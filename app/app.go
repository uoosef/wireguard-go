package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/netip"
	"path"

	"github.com/bepass-org/warp-plus/iputils"
	"github.com/bepass-org/warp-plus/psiphon"
	"github.com/bepass-org/warp-plus/warp"
	"github.com/bepass-org/warp-plus/wiresocks"
	"github.com/go-ini/ini"
)

const singleMTU = 1330
const doubleMTU = 1280 // minimum mtu for IPv6, may cause frag reassembly somewhere
const connTestEndpoint = "http://1.1.1.1:80/"

type WarpOptions struct {
	Bind            netip.AddrPort
	Endpoint        string
	License         string
	DnsAddr         netip.Addr
	Psiphon         *PsiphonOptions
	Gool            bool
	Scan            *wiresocks.ScanOptions
	CacheDir        string
	Tun             bool
	FwMark          uint32
	WireguardConfig string
}

type PsiphonOptions struct {
	Country string
}

func RunWarp(ctx context.Context, l *slog.Logger, opts WarpOptions) error {
	if opts.WireguardConfig != "" {
		if err := runWireguard(ctx, l, opts); err != nil {
			return err
		}

		return nil
	}

	if opts.Psiphon != nil && opts.Gool {
		return errors.New("can't use psiphon and gool at the same time")
	}

	if opts.Psiphon != nil && opts.Psiphon.Country == "" {
		return errors.New("must provide country for psiphon")
	}

	if opts.Psiphon != nil && opts.Tun {
		return errors.New("can't use psiphon and tun at the same time")
	}

	// create identities
	if err := createPrimaryAndSecondaryIdentities(l.With("subsystem", "warp/account"), opts); err != nil {
		return err
	}

	// Decide Working Scenario
	endpoints := []string{opts.Endpoint, opts.Endpoint}

	if opts.Scan != nil {
		cfg, err := ini.Load(path.Join(opts.CacheDir, "primary", "wgcf-profile.ini"))
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		// Reading the private key from the 'Interface' section
		opts.Scan.PrivateKey = cfg.Section("Interface").Key("PrivateKey").String()

		// Reading the public key from the 'Peer' section
		opts.Scan.PublicKey = cfg.Section("Peer").Key("PublicKey").String()

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
		warpErr = runWarpWithPsiphon(ctx, l, opts, endpoints[0])
	case opts.Gool:
		l.Info("running in warp-in-warp (gool) mode")
		// run warp in warp
		warpErr = runWarpInWarp(ctx, l, opts, endpoints)
	default:
		l.Info("running in normal warp mode")
		// just run primary warp on bindAddress
		warpErr = runWarp(ctx, l, opts, endpoints[0])
	}

	return warpErr
}

func runWireguard(ctx context.Context, l *slog.Logger, opts WarpOptions) error {
	conf, err := wiresocks.ParseConfig(opts.WireguardConfig)
	if err != nil {
		return err
	}

	// Set up MTU
	conf.Interface.MTU = singleMTU
	// Set up DNS Address
	conf.Interface.DNS = []netip.Addr{opts.DnsAddr}

	// Enable trick and keepalive on all peers in config
	for i, peer := range conf.Peers {
		peer.Trick = true
		peer.KeepAlive = 3

		// Try resolving if the endpoint is a domain
		addr, err := iputils.ParseResolveAddressPort(peer.Endpoint, false)
		if err == nil {
			peer.Endpoint = addr.String()
		}

		conf.Peers[i] = peer
	}

	if opts.Tun {
		// Create a new tun interface
		tunDev, err := newNormalTun([]netip.Addr{opts.DnsAddr})
		if err != nil {
			return err
		}

		// Establish wireguard tunnel on tun interface
		if err := establishWireguard(l, conf, tunDev, true, opts.FwMark); err != nil {
			return err
		}
		l.Info("serving tun", "interface", "warp0")
		return nil
	}

	// Create userspace tun network stack
	tunDev, tnet, err := newUsermodeTun(conf)
	if err != nil {
		return err
	}

	// Establish wireguard on userspace stack
	if err := establishWireguard(l, conf, tunDev, false, opts.FwMark); err != nil {
		return err
	}

	// Test wireguard connectivity
	if err := usermodeTunTest(ctx, l, tnet); err != nil {
		return err
	}

	// Run a proxy on the userspace stack
	_, err = wiresocks.StartProxy(ctx, l, tnet, opts.Bind)
	if err != nil {
		return err
	}

	l.Info("serving proxy", "address", opts.Bind)

	return nil
}

func runWarp(ctx context.Context, l *slog.Logger, opts WarpOptions, endpoint string) error {
	// Set up primary/outer warp config
	conf, err := wiresocks.ParseConfig(path.Join(opts.CacheDir, "primary", "wgcf-profile.ini"))
	if err != nil {
		return err
	}

	// Set up MTU
	conf.Interface.MTU = singleMTU
	// Set up DNS Address
	conf.Interface.DNS = []netip.Addr{opts.DnsAddr}

	// Enable trick and keepalive on all peers in config
	for i, peer := range conf.Peers {
		peer.Endpoint = endpoint
		peer.Trick = true
		peer.KeepAlive = 3
		conf.Peers[i] = peer
	}

	if opts.Tun {
		// Create a new tun interface
		tunDev, err := newNormalTun([]netip.Addr{opts.DnsAddr})
		if err != nil {
			return err
		}

		// Establish wireguard tunnel on tun interface
		if err := establishWireguard(l, conf, tunDev, true, opts.FwMark); err != nil {
			return err
		}
		l.Info("serving tun", "interface", "warp0")
		return nil
	}

	// Create userspace tun network stack
	tunDev, tnet, err := newUsermodeTun(conf)
	if err != nil {
		return err
	}

	// Establish wireguard on userspace stack
	if err := establishWireguard(l, conf, tunDev, false, opts.FwMark); err != nil {
		return err
	}

	// Test wireguard connectivity
	if err := usermodeTunTest(ctx, l, tnet); err != nil {
		return err
	}

	// Run a proxy on the userspace stack
	_, err = wiresocks.StartProxy(ctx, l, tnet, opts.Bind)
	if err != nil {
		return err
	}

	l.Info("serving proxy", "address", opts.Bind)
	return nil
}

func runWarpInWarp(ctx context.Context, l *slog.Logger, opts WarpOptions, endpoints []string) error {
	// Set up primary/outer warp config
	conf, err := wiresocks.ParseConfig(path.Join(opts.CacheDir, "primary", "wgcf-profile.ini"))
	if err != nil {
		return err
	}

	// Set up MTU
	conf.Interface.MTU = singleMTU
	// Set up DNS Address
	conf.Interface.DNS = []netip.Addr{opts.DnsAddr}

	// Enable trick and keepalive on all peers in config
	for i, peer := range conf.Peers {
		peer.Endpoint = endpoints[0]
		peer.Trick = true
		peer.KeepAlive = 3
		conf.Peers[i] = peer
	}

	// Create userspace tun network stack
	tunDev, tnet, err := newUsermodeTun(conf)
	if err != nil {
		return err
	}

	// Establish wireguard on userspace stack and bind the wireguard sockets to the default interface and apply
	if err := establishWireguard(l.With("gool", "outer"), conf, tunDev, opts.Tun, opts.FwMark); err != nil {
		return err
	}

	// Test wireguard connectivity
	if err := usermodeTunTest(ctx, l, tnet); err != nil {
		return err
	}

	// Create a UDP port forward between localhost and the remote endpoint
	addr, err := wiresocks.NewVtunUDPForwarder(ctx, netip.MustParseAddrPort("127.0.0.1:0"), endpoints[0], tnet, singleMTU)
	if err != nil {
		return err
	}

	// Set up secondary/inner warp config
	conf, err = wiresocks.ParseConfig(path.Join(opts.CacheDir, "secondary", "wgcf-profile.ini"))
	if err != nil {
		return err
	}

	// Set up MTU
	conf.Interface.MTU = doubleMTU
	// Set up DNS Address
	conf.Interface.DNS = []netip.Addr{opts.DnsAddr}

	// Enable keepalive on all peers in config
	for i, peer := range conf.Peers {
		peer.Endpoint = addr.String()
		peer.KeepAlive = 10
		conf.Peers[i] = peer
	}

	if opts.Tun {
		// Create a new tun interface
		tunDev, err := newNormalTun([]netip.Addr{opts.DnsAddr})
		if err != nil {
			return err
		}

		// Establish wireguard tunnel on tun interface but don't bind
		// wireguard sockets to default interface and don't apply fwmark.
		if err := establishWireguard(l.With("gool", "inner"), conf, tunDev, false, opts.FwMark); err != nil {
			return err
		}
		l.Info("serving tun", "interface", "warp0")
		return nil
	}

	// Create userspace tun network stack
	tunDev, tnet, err = newUsermodeTun(conf)
	if err != nil {
		return err
	}

	// Establish wireguard on userspace stack
	if err := establishWireguard(l.With("gool", "inner"), conf, tunDev, false, opts.FwMark); err != nil {
		return err
	}

	// Test wireguard connectivity
	if err := usermodeTunTest(ctx, l, tnet); err != nil {
		return err
	}

	_, err = wiresocks.StartProxy(ctx, l, tnet, opts.Bind)
	if err != nil {
		return err
	}

	l.Info("serving proxy", "address", opts.Bind)
	return nil
}

func runWarpWithPsiphon(ctx context.Context, l *slog.Logger, opts WarpOptions, endpoint string) error {
	// Set up primary/outer warp config
	conf, err := wiresocks.ParseConfig(path.Join(opts.CacheDir, "primary", "wgcf-profile.ini"))
	if err != nil {
		return err
	}

	// Set up MTU
	conf.Interface.MTU = singleMTU
	// Set up DNS Address
	conf.Interface.DNS = []netip.Addr{opts.DnsAddr}

	// Enable trick and keepalive on all peers in config
	for i, peer := range conf.Peers {
		peer.Endpoint = endpoint
		peer.Trick = true
		peer.KeepAlive = 3
		conf.Peers[i] = peer
	}

	// Create userspace tun network stack
	tunDev, tnet, err := newUsermodeTun(conf)
	if err != nil {
		return err
	}

	// Establish wireguard on userspace stack
	if err := establishWireguard(l, conf, tunDev, false, opts.FwMark); err != nil {
		return err
	}

	// Test wireguard connectivity
	if err := usermodeTunTest(ctx, l, tnet); err != nil {
		return err
	}

	// Run a proxy on the userspace stack
	warpBind, err := wiresocks.StartProxy(ctx, l, tnet, netip.MustParseAddrPort("127.0.0.1:0"))
	if err != nil {
		return err
	}

	// run psiphon
	err = psiphon.RunPsiphon(ctx, l.With("subsystem", "psiphon"), warpBind.String(), opts.CacheDir, opts.Bind.String(), opts.Psiphon.Country)
	if err != nil {
		return fmt.Errorf("unable to run psiphon %w", err)
	}

	l.Info("serving proxy", "address", opts.Bind)
	return nil
}

func createPrimaryAndSecondaryIdentities(l *slog.Logger, opts WarpOptions) error {
	// make primary identity
	err := warp.LoadOrCreateIdentity(l, path.Join(opts.CacheDir, "primary"), opts.License)
	if err != nil {
		l.Error("couldn't load primary warp identity")
		return err
	}

	// make secondary
	err = warp.LoadOrCreateIdentity(l, path.Join(opts.CacheDir, "secondary"), opts.License)
	if err != nil {
		l.Error("couldn't load secondary warp identity")
		return err
	}

	return nil
}
