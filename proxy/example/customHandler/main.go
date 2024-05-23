package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/bepass-org/warp-plus/proxy/pkg/mixed"
	"github.com/bepass-org/warp-plus/proxy/pkg/routing"
	"github.com/bepass-org/warp-plus/proxy/pkg/statute"
)

func main() {
	proxy := mixed.NewProxy(
		mixed.WithBindAddress("127.0.0.1:1080"),
		mixed.WithUserHandler(generalHandler),
		mixed.WithRoutingRuleConfig( // optional
			routing.NewRoutingRuleConfig(
				routing.RoutingRuleRawConfig{
					RouteRuleType:    routing.DirectRoutingRule,
					IpList:           []string{"172.16.0.1"},
					CidrList:         []string{"127.0.0.1/24", "10.0.0.0/8", "192.168.0.0/16"},
					DomainList:       []string{"example.com"},
					DomainRegexpList: []string{`.*\.ir$`},
				},
			),
		),
	)
	_ = proxy.ListenAndServe()
}

func generalHandler(req *statute.ProxyRequest) error {
	fmt.Println("handling request to", req.Destination)
	conn, err := net.Dial(req.Network, req.Destination)
	if err != nil {
		return err
	}
	go func() {
		_, err := io.Copy(conn, req.Conn)
		if err != nil {
			log.Println(err)
		}
	}()
	_, err = io.Copy(req.Conn, conn)
	return err
}
