package proxy

import (
	"errors"
	"fmt"
	"log"

	"github.com/miekg/dns"
	"github.com/vmorsell/go-dns-proxy/cache"
)

var (
	ErrNoQuestion = fmt.Errorf("no question in message")
)

func (p *proxy) Resolve(r *dns.Msg) (*dns.Msg, ResolveStrategy, error) {
	if len(r.Question) == 0 {
		return nil, NONE, ErrNoQuestion
	}

	for _, q := range r.Question {
		if p.blocker.IsHostBlocked(q.Name[:len(q.Name)-1]) {
			m := p.blocker.GetBlockedReply(r)
			return m, BLOCKER, nil
		}
	}

	res, err := p.resolveFromCache(r)
	if err != nil {
		if errors.Is(err, cache.ErrNotFound) {
			log.Print("Cache not found")
		}
	}
	if res != nil {
		return res, CACHE, nil
	}

	res, err = p.resolveFromServer(r)
	if err != nil {
		return nil, SERVER, fmt.Errorf("resolve from server: %w", err)
	}

	return res, SERVER, nil
}

func (p *proxy) resolveFromCache(r *dns.Msg) (*dns.Msg, error) {
	key, err := p.cache.Key(r)
	if err != nil {
		return nil, fmt.Errorf("key: %w", err)
	}

	res, err := p.cache.Get(key)
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}
	return res, nil
}

func (p *proxy) resolveFromServer(r *dns.Msg) (*dns.Msg, error) {
	res, err := dns.Exchange(r, p.server)
	if err != nil {
		return nil, fmt.Errorf("exchange: %w", err)
	}
	return res, nil
}
