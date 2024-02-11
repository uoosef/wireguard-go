package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/bepass-org/wireguard-go/app"
)

type Flags struct {
	Verbose        bool
	BindAddress    string
	Endpoint       string
	License        string
	Country        string
	PsiphonEnabled bool
	Gool           bool
	Scan           bool
}

var validFlags = map[string]bool{
	"-v":       true,
	"-b":       true,
	"-e":       true,
	"-k":       true,
	"-country": true,
	"-cfon":    true,
	"-gool":    true,
	"-scan":    true,
}

func newFlags() *Flags {
	return &Flags{}
}

func (f *Flags) setup() {
	flag.BoolVar(&f.Verbose, "v", false, "verbose")
	flag.StringVar(&f.BindAddress, "b", "127.0.0.1:8086", "socks bind address")
	flag.StringVar(&f.Endpoint, "e", "notset", "warp clean IP")
	flag.StringVar(&f.License, "k", "notset", "license key")
	flag.StringVar(&f.Country, "country", "", "psiphon country code in ISO 3166-1 alpha-2 format")
	flag.BoolVar(&f.PsiphonEnabled, "cfon", false, "enable Psiphon over warp")
	flag.BoolVar(&f.Gool, "gool", false, "enable warp gooling")
	flag.BoolVar(&f.Scan, "scan", false, "enable warp scanner(experimental)")

	flag.Usage = usage
	flag.Parse()
}

var validCountryCodes = map[string]bool{
	"AT": true,
	"BE": true,
	"BG": true,
	"BR": true,
	"CA": true,
	"CH": true,
	"CZ": true,
	"DE": true,
	"DK": true,
	"EE": true,
	"ES": true,
	"FI": true,
	"FR": true,
	"GB": true,
	"HU": true,
	"IE": true,
	"IN": true,
	"IT": true,
	"JP": true,
	"LV": true,
	"NL": true,
	"NO": true,
	"PL": true,
	"RO": true,
	"RS": true,
	"SE": true,
	"SG": true,
	"SK": true,
	"UA": true,
	"US": true,
}

func usage() {
	log.Println("./warp-plus-go [-v] [-b addr:port] [-c config-file-path] [-e warp-ip] [-k license-key] [-country country-code] [-cfon] [-gool]")
	flag.PrintDefaults()
}

func validateFlags(f *Flags) error {
	if _, err := net.ResolveTCPAddr("tcp", f.BindAddress); err != nil {
		return fmt.Errorf("invalid bindAddress format: %s", f.BindAddress)
	}

	if ip := net.ParseIP(f.Endpoint); ip == nil {
		return fmt.Errorf("invalid warp clean IP: %s", f.Endpoint)
	}

	if f.PsiphonEnabled && f.Country == "" {
		return fmt.Errorf("if Psiphon is enabled, country code must be provided")
	}

	if !validCountryCodes[f.Country] {
		validCountries := make([]string, 0, len(validCountryCodes))

		for code, _ := range validCountryCodes {
			validCountries = append(validCountries, code)
		}

		return fmt.Errorf("invalid country code: %s. Valid country codes: $s", f.Country, validCountries)
	}

	return nil
}

func main() {
	// Check for unexpected flags
	for _, arg := range os.Args[1:] {
		if !validFlags[arg] {
			log.Fatalf("Invalid flag: %s", arg)
		}
	}

	flags := newFlags()
	flags.setup()

	if err := validateFlags(flags); err != nil {
		log.Fatalf("Validatrion error: %v", err)
	}

	sigchan := make(chan os.Signal)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		err := app.RunWarp(flags.PsiphonEnabled, flags.Gool, flags.Scan, flags.Verbose, flags.Country, flags.BindAddress, flags.Endpoint, flags.License, ctx)
		if err != nil {
			log.Fatal(err)
		}
	}()

	<-sigchan
	cancel()
}
