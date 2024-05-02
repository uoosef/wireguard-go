package wiresocks

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"

	"github.com/bepass-org/warp-plus/wireguard/conn"
	"github.com/bepass-org/warp-plus/wireguard/device"
	"github.com/bepass-org/warp-plus/wireguard/tun/netstack"
	"github.com/things-go/go-socks5/bufferpool"
)

// StartWireguard creates a tun interface on netstack given a configuration
func StartWireguard(ctx context.Context, l *slog.Logger, conf *Configuration) (*VirtualTun, error) {
	var request bytes.Buffer

	request.WriteString(fmt.Sprintf("private_key=%s\n", conf.Interface.PrivateKey))

	for _, peer := range conf.Peers {
		request.WriteString(fmt.Sprintf("public_key=%s\n", peer.PublicKey))
		request.WriteString(fmt.Sprintf("persistent_keepalive_interval=%d\n", peer.KeepAlive))
		request.WriteString(fmt.Sprintf("preshared_key=%s\n", peer.PreSharedKey))
		request.WriteString(fmt.Sprintf("endpoint=%s\n", peer.Endpoint))
		request.WriteString(fmt.Sprintf("trick=%t\n", peer.Trick))

		for _, cidr := range peer.AllowedIPs {
			request.WriteString(fmt.Sprintf("allowed_ip=%s\n", cidr))
		}
	}

	tun, tnet, err := netstack.CreateNetTUN(conf.Interface.Addresses, conf.Interface.DNS, conf.Interface.MTU)
	if err != nil {
		return nil, err
	}

	dev := device.NewDevice(tun, conn.NewDefaultBind(), device.NewSLogger(l.With("subsystem", "wireguard-go")))
	err = dev.IpcSet(request.String())
	if err != nil {
		return nil, err
	}

	err = dev.Up()
	if err != nil {
		return nil, err
	}

	return &VirtualTun{
		Tnet:   tnet,
		Logger: l.With("subsystem", "vtun"),
		Dev:    dev,
		Ctx:    ctx,
		pool:   bufferpool.NewPool(256 * 1024),
	}, nil
}
