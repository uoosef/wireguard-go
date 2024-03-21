package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/netip"
	"os"
	"path/filepath"

	"github.com/bepass-org/warp-plus/psiphon"
	"github.com/bepass-org/warp-plus/warp"
	"github.com/bepass-org/warp-plus/wiresocks"
)

const singleMTU = 1400
const doubleMTU = 1320

type WarpOptions struct {
	Bind     netip.AddrPort
	Endpoint string
	License  string
	Psiphon  *PsiphonOptions
	Gool     bool
	Scan     *wiresocks.ScanOptions
}

type PsiphonOptions struct {
	Country string
}

func RunWarp(ctx context.Context, l *slog.Logger, opts WarpOptions) error {
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
	l.Debug("'primary' and 'secondary' directories are ready")

	// Change the current working directory to 'stuff'
	if err := os.Chdir("stuff"); err != nil {
		return fmt.Errorf("error changing to 'stuff' directory: %w", err)
	}
	l.Debug("Changed working directory to 'stuff'")

	// create identities
	if err := createPrimaryAndSecondaryIdentities(l.With("subsystem", "warp/account"), opts.License); err != nil {
		return err
	}

	// Decide Working Scenario
	endpoints := []string{opts.Endpoint, opts.Endpoint}

	if opts.Scan != nil {
		res, err := wiresocks.RunScan(ctx, l, *opts.Scan)
		if err != nil {
			return err
		}

		l.Info("scan results", "endpoints", res)

		endpoints = make([]string, len(res))
		for i := 0; i < len(res); i++ {
			endpoints[i] = res[i].AddrPort.String()
		}
	}
	l.Info("using warp endpoints", "endpoints", endpoints)

	var warpErr error
	switch {
	case opts.Psiphon != nil:
		l.Info("running in Psiphon (cfon) mode")
		// run primary warp on a random tcp port and run psiphon on bind address
		warpErr = runWarpWithPsiphon(ctx, l, opts.Bind, endpoints[0], opts.Psiphon.Country)
	case opts.Gool:
		l.Info("running in warp-in-warp (gool) mode")
		// run warp in warp
		warpErr = runWarpInWarp(ctx, l, opts.Bind, endpoints)
	default:
		l.Info("running in normal warp mode")
		// just run primary warp on bindAddress
		warpErr = runWarp(ctx, l, opts.Bind, endpoints[0])
	}

	return warpErr
}

func runWarp(ctx context.Context, l *slog.Logger, bind netip.AddrPort, endpoint string) error {
	conf, err := wiresocks.ParseConfig("./primary/wgcf-profile.ini", endpoint)
	if err != nil {
		return err
	}
	conf.Interface.MTU = singleMTU

	for i, peer := range conf.Peers {
		peer.Trick = true
		peer.KeepAlive = 3
		conf.Peers[i] = peer
	}

	tnet, err := wiresocks.StartWireguard(ctx, l, conf)
	if err != nil {
		return err
	}

	tnet.StartProxy(bind)
	l.Info("serving proxy", "address", bind)

	return nil
}

func runWarpWithPsiphon(ctx context.Context, l *slog.Logger, bind netip.AddrPort, endpoint string, country string) error {
	conf, err := wiresocks.ParseConfig("./primary/wgcf-profile.ini", endpoint)
	if err != nil {
		return err
	}
	conf.Interface.MTU = singleMTU

	for i, peer := range conf.Peers {
		peer.Trick = true
		peer.KeepAlive = 3
		conf.Peers[i] = peer
	}

	tnet, err := wiresocks.StartWireguard(ctx, l, conf)
	if err != nil {
		return err
	}

	warpBind, err := tnet.StartProxy(netip.MustParseAddrPort("127.0.0.1:0"))
	if err != nil {
		return err
	}

	// run psiphon
	err = psiphon.RunPsiphon(ctx, l.With("subsystem", "psiphon"), warpBind.String(), bind.String(), country)
	if err != nil {
		return fmt.Errorf("unable to run psiphon %w", err)
	}

	l.Info("serving proxy", "address", bind)

	return nil
}

func runWarpInWarp(ctx context.Context, l *slog.Logger, bind netip.AddrPort, endpoints []string) error {
	// Run outer warp
	conf, err := wiresocks.ParseConfig("./primary/wgcf-profile.ini", endpoints[0])
	if err != nil {
		return err
	}
	conf.Interface.MTU = singleMTU

	for i, peer := range conf.Peers {
		peer.Trick = true
		peer.KeepAlive = 3
		conf.Peers[i] = peer
	}

	tnet, err := wiresocks.StartWireguard(ctx, l.With("gool", "outer"), conf)
	if err != nil {
		return err
	}

	// Create a UDP port forward between localhost and the remote endpoint
	addr, err := wiresocks.NewVtunUDPForwarder(ctx, netip.MustParseAddrPort("127.0.0.1:0"), endpoints[1], tnet, singleMTU)
	if err != nil {
		return err
	}

	// Run inner warp
	conf, err = wiresocks.ParseConfig("./secondary/wgcf-profile.ini", addr.String())
	if err != nil {
		return err
	}
	conf.Interface.MTU = doubleMTU

	for i, peer := range conf.Peers {
		peer.KeepAlive = 10
		conf.Peers[i] = peer
	}

	tnet, err = wiresocks.StartWireguard(ctx, l.With("gool", "inner"), conf)
	if err != nil {
		return err
	}

	_, err = tnet.StartProxy(bind)
	if err != nil {
		return err
	}

	l.Info("serving proxy", "address", bind)
	return nil
}

func createPrimaryAndSecondaryIdentities(l *slog.Logger, license string) error {
	// make primary identity
	warp.UpdatePath("./primary")
	if !warp.CheckProfileExists(license) {
		err := warp.LoadOrCreateIdentity(l, license)
		if err != nil {
			return err
		}
	}
	// make secondary
	warp.UpdatePath("./secondary")
	if !warp.CheckProfileExists(license) {
		err := warp.LoadOrCreateIdentity(l, license)
		if err != nil {
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
