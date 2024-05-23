package routing

import (
	"testing"
)

func TestRoutingRuleContains(t *testing.T) {

	testRoutingConfig := NewRoutingRuleConfig(
		RoutingRuleRawConfig{
			RouteRuleType:    DirectRoutingRule,
			IpList:           []string{"192.168.1.1"},
			CidrList:         []string{"10.0.0.0/24"},
			DomainList:       []string{"example.com"},
			DomainRegexpList: []string{`^.*\.example\.com$`, `.*\.digikala\.com$`},
		},
	)
	// Define test cases
	testCases := []struct {
		name        string
		routeConfig *RoutingRuleConfig
		destination string
		expected    bool
	}{
		// Test cases for IP addresses
		{
			name:        "IP in IP list",
			routeConfig: testRoutingConfig,
			destination: "192.168.1.1",
			expected:    true,
		},
		{
			name:        "IP in IP list",
			routeConfig: testRoutingConfig,
			destination: "192.168.1.1:443",
			expected:    true,
		},
		{
			name:        "IP in IP list",
			routeConfig: testRoutingConfig,
			destination: "192.168.1.1:80",
			expected:    true,
		},
		{
			name:        "IP in CIDR list",
			routeConfig: testRoutingConfig,
			destination: "10.0.0.5",
			expected:    true,
		},
		{
			name:        "IP in CIDR list",
			routeConfig: testRoutingConfig,
			destination: "10.0.0.5:443",
			expected:    true,
		},
		{
			name:        "IP not in lists",
			routeConfig: testRoutingConfig,
			destination: "172.16.0.1",
			expected:    false,
		},
		{
			name:        "IP not in lists",
			routeConfig: testRoutingConfig,
			destination: "172.16.0.1:80",
			expected:    false,
		},
		{
			name:        "IP not in lists",
			routeConfig: testRoutingConfig,
			destination: "172.16.0.1",
			expected:    false,
		},
		{
			name:        "Invalid IP",
			routeConfig: testRoutingConfig,
			destination: "256.256.256.256",
			expected:    false,
		},
		// Test cases for domain names
		{
			name:        "Domain in domain list",
			routeConfig: testRoutingConfig,
			destination: "example.com",
			expected:    true,
		},
		{
			name:        "Domain matching regex",
			routeConfig: testRoutingConfig,
			destination: "subdomain.example.com",
			expected:    true,
		},
		{
			name:        "Domain not in lists",
			routeConfig: testRoutingConfig,
			destination: "example.net:443",
			expected:    false,
		},
		{
			name:        "Invalid domain",
			routeConfig: testRoutingConfig,
			destination: "invalid",
			expected:    false,
		},
		// Test cases for regular expressions
		{
			name:        "Destination matches regex",
			routeConfig: testRoutingConfig,
			destination: "prefix.example.com",
			expected:    true,
		},
		{
			name:        "Destination matches regex",
			routeConfig: testRoutingConfig,
			destination: "a.digikala.com",
			expected:    true,
		},
		{
			name:        "Destination matches regex",
			routeConfig: testRoutingConfig,
			destination: "test.example.com:443",
			expected:    true,
		},
		{
			name:        "Destination does not match regex",
			routeConfig: testRoutingConfig,
			destination: "digikala.com:80",
			expected:    false,
		},
		{
			name:        "Destination does not match regex",
			routeConfig: testRoutingConfig,
			destination: "test.digikala.com.invalid.com",
			expected:    false,
		},
		{
			name:        "Destination does not match regex",
			routeConfig: testRoutingConfig,
			destination: "anotherexample.com",
			expected:    false,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.routeConfig.Contains(tc.destination)
			if actual != tc.expected {
				t.Errorf("Expected %v, got %v for destination %s", tc.expected, actual, tc.destination)
			}
		})
	}
}
