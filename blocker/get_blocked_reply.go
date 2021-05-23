package blocker

import (
	"net"

	"github.com/miekg/dns"
)

func (b *blocker) GetBlockedReply(r *dns.Msg) *dns.Msg {
	m := new(dns.Msg)
	m.SetReply(r)

	ip := net.ParseIP("0.0.0.0")
	rrHeader := dns.RR_Header{
		Name:   r.Question[0].Name,
		Rrtype: dns.TypeA,
		Class:  dns.ClassINET,
		Ttl:    3600,
	}
	a := &dns.A{Hdr: rrHeader, A: ip}
	m.Answer = append(m.Answer, a)
	return m
}
