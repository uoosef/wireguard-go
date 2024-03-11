package ping

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"time"

	"github.com/bepass-org/wireguard-go/ipscanner/internal/statute"
)

type HttpPingResult struct {
	Time   int
	Proto  string
	Status int
	Length int
	Err    error
	IP     netip.Addr
}

func (h *HttpPingResult) Result() int {
	return h.Time
}

func (h *HttpPingResult) Error() error {
	return h.Err
}

func (h *HttpPingResult) String() string {
	if h.Err != nil {
		return fmt.Sprintf("%s", h.Err)
	} else {
		return fmt.Sprintf("%s: protocol=%s, status=%d, length=%d, time=%d ms", h.IP.String(), h.Proto, h.Status, h.Length, h.Time)
	}
}

type HttpPing struct {
	Method string
	URL    string
	IP     netip.Addr

	opts statute.ScannerOptions
}

func (h *HttpPing) Ping() statute.IPingResult {
	return h.PingContext(context.Background())
}

func (h *HttpPing) PingContext(ctx context.Context) statute.IPingResult {
	u, err := url.Parse(h.URL)
	if err != nil {
		return h.errorResult(err)
	}
	orighost := u.Host
	port := u.Port()
	ip := statute.CloneIP(h.IP)
	if !ip.IsValid() {
		return h.errorResult(fmt.Errorf("no IP specified"))
	}
	ipstr := ip.String()
	if statute.IsIPv6(ip) {
		ipstr = fmt.Sprintf("[%s]", ipstr)
	}
	targetAddr := net.JoinHostPort(ipstr, port)

	req, err := http.NewRequestWithContext(ctx, h.Method, h.URL, nil)
	if err != nil {
		return h.errorResult(err)
	}
	ua := "httping"
	if h.opts.UserAgent != "" {
		ua = h.opts.UserAgent
	}
	req.Header.Set("User-Agent", ua)
	if h.opts.Referrer != "" {
		req.Header.Set("Referer", h.opts.Referrer)
	}
	req.Host = orighost

	client := h.opts.HttpClientFunc(h.opts.RawDialerFunc, h.opts.TLSDialerFunc, h.opts.QuicDialerFunc, targetAddr)

	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	t0 := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return h.errorResult(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return h.errorResult(err)
	}
	return &HttpPingResult{int(time.Since(t0).Milliseconds()), resp.Proto, resp.StatusCode, len(body), nil, ip}
}

func (h *HttpPing) errorResult(err error) *HttpPingResult {
	r := &HttpPingResult{}
	r.Err = err
	return r
}

func NewHttpPing(ip netip.Addr, method, url string, opts *statute.ScannerOptions) *HttpPing {
	return &HttpPing{
		IP:     ip,
		Method: method,
		URL:    url,

		opts: *opts,
	}
}

var (
	_ statute.IPing       = (*HttpPing)(nil)
	_ statute.IPingResult = (*HttpPingResult)(nil)
)
