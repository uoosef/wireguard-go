package ping

import (
	"context"
	"fmt"
	"net"
	"net/netip"
	"time"

	"github.com/bepass-org/wireguard-go/ipscanner/internal/statute"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
)

type QuicPingResult struct {
	Time        int
	Err         error
	IP          netip.Addr
	QUICVersion uint32
	TLSVersion  uint16
}

func (h *QuicPingResult) Result() int {
	return h.Time
}

func (h *QuicPingResult) Error() error {
	return h.Err
}

func (h *QuicPingResult) String() string {
	if h.Err != nil {
		return fmt.Sprintf("%s", h.Err)
	} else {
		return fmt.Sprintf("%s: quic=%s, tls=%s, time=%d ms", h.IP.String(), quic.VersionNumber(h.QUICVersion).String(), statute.TlsVersionToString(h.TLSVersion), h.Time)
	}
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
	ip := statute.CloneIP(h.IP)
	if !ip.IsValid() {
		return h.errorResult(fmt.Errorf("no IP specified"))
	}
	addr := net.JoinHostPort(ip.String(), fmt.Sprint(h.Port))

	t0 := time.Now()
	conn, err := h.opts.QuicDialerFunc(ctx, addr, nil, nil)
	if err != nil {
		return h.errorResult(err)
	}

	defer conn.CloseWithError(quic.ApplicationErrorCode(uint64(http3.ErrCodeNoError)), "")
	return &QuicPingResult{int(time.Since(t0).Milliseconds()), nil, ip, uint32(conn.ConnectionState().Version), conn.ConnectionState().TLS.Version}
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
