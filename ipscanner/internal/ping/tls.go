package ping

import (
	"context"
	"fmt"
	"net"
	"net/netip"
	"strconv"
	"time"

	"github.com/bepass-org/wireguard-go/ipscanner/internal/statute"
)

type TlsPingResult struct {
	Time       int
	TLSVersion uint16
	Err        error
	IP         netip.Addr
}

func (t *TlsPingResult) Result() int {
	return t.Time
}

func (t *TlsPingResult) Error() error {
	return t.Err
}

func (t *TlsPingResult) String() string {
	if t.Err != nil {
		return fmt.Sprintf("%s", t.Err)
	} else {
		return fmt.Sprintf("%s: protocol=%s, time=%d ms", t.IP.String(), statute.TlsVersionToString(t.TLSVersion), t.Result())
	}
}

type TlsPing struct {
	Host string
	Port uint16
	IP   netip.Addr

	opts *statute.ScannerOptions
}

func (t *TlsPing) Ping() statute.IPingResult {
	return t.PingContext(context.Background())
}

func (t *TlsPing) PingContext(ctx context.Context) statute.IPingResult {
	ip := statute.CloneIP(t.IP)

	if !ip.IsValid() {
		return t.errorResult(fmt.Errorf("no IP specified"))
	}

	t0 := time.Now()
	client, err := t.opts.TLSDialerFunc(ctx, "tcp", net.JoinHostPort(ip.String(), strconv.FormatUint(uint64(t.Port), 10)))
	if err != nil {
		return t.errorResult(err)
	}
	defer client.Close()
	return &TlsPingResult{int(time.Since(t0).Milliseconds()), t.opts.TlsVersion, nil, ip}
}

func NewTlsPing(ip netip.Addr, host string, port uint16, opts *statute.ScannerOptions) *TlsPing {
	return &TlsPing{
		IP:   ip,
		Host: host,
		Port: port,
		opts: opts,
	}
}

func (t *TlsPing) errorResult(err error) *TlsPingResult {
	r := &TlsPingResult{}
	r.Err = err
	return r
}

var (
	_ statute.IPing       = (*TlsPing)(nil)
	_ statute.IPingResult = (*TlsPingResult)(nil)
)
