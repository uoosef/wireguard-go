package main

import (
	"context"
	"log/slog"
	"net"
	"net/netip"
	"os"
	"time"

	"github.com/bepass-org/wireguard-go/ipscanner"
	"github.com/bepass-org/wireguard-go/ipscanner/internal/statute"
	"github.com/bepass-org/wireguard-go/warp"
	"github.com/fatih/color"
	"github.com/rodaine/table"
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

func RunScan(privKey, pubKey string) (result []statute.IPInfo) {
	// new scanner
	scanner := ipscanner.NewScanner(
		ipscanner.WithLogger(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))),
		ipscanner.WithWarpPing(),
		ipscanner.WithWarpPrivateKey(privKey),
		ipscanner.WithWarpPeerPublicKey(pubKey),
		ipscanner.WithUseIPv6(canConnectIPv6(googlev6DNSAddr80)),
		ipscanner.WithUseIPv4(true),
		ipscanner.WithMaxDesirableRTT(500*time.Millisecond),
		ipscanner.WithCidrList(warp.WarpPrefixes()),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
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
			return
		}

		select {
		case <-ctx.Done():
			// Context is done
			return
		case <-t.C:
			// Prevent the loop from spinning too fast
			continue
		}
	}
}

func main() {
	result := RunScan(privKey, pubKey)

	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Address", "RTT (ping)", "Time")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, info := range result {
		tbl.AddRow(info.AddrPort, info.RTT, info.CreatedAt)
	}

	tbl.Print()
}
