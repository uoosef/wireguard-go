package main

import "github.com/bepass-org/wireguard-go/ipscanner"

func main() {
	// new scanner
	scanner := ipscanner.NewScanner(
		ipscanner.WithHTTPPing(),
		ipscanner.WithUseIPv6(true),
	)
	go scanner.Run()
	select {}
}
