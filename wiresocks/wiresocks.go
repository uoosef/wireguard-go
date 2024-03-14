package wiresocks

import (
	"bytes"
	"context"
	"fmt"
	"net/netip"

	"github.com/bepass-org/wireguard-go/conn"
	"github.com/bepass-org/wireguard-go/device"
	"github.com/bepass-org/wireguard-go/tun/netstack"
)

// DeviceSetting contains the parameters for setting up a tun interface
type DeviceSetting struct {
	ipcRequest string
	dns        []netip.Addr
	deviceAddr []netip.Addr
	mtu        int
}

// serialize the config into an IPC request and DeviceSetting
func createIPCRequest(conf *Configuration) (*DeviceSetting, error) {
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

	setting := &DeviceSetting{ipcRequest: request.String(), dns: conf.Interface.DNS, deviceAddr: conf.Interface.Addresses, mtu: conf.Interface.MTU}
	return setting, nil
}

// StartWireguard creates a tun interface on netstack given a configuration
func StartWireguard(ctx context.Context, conf *Configuration, verbose bool) (*VirtualTun, error) {
	setting, err := createIPCRequest(conf)
	if err != nil {
		return nil, err
	}

	tun, tnet, err := netstack.CreateNetTUN(setting.deviceAddr, setting.dns, setting.mtu)
	if err != nil {
		return nil, err
	}

	logLevel := device.LogLevelVerbose
	if !verbose {
		logLevel = device.LogLevelSilent
	}

	dev := device.NewDevice(tun, conn.NewDefaultBind(), device.NewLogger(logLevel, ""))
	err = dev.IpcSet(setting.ipcRequest)
	if err != nil {
		return nil, err
	}

	err = dev.Up()
	if err != nil {
		return nil, err
	}

	return &VirtualTun{
		Tnet:      tnet,
		SystemDNS: len(setting.dns) == 0,
		Verbose:   verbose,
		Logger: DefaultLogger{
			verbose: verbose,
		},
		Dev: dev,
		Ctx: ctx,
	}, nil
}
