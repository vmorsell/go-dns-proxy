package cache

import (
	"testing"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	tests := []struct {
		name string
		msgs map[string]*dns.Msg
		key  string
		res  *dns.Msg
		err  error
	}{
		{
			name: "not found",
			msgs: map[string]*dns.Msg{
				"other.se": {Question: []dns.Question{
					{Name: "other.se"},
				}},
			},
			key: "x",
			err: ErrNotFound{"x"},
		},
		{
			name: "ok",
			msgs: map[string]*dns.Msg{
				"x": {Question: []dns.Question{
					{Name: "x"},
				}},
			},
			key: "x",
			res: &dns.Msg{Question: []dns.Question{
				{Name: "x"},
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cache{
				msgs: tt.msgs,
			}
			res, err := c.Get(tt.key)
			require.Equal(t, tt.res, res)
			require.Equal(t, tt.err, err)
		})
	}
}
