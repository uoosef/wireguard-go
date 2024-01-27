package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/uoosef/wireguard-go/device"
	"github.com/uoosef/wireguard-go/psiphon"
	"github.com/uoosef/wireguard-go/warp"
	"github.com/uoosef/wireguard-go/wiresocks"
)

func usage() {
	log.Println("Usage: wiresocks [-v] [-b addr:port] [-l license] <config file path>")
	flag.PrintDefaults()
}

func main() {
	var (
		verbose        = flag.Bool("v", false, "verbose")
		bindAddress    = flag.String("b", "127.0.0.1:8086", "socks bind address")
		configFile     = flag.String("c", "./wgcf-profile.ini", "ini config file path")
		endpoint       = flag.String("e", "notset", "warp clean ip")
		license        = flag.String("k", "notset", "license key")
		hostbind       = flag.Bool("h", false, "bind to 0.0.0.0 instead of local host")
		country        = flag.String("country", "", "psiphon country code in ISO 3166-1 alpha-2 format")
		psiphonEnabled = flag.Bool("cfon", false, "enable psiphonEnabled over warp")
		port           = flag.Int("p", 0, "port for wiresocks to listen on")
		pbind          = "127.0.0.1:8086"
	)

	flag.Usage = usage
	flag.Parse()

	wiresocks.Verbose = *verbose

	if *psiphonEnabled {
		pbind = *bindAddress
		randomBind, err := findFreePort(*hostbind, *port)
		if err != nil {
			log.Fatal(err)
		}
		bindAddress = &randomBind
	}

	// check if wgcf-profile.conf exists
	if !warp.CheckProfileExists(*license) {
		if *license == "notset" {
			*license = ""
		}
		err := warp.LoadOrCreateIdentity(*license)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
	}

	conf, err := wiresocks.ParseConfig(*configFile, *endpoint)
	if err != nil {
		log.Fatal(err)
	}

	logLevel := device.LogLevelVerbose
	if !*verbose {
		logLevel = device.LogLevelSilent
	}

	// Setup channel to listen for interrupt signal (Ctrl+C)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	tnet, err := wiresocks.StartWireguard(conf.Device, logLevel)
	if err != nil {
		log.Fatal(err)
	}

	go wiresocks.StartProxy(tnet, *bindAddress)

	if *psiphonEnabled {
		psiphon.RunPsiphon(*bindAddress, pbind, *country)
	} else {
		log.Println("Wiresocks started successfully")
	}

	// Wait for interrupt signal
	<-sigChan

	log.Println("Bye!")
}

func findFreePort(hostflag bool, port int) (string, error) {
	// Listen on TCP port 0, which tells the OS to pick a free port.

	var listener net.Listener
	var err error
	var address string = "127.0.0.1:" + strconv.Itoa(port)

	if hostflag {
		address = "0.0.0.0:" + strconv.Itoa(port)
	}

	listener, err = net.Listen("tcp", address)

	if err != nil {
		return "", err // Return error if unable to listen on a port
	}
	defer listener.Close() // Ensure the listener is closed when the function returns

	// Get the port from the listener's address
	addr := listener.Addr().String()

	return addr, nil
}
