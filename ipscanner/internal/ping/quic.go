package ping

import (
	"context"
	"errors"
	"fmt"
	"net/netip"
	"time"

	"github.com/bepass-org/warp-plus/ipscanner/internal/statute"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
)

type QuicPingResult struct {
	AddrPort    netip.AddrPort
	QUICVersion quic.VersionNumber
	TLSVersion  uint16
	RTT         time.Duration
	Err         error
}

func (h *QuicPingResult) Result() statute.IPInfo {
	return statute.IPInfo{AddrPort: h.AddrPort, RTT: h.RTT, CreatedAt: time.Now()}
}

func (h *QuicPingResult) Error() error {
	return h.Err
}

func (h *QuicPingResult) String() string {
	if h.Err != nil {
		return fmt.Sprintf("%s", h.Err)
	}

	return fmt.Sprintf("%s: quic=%s, tls=%s, time=%d ms", h.AddrPort, quic.VersionNumber(h.QUICVersion), statute.TlsVersionToString(h.TLSVersion), h.RTT)
}

type QuicPing struct {
	Host string
	Port uint16
	IP   netip.Addr

	opts statute.ScannerOptions
}

func (h *QuicPing) Ping() statute.IPingResult {
	return h.PingContext(context.Background())
}

func (h *QuicPing) PingContext(ctx context.Context) statute.IPingResult {
	if !h.IP.IsValid() {
		return h.errorResult(errors.New("no IP specified"))
	}

	addr := netip.AddrPortFrom(h.IP, h.Port)

	t0 := time.Now()
	conn, err := h.opts.QuicDialerFunc(ctx, addr.String(), nil, nil)
	if err != nil {
		return h.errorResult(err)
	}

	res := QuicPingResult{
		AddrPort:    addr,
		RTT:         time.Since(t0),
		QUICVersion: conn.ConnectionState().Version,
		TLSVersion:  conn.ConnectionState().TLS.Version,
		Err:         nil,
	}

	defer conn.CloseWithError(quic.ApplicationErrorCode(uint64(http3.ErrCodeNoError)), "")
	return &res
}

func NewQuicPing(ip netip.Addr, host string, port uint16, opts *statute.ScannerOptions) *QuicPing {
	return &QuicPing{
		IP:   ip,
		Host: host,
		Port: port,
		opts: *opts,
	}
}

func (h *QuicPing) errorResult(err error) *QuicPingResult {
	r := &QuicPingResult{}
	r.Err = err
	return r
}

var (
	_ statute.IPing       = (*QuicPing)(nil)
	_ statute.IPingResult = (*QuicPingResult)(nil)
)
