package warp

import (
	"errors"
	"fmt"
	"io"
	"math/big"
	"math/rand"
	"net"
	"net/netip"
	"time"

	tls "github.com/refraction-networking/utls"
)

var previousIP string

// Dialer is a struct that holds various options for custom dialing.
type Dialer struct {
}

const (
	extensionServerName   uint16 = 0x0
	utlsExtensionSNICurve uint16 = 0x15
)

func hostnameInSNI(name string) string {
	return name
}

// SNIExtension implements server_name (0)
type SNIExtension struct {
	*tls.GenericExtension
	ServerName string // not an array because go crypto/tls doesn't support multiple SNIs
}

// Len returns the length of the SNIExtension.
func (e *SNIExtension) Len() int {
	// Literal IP addresses, absolute FQDNs, and empty strings are not permitted as SNI values.
	// See RFC 6066, Section 3.
	hostName := hostnameInSNI(e.ServerName)
	if len(hostName) == 0 {
		return 0
	}
	return 4 + 2 + 1 + 2 + len(hostName)
}

// Read reads the SNIExtension.
func (e *SNIExtension) Read(b []byte) (int, error) {
	// Literal IP addresses, absolute FQDNs, and empty strings are not permitted as SNI values.
	// See RFC 6066, Section 3.
	hostName := hostnameInSNI(e.ServerName)
	if len(hostName) == 0 {
		return 0, io.EOF
	}
	if len(b) < e.Len() {
		return 0, io.ErrShortBuffer
	}
	// RFC 3546, section 3.1
	b[0] = byte(extensionServerName >> 8)
	b[1] = byte(extensionServerName)
	b[2] = byte((len(hostName) + 5) >> 8)
	b[3] = byte(len(hostName) + 5)
	b[4] = byte((len(hostName) + 3) >> 8)
	b[5] = byte(len(hostName) + 3)
	// b[6] Server Name Type: host_name (0)
	b[7] = byte(len(hostName) >> 8)
	b[8] = byte(len(hostName))
	copy(b[9:], hostName)
	return e.Len(), io.EOF
}

// SNICurveExtension implements SNICurve (0x15) extension
type SNICurveExtension struct {
	*tls.GenericExtension
	SNICurveLen int
	WillPad     bool // set false to disable extension
}

// Len returns the length of the SNICurveExtension.
func (e *SNICurveExtension) Len() int {
	if e.WillPad {
		return 4 + e.SNICurveLen
	}
	return 0
}

// Read reads the SNICurveExtension.
func (e *SNICurveExtension) Read(b []byte) (n int, err error) {
	if !e.WillPad {
		return 0, io.EOF
	}
	if len(b) < e.Len() {
		return 0, io.ErrShortBuffer
	}
	// https://tools.ietf.org/html/rfc7627
	b[0] = byte(utlsExtensionSNICurve >> 8)
	b[1] = byte(utlsExtensionSNICurve)
	b[2] = byte(e.SNICurveLen >> 8)
	b[3] = byte(e.SNICurveLen)
	y := make([]byte, 1200)
	copy(b[4:], y)
	return e.Len(), io.EOF
}

// makeTLSHelloPacketWithSNICurve creates a TLS hello packet with SNICurve.
func (d *Dialer) makeTLSHelloPacketWithSNICurve(plainConn net.Conn, config *tls.Config, sni string) (*tls.UConn, error) {
	SNICurveSize := 1200

	utlsConn := tls.UClient(plainConn, config, tls.HelloCustom)
	spec := tls.ClientHelloSpec{
		TLSVersMax: tls.VersionTLS12,
		TLSVersMin: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.GREASE_PLACEHOLDER,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_AES_128_GCM_SHA256, // tls 1.3
			tls.FAKE_TLS_DHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
		Extensions: []tls.TLSExtension{
			&SNICurveExtension{
				SNICurveLen: SNICurveSize,
				WillPad:     true,
			},
			&tls.SupportedCurvesExtension{Curves: []tls.CurveID{tls.X25519, tls.CurveP256}},
			&tls.SupportedPointsExtension{SupportedPoints: []byte{0}}, // uncompressed
			&tls.SessionTicketExtension{},
			&tls.ALPNExtension{AlpnProtocols: []string{"http/1.1"}},
			&tls.SignatureAlgorithmsExtension{SupportedSignatureAlgorithms: []tls.SignatureScheme{
				tls.ECDSAWithP256AndSHA256,
				tls.ECDSAWithP384AndSHA384,
				tls.ECDSAWithP521AndSHA512,
				tls.PSSWithSHA256,
				tls.PSSWithSHA384,
				tls.PSSWithSHA512,
				tls.PKCS1WithSHA256,
				tls.PKCS1WithSHA384,
				tls.PKCS1WithSHA512,
				tls.ECDSAWithSHA1,
				tls.PKCS1WithSHA1}},
			&tls.KeyShareExtension{KeyShares: []tls.KeyShare{
				{Group: tls.CurveID(tls.GREASE_PLACEHOLDER), Data: []byte{0}},
				{Group: tls.X25519},
			}},
			&tls.PSKKeyExchangeModesExtension{Modes: []uint8{1}}, // pskModeDHE
			&SNIExtension{
				ServerName: sni,
			},
		},
		GetSessionID: nil,
	}
	err := utlsConn.ApplyPreset(&spec)

	if err != nil {
		return nil, fmt.Errorf("uTlsConn.Handshake() error: %+v", err)
	}

	err = utlsConn.Handshake()

	if err != nil {
		return nil, fmt.Errorf("uTlsConn.Handshake() error: %+v", err)
	}

	return utlsConn, nil
}

// RandomIPFromPrefix returns a random IP from the provided CIDR prefix.
// Supports IPv4 and IPv6. Does not support mapped inputs.
func RandomIPFromPrefix(cidr netip.Prefix) (netip.Addr, error) {
	startingAddress := cidr.Masked().Addr()
	if startingAddress.Is4In6() {
		return netip.Addr{}, errors.New("mapped v4 addresses not supported")
	}

	prefixLen := cidr.Bits()
	if prefixLen == -1 {
		return netip.Addr{}, fmt.Errorf("invalid cidr: %s", cidr)
	}

	// Initialise rand number generator
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Find the bit length of the Host portion of the provided CIDR
	// prefix
	hostLen := big.NewInt(int64(startingAddress.BitLen() - prefixLen))

	// Find the max value for our random number
	max := new(big.Int).Exp(big.NewInt(2), hostLen, nil)

	// Generate the random number
	randInt := new(big.Int).Rand(rng, max)

	// Get the first address in the CIDR prefix in 16-bytes form
	startingAddress16 := startingAddress.As16()

	// Convert the first address into a decimal number
	startingAddressInt := new(big.Int).SetBytes(startingAddress16[:])

	// Add the random number to the decimal form of the starting address
	// to get a random address in the desired range
	randomAddressInt := new(big.Int).Add(startingAddressInt, randInt)

	// Convert the random address from decimal form back into netip.Addr
	randomAddress, ok := netip.AddrFromSlice(randomAddressInt.FillBytes(make([]byte, 16)))
	if !ok {
		return netip.Addr{}, fmt.Errorf("failed to generate random IP from CIDR: %s", cidr)
	}

	// Unmap any mapped v4 addresses before return
	return randomAddress.Unmap(), nil
}

// TLSDial dials a TLS connection.
func (d *Dialer) TLSDial(plainDialer *net.Dialer, network, addr string) (net.Conn, error) {
	sni, _, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	ip, err := RandomIPFromPrefix(netip.MustParsePrefix("141.101.113.0/24"))
	if err != nil {
		return nil, err
	}
	plainConn, err := plainDialer.Dial(network, ip.String()+":443")
	if err != nil {
		return nil, err
	}

	config := tls.Config{
		ServerName:         sni,
		InsecureSkipVerify: true,
		NextProtos:         nil,
		MinVersion:         tls.VersionTLS10,
	}

	utlsConn, handshakeErr := d.makeTLSHelloPacketWithSNICurve(plainConn, &config, sni)
	if handshakeErr != nil {
		_ = plainConn.Close()
		fmt.Println(handshakeErr)
		return nil, handshakeErr
	}
	return utlsConn, nil
}
