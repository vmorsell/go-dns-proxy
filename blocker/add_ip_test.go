package blocker

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddIP(t *testing.T) {
	ip := "1.1.1.1"
	want := map[string]struct{}{
		ip: {},
	}

	b := &blocker{
		ips: make(map[string]struct{}),
	}
	err := b.AddIP(ip)
	require.Nil(t, err)
	require.Equal(t, want, b.ips)
}
