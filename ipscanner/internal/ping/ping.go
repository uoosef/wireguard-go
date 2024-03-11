package ping

import (
	"fmt"
	"net/netip"

	"github.com/bepass-org/wireguard-go/ipscanner/internal/statute"
)

type Ping struct {
	Options *statute.ScannerOptions
}

// DoPing performs a ping on the given IP address.
func (p *Ping) DoPing(ip netip.Addr) (int, error) {
	var sum, ops, hp, tp int
	var err error
	if p.Options.SelectedOps&statute.HTTPPing > 0 {
		hp, err = p.httpPing(ip)
		if err != nil {
			return 0, err
		}
		ops++
		sum += hp
	}
	if p.Options.SelectedOps&statute.TLSPing > 0 {
		tp, err = p.tlsPing(ip)
		if err != nil {
			return 0, err
		}
		ops++
		sum += tp
	}
	if p.Options.SelectedOps&statute.TCPPing > 0 {
		tp, err = p.tcpPing(ip)
		if err != nil {
			return 0, err
		}
		ops++
		sum += tp
	}
	if p.Options.SelectedOps&statute.QUICPing > 0 {
		tp, err = p.quicPing(ip)
		if err != nil {
			return 0, err
		}
		ops++
		sum += tp
	}
	if p.Options.SelectedOps&statute.WARPPing > 0 {
		tp, err = p.warpPing(ip)
		if err != nil {
			return 0, err
		}
		ops++
		sum += tp
	}
	if ops == 0 {
		return 99, nil
	}
	return sum / ops, nil
}

func (p *Ping) httpPing(ip netip.Addr) (int, error) {
	return p.calc(
		NewHttpPing(
			ip,
			"GET",
			fmt.Sprintf(
				"https://%s:%d%s",
				p.Options.Hostname,
				p.Options.Port,
				p.Options.HTTPPath,
			),
			p.Options,
		),
	)
}

func (p *Ping) warpPing(ip netip.Addr) (int, error) {
	return p.calc(
		NewWarpPing(
			ip,
			p.Options,
		),
	)
}

func (p *Ping) tlsPing(ip netip.Addr) (int, error) {
	return p.calc(
		NewTlsPing(ip, p.Options.Hostname, p.Options.Port, p.Options),
	)
}

func (p *Ping) tcpPing(ip netip.Addr) (int, error) {
	return p.calc(
		NewTcpPing(ip, p.Options.Hostname, p.Options.Port, p.Options),
	)
}

func (p *Ping) quicPing(ip netip.Addr) (int, error) {
	return p.calc(
		NewQuicPing(ip, p.Options.Hostname, p.Options.Port, p.Options),
	)
}

func (p *Ping) calc(tp statute.IPing) (int, error) {
	pr := tp.Ping()
	err := pr.Error()
	if err != nil {
		return 0, err
	}
	return pr.Result(), nil
}
