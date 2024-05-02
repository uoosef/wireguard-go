package wiresocks

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/netip"

	"github.com/go-ini/ini"
)

type PeerConfig struct {
	PublicKey    string
	PreSharedKey string
	Endpoint     string
	KeepAlive    int
	AllowedIPs   []netip.Prefix
	Trick        bool
}

type InterfaceConfig struct {
	PrivateKey string
	Addresses  []netip.Addr
	DNS        []netip.Addr
	MTU        int
}

type Configuration struct {
	Interface *InterfaceConfig
	Peers     []PeerConfig
}

func encodeBase64ToHex(key string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return "", fmt.Errorf("invalid base64 string: %s", key)
	}
	if len(decoded) != 32 {
		return "", fmt.Errorf("key should be 32 bytes: %s", key)
	}
	return hex.EncodeToString(decoded), nil
}

// ParseInterface parses the [Interface] section
func ParseInterface(cfg *ini.File) (InterfaceConfig, error) {
	device := InterfaceConfig{}
	interfaces, err := cfg.SectionsByName("Interface")
	if len(interfaces) != 1 || err != nil {
		return InterfaceConfig{}, errors.New("only one [Interface] is expected")
	}
	iface := interfaces[0]

	key := iface.Key("Address")
	if key == nil {
		return InterfaceConfig{}, nil
	}

	var addresses []netip.Addr
	for _, str := range key.StringsWithShadows(",") {
		prefix, err := netip.ParsePrefix(str)
		if err != nil {
			return InterfaceConfig{}, err
		}

		addresses = append(addresses, prefix.Addr())
	}
	device.Addresses = addresses

	key = iface.Key("PrivateKey")
	if key == nil {
		return InterfaceConfig{}, errors.New("PrivateKey should not be empty")
	}

	privateKeyHex, err := encodeBase64ToHex(key.String())
	if err != nil {
		return InterfaceConfig{}, err
	}
	device.PrivateKey = privateKeyHex

	key = iface.Key("DNS")
	if key == nil {
		return InterfaceConfig{}, nil
	}

	addresses = []netip.Addr{}
	for _, str := range key.StringsWithShadows(",") {
		ip, err := netip.ParseAddr(str)
		if err != nil {
			return InterfaceConfig{}, err
		}
		addresses = append(addresses, ip)
	}
	device.DNS = addresses

	if sectionKey, err := iface.GetKey("MTU"); err == nil {
		value, err := sectionKey.Int()
		if err != nil {
			return InterfaceConfig{}, err
		}
		device.MTU = value
	}

	return device, nil
}

// ParsePeers parses the [Peer] section and extract the information into `peers`
func ParsePeers(cfg *ini.File) ([]PeerConfig, error) {
	sections, err := cfg.SectionsByName("Peer")
	if len(sections) < 1 || err != nil {
		return nil, errors.New("at least one [Peer] is expected")
	}

	peers := make([]PeerConfig, len(sections))
	for i, section := range sections {
		peer := PeerConfig{
			PreSharedKey: "0000000000000000000000000000000000000000000000000000000000000000",
			KeepAlive:    0,
		}

		if sectionKey, err := section.GetKey("PublicKey"); err == nil {
			value, err := encodeBase64ToHex(sectionKey.String())
			if err != nil {
				return nil, err
			}
			peer.PublicKey = value
		}

		if sectionKey, err := section.GetKey("PreSharedKey"); err == nil {
			value, err := encodeBase64ToHex(sectionKey.String())
			if err != nil {
				return nil, err
			}
			peer.PreSharedKey = value
		}

		if sectionKey, err := section.GetKey("PersistentKeepalive"); err == nil {
			value, err := sectionKey.Int()
			if err != nil {
				return nil, err
			}
			peer.KeepAlive = value
		}

		if sectionKey, err := section.GetKey("AllowedIPs"); err == nil {
			var ips []netip.Prefix
			for _, str := range sectionKey.StringsWithShadows(",") {
				prefix, err := netip.ParsePrefix(str)
				if err != nil {
					return nil, err
				}
				ips = append(ips, prefix)
			}
			peer.AllowedIPs = ips
		}

		if sectionKey, err := section.GetKey("Endpoint"); err == nil {
			peer.Endpoint = sectionKey.String()
		}

		if sectionKey, err := section.GetKey("Trick"); err == nil {
			value, err := sectionKey.Bool()
			if err != nil {
				return nil, err
			}
			peer.Trick = value
		}

		peers[i] = peer
	}

	return peers, nil
}

// ParseConfig takes the path of a configuration file and parses it into Configuration
func ParseConfig(path string) (*Configuration, error) {
	iniOpt := ini.LoadOptions{
		Insensitive:            true,
		AllowShadows:           true,
		AllowNonUniqueSections: true,
	}

	cfg, err := ini.LoadSources(iniOpt, path)
	if err != nil {
		return nil, err
	}

	iface, err := ParseInterface(cfg)
	if err != nil {
		return nil, err
	}

	peers, err := ParsePeers(cfg)
	if err != nil {
		return nil, err
	}

	return &Configuration{Interface: &iface, Peers: peers}, nil
}
