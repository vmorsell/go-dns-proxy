package cache

import "github.com/miekg/dns"

func (c *cache) Set(key string, msg *dns.Msg) error {
	c.mu.Lock()
	c.msgs[key] = msg
	c.mu.Unlock()
	return nil
}
