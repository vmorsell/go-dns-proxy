package blocker

import (
	"bufio"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddHostsFile(t *testing.T) {
	tests := []struct {
		name string
		s    *bufio.Scanner
		n    int
		err  error
	}{
		{
			name: "scanner is nil",
			err:  errScannerIsNil,
		},
		{
			name: "ok",
			s:    bufio.NewScanner(strings.NewReader("0.0.0.0 host.se")),
			n:    1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &blocker{
				ips:   make(map[string]struct{}),
				hosts: make(map[string]struct{}),
			}
			n, err := b.addHostsFile(tt.s)
			require.Equal(t, tt.n, n)
			require.Equal(t, tt.err, err)
		})
	}
}

func TestParseLine(t *testing.T) {
	tests := []struct {
		name string
		line string
		host string
		err  error
	}{
		{
			name: "full line commented",
			line: "   # 1.2.3.4 host.se",
			err:  errNoHostInLine,
		},
		{
			name: "host commented",
			line: "  1.2.3.4  #host.se",
			err:  errNoHostInLine,
		},
		{
			name: "ok - with leading IP address",
			line: "1.2.3.4 host.se",
			host: "host.se",
		},
		{
			name: "ok - without leading IP address",
			line: "host.se",
			host: "host.se",
		},
		{
			name: "ok - with leading and trailing whitespaces",
			line: "\t 1.2.3.4 \thost.se  \n",
			host: "host.se",
		},
		{
			name: "ok - with comment after host",
			line: "1.2.3.4 host.se#comment",
			host: "host.se",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			host, err := parseLine(tt.line)
			require.Equal(t, tt.host, host)
			require.Equal(t, tt.err, err)
		})
	}
}
