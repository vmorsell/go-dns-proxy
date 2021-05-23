package cache

import (
	"fmt"

	"github.com/miekg/dns"
)

var ErrNotFound = fmt.Errorf("not found")

func (c *cache) Get(key string) (*dns.Msg, error) {
	c.mu.RLock()
	res, ok := c.msgs[key]
	c.mu.RUnlock()
	if !ok {
		return nil, ErrNotFound
	}
	return res, nil
}
