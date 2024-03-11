package engine

import (
	"context"
	"fmt"
	"net/netip"
	"strings"
	"time"

	"github.com/bepass-org/wireguard-go/ipscanner/internal/iterator"
	"github.com/bepass-org/wireguard-go/ipscanner/internal/ping"
	"github.com/bepass-org/wireguard-go/ipscanner/internal/statute"
)

type Engine struct {
	generator  *iterator.IpGenerator
	ipQueue    *IPQueue
	ctx        context.Context
	cancelFunc context.CancelFunc
	ping       func(netip.Addr) (int, error)
	statute.Logger
}

func NewScannerEngine(opts *statute.ScannerOptions, ctx ...context.Context) *Engine {
	queue := NewIPQueue(opts)
	var contextToUse context.Context
	var cancel context.CancelFunc

	if len(ctx) > 0 {
		contextToUse = ctx[0]
	} else {
		contextToUse, cancel = context.WithCancel(context.Background())
	}
	p := ping.Ping{
		Options: opts,
	}
	return &Engine{
		ipQueue:    queue,
		ctx:        contextToUse,
		cancelFunc: cancel,
		ping:       p.DoPing,
		generator:  iterator.NewIterator(opts),
		Logger:     opts.Logger,
	}
}

func (e *Engine) GetAvailableIPs(desc bool) []netip.Addr {
	if e.ipQueue != nil {
		return e.ipQueue.AvailableIPs(desc)
	}
	return nil
}

func (e *Engine) Run() {
	for {
		select {
		case <-e.ctx.Done():
			fmt.Println("Context Done!")
			return
		case <-e.ipQueue.available:
			e.Logger.Debug("New Scanning Round Started")
			batch, err := e.generator.NextBatch()
			if err != nil {
				e.Logger.Error("Error while generating IP: %v", err)
				// in case of disastrous error, to prevent resource draining wait for 2 seconds and try again
				time.Sleep(2 * time.Second)
				continue
			}
			for _, ip := range batch {
				select {
				case <-e.ctx.Done():
					fmt.Println("Context Done!")
					return
				default:
					e.Logger.Debug("Pinging IP: %s", ip)
					if rtt, err := e.ping(ip); err == nil {
						ipInfo := statute.IPInfo{
							IP:        ip,
							RTT:       rtt,
							CreatedAt: time.Now(),
						}
						e.Logger.Debug("IP: %s, RTT: %d", ip, rtt)
						e.ipQueue.Enqueue(ipInfo)
					} else {
						// if timeout error
						if strings.Contains(err.Error(), ": i/o timeout") {
							e.Logger.Debug("Timeout Error: %s", ip)
							continue
						}
						e.Logger.Error("Error while pinging IP: %s, Error: %v", ip, err)
					}
				}
			}
		default:
			e.Logger.Debug("Engine: call the expire function")
			e.ipQueue.Expire()
			time.Sleep(200 * time.Millisecond)
		}
	}
}

func (e *Engine) Cancel() {
	e.cancelFunc()
}
