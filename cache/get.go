package cache

import (
	"fmt"

	"github.com/miekg/dns"
)

var (
	errNotFound = func(key string) error { return fmt.Errorf("not found: %s", key) }
)

func (c *cache) Get(key string) (*dns.Msg, error) {
	c.mu.RLock()
	res, ok := c.msgs[key]
	c.mu.RUnlock()
	if !ok {
		return nil, errNotFound(key)
	}
	return res, nil
}
