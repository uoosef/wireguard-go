//go:build windows

package app

import (
	"errors"
	"fmt"
	"net"
	"net/netip"

	"github.com/bepass-org/warp-plus/wireguard/conn"
	"github.com/bepass-org/warp-plus/wireguard/device"
	"github.com/bepass-org/warp-plus/wireguard/tun"
	wgtun "github.com/bepass-org/warp-plus/wireguard/tun"
	"golang.org/x/sys/windows"
	"golang.zx2c4.com/wireguard/windows/tunnel/winipcfg"
)

const wintunGUID = "c33d325f-20cd-44e5-998c-19b0c15b4df1"
const family4 = winipcfg.AddressFamily(windows.AF_INET)
const family6 = winipcfg.AddressFamily(windows.AF_INET6)

func newNormalTun(dns []netip.Addr) (wgtun.Device, error) {
	guid, _ := windows.GUIDFromString(wintunGUID)
	tunDev, err := wgtun.CreateTUNWithRequestedGUID("warp0", &guid, 1280)
	if err != nil {
		return nil, err
	}

	nativeTunDevice := tunDev.(*tun.NativeTun)
	luid := winipcfg.LUID(nativeTunDevice.LUID())

	err = luid.SetIPAddressesForFamily(family4, []netip.Prefix{netip.MustParsePrefix("172.16.0.2/24")})
	if err != nil {
		return nil, err
	}

	// Set this to break IPv6 and prevent leaks. TODO: fix windows ipv6 tun
	err = luid.SetIPAddressesForFamily(family6, []netip.Prefix{netip.MustParsePrefix("fd12:3456:789a:1::1/128")})
	if err != nil {
		return nil, err
	}

tryAgain4:
	err = luid.SetRoutesForFamily(family4, []*winipcfg.RouteData{{Destination: netip.MustParsePrefix("0.0.0.0/0"), NextHop: netip.IPv4Unspecified(), Metric: 0}})
	if err != nil && err == windows.ERROR_NOT_FOUND {
		goto tryAgain4
	} else if err != nil {
		return nil, err
	}

	var ipif *winipcfg.MibIPInterfaceRow
	ipif, err = luid.IPInterface(family4)
	if err != nil {
		return nil, err
	}
	ipif.ForwardingEnabled = true
	ipif.RouterDiscoveryBehavior = winipcfg.RouterDiscoveryDisabled
	ipif.DadTransmits = 0
	ipif.ManagedAddressConfigurationSupported = false
	ipif.OtherStatefulConfigurationSupported = false
	ipif.NLMTU = uint32(1280)
	ipif.UseAutomaticMetric = false
	ipif.Metric = 0

	err = ipif.Set()
	if err != nil && err == windows.ERROR_NOT_FOUND {
		goto tryAgain4
	} else if err != nil {
		return nil, fmt.Errorf("unable to set metric and MTU: %w", err)
	}

tryAgain6:
	err = luid.SetRoutesForFamily(family6, []*winipcfg.RouteData{{Destination: netip.MustParsePrefix("::/0"), NextHop: netip.IPv6Unspecified(), Metric: 0}})
	if err != nil && err == windows.ERROR_NOT_FOUND {
		goto tryAgain6
	} else if err != nil {
		return nil, err
	}

	var ipif6 *winipcfg.MibIPInterfaceRow
	ipif6, err = luid.IPInterface(family6)
	if err != nil {
		return nil, err
	}
	ipif6.RouterDiscoveryBehavior = winipcfg.RouterDiscoveryDisabled
	ipif6.DadTransmits = 0
	ipif6.ManagedAddressConfigurationSupported = false
	ipif6.OtherStatefulConfigurationSupported = false
	ipif6.NLMTU = uint32(1280)
	ipif6.UseAutomaticMetric = false
	ipif6.Metric = 0

	err = ipif6.Set()
	if err != nil && err == windows.ERROR_NOT_FOUND {
		goto tryAgain6
	} else if err != nil {
		return nil, fmt.Errorf("unable to set metric and MTU: %w", err)
	}

	err = luid.SetDNS(family4, dns, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to set DNS: %w", err)
	}

	return tunDev, nil

}

func getAutoDetectInterfaceByFamily(family winipcfg.AddressFamily) (string, error) {
	interfaces, err := winipcfg.GetAdaptersAddresses(family, winipcfg.GAAFlagIncludeGateways)
	if err != nil {
		return "", fmt.Errorf("get default interface failure. %w", err)
	}

	var destination netip.Prefix
	if family == family4 {
		destination = netip.PrefixFrom(netip.IPv4Unspecified(), 0)
	} else {
		destination = netip.PrefixFrom(netip.IPv6Unspecified(), 0)
	}

	for _, ifaceM := range interfaces {
		if ifaceM.OperStatus != winipcfg.IfOperStatusUp {
			continue
		}

		ifname := ifaceM.FriendlyName()

		if ifname == "warp0" {
			continue
		}

		for gatewayAddress := ifaceM.FirstGatewayAddress; gatewayAddress != nil; gatewayAddress = gatewayAddress.Next {
			nextHop, _ := netip.AddrFromSlice(gatewayAddress.Address.IP())

			if _, err = ifaceM.LUID.Route(destination, nextHop.Unmap()); err == nil {
				return ifname, nil
			}
		}
	}

	return "", errors.New("interface not found")
}

func bindToIface(dev *device.Device) error {
	ifaceName, err := getAutoDetectInterfaceByFamily(winipcfg.AddressFamily(family4))
	if err != nil {
		return err
	}

	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		return err
	}

	bind, ok := dev.Bind().(conn.BindSocketToInterface)
	if !ok {
		return errors.New("failed to cast to bindsockettointerface")
	}

	if err := bind.BindSocketToInterface4(uint32(iface.Index), false); err != nil {
		return err
	}

	if err := bind.BindSocketToInterface6(uint32(iface.Index), false); err != nil {
		return err
	}

	return nil
}
