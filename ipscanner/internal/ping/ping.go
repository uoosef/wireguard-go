package ping

import (
	"errors"
	"fmt"
	"net/netip"

	"github.com/bepass-org/warp-plus/ipscanner/internal/statute"
)

type Ping struct {
	Options *statute.ScannerOptions
}

// DoPing performs a ping on the given IP address.
func (p *Ping) DoPing(ip netip.Addr) (statute.IPInfo, error) {
	if p.Options.SelectedOps&statute.HTTPPing > 0 {
		res, err := p.httpPing(ip)
		if err != nil {
			return statute.IPInfo{}, err
		}

		return res, nil
	}
	if p.Options.SelectedOps&statute.TLSPing > 0 {
		res, err := p.tlsPing(ip)
		if err != nil {
			return statute.IPInfo{}, err
		}

		return res, nil
	}
	if p.Options.SelectedOps&statute.TCPPing > 0 {
		res, err := p.tcpPing(ip)
		if err != nil {
			return statute.IPInfo{}, err
		}

		return res, nil
	}
	if p.Options.SelectedOps&statute.QUICPing > 0 {
		res, err := p.quicPing(ip)
		if err != nil {
			return statute.IPInfo{}, err
		}

		return res, nil
	}
	if p.Options.SelectedOps&statute.WARPPing > 0 {
		res, err := p.warpPing(ip)
		if err != nil {
			return statute.IPInfo{}, err
		}

		return res, nil
	}

	return statute.IPInfo{}, errors.New("no ping operation selected")
}

func (p *Ping) httpPing(ip netip.Addr) (statute.IPInfo, error) {
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

func (p *Ping) warpPing(ip netip.Addr) (statute.IPInfo, error) {
	return p.calc(NewWarpPing(ip, p.Options))
}

func (p *Ping) tlsPing(ip netip.Addr) (statute.IPInfo, error) {
	return p.calc(
		NewTlsPing(ip, p.Options.Hostname, p.Options.Port, p.Options),
	)
}

func (p *Ping) tcpPing(ip netip.Addr) (statute.IPInfo, error) {
	return p.calc(
		NewTcpPing(ip, p.Options.Hostname, p.Options.Port, p.Options),
	)
}

func (p *Ping) quicPing(ip netip.Addr) (statute.IPInfo, error) {
	return p.calc(
		NewQuicPing(ip, p.Options.Hostname, p.Options.Port, p.Options),
	)
}

func (p *Ping) calc(tp statute.IPing) (statute.IPInfo, error) {
	pr := tp.Ping()
	err := pr.Error()
	if err != nil {
		return statute.IPInfo{}, err
	}
	return pr.Result(), nil
}
