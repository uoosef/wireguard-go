package main

import (
	"github.com/bepass-org/warp-plus/proxy/pkg/mixed"
)

func main() {
	proxy := mixed.NewProxy()
	_ = proxy.ListenAndServe()
}
