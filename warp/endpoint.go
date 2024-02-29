package warp

import (
	"math/rand"
	"net/netip"
	"time"
)

func WarpPrefixes() []netip.Prefix {
	return []netip.Prefix{
		netip.MustParsePrefix("162.159.192.0/24"),
		netip.MustParsePrefix("162.159.193.0/24"),
		netip.MustParsePrefix("162.159.195.0/24"),
		netip.MustParsePrefix("188.114.96.0/24"),
		netip.MustParsePrefix("188.114.97.0/24"),
		netip.MustParsePrefix("188.114.98.0/24"),
		netip.MustParsePrefix("188.114.99.0/24"),
		netip.MustParsePrefix("2606:4700:d0::/48"),
		netip.MustParsePrefix("2606:4700:d1::/48"),
	}
}

func RandomWarpEndpoint() (netip.AddrPort, error) {
	ports := []int{500, 854, 859, 864, 878, 880, 890, 891, 894, 903, 908, 928, 934, 939, 942,
		943, 945, 946, 955, 968, 987, 988, 1002, 1010, 1014, 1018, 1070, 1074, 1180, 1387, 1701,
		1843, 2371, 2408, 2506, 3138, 3476, 3581, 3854, 4177, 4198, 4233, 4500, 5279,
		5956, 7103, 7152, 7156, 7281, 7559, 8319, 8742, 8854, 8886}

	// Seed the random number generator
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Pick a random port number
	randomPort := uint16(ports[rng.Intn(len(ports))])

	cidrs := WarpPrefixes()

	randomIP, err := RandomIPFromPrefix(cidrs[rng.Intn(len(cidrs))])
	if err != nil {
		return netip.AddrPort{}, err
	}

	return netip.AddrPortFrom(randomIP, randomPort), nil
}
