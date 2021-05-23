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
		host string
		err  error
	}{
		{
			name: "scanner is nil",
			err:  errScannerIsNil,
		},
		{
			name: "ok",
			s:    bufio.NewScanner(strings.NewReader("0.0.0.0 host.se")),
			host: "host.se",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &blocker{
				ips:   make(map[string]struct{}),
				hosts: make(map[string]struct{}),
			}
			err := b.addHostsFile(tt.s)
			require.Equal(t, tt.err, err)

			if tt.host != "" {
				res := b.IsHostBlocked(tt.host)
				require.True(t, res)
			}
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
			name: "pattern not found",
			line: "x",
			err:  errNoHostInLine,
		},
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
			name: "ok",
			line: "1.2.3.4 host.se",
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
