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
}

func NewProxy(server string) *Proxy {
	return &Proxy{
		server: server,
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
	res, err := p.resolveFromServer(r)
	if err != nil {
		log.Printf("exchange: %s", err.Error())
	}
	for _, rec := range res.Answer {
		log.Println(rec)
	}

	res.SetReply(r)
	w.WriteMsg(res)

	log.Printf("resolved in %d ms", time.Since(start).Milliseconds())
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

func (p *Proxy) resolveFromServer(r *dns.Msg) (*dns.Msg, error) {
	log.Printf("resolving from %s", p.server)
	res, err := dns.Exchange(r, p.server)
	if err != nil {
		return nil, fmt.Errorf("exchange: %w", err)
	}
	return res, nil
}
