package cache

import (
	"fmt"

	"github.com/miekg/dns"
)

type ErrNotFound struct {
	Key string
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("%s not found", e.Key)
}

func (c *cache) Get(key string) (*dns.Msg, error) {
	c.mu.RLock()
	res, ok := c.msgs[key]
	c.mu.RUnlock()
	if !ok {
		return nil, ErrNotFound{key}
	}
	return res, nil
}
