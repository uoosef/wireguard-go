//go:build !windows

package app

import (
	"net/netip"

	"github.com/bepass-org/warp-plus/wireguard/device"
	wgtun "github.com/bepass-org/warp-plus/wireguard/tun"
)

func newNormalTun(_ []netip.Addr) (wgtun.Device, error) {
	tunDev, err := wgtun.CreateTUN("warp0", 1280)
	if err != nil {
		return nil, err
	}
	return tunDev, nil
}

func bindToIface(_ *device.Device) error {
	return nil
}
