package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/bepass-org/wireguard-go/app"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func usage() {
	log.Println("Usage: wiresocks [-v] [-b addr:port] [-l license] <config file path>")
	flag.PrintDefaults()
}

var currentVersion = "v1.0.3-rc2"

func main() {
	var (
		verbose        = flag.Bool("v", false, "verbose")
		version        = flag.Bool("version", false, "Show current version")
		bindAddress    = flag.String("b", "127.0.0.1:8086", "socks bind address")
		endpoint       = flag.String("e", "notset", "warp clean ip")
		license        = flag.String("k", "notset", "license key")
		country        = flag.String("country", "", "psiphon country code in ISO 3166-1 alpha-2 format")
		psiphonEnabled = flag.Bool("cfon", false, "enable psiphonEnabled over warp")
		gool           = flag.Bool("gool", false, "enable warp gooling")
		scan           = flag.Bool("scan", false, "enable warp scanner(experimental)")
		rtt            = flag.Int("rtt", 1000, "scanner rtt threshold, default 1000")
	)

	flag.Usage = usage
	flag.Parse()

	if *version {
		fmt.Printf("Current Version: %v \n", currentVersion)
		os.Exit(0)
	}

	sigchan := make(chan os.Signal)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		err := app.RunWarp(*psiphonEnabled, *gool, *scan, *verbose, *country, *bindAddress, *endpoint, *license, ctx, *rtt)
		if err != nil {
			log.Fatal(err)
		}
	}()

	<-sigchan
	cancel()
}
