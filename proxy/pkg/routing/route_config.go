package routing

import (
	"errors"
	"fmt"
	"net"
	"regexp"
)

type RoutingRuleType string

const (
	ProxyRoutingRule  RoutingRuleType = "Proxy"
	DirectRoutingRule RoutingRuleType = "Direct"
	BlockRoutingRule  RoutingRuleType = "Block"
)
const domainRegexPattern = `^([A-Za-z0-9-]{1,63}\.)+[A-Za-z]{2,6}$`

var ErrBlockRoutingRule = errors.New("block routing rule")

var domainRegex = regexp.MustCompile(domainRegexPattern)

type RoutingRuleConfig struct {
	routeRuleType    RoutingRuleType
	ipList           []*net.IP
	cidrList         []*net.IPNet
	domainList       []string
	domainRegexpList []*regexp.Regexp
}

type RoutingRuleRawConfig struct {
	RouteRuleType    RoutingRuleType `json:"type"`
	IpList           []string        `json:"ip_list"`
	CidrList         []string        `json:"cidr_list"`
	DomainList       []string        `json:"domain_list"`
	DomainRegexpList []string        `json:"domain_regexp_list"`
}

func NewRoutingRuleConfig(
	routingRuleRawConfig RoutingRuleRawConfig,
) *RoutingRuleConfig {

	if !(routingRuleRawConfig.RouteRuleType == ProxyRoutingRule ||
		routingRuleRawConfig.RouteRuleType == DirectRoutingRule ||
		routingRuleRawConfig.RouteRuleType == BlockRoutingRule) {

		panic(fmt.Sprintf(
			"invalid routing type: %s, use these types: %s,%s,%s",
			routingRuleRawConfig.RouteRuleType,
			ProxyRoutingRule, DirectRoutingRule,
			BlockRoutingRule,
		))
	}

	rc := &RoutingRuleConfig{routeRuleType: routingRuleRawConfig.RouteRuleType}

	// Parse IP addresses
	for _, ipStr := range routingRuleRawConfig.IpList {
		ip := net.ParseIP(ipStr)
		if ip != nil {
			rc.ipList = append(rc.ipList, &ip)
		}
	}

	// Parse CIDR blocks
	for _, cidrStr := range routingRuleRawConfig.CidrList {
		_, cidr, err := net.ParseCIDR(cidrStr)
		if err == nil {
			rc.cidrList = append(rc.cidrList, cidr)
		}
	}

	// Parse domain names
	for _, domainStr := range routingRuleRawConfig.DomainList {
		if host, _, err := net.SplitHostPort(domainStr); err == nil {
			domainStr = host
		}
		rc.domainList = append(rc.domainList, domainStr)
	}

	// Compile regular expressions
	for _, regexStr := range routingRuleRawConfig.DomainRegexpList {
		if compiledRegex, err := regexp.Compile(regexStr); err == nil {
			rc.domainRegexpList = append(rc.domainRegexpList, compiledRegex)
		} else {
			fmt.Println(err)
		}
	}

	return rc
}

func GetRoutingRule(routeConfigList []*RoutingRuleConfig, Destination string) RoutingRuleType {
	for _, routeConfig := range routeConfigList {
		if routeConfig.Contains(Destination) {
			return routeConfig.routeRuleType
		}
	}
	return ""
}

func (rc *RoutingRuleConfig) Contains(Destination string) bool {

	if Host, _, err := net.SplitHostPort(Destination); err == nil {
		Destination = Host
	}

	if ip := net.ParseIP(Destination); ip != nil {
		return rc.containsIP(ip)
	}

	if domainRegex.MatchString(Destination) {
		return rc.containsDomain(Destination) || rc.containsDomainRegexp(Destination)
	}

	return false
}

func (rc *RoutingRuleConfig) containsIP(Destination net.IP) bool {
	for _, netIP := range rc.ipList {
		if netIP.Equal(Destination) {
			return true
		}
	}
	for _, cidr := range rc.cidrList {
		if cidr.Contains(Destination) {
			return true
		}
	}
	return false
}

func (rc *RoutingRuleConfig) containsDomain(destinationHost string) bool {
	for _, domain := range rc.domainList {
		if domain == destinationHost {
			return true
		}
	}
	return false
}

func (rc *RoutingRuleConfig) containsDomainRegexp(Destination string) bool {
	for _, domainRegexp := range rc.domainRegexpList {
		if domainRegexp.MatchString(Destination) {
			return true
		}
	}
	return false
}
