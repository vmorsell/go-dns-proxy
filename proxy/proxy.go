package proxy

import (
	"github.com/miekg/dns"
	"github.com/vmorsell/go-dns-proxy/blocker"
	"github.com/vmorsell/go-dns-proxy/cache"
)

type Proxy interface {
	Listen(ports []Port)
	Resolve(r *dns.Msg) (*dns.Msg, ResolveStrategy, error)
}

type proxy struct {
	server  string
	cache   cache.Cache
	blocker blocker.Blocker
}

func New(server string, cache cache.Cache, blocker blocker.Blocker) Proxy {
	return &proxy{
		server:  server,
		cache:   cache,
		blocker: blocker,
	}
}
