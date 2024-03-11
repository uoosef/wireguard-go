package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/netip"
	"os"
	"path/filepath"
	"time"

	"github.com/bepass-org/wireguard-go/psiphon"
	"github.com/bepass-org/wireguard-go/warp"
	"github.com/bepass-org/wireguard-go/wiresocks"
)

type WarpOptions struct {
	LogLevel string
	Bind     netip.AddrPort
	Endpoint string
	License  string
	Psiphon  *PsiphonOptions
	Gool     bool
	Scan     *ScanOptions
}

type PsiphonOptions struct {
	Country string
}

type ScanOptions struct {
	MaxRTT time.Duration
}

func RunWarp(ctx context.Context, opts WarpOptions) error {
	if opts.Psiphon != nil && opts.Gool {
		return errors.New("can't use psiphon and gool at the same time")
	}

	if opts.Psiphon != nil && opts.Psiphon.Country == "" {
		return errors.New("must provide country for psiphon")
	}

	// create necessary file structures
	if err := makeDirs(); err != nil {
		return err
	}
	log.Println("'primary' and 'secondary' directories are ready")

	// Change the current working directory to 'stuff'
	if err := os.Chdir("stuff"); err != nil {
		return fmt.Errorf("error changing to 'stuff' directory: %w", err)
	}
	log.Println("Changed working directory to 'stuff'")

	// create identities
	if err := createPrimaryAndSecondaryIdentities(opts.License); err != nil {
		return err
	}

	// Decide Working Scenario
	endpoints := []string{opts.Endpoint, opts.Endpoint}

	if opts.Scan != nil {
		res, err := wiresocks.RunScan(ctx, opts.Scan.MaxRTT)
		if err != nil {
			return err
		}
		endpoints = make([]string, len(res))
		for i := 0; i < len(res); i++ {
			endpoints[i] = res[i].String()
		}
	}

	var warpErr error
	switch {
	case opts.Psiphon != nil:
		// run primary warp on a random tcp port and run psiphon on bind address
		warpErr = runWarpWithPsiphon(ctx, opts.Bind, endpoints, opts.Psiphon.Country, opts.LogLevel == "debug")
	case opts.Gool:
		// run warp in warp
		warpErr = runWarpInWarp(ctx, opts.Bind, endpoints, opts.LogLevel == "debug")
	default:
		// just run primary warp on bindAddress
		_, _, warpErr = runWarp(ctx, opts.Bind, endpoints, "./primary/wgcf-profile.ini", opts.LogLevel == "debug", true, true)
	}

	return warpErr
}

func runWarp(ctx context.Context, bind netip.AddrPort, endpoints []string, confPath string, verbose, startProxy bool, trick bool) (*wiresocks.VirtualTun, int, error) {
	conf, err := wiresocks.ParseConfig(confPath, endpoints[0])
	if err != nil {
		log.Println(err)
		return nil, 0, err
	}

	if trick {
		conf.Device.Trick = trick
	}

	tnet, err := wiresocks.StartWireguard(ctx, conf.Device, verbose)
	if err != nil {
		log.Println(err)
		return nil, 0, err
	}

	if startProxy {
		tnet.StartProxy(bind)
	}

	return tnet, conf.Device.MTU, nil
}

func runWarpWithPsiphon(ctx context.Context, bind netip.AddrPort, endpoints []string, country string, verbose bool) error {
	// make a random bind address for warp
	warpBindAddress, err := findFreePort("tcp")
	if err != nil {
		log.Println("There are no free tcp ports on Device!")
		return err
	}

	_, _, err = runWarp(ctx, warpBindAddress, endpoints, "./primary/wgcf-profile.ini", verbose, true, true)
	if err != nil {
		return err
	}

	// run psiphon
	err = psiphon.RunPsiphon(warpBindAddress.String(), bind.String(), country, ctx)
	if err != nil {
		log.Printf("unable to run psiphon %v", err)
		return fmt.Errorf("unable to run psiphon %w", err)
	}

	log.Printf("Serving on %s", bind)

	return nil
}

func runWarpInWarp(ctx context.Context, bind netip.AddrPort, endpoints []string, verbose bool) error {
	// run secondary warp
	vTUN, mtu, err := runWarp(ctx, netip.AddrPort{}, endpoints, "./secondary/wgcf-profile.ini", verbose, false, true)
	if err != nil {
		return err
	}

	// run virtual endpoint
	virtualEndpointBindAddress, err := findFreePort("udp")
	if err != nil {
		log.Println("There are no free udp ports on Device!")
		return err
	}
	addr := endpoints[1]
	err = wiresocks.NewVtunUDPForwarder(virtualEndpointBindAddress.String(), addr, vTUN, mtu+100, ctx)
	if err != nil {
		log.Println(err)
		return err
	}

	// run primary warp
	_, _, err = runWarp(ctx, bind, []string{virtualEndpointBindAddress.String()}, "./primary/wgcf-profile.ini", verbose, true, false)
	if err != nil {
		return err
	}
	return nil
}

func findFreePort(network string) (netip.AddrPort, error) {
	if network == "udp" {
		addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
		if err != nil {
			return netip.AddrPort{}, err
		}

		conn, err := net.ListenUDP("udp", addr)
		if err != nil {
			return netip.AddrPort{}, err
		}
		defer conn.Close()

		return netip.MustParseAddrPort(conn.LocalAddr().String()), nil
	}
	// Listen on TCP port 0, which tells the OS to pick a free port.
	listener, err := net.Listen(network, "127.0.0.1:0")
	if err != nil {
		return netip.AddrPort{}, err // Return error if unable to listen on a port
	}
	defer listener.Close() // Ensure the listener is closed when the function returns

	// Get the port from the listener's address
	return netip.MustParseAddrPort(listener.Addr().String()), nil
}

func createPrimaryAndSecondaryIdentities(license string) error {
	// make primary identity
	warp.UpdatePath("./primary")
	if !warp.CheckProfileExists(license) {
		err := warp.LoadOrCreateIdentity(license)
		if err != nil {
			log.Printf("error: %v", err)
			return err
		}
	}
	// make secondary
	warp.UpdatePath("./secondary")
	if !warp.CheckProfileExists(license) {
		err := warp.LoadOrCreateIdentity(license)
		if err != nil {
			log.Printf("error: %v", err)
			return err
		}
	}
	return nil
}

func makeDirs() error {
	stuffDir := "stuff"
	primaryDir := "primary"
	secondaryDir := "secondary"

	// Check if 'stuff' directory exists, if not create it
	if _, err := os.Stat(stuffDir); os.IsNotExist(err) {
		if err := os.Mkdir(stuffDir, 0o755); err != nil {
			return fmt.Errorf("error creating 'stuff' directory: %w", err)
		}
	}

	// Create 'primary' and 'secondary' directories if they don't exist
	for _, dir := range []string{primaryDir, secondaryDir} {
		if _, err := os.Stat(filepath.Join(stuffDir, dir)); os.IsNotExist(err) {
			if err := os.Mkdir(filepath.Join(stuffDir, dir), 0o755); err != nil {
				return fmt.Errorf("error creating '%s' directory: %w", dir, err)
			}
		}
	}

	return nil
}
