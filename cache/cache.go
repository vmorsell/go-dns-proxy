package cache

import (
	"sync"

	"github.com/miekg/dns"
)

type Cache interface {
	Get(key string) (*dns.Msg, error)
	Set(key string, msg *dns.Msg) error
}

type cache struct {
	msgs map[string]*dns.Msg
	mu   sync.RWMutex
}

func New() Cache {
	return &cache{
		msgs: make(map[string]*dns.Msg),
	}
}
