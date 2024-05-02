package app

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/bepass-org/warp-plus/wireguard/conn"
	"github.com/bepass-org/warp-plus/wireguard/device"
	wgtun "github.com/bepass-org/warp-plus/wireguard/tun"
	"github.com/bepass-org/warp-plus/wireguard/tun/netstack"
	"github.com/bepass-org/warp-plus/wiresocks"
)

func newUsermodeTun(conf *wiresocks.Configuration) (wgtun.Device, *netstack.Net, error) {
	tunDev, tnet, err := netstack.CreateNetTUN(conf.Interface.Addresses, conf.Interface.DNS, conf.Interface.MTU)
	if err != nil {
		return nil, nil, err
	}

	return tunDev, tnet, nil
}

func usermodeTunTest(ctx context.Context, l *slog.Logger, tnet *netstack.Net) error {
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(30*time.Second))
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		client := http.Client{Transport: &http.Transport{
			DialContext:           tnet.DialContext,
			ResponseHeaderTimeout: 5 * time.Second,
		}}
		resp, err := client.Get(connTestEndpoint)
		if err != nil {
			l.Error("connection test failed", "error", err.Error())
			continue
		}
		_, err = io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			l.Error("connection test failed", "error", err.Error())
			continue
		}

		l.Info("connection test successful")
		break
	}

	return nil
}

func establishWireguard(l *slog.Logger, conf *wiresocks.Configuration, tunDev wgtun.Device, bind bool, fwmark uint32) error {
	// create the IPC message to establish the wireguard conn
	var request bytes.Buffer

	request.WriteString(fmt.Sprintf("private_key=%s\n", conf.Interface.PrivateKey))
	if fwmark != 0 {
		request.WriteString(fmt.Sprintf("fwmark=%d\n", fwmark))
	}

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

	dev := device.NewDevice(
		tunDev,
		conn.NewDefaultBind(),
		device.NewSLogger(l.With("subsystem", "wireguard-go")),
	)

	if err := dev.IpcSet(request.String()); err != nil {
		return err
	}

	if err := dev.Up(); err != nil {
		return err
	}

	if bind {
		if err := bindToIface(dev); err != nil {
			return err
		}
	}

	return nil
}
