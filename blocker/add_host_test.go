package blocker

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddHost(t *testing.T) {
	host := "upper.st."
	want := map[string]struct{}{
		host: {},
	}

	b := &blocker{
		hosts: make(map[string]struct{}),
	}
	err := b.AddHost(host)
	require.Nil(t, err)
	require.Equal(t, want, b.hosts)
}
