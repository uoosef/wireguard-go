package wiresocks

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/bepass-org/warp-plus/ipscanner"
	"github.com/bepass-org/warp-plus/warp"
	"github.com/go-ini/ini"
)

type ScanOptions struct {
	V4     bool
	V6     bool
	MaxRTT time.Duration
}

func RunScan(ctx context.Context, l *slog.Logger, opts ScanOptions) (result []ipscanner.IPInfo, err error) {
	cfg, err := ini.Load("./stuff/primary/wgcf-profile.ini")
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Reading the private key from the 'Interface' section
	privateKey := cfg.Section("Interface").Key("PrivateKey").String()

	// Reading the public key from the 'Peer' section
	publicKey := cfg.Section("Peer").Key("PublicKey").String()

	// new scanner
	scanner := ipscanner.NewScanner(
		ipscanner.WithLogger(l.With(slog.String("subsystem", "scanner"))),
		ipscanner.WithWarpPing(),
		ipscanner.WithWarpPrivateKey(privateKey),
		ipscanner.WithWarpPeerPublicKey(publicKey),
		ipscanner.WithUseIPv4(opts.V4),
		ipscanner.WithUseIPv6(opts.V6),
		ipscanner.WithMaxDesirableRTT(opts.MaxRTT),
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
