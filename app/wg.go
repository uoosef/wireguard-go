package app

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/bepass-org/warp-plus/wireguard/conn"
	"github.com/bepass-org/warp-plus/wireguard/device"
	wgtun "github.com/bepass-org/warp-plus/wireguard/tun"
	"github.com/bepass-org/warp-plus/wireguard/tun/netstack"
	"github.com/bepass-org/warp-plus/wiresocks"
)

const connTestEndpoint = "http://1.1.1.1/cdn-cgi/trace"

func usermodeTunTest(ctx context.Context, l *slog.Logger, tnet *netstack.Net) error {
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(5*time.Second))
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
		resp, err := client.Head(connTestEndpoint)
		if err != nil {
			l.Error("connection test failed")
			continue
		}
		if resp.StatusCode != http.StatusOK {
			l.Error("connection test failed")
			continue
		}

		l.Info("connection test successful")
		break
	}

	return nil
}

func waitHandshake(ctx context.Context, l *slog.Logger, dev *device.Device) error {
	lastHandshakeSecs := "0"
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		get, err := dev.IpcGet()
		if err != nil {
			continue
		}
		scanner := bufio.NewScanner(strings.NewReader(get))
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				break
			}

			key, value, ok := strings.Cut(line, "=")
			if !ok {
				continue
			}

			if key == "last_handshake_time_sec" {
				lastHandshakeSecs = value
				break
			}
		}
		if lastHandshakeSecs != "0" {
			l.Debug("handshake complete")
			break
		}

		l.Debug("waiting on handshake")
		time.Sleep(1 * time.Second)
	}

	return nil
}

func establishWireguard(l *slog.Logger, conf *wiresocks.Configuration, tunDev wgtun.Device, fwmark uint32, t string) error {
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
		request.WriteString(fmt.Sprintf("trick=%s\n", t))
		request.WriteString(fmt.Sprintf("reserved=%d,%d,%d\n", peer.Reserved[0], peer.Reserved[1], peer.Reserved[2]))

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

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(15*time.Second))
	defer cancel()
	if err := waitHandshake(ctx, l, dev); err != nil {
		dev.BindClose()
		dev.Close()
		return err
	}

	return nil
}
