package blocker

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsIPBlocked(t *testing.T) {
	tests := []struct {
		name string
		ips  map[string]struct{}
		ip   string
		res  bool
	}{
		{
			name: "not blocked",
			ips:  make(map[string]struct{}),
			ip:   "1.1.1.1",
			res:  false,
		},
		{
			name: "blocked",
			ips: map[string]struct{}{
				"1.1.1.1": {},
			},
			ip:  "1.1.1.1",
			res: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &blocker{
				ips: tt.ips,
			}

			res := b.IsIPBlocked(tt.ip)
			require.Equal(t, tt.res, res)
		})
	}
}
