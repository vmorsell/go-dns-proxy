package blocker

import (
	"sync"

	"github.com/miekg/dns"
)

type Blocker interface {
	AddIP(ip string) error
	AddHost(host string) error
	AddHostsFiles(urls ...string) error
	IsIPBlocked(ip string) bool
	IsHostBlocked(host string) bool
	GetBlockedReply(r *dns.Msg) *dns.Msg
}

type blocker struct {
	ips   map[string]struct{}
	hosts map[string]struct{}
	mu    sync.RWMutex
}

func New() Blocker {
	return &blocker{
		ips:   make(map[string]struct{}),
		hosts: make(map[string]struct{}),
	}
}
