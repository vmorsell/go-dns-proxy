package cache

import (
	"testing"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/require"
)

func TestKey(t *testing.T) {
	tests := []struct {
		name string
		r    *dns.Msg
		res  string
		err  error
	}{
		{
			name: "no questions in message",
			r:    &dns.Msg{},
			err:  ErrNoQuestions,
		},
		{
			name: "ok - one question",
			r: &dns.Msg{
				Question: []dns.Question{
					{Name: "x"},
				},
			},
			res: "x",
		},
		{
			name: "ok - multiple questions",
			r: &dns.Msg{
				Question: []dns.Question{
					{Name: "x"},
					{Name: "y"},
				},
			},
			res: "x,y",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cache{}
			res, err := c.Key(tt.r)
			require.Equal(t, tt.res, res)
			require.Equal(t, tt.err, err)
		})
	}
}
