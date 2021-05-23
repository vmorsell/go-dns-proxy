package proxy

import (
	"fmt"
	"testing"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/require"
	"github.com/vmorsell/go-dns-proxy/blocker"
	"github.com/vmorsell/go-dns-proxy/cache"
)

func TestResolve(t *testing.T) {
	t.Run("no questions in message", func(t *testing.T) {
		r := &dns.Msg{}

		p := &proxy{}
		_, _, err := p.Resolve(r)
		require.Equal(t, ErrNoQuestion, err)
	})

	t.Run("blocked", func(t *testing.T) {
		host := "blocked.se"
		r := &dns.Msg{
			Question: []dns.Question{
				{Name: fmt.Sprintf("%s.", host)},
			},
		}

		p := &proxy{
			blocker: blocker.New(),
		}
		p.blocker.AddHost(host)
		_, strat, _ := p.Resolve(r)
		require.Equal(t, BLOCKER, strat)
	})
	t.Run("cached", func(t *testing.T) {
		r := &dns.Msg{
			Question: []dns.Question{
				{Name: "test.se."},
			},
		}

		p := &proxy{
			blocker: blocker.New(),
			cache:   cache.New(),
		}
		key, _ := p.cache.Key(r)
		p.cache.Set(key, r)
		_, strat, _ := p.Resolve(r)
		require.Equal(t, CACHE, strat)
	})

	t.Run("server", func(t *testing.T) {
		r := &dns.Msg{
			Question: []dns.Question{
				{Name: "upper.st."},
			},
		}

		p := New("1.1.1.1:53", cache.New(), blocker.New())
		_, strat, _ := p.Resolve(r)
		require.Equal(t, SERVER, strat)
	})
}

func TestResolveFromCache(t *testing.T) {
	t.Run("not found", func(t *testing.T) {
		c := cache.New()
		p := &proxy{
			cache: c,
		}

		key := "x"

		res, err := p.cache.Get(key)
		require.Nil(t, res)
		require.Equal(t, cache.ErrNotFound, err)
	})

	t.Run("not found", func(t *testing.T) {
		c := cache.New()
		p := &proxy{
			cache: c,
		}

		key := "x"
		msg := &dns.Msg{
			Question: []dns.Question{
				{Name: key},
			},
		}

		err := p.cache.Set(key, msg)
		require.Nil(t, err)

		res, err := p.cache.Get(key)
		require.Equal(t, msg, res)
		require.Nil(t, err)
	})
}
