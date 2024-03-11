package ipscanner

import (
	"crypto/tls"
	"net/netip"
	"time"

	"github.com/bepass-org/wireguard-go/ipscanner/internal/engine"
	"github.com/bepass-org/wireguard-go/ipscanner/internal/statute"
)

type IPScanner struct {
	options statute.ScannerOptions
	logger  statute.Logger
	engine  *engine.Engine
	// onChange func([]netip.Addr)
}

func NewScanner(options ...Option) *IPScanner {
	p := &IPScanner{
		options: statute.ScannerOptions{
			UseIPv4:               true,
			UseIPv6:               true,
			CidrList:              statute.DefaultCFRanges(),
			SelectedOps:           0,
			Logger:                statute.DefaultLogger{},
			InsecureSkipVerify:    true,
			RawDialerFunc:         statute.DefaultDialerFunc,
			TLSDialerFunc:         statute.DefaultTLSDialerFunc,
			QuicDialerFunc:        statute.DefaultQuicDialerFunc,
			HttpClientFunc:        statute.DefaultHTTPClientFunc,
			UseHTTP3:              false,
			UseHTTP2:              false,
			DisableCompression:    false,
			HTTPPath:              "/",
			Referrer:              "",
			UserAgent:             "Chrome/80.0.3987.149",
			Hostname:              "www.cloudflare.com",
			WarpPresharedKey:      "",
			WarpPeerPublicKey:     "",
			WarpPrivateKey:        "",
			Port:                  443,
			IPQueueSize:           8,
			MaxDesirableRTT:       400,
			IPQueueTTL:            30 * time.Second,
			IPQueueChangeCallback: statute.DefaultIPQueueChangeCallback,
			ConnectionTimeout:     1 * time.Second,
			HandshakeTimeout:      1 * time.Second,
			TlsVersion:            tls.VersionTLS13,
		},
		logger: statute.DefaultLogger{},
	}

	for _, option := range options {
		option(p)
	}

	return p
}

type Option func(*IPScanner)

func WithUseIPv4(useIPv4 bool) Option {
	return func(i *IPScanner) {
		i.options.UseIPv4 = useIPv4
	}
}

func WithUseIPv6(useIPv6 bool) Option {
	return func(i *IPScanner) {
		i.options.UseIPv6 = useIPv6
	}
}

func WithDialer(d statute.TDialerFunc) Option {
	return func(i *IPScanner) {
		i.options.RawDialerFunc = d
	}
}

func WithTLSDialer(t statute.TDialerFunc) Option {
	return func(i *IPScanner) {
		i.options.TLSDialerFunc = t
	}
}

func WithQuicDialer(q statute.TQuicDialerFunc) Option {
	return func(i *IPScanner) {
		i.options.QuicDialerFunc = q
	}
}

func WithHttpClientFunc(h statute.THTTPClientFunc) Option {
	return func(i *IPScanner) {
		i.options.HttpClientFunc = h
	}
}

func WithUseHTTP3(useHTTP3 bool) Option {
	return func(i *IPScanner) {
		i.options.UseHTTP3 = useHTTP3
	}
}

func WithUseHTTP2(useHTTP2 bool) Option {
	return func(i *IPScanner) {
		i.options.UseHTTP2 = useHTTP2
	}
}

func WithDisableCompression(disableCompression bool) Option {
	return func(i *IPScanner) {
		i.options.DisableCompression = disableCompression
	}
}

func WithHttpPath(path string) Option {
	return func(i *IPScanner) {
		i.options.HTTPPath = path
	}
}

func WithReferrer(referrer string) Option {
	return func(i *IPScanner) {
		i.options.Referrer = referrer
	}
}

func WithUserAgent(userAgent string) Option {
	return func(i *IPScanner) {
		i.options.UserAgent = userAgent
	}
}

func WithLogger(logger statute.Logger) Option {
	return func(i *IPScanner) {
		i.options.Logger = logger
	}
}

func WithInsecureSkipVerify(insecureSkipVerify bool) Option {
	return func(i *IPScanner) {
		i.options.InsecureSkipVerify = insecureSkipVerify
	}
}

func WithHostname(hostname string) Option {
	return func(i *IPScanner) {
		i.options.Hostname = hostname
	}
}

func WithPort(port uint16) Option {
	return func(i *IPScanner) {
		i.options.Port = port
	}
}

func WithCidrList(cidrList []netip.Prefix) Option {
	return func(i *IPScanner) {
		i.options.CidrList = cidrList
	}
}

func WithHTTPPing() Option {
	return func(i *IPScanner) {
		i.options.SelectedOps |= statute.HTTPPing
	}
}

func WithWarpPing() Option {
	return func(i *IPScanner) {
		i.options.SelectedOps |= statute.WARPPing
	}
}

func WithQUICPing() Option {
	return func(i *IPScanner) {
		i.options.SelectedOps |= statute.QUICPing
	}
}

func WithTCPPing() Option {
	return func(i *IPScanner) {
		i.options.SelectedOps |= statute.TCPPing
	}
}

func WithTLSPing() Option {
	return func(i *IPScanner) {
		i.options.SelectedOps |= statute.TLSPing
	}
}

func WithIPQueueSize(size int) Option {
	return func(i *IPScanner) {
		i.options.IPQueueSize = size
	}
}

func WithMaxDesirableRTT(threshold int) Option {
	return func(i *IPScanner) {
		i.options.MaxDesirableRTT = threshold
	}
}

func WithIPQueueTTL(ttl time.Duration) Option {
	return func(i *IPScanner) {
		i.options.IPQueueTTL = ttl
	}
}

func WithIPQueueChangeCallback(callback statute.TIPQueueChangeCallback) Option {
	return func(i *IPScanner) {
		i.options.IPQueueChangeCallback = callback
	}
}

func WithConnectionTimeout(timeout time.Duration) Option {
	return func(i *IPScanner) {
		i.options.ConnectionTimeout = timeout
	}
}

func WithHandshakeTimeout(timeout time.Duration) Option {
	return func(i *IPScanner) {
		i.options.HandshakeTimeout = timeout
	}
}

func WithTlsVersion(version uint16) Option {
	return func(i *IPScanner) {
		i.options.TlsVersion = version
	}
}

func WithWarpPrivateKey(privateKey string) Option {
	return func(i *IPScanner) {
		i.options.WarpPrivateKey = privateKey
	}
}

func WithWarpPeerPublicKey(peerPublicKey string) Option {
	return func(i *IPScanner) {
		i.options.WarpPeerPublicKey = peerPublicKey
	}
}

func WithWarpPreSharedKey(presharedKey string) Option {
	return func(i *IPScanner) {
		i.options.WarpPresharedKey = presharedKey
	}
}

func (i *IPScanner) SetIPQueueChangeCallback(callback statute.TIPQueueChangeCallback) {
	i.options.IPQueueChangeCallback = callback
}

// run engine and in case of new event call onChange callback also if it gets canceled with context
// cancel all operations

func (i *IPScanner) Run() {
	statute.FinalOptions = &i.options
	if !i.options.UseIPv4 && !i.options.UseIPv6 {
		i.logger.Error("Fatal: both IPv4 and IPv6 are disabled, nothing to do")
		return
	}
	i.engine = engine.NewScannerEngine(&i.options)
	go i.engine.Run()
}

func (i *IPScanner) Stop() {
	i.engine.Cancel()
}

func (i *IPScanner) GetAvailableIPs() []netip.Addr {
	if i.engine != nil {
		return i.engine.GetAvailableIPs(false)
	}
	return nil
}

type IPInfo = statute.IPInfo
