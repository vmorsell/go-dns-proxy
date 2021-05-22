package main

import "github.com/miekg/dns"

type Cache struct {
	responses map[string]*dns.Msg
}

func NewCache() *Cache {
	return &Cache{
		responses: make(map[string]*dns.Msg),
	}
}

func (c *Cache) query(key string) *dns.Msg {
	return c.responses[key]
}

func (c *Cache) add(key string, r *dns.Msg) {
	c.responses[key] = r
}
