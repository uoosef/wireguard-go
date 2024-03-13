package wiresocks

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/netip"
	"time"

	"github.com/bepass-org/wireguard-go/ipscanner"
	"github.com/bepass-org/wireguard-go/warp"
	"github.com/go-ini/ini"
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

func RunScan(ctx context.Context, rtt time.Duration) (result []ipscanner.IPInfo, err error) {
	cfg, err := ini.Load("./primary/wgcf-profile.ini")
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Reading the private key from the 'Interface' section
	privateKey := cfg.Section("Interface").Key("PrivateKey").String()

	// Reading the public key from the 'Peer' section
	publicKey := cfg.Section("Peer").Key("PublicKey").String()

	// new scanner
	scanner := ipscanner.NewScanner(
		ipscanner.WithWarpPing(),
		ipscanner.WithWarpPrivateKey(privateKey),
		ipscanner.WithWarpPeerPublicKey(publicKey),
		ipscanner.WithUseIPv6(canConnectIPv6(netip.MustParseAddrPort("[2001:4860:4860::8888]:80"))),
		ipscanner.WithUseIPv4(true),
		ipscanner.WithMaxDesirableRTT(rtt),
		ipscanner.WithCidrList(warp.WarpPrefixes()),
	)

	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	scanner.Run(ctx)

	t := time.NewTicker(1 * time.Second)
	defer t.Stop()

	for {
		ipList := scanner.GetAvailableIPs()
		if len(ipList) > 1 {
			for i := 0; i < 2; i++ {
				result = append(result, ipList[i])
			}
			return result, nil
		}

		select {
		case <-ctx.Done():
			// Context is done - canceled externally
			return nil, errors.New("user canceled the operation")
		case <-t.C:
			// Prevent the loop from spinning too fast
			continue
		}
	}
}
