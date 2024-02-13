package wiresocks

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/bepass-org/wireguard-go/warp"
	"math/rand"
	"net"
	"strings"
	"time"

	"github.com/go-ini/ini"

	"net/netip"
)

type PeerConfig struct {
	PublicKey    string
	PreSharedKey string
	Endpoint     *string
	KeepAlive    int
	AllowedIPs   []netip.Prefix
}

// DeviceConfig contains the information to initiate a wireguard connection
type DeviceConfig struct {
	SecretKey  string
	Endpoint   []netip.Addr
	Peers      []PeerConfig
	DNS        []netip.Addr
	MTU        int
	ListenPort *int
}

type Configuration struct {
	Device *DeviceConfig
}

var (
	dnsAddresses = []string{"8.8.8.8", "8.8.4.4"}
	dc           = 0
)

func parseString(section *ini.Section, keyName string) (string, error) {
	key := section.Key(strings.ToLower(keyName))
	if key == nil {
		return "", errors.New(keyName + " should not be empty")
	}
	return key.String(), nil
}

func parseBase64KeyToHex(section *ini.Section, keyName string) (string, error) {
	key, err := parseString(section, keyName)
	if err != nil {
		return "", err
	}
	result, err := encodeBase64ToHex(key)
	if err != nil {
		return result, err
	}

	return result, nil
}

func encodeBase64ToHex(key string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return "", errors.New("invalid base64 string: " + key)
	}
	if len(decoded) != 32 {
		return "", errors.New("key should be 32 bytes: " + key)
	}
	return hex.EncodeToString(decoded), nil
}

func parseNetIP(section *ini.Section, keyName string) ([]netip.Addr, error) {
	key := section.Key(keyName)
	if key == nil {
		return []netip.Addr{}, nil
	}

	var ips []netip.Addr
	for _, str := range key.StringsWithShadows(",") {
		str = strings.TrimSpace(str)
		if str == "1.1.1.1" {
			str = dnsAddresses[dc%len(dnsAddresses)]
			dc++
		}
		ip, err := netip.ParseAddr(str)
		if err != nil {
			return nil, err
		}
		ips = append(ips, ip)
	}
	return ips, nil
}

func parseCIDRNetIP(section *ini.Section, keyName string) ([]netip.Addr, error) {
	key := section.Key(keyName)
	if key == nil {
		return []netip.Addr{}, nil
	}

	var ips []netip.Addr
	for _, str := range key.StringsWithShadows(",") {
		prefix, err := netip.ParsePrefix(str)
		if err != nil {
			return nil, err
		}

		addr := prefix.Addr()
		ips = append(ips, addr)
	}
	return ips, nil
}

func parseAllowedIPs(section *ini.Section) ([]netip.Prefix, error) {
	key := section.Key("AllowedIPs")
	if key == nil {
		return []netip.Prefix{}, nil
	}

	var ips []netip.Prefix
	for _, str := range key.StringsWithShadows(",") {
		prefix, err := netip.ParsePrefix(str)
		if err != nil {
			return nil, err
		}

		ips = append(ips, prefix)
	}
	return ips, nil
}

func resolveIP(ip string) (*net.IPAddr, error) {
	return net.ResolveIPAddr("ip", ip)
}

func ResolveIPPAndPort(addr string) (string, error) {
	if addr == "engage.cloudflareclient.com:2408" {
		// Define your specific list of port numbers
		ports := []int{500, 854, 859, 864, 878, 880, 890, 891, 894, 903, 908, 928, 934, 939, 942,
			943, 945, 946, 955, 968, 987, 988, 1002, 1010, 1014, 1018, 1070, 1074, 1180, 1387, 1701,
			1843, 2371, 2408, 2506, 3138, 3476, 3581, 3854, 4177, 4198, 4233, 4500, 5279,
			5956, 7103, 7152, 7156, 7281, 7559, 8319, 8742, 8854, 8886}

		// Seed the random number generator
		rand.Seed(time.Now().UnixNano())

		// Pick a random port number
		randomPort := ports[rand.Intn(len(ports))]

		cidrs := []string{
			"162.159.195.0/24", "188.114.96.0/24", "162.159.192.0/24",
			"188.114.97.0/24", "188.114.99.0/24", "188.114.98.0/24",
		}

		ip, err := warp.RandomIPFromRange(cidrs[rand.Intn(len(cidrs))])
		if err == nil {
			return fmt.Sprintf("%s:%d", ip.String(), randomPort), nil
		}
	}
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return "", err
	}

	ip, err := resolveIP(host)
	if err != nil {
		return "", err
	}
	return net.JoinHostPort(ip.String(), port), nil
}

// ParseInterface parses the [Interface] section and extract the information into `device`
func ParseInterface(cfg *ini.File, device *DeviceConfig) error {
	sections, err := cfg.SectionsByName("Interface")
	if len(sections) != 1 || err != nil {
		return errors.New("one and only one [Interface] is expected")
	}
	section := sections[0]

	address, err := parseCIDRNetIP(section, "Address")
	if err != nil {
		return err
	}

	device.Endpoint = address

	privKey, err := parseBase64KeyToHex(section, "PrivateKey")
	if err != nil {
		return err
	}
	device.SecretKey = privKey

	dns, err := parseNetIP(section, "DNS")
	if err != nil {
		return err
	}
	device.DNS = dns

	if sectionKey, err := section.GetKey("MTU"); err == nil {
		value, err := sectionKey.Int()
		if err != nil {
			return err
		}
		device.MTU = value
		if dc == 0 {
			device.MTU = value + 100
		}
	}

	if sectionKey, err := section.GetKey("ListenPort"); err == nil {
		value, err := sectionKey.Int()
		if err != nil {
			return err
		}
		device.ListenPort = &value
	}

	return nil
}

// ParsePeers parses the [Peer] section and extract the information into `peers`
func ParsePeers(cfg *ini.File, peers *[]PeerConfig, endpoint string) error {
	sections, err := cfg.SectionsByName("Peer")
	if len(sections) < 1 || err != nil {
		return errors.New("at least one [Peer] is expected")
	}

	for _, section := range sections {
		peer := PeerConfig{
			PreSharedKey: "0000000000000000000000000000000000000000000000000000000000000000",
			KeepAlive:    0,
		}

		decoded, err := parseBase64KeyToHex(section, "PublicKey")
		if err != nil {
			return err
		}
		peer.PublicKey = decoded

		if sectionKey, err := section.GetKey("PreSharedKey"); err == nil {
			value, err := encodeBase64ToHex(sectionKey.String())
			if err != nil {
				return err
			}
			peer.PreSharedKey = value
		}

		if sectionKey, err := section.GetKey("Endpoint"); err == nil {
			value := sectionKey.String()
			if endpoint != "notset" {
				peer.Endpoint = &endpoint
			} else {
				decoded, err = ResolveIPPAndPort(strings.ToLower(value))
				if err != nil {
					return err
				}
				peer.Endpoint = &decoded
			}
		}

		if sectionKey, err := section.GetKey("PersistentKeepalive"); err == nil {
			value, err := sectionKey.Int()
			if err != nil {
				return err
			}
			peer.KeepAlive = value
		}

		peer.AllowedIPs, err = parseAllowedIPs(section)
		if err != nil {
			return err
		}

		*peers = append(*peers, peer)
	}
	return nil
}

// ParseConfig takes the path of a configuration file and parses it into Configuration
func ParseConfig(path string, endpoint string) (*Configuration, error) {
	iniOpt := ini.LoadOptions{
		Insensitive:            true,
		AllowShadows:           true,
		AllowNonUniqueSections: true,
	}

	cfg, err := ini.LoadSources(iniOpt, path)
	if err != nil {
		return nil, err
	}

	device := &DeviceConfig{
		MTU: 1420,
	}

	root := cfg.Section("")
	wgConf, err := root.GetKey("WGConfig")
	wgCfg := cfg
	if err == nil {
		wgCfg, err = ini.LoadSources(iniOpt, wgConf.String())
		if err != nil {
			return nil, err
		}
	}

	err = ParseInterface(wgCfg, device)
	if err != nil {
		return nil, err
	}

	err = ParsePeers(wgCfg, &device.Peers, endpoint)
	if err != nil {
		return nil, err
	}

	return &Configuration{
		Device: device,
	}, nil
}
