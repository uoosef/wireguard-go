package main

import (
	"fmt"
	"net"
	"net/netip"
	"time"

	"github.com/bepass-org/wireguard-go/ipscanner"
	"github.com/bepass-org/wireguard-go/warp"
)

var (
	privKey           = "yGXeX7gMyUIZmK5QIgC7+XX5USUSskQvBYiQ6LdkiXI="
	pubKey            = "bmXOC+F1FxEMF9dyiK2H5/1SUtzH0JuVo51h2wPfgyo="
	googlev6DNSAddr80 = netip.MustParseAddrPort("[2001:4860:4860::8888]:80")
)

func canConnectIPv6(remoteAddr netip.AddrPort) bool {
	dialer := net.Dialer{
		Timeout: 5 * time.Second,
	}

	conn, err := dialer.Dial("tcp6", remoteAddr.String())
	if err != nil {
		return false
	}
	defer conn.Close()

	return true
}

func RunScan(privKey, pubKey string) (result []netip.AddrPort) {
	// new scanner
	scanner := ipscanner.NewScanner(
		ipscanner.WithWarpPing(),
		ipscanner.WithWarpPrivateKey(privKey),
		ipscanner.WithWarpPeerPublicKey(pubKey),
		ipscanner.WithUseIPv6(canConnectIPv6(googlev6DNSAddr80)),
		ipscanner.WithUseIPv4(true),
		ipscanner.WithMaxDesirableRTT(500),
		ipscanner.WithCidrList(warp.WarpPrefixes()),
	)
	scanner.Run()
	var ipList []netip.Addr
	for {
		ipList = scanner.GetAvailableIPs()
		if len(ipList) > 2 {
			scanner.Stop()
			break
		}
		time.Sleep(1 * time.Second)
	}
	for i := 0; i < 2; i++ {
		result = append(result, netip.AddrPortFrom(ipList[i], warp.RandomWarpPort()))
	}
	return
}

func main() {
	fmt.Println(RunScan(privKey, pubKey))
	time.Sleep(10 * time.Second)
}
