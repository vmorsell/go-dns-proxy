package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/miekg/dns"
)

const (
	port      = 53
	dnsServer = "1.1.1.1:53"
)

func main() {
	cache := NewCache()
	proxy := NewProxy(dnsServer, cache)

	dns.HandleFunc(".", proxy.handler)
	go proxy.start("tcp", port)
	go proxy.start("udp", port)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s := <-sig
	log.Printf("received %v signal", s)
}
