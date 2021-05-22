package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/miekg/dns"
)

const (
	port      = 53
	dnsServer = "1.1.1.1:53"
)

func handler(w dns.ResponseWriter, r *dns.Msg) {
	if r.Opcode != dns.OpcodeQuery {
		log.Printf("unsupported opcode: %v", r.Opcode)
		return
	}

	var names []string
	for _, q := range r.Question {
		names = append(names, q.Name)
	}
	log.Printf("resolving: %s", strings.Join(names, ", "))

	res, err := dns.Exchange(r, dnsServer)
	if err != nil {
		log.Printf("exchange: %s", err.Error())
	}
	for _, rec := range res.Answer {
		log.Println(rec)
	}

	res.SetReply(r)
	w.WriteMsg(res)
}

func main() {
	dns.HandleFunc(".", handler)
	go run("tcp", port)
	go run("udp", port)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s := <-sig
	log.Printf("received %v signal", s)
}

func run(protocol string, port int) {
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
