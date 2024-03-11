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

type TcpPingResult struct {
	Time int
	Err  error
	IP   netip.Addr
}

func (tp *TcpPingResult) Result() int {
	return tp.Time
}

func (tp *TcpPingResult) Error() error {
	return tp.Err
}

func (tp *TcpPingResult) String() string {
	if tp.Err != nil {
		return fmt.Sprintf("%s", tp.Err)
	} else {
		return fmt.Sprintf("%s: time=%d ms", tp.IP.String(), tp.Time)
	}
}

type TcpPing struct {
	host string
	Port uint16
	ip   netip.Addr

	opts statute.ScannerOptions
}

func (tp *TcpPing) SetHost(host string) {
	tp.host = host
	tp.ip, _ = netip.ParseAddr(host)
}

func (tp *TcpPing) Host() string {
	return tp.host
}

func (tp *TcpPing) Ping() statute.IPingResult {
	return tp.PingContext(context.Background())
}

func (tp *TcpPing) PingContext(ctx context.Context) statute.IPingResult {
	ip := statute.CloneIP(tp.ip)
	if !ip.IsValid() {
		return &TcpPingResult{0, fmt.Errorf("no IP specified"), netip.Addr{}}
	}
	t0 := time.Now()
	conn, err := tp.opts.RawDialerFunc(ctx, "tcp", net.JoinHostPort(ip.String(), strconv.FormatUint(uint64(tp.Port), 10)))
	if err != nil {
		return &TcpPingResult{0, err, netip.Addr{}}
	}
	defer conn.Close()
	return &TcpPingResult{int(time.Since(t0).Milliseconds()), nil, ip}
}

func NewTcpPing(ip netip.Addr, host string, port uint16, opts *statute.ScannerOptions) *TcpPing {
	return &TcpPing{
		host: host,
		Port: port,
		ip:   ip,

		opts: *opts,
	}
}

var (
	_ statute.IPing       = (*TcpPing)(nil)
	_ statute.IPingResult = (*TcpPingResult)(nil)
)
