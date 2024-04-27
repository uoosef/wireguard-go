package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/netip"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	_ "net/http/pprof"

	"github.com/adrg/xdg"
	"github.com/bepass-org/warp-plus/app"
	"github.com/bepass-org/warp-plus/warp"
	"github.com/bepass-org/warp-plus/wiresocks"

	"github.com/carlmjohnson/versioninfo"
	"github.com/peterbourgon/ff/v4"
	"github.com/peterbourgon/ff/v4/ffhelp"
	"github.com/peterbourgon/ff/v4/ffjson"
)

const appName = "warp-plus"

var psiphonCountries = []string{
	"AT",
	"BE",
	"BG",
	"BR",
	"CA",
	"CH",
	"CZ",
	"DE",
	"DK",
	"EE",
	"ES",
	"FI",
	"FR",
	"GB",
	"HR",
	"HU",
	"IE",
	"IN",
	"IT",
	"JP",
	"LV",
	"NL",
	"NO",
	"PL",
	"PT",
	"RO",
	"RS",
	"SE",
	"SG",
	"SK",
	"UA",
	"US",
}

var version string = ""

func main() {
	fs := ff.NewFlagSet(appName)
	var (
		v4       = fs.BoolShort('4', "only use IPv4 for random warp endpoint")
		v6       = fs.BoolShort('6', "only use IPv6 for random warp endpoint")
		verbose  = fs.Bool('v', "verbose", "enable verbose logging")
		bind     = fs.String('b', "bind", "127.0.0.1:8086", "socks bind address")
		endpoint = fs.String('e', "endpoint", "", "warp endpoint")
		key      = fs.String('k', "key", "", "warp key")
		gool     = fs.BoolLong("gool", "enable gool mode (warp in warp)")
		psiphon  = fs.BoolLong("cfon", "enable psiphon mode (must provide country as well)")
		country  = fs.StringEnumLong("country", fmt.Sprintf("psiphon country code (valid values: %s)", psiphonCountries), psiphonCountries...)
		scan     = fs.BoolLong("scan", "enable warp scanning")
		rtt      = fs.DurationLong("rtt", 1000*time.Millisecond, "scanner rtt limit")
		cacheDir = fs.StringLong("cache-dir", "", "directory to store generated profiles")
		tun      = fs.BoolLong("tun", "enable tun interface")
		fwmark   = fs.UintLong("fwmark", 0x1375, "set linux firewall mark for tun mode")
		_        = fs.String('c', "config", "", "path to config file")
		verFlag  = fs.BoolLong("version", "displays version number")
	)

	err := ff.Parse(
		fs,
		os.Args[1:],
		ff.WithConfigFileFlag("config"),
		ff.WithConfigFileParser(ffjson.Parse),
	)
	switch {
	case errors.Is(err, ff.ErrHelp):
		fmt.Fprintf(os.Stderr, "%s\n", ffhelp.Flags(fs))
		os.Exit(0)
	case err != nil:
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if *verFlag {
		if version == "" {
			version = versioninfo.Short()
		}
		fmt.Fprintf(os.Stderr, "%s\n", version)
		os.Exit(0)
	}

	l := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	if *verbose {
		l = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	if *psiphon && *gool {
		fatal(l, errors.New("can't use cfon and gool at the same time"))
	}

	if *v4 && *v6 {
		fatal(l, errors.New("can't force v4 and v6 at the same time"))
	}

	if !*v4 && !*v6 {
		*v4, *v6 = true, true
	}

	bindAddrPort, err := netip.ParseAddrPort(*bind)
	if err != nil {
		fatal(l, fmt.Errorf("invalid bind address: %w", err))
	}

	opts := app.WarpOptions{
		Bind:     bindAddrPort,
		Endpoint: *endpoint,
		License:  *key,
		Gool:     *gool,
	}

	switch {
	case *cacheDir != "":
		opts.CacheDir = *cacheDir
	case xdg.CacheHome != "":
		opts.CacheDir = path.Join(xdg.CacheHome, appName)
	case os.Getenv("HOME") != "":
		opts.CacheDir = path.Join(os.Getenv("HOME"), ".cache", appName)
	default:
		opts.CacheDir = "warp_plus_cache"
	}

	if *psiphon {
		l.Info("psiphon mode enabled", "country", *country)
		opts.Psiphon = &app.PsiphonOptions{Country: *country}
	}

	if *scan {
		l.Info("scanner mode enabled", "max-rtt", rtt)
		opts.Scan = &wiresocks.ScanOptions{V4: *v4, V6: *v6, MaxRTT: *rtt}
	}

	if *tun {
		l.Info("tun mode enabled")
		opts.Tun = &app.TunOptions{FwMark: uint32(*fwmark)}
	}

	// If the endpoint is not set, choose a random warp endpoint
	if opts.Endpoint == "" {
		addrPort, err := warp.RandomWarpEndpoint(*v4, *v6)
		if err != nil {
			fatal(l, err)
		}
		opts.Endpoint = addrPort.String()
	}

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	go func() {
		if err := app.RunWarp(ctx, l, opts); err != nil {
			fatal(l, err)
		}
	}()

	<-ctx.Done()
}

func fatal(l *slog.Logger, err error) {
	l.Error(err.Error())
	os.Exit(1)
}
