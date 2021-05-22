package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/miekg/dns"
)

type Proxy struct {
	server string
	cache  *Cache
}

func NewProxy(server string, cache *Cache) *Proxy {
	return &Proxy{
		server: server,
		cache:  cache,
	}
}

func (p *Proxy) handler(w dns.ResponseWriter, r *dns.Msg) {
	start := time.Now()

	if r.Opcode != dns.OpcodeQuery {
		log.Printf("unsupported opcode: %v", r.Opcode)
		return
	}

	var names []string
	for _, q := range r.Question {
		names = append(names, q.Name)
	}
	log.Printf("resolving: %s", strings.Join(names, ", "))

	res, strat, err := p.resolve(r)
	if err != nil {
		log.Printf("resolve: %s", err.Error())
		return
	}

	if strat != CACHE {
		p.cache.add(getCacheKey(r), res)
	}

	for _, rec := range res.Answer {
		log.Println(rec)
	}

	res.SetReply(r)
	w.WriteMsg(res)

	log.Printf("resolved from %s in %d ms", strat.String(), time.Since(start).Milliseconds())
}

func (p *Proxy) start(protocol string, port int) {
	srv := &dns.Server{
		Addr: fmt.Sprintf(":%d", port),
		Net:  protocol,
	}
	defer srv.Shutdown()

	log.Printf("%s: listening on port %d", protocol, port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("%s service: %s", protocol, err.Error())
	}
}

func (p *Proxy) resolve(r *dns.Msg) (*dns.Msg, ResolveStrategy, error) {
	if res := p.resolveFromCache(r); res != nil {
		return res, CACHE, nil
	}

	res, err := p.resolveFromServer(r)
	if err != nil {
		return nil, SERVER, fmt.Errorf("resolve from server: %w", err)
	}
	return res, SERVER, nil
}

func getCacheKey(r *dns.Msg) string {
	return r.Question[0].Name
}

func (p *Proxy) resolveFromCache(r *dns.Msg) *dns.Msg {
	if len(r.Question) != 1 {
		return nil
	}
	log.Printf("checking cache...")
	return p.cache.query(getCacheKey(r))
}

func (p *Proxy) resolveFromServer(r *dns.Msg) (*dns.Msg, error) {
	log.Printf("resolving from %s...", p.server)
	res, err := dns.Exchange(r, p.server)
	if err != nil {
		return nil, fmt.Errorf("exchange: %w", err)
	}
	return res, nil
}
