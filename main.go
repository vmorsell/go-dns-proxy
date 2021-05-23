package main

import (
	"github.com/vmorsell/go-dns-proxy/blocker"
	"github.com/vmorsell/go-dns-proxy/cache"
	"github.com/vmorsell/go-dns-proxy/proxy"
)

const (
	port      = 53
	dnsServer = "1.1.1.1:53"
)

func main() {
	b := blocker.New()
	c := cache.New()
	p := proxy.New(dnsServer, c, b)

	p.Listen([]proxy.Port{
		{Number: port, Protocol: "tcp"},
		{Number: port, Protocol: "udp"},
	})
}
