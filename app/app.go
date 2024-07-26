package app

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"net/netip"
	"path"

	"github.com/bepass-org/warp-plus/iputils"
	"github.com/bepass-org/warp-plus/psiphon"
	"github.com/bepass-org/warp-plus/warp"
	"github.com/bepass-org/warp-plus/wiresocks"
)

const singleMTU = 1330
const doubleMTU = 1280 // minimum mtu for IPv6, may cause frag reassembly somewhere

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
	Reserved        string
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

	// Decide Working Scenario
	endpoints := []string{opts.Endpoint, opts.Endpoint}

	if opts.Scan != nil {
		// make primary identity
		ident, err := warp.LoadOrCreateIdentity(l, path.Join(opts.CacheDir, "primary"), opts.License)
		if err != nil {
			l.Error("couldn't load primary warp identity")
			return err
		}

		// Reading the private key from the 'Interface' section
		opts.Scan.PrivateKey = ident.PrivateKey

		// Reading the public key from the 'Peer' section
		opts.Scan.PublicKey = ident.Config.Peers[0].PublicKey

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

	// // Test wireguard connectivity
	// if err := usermodeTunTest(ctx, l, tnet); err != nil {
	// 	return err
	// }

	// Run a proxy on the userspace stack
	_, err = wiresocks.StartProxy(ctx, l, tnet, opts.Bind)
	if err != nil {
		return err
	}

	l.Info("serving proxy", "address", opts.Bind)

	return nil
}

func runWarp(ctx context.Context, l *slog.Logger, opts WarpOptions, endpoint string) error {
	// make primary identity
	ident, err := warp.LoadOrCreateIdentity(l, path.Join(opts.CacheDir, "primary"), opts.License)
	if err != nil {
		l.Error("couldn't load primary warp identity")
		return err
	}

	conf := generateWireguardConfig(ident)

	// Set up MTU
	conf.Interface.MTU = singleMTU
	// Set up DNS Address
	conf.Interface.DNS = []netip.Addr{opts.DnsAddr}

	// Enable trick and keepalive on all peers in config
	for i, peer := range conf.Peers {
		peer.Endpoint = endpoint
		peer.Trick = true
		peer.KeepAlive = 3

		if opts.Reserved != "" {
			r, err := wiresocks.ParseReserved(opts.Reserved)
			if err != nil {
				return err
			}
			peer.Reserved = r
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
		if err := establishWireguard(l, &conf, tunDev, true, opts.FwMark); err != nil {
			return err
		}
		l.Info("serving tun", "interface", "warp0")
		return nil
	}

	// Create userspace tun network stack
	tunDev, tnet, err := newUsermodeTun(&conf)
	if err != nil {
		return err
	}

	// Establish wireguard on userspace stack
	if err := establishWireguard(l, &conf, tunDev, false, opts.FwMark); err != nil {
		return err
	}

	// // Test wireguard connectivity
	// if err := usermodeTunTest(ctx, l, tnet); err != nil {
	// 	return err
	// }

	// Run a proxy on the userspace stack
	_, err = wiresocks.StartProxy(ctx, l, tnet, opts.Bind)
	if err != nil {
		return err
	}

	l.Info("serving proxy", "address", opts.Bind)
	return nil
}

func runWarpInWarp(ctx context.Context, l *slog.Logger, opts WarpOptions, endpoints []string) error {
	// make primary identity
	ident1, err := warp.LoadOrCreateIdentity(l, path.Join(opts.CacheDir, "primary"), opts.License)
	if err != nil {
		l.Error("couldn't load primary warp identity")
		return err
	}

	conf := generateWireguardConfig(ident1)

	// Set up MTU
	conf.Interface.MTU = singleMTU
	// Set up DNS Address
	conf.Interface.DNS = []netip.Addr{opts.DnsAddr}

	// Enable trick and keepalive on all peers in config
	for i, peer := range conf.Peers {
		peer.Endpoint = endpoints[0]
		peer.Trick = true
		peer.KeepAlive = 3

		if opts.Reserved != "" {
			r, err := wiresocks.ParseReserved(opts.Reserved)
			if err != nil {
				return err
			}
			peer.Reserved = r
		}

		conf.Peers[i] = peer
	}

	// Create userspace tun network stack
	tunDev, tnet, err := newUsermodeTun(&conf)
	if err != nil {
		return err
	}

	// Establish wireguard on userspace stack and bind the wireguard sockets to the default interface and apply
	if err := establishWireguard(l.With("gool", "outer"), &conf, tunDev, opts.Tun, opts.FwMark); err != nil {
		return err
	}

	// // Test wireguard connectivity
	// if err := usermodeTunTest(ctx, l, tnet); err != nil {
	// 	return err
	// }

	// Create a UDP port forward between localhost and the remote endpoint
	addr, err := wiresocks.NewVtunUDPForwarder(ctx, netip.MustParseAddrPort("127.0.0.1:0"), endpoints[0], tnet, singleMTU)
	if err != nil {
		return err
	}

	// make secondary
	ident2, err := warp.LoadOrCreateIdentity(l, path.Join(opts.CacheDir, "secondary"), opts.License)
	if err != nil {
		l.Error("couldn't load secondary warp identity")
		return err
	}

	conf = generateWireguardConfig(ident2)

	// Set up MTU
	conf.Interface.MTU = doubleMTU
	// Set up DNS Address
	conf.Interface.DNS = []netip.Addr{opts.DnsAddr}

	// Enable keepalive on all peers in config
	for i, peer := range conf.Peers {
		peer.Endpoint = addr.String()
		peer.KeepAlive = 10

		if opts.Reserved != "" {
			r, err := wiresocks.ParseReserved(opts.Reserved)
			if err != nil {
				return err
			}
			peer.Reserved = r
		}

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
		if err := establishWireguard(l.With("gool", "inner"), &conf, tunDev, false, opts.FwMark); err != nil {
			return err
		}
		l.Info("serving tun", "interface", "warp0")
		return nil
	}

	// Create userspace tun network stack
	tunDev, tnet, err = newUsermodeTun(&conf)
	if err != nil {
		return err
	}

	// Establish wireguard on userspace stack
	if err := establishWireguard(l.With("gool", "inner"), &conf, tunDev, false, opts.FwMark); err != nil {
		return err
	}

	// // Test wireguard connectivity
	// if err := usermodeTunTest(ctx, l, tnet); err != nil {
	// 	return err
	// }

	_, err = wiresocks.StartProxy(ctx, l, tnet, opts.Bind)
	if err != nil {
		return err
	}

	l.Info("serving proxy", "address", opts.Bind)
	return nil
}

func runWarpWithPsiphon(ctx context.Context, l *slog.Logger, opts WarpOptions, endpoint string) error {
	// make primary identity
	ident, err := warp.LoadOrCreateIdentity(l, path.Join(opts.CacheDir, "primary"), opts.License)
	if err != nil {
		l.Error("couldn't load primary warp identity")
		return err
	}

	conf := generateWireguardConfig(ident)

	// Set up MTU
	conf.Interface.MTU = singleMTU
	// Set up DNS Address
	conf.Interface.DNS = []netip.Addr{opts.DnsAddr}

	// Enable trick and keepalive on all peers in config
	for i, peer := range conf.Peers {
		peer.Endpoint = endpoint
		peer.Trick = true
		peer.KeepAlive = 3

		if opts.Reserved != "" {
			r, err := wiresocks.ParseReserved(opts.Reserved)
			if err != nil {
				return err
			}
			peer.Reserved = r
		}

		conf.Peers[i] = peer
	}

	// Create userspace tun network stack
	tunDev, tnet, err := newUsermodeTun(&conf)
	if err != nil {
		return err
	}

	// Establish wireguard on userspace stack
	if err := establishWireguard(l, &conf, tunDev, false, opts.FwMark); err != nil {
		return err
	}

	// // Test wireguard connectivity
	// if err := usermodeTunTest(ctx, l, tnet); err != nil {
	// 	return err
	// }

	// Run a proxy on the userspace stack
	warpBind, err := wiresocks.StartProxy(ctx, l, tnet, netip.MustParseAddrPort("127.0.0.1:0"))
	if err != nil {
		return err
	}

	// run psiphon
	err = psiphon.RunPsiphon(ctx, l.With("subsystem", "psiphon"), warpBind, opts.CacheDir, opts.Bind, opts.Psiphon.Country)
	if err != nil {
		return fmt.Errorf("unable to run psiphon %w", err)
	}

	l.Info("serving proxy", "address", opts.Bind)
	return nil
}

func generateWireguardConfig(i *warp.Identity) wiresocks.Configuration {
	priv, _ := wiresocks.EncodeBase64ToHex(i.PrivateKey)
	pub, _ := wiresocks.EncodeBase64ToHex(i.Config.Peers[0].PublicKey)
	clientID, _ := base64.StdEncoding.DecodeString(i.Config.ClientID)
	return wiresocks.Configuration{
		Interface: &wiresocks.InterfaceConfig{
			PrivateKey: priv,
			Addresses: []netip.Addr{
				netip.MustParseAddr(i.Config.Interface.Addresses.V4),
				netip.MustParseAddr(i.Config.Interface.Addresses.V6),
			},
		},
		Peers: []wiresocks.PeerConfig{{
			PublicKey:    pub,
			PreSharedKey: "0000000000000000000000000000000000000000000000000000000000000000",
			AllowedIPs: []netip.Prefix{
				netip.MustParsePrefix("0.0.0.0/0"),
				netip.MustParsePrefix("::/0"),
			},
			Endpoint: i.Config.Peers[0].Endpoint.Host,
			Reserved: [3]byte{clientID[0], clientID[1], clientID[2]},
		}},
	}
}
