package cache

import (
	"testing"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/require"
)

func TestSet(t *testing.T) {
	key := "x"
	msg := &dns.Msg{
		Question: []dns.Question{
			{Name: key},
		},
	}
	want := map[string]*dns.Msg{
		key: msg,
	}

	c := &cache{
		msgs: make(map[string]*dns.Msg),
	}
	err := c.Set(key, msg)
	require.Equal(t, c.msgs, want)
	require.Nil(t, err)
}
