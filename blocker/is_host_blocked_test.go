package blocker

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsHostBlocked(t *testing.T) {
	tests := []struct {
		name  string
		hosts map[string]struct{}
		host  string
		res   bool
	}{
		{
			name:  "not blocked",
			hosts: make(map[string]struct{}),
			host:  "bogus.com.",
			res:   false,
		},
		{
			name: "blocked",
			hosts: map[string]struct{}{
				"bogus.com.": {},
			},
			host: "bogus.com.",
			res:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &blocker{
				hosts: tt.hosts,
			}

			res := b.IsHostBlocked(tt.host)
			require.Equal(t, tt.res, res)
		})
	}
}
