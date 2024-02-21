package app

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/bepass-org/wireguard-go/psiphon"
	"github.com/bepass-org/wireguard-go/warp"
	"github.com/bepass-org/wireguard-go/wiresocks"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"
)

func RunWarp(psiphonEnabled, gool, scan, verbose bool, country, bindAddress, endpoint, license string, ctx context.Context, rtThreshold int) error {
	// check if user input is not correct
	if (psiphonEnabled && gool) || (!psiphonEnabled && country != "") {
		log.Println("Wrong combination of flags!")
		flag.Usage()
		return errors.New("wrong command")
	}

	//create necessary file structures
	if err := makeDirs(); err != nil {
		return err
	}

	// Change the current working directory to 'stuff'
	if err := os.Chdir("stuff"); err != nil {
		log.Printf("Error changing to 'stuff' directory: %v\n", err)
		return fmt.Errorf("Error changing to 'stuff' directory: %v\n", err)
	}
	log.Println("Changed working directory to 'stuff'")
	defer func() {
		// back where you where
		if err := os.Chdir(".."); err != nil {
			log.Fatal("Error changing to 'main' directory:", err)
		}
	}()

	//create identities
	if err := createPrimaryAndSecondaryIdentities(license); err != nil {
		return err
	}

	//Decide Working Scenario
	endpoints := []string{endpoint, endpoint}

	if scan {
		var err error
		endpoints, err = wiresocks.RunScan(&ctx, rtThreshold)
		if err != nil {
			return err
		}
	}

	if !psiphonEnabled && !gool {
		// just run primary warp on bindAddress
		_, _, err := runWarp(bindAddress, endpoints, "./primary/wgcf-profile.ini", verbose, true, ctx, true)
		return err
	} else if psiphonEnabled && !gool {
		// run primary warp on a random tcp port and run psiphon on bind address
		return runWarpWithPsiphon(bindAddress, endpoints, country, verbose, ctx)
	} else if !psiphonEnabled && gool {
		// run warp in warp
		return runWarpInWarp(bindAddress, endpoints, verbose, ctx)
	}

	return fmt.Errorf("unknown error, it seems core related issue")
}

func runWarp(bindAddress string, endpoints []string, confPath string, verbose, startProxy bool, ctx context.Context, showServing bool) (*wiresocks.VirtualTun, int, error) {
	conf, err := wiresocks.ParseConfig(confPath, endpoints[0])
	if err != nil {
		log.Println(err)
		return nil, 0, err
	}

	tnet, err := wiresocks.StartWireguard(conf.Device, verbose, ctx)
	if err != nil {
		log.Println(err)
		return nil, 0, err
	}

	if startProxy {
		tnet.StartProxy(bindAddress)
	}

	if showServing {
		log.Printf("Serving on %s\n", bindAddress)
	}

	return tnet, conf.Device.MTU, nil
}

func runWarpWithPsiphon(bindAddress string, endpoints []string, country string, verbose bool, ctx context.Context) error {
	// make a random bind address for warp
	warpBindAddress, err := findFreePort("tcp")
	if err != nil {
		log.Println("There are no free tcp ports on Device!")
		return err
	}

	_, _, err = runWarp(warpBindAddress, endpoints, "./primary/wgcf-profile.ini", verbose, true, ctx, false)
	if err != nil {
		return err
	}

	// run psiphon
	err = psiphon.RunPsiphon(warpBindAddress, bindAddress, country, ctx)
	if err != nil {
		log.Printf("unable to run psiphon %v", err)
		return fmt.Errorf("unable to run psiphon %v", err)
	}

	log.Printf("Serving on %s\n", bindAddress)

	return nil
}

func runWarpInWarp(bindAddress string, endpoints []string, verbose bool, ctx context.Context) error {
	// run secondary warp
	vTUN, mtu, err := runWarp("", endpoints, "./secondary/wgcf-profile.ini", verbose, false, ctx, false)
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
	if addr == "notset" {
		addr, _ = wiresocks.ResolveIPPAndPort("engage.cloudflareclient.com:2408")
	}
	err = wiresocks.NewVtunUDPForwarder(virtualEndpointBindAddress, addr, vTUN, mtu+100, ctx)
	if err != nil {
		log.Println(err)
		return err
	}

	// run primary warp
	_, _, err = runWarp(bindAddress, []string{virtualEndpointBindAddress}, "./primary/wgcf-profile.ini", verbose, true, ctx, true)
	if err != nil {
		return err
	}
	return nil
}

func findFreePort(network string) (string, error) {
	if network == "udp" {
		addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
		if err != nil {
			return "", err
		}

		conn, err := net.ListenUDP("udp", addr)
		if err != nil {
			return "", err
		}
		defer conn.Close()

		return conn.LocalAddr().(*net.UDPAddr).String(), nil
	}
	// Listen on TCP port 0, which tells the OS to pick a free port.
	listener, err := net.Listen(network, "127.0.0.1:0")
	if err != nil {
		return "", err // Return error if unable to listen on a port
	}
	defer listener.Close() // Ensure the listener is closed when the function returns

	// Get the port from the listener's address
	addr := listener.Addr().String()

	return addr, nil
}

func createPrimaryAndSecondaryIdentities(license string) error {
	// make primary identity
	_license := license
	if _license == "notset" {
		_license = ""
	}
	warp.UpdatePath("./primary")
	if !warp.CheckProfileExists(license) {
		err := warp.LoadOrCreateIdentity(_license)
		if err != nil {
			log.Printf("error: %v", err)
			return fmt.Errorf("error: %v", err)
		}
	}
	// make secondary
	warp.UpdatePath("./secondary")
	if !warp.CheckProfileExists(license) {
		err := warp.LoadOrCreateIdentity(_license)
		if err != nil {
			log.Printf("error: %v", err)
			return fmt.Errorf("error: %v", err)
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
		log.Println("'stuff' directory does not exist, creating it...")
		if err := os.Mkdir(stuffDir, 0755); err != nil {
			log.Println("Error creating 'stuff' directory:", err)
			return errors.New("Error creating 'stuff' directory:" + err.Error())
		}
	}

	// Create 'primary' and 'secondary' directories if they don't exist
	for _, dir := range []string{primaryDir, secondaryDir} {
		if _, err := os.Stat(filepath.Join(stuffDir, dir)); os.IsNotExist(err) {
			log.Printf("Creating '%s' directory...\n", dir)
			if err := os.Mkdir(filepath.Join(stuffDir, dir), 0755); err != nil {
				log.Printf("Error creating '%s' directory: %v\n", dir, err)
				return fmt.Errorf("Error creating '%s' directory: %v\n", dir, err)
			}
		}
	}
	log.Println("'primary' and 'secondary' directories are ready")
	return nil
}
func isPortOpen(address string, timeout time.Duration) bool {
	// Try to establish a connection
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return false
	}
	defer conn.Close()

	return true
}

func waitForPortToGetsOpenOrTimeout(addressToCheck string) {
	timeout := 5 * time.Second
	checkInterval := 500 * time.Millisecond

	// Set a deadline for when to stop checking
	deadline := time.Now().Add(timeout)

	for {
		if time.Now().After(deadline) {
			log.Fatalf("Timeout reached, port %s is not open", addressToCheck)
		}

		if isPortOpen(addressToCheck, checkInterval) {
			log.Printf("Port %s is now open", addressToCheck)
			break
		}

		time.Sleep(checkInterval)
	}
}
