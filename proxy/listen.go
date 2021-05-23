package proxy

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/miekg/dns"
)

type Port struct {
	Number   int
	Protocol string
}

func (p *proxy) Listen(ports []Port) {
	dns.HandleFunc(".", p.Handler)
	for _, p := range ports {
		go listen(p)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s := <-sig
	log.Printf("received %v signal, exiting", s)
}

func listen(port Port) {
	srv := &dns.Server{
		Addr: fmt.Sprintf(":%d", port.Number),
		Net:  port.Protocol,
	}
	defer srv.Shutdown()

	log.Printf("%s: listening on port %d", port.Protocol, port.Number)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("%s: %s", port.Protocol, err.Error())
	}
}
