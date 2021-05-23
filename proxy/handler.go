package proxy

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/miekg/dns"
)

func (p *proxy) Handler(w dns.ResponseWriter, r *dns.Msg) {
	start := time.Now()

	var hosts []string
	for _, q := range r.Question {
		hosts = append(hosts, q.Name)
	}
	log.Printf("Resolving %s", strings.Join(hosts, ", "))

	res, strat, err := p.handle(r)
	if err != nil {
		log.Printf("Error: %s", err.Error())
		return
	}

	for _, a := range res.Answer {
		log.Print(a)
	}

	res.SetReply(r)
	w.WriteMsg(res)

	log.Printf("Request resolved using %s in %d ms", strat, time.Since(start).Milliseconds())
}

var ErrUnsupportedOperation = fmt.Errorf("unsupported operation")

func (p *proxy) handle(r *dns.Msg) (*dns.Msg, ResolveStrategy, error) {
	if r.Opcode != dns.OpcodeQuery {
		return nil, NONE, ErrUnsupportedOperation
	}

	res, strat, err := p.Resolve(r)
	if err != nil {
		return nil, strat, fmt.Errorf("resolve: %w", err)
	}

	if strat == CACHE {
		return res, strat, nil
	}

	key, err := p.cache.Key(r)
	if err != nil {
		log.Printf("cache key: %s", err.Error())
		return res, strat, nil
	}

	if err := p.cache.Set(key, res); err != nil {
		log.Printf("cache: %s", err.Error())
		return res, SERVER, nil
	}

	return res, strat, nil
}
