package engine

import (
	"context"
	"log/slog"
	"net/netip"
	"time"

	"github.com/bepass-org/warp-plus/ipscanner/internal/iterator"
	"github.com/bepass-org/warp-plus/ipscanner/internal/ping"
	"github.com/bepass-org/warp-plus/ipscanner/internal/statute"
)

type Engine struct {
	generator *iterator.IpGenerator
	ipQueue   *IPQueue
	ping      func(netip.Addr) (statute.IPInfo, error)
	log       *slog.Logger
}

func NewScannerEngine(opts *statute.ScannerOptions) *Engine {
	queue := NewIPQueue(opts)

	p := ping.Ping{
		Options: opts,
	}
	return &Engine{
		ipQueue:   queue,
		ping:      p.DoPing,
		generator: iterator.NewIterator(opts),
		log:       opts.Logger.With(slog.String("subsystem", "scanner/engine")),
	}
}

func (e *Engine) GetAvailableIPs(desc bool) []statute.IPInfo {
	if e.ipQueue != nil {
		return e.ipQueue.AvailableIPs(desc)
	}
	return nil
}

func (e *Engine) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-e.ipQueue.available:
			e.log.Debug("Started new scanning round")
			batch, err := e.generator.NextBatch()
			if err != nil {
				e.log.Error("Error while generating IP: %v", err)
				// in case of disastrous error, to prevent resource draining wait for 2 seconds and try again
				time.Sleep(2 * time.Second)
				continue
			}
			for _, ip := range batch {
				select {
				case <-ctx.Done():
					return
				default:
					e.log.Debug("pinging IP", "addr", ip)
					if ipInfo, err := e.ping(ip); err == nil {
						e.log.Debug("ping success", "addr", ipInfo.AddrPort, "rtt", ipInfo.RTT)
						e.ipQueue.Enqueue(ipInfo)
					} else {
						e.log.Error("ping error", "addr", ip, "error", err)
					}
				}
			}
		default:
			e.log.Debug("calling expire")
			e.ipQueue.Expire()
			time.Sleep(200 * time.Millisecond)
		}
	}
}
