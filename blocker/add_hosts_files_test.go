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
		ip   *string
		host *string
		err  error
	}{
		{
			name: "scanner is nil",
			err:  errScannerIsNil,
		},
		{
			name: "ok",
			s:    bufio.NewScanner(strings.NewReader("1.2.3.4 host.se")),
			ip:   strPtr("1.2.3.4"),
			host: strPtr("host.se"),
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

			if tt.ip != nil {
				res := b.IsIPBlocked(*tt.ip)
				require.True(t, res)
			}

			if tt.host != nil {
				res := b.IsHostBlocked(*tt.host)
				require.True(t, res)
			}
		})
	}
}

func TestParseLine(t *testing.T) {
	tests := []struct {
		name string
		line string
		ip   *string
		host *string
		err  error
	}{
		{
			name: "pattern not found",
			line: "x",
			err:  ErrPatternNotFound,
		},
		{
			name: "full line commented",
			line: "   # 1.2.3.4 host.se",
			err:  ErrPatternNotFound,
		},
		{
			name: "host commented",
			line: "  1.2.3.4  #host.se",
			err:  ErrPatternNotFound,
		},
		{
			name: "ip ignored",
			line: "0.0.0.0 host.se",
			host: strPtr("host.se"),
		},
		{
			name: "ok",
			line: "1.2.3.4 host.se",
			ip:   strPtr("1.2.3.4"),
			host: strPtr("host.se"),
		},
		{
			name: "ok - with leading and trailing whitespaces",
			line: "\t 1.2.3.4 \thost.se  \n",
			ip:   strPtr("1.2.3.4"),
			host: strPtr("host.se"),
		},
		{
			name: "ok - with comment after host",
			line: "1.2.3.4 host.se#comment",
			ip:   strPtr("1.2.3.4"),
			host: strPtr("host.se"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip, host, err := parseLine(tt.line)
			require.Equal(t, tt.ip, ip)
			require.Equal(t, tt.host, host)
			require.Equal(t, tt.err, err)
		})
	}
}

func TestIPIsIgnored(t *testing.T) {
	tests := []struct {
		name string
		ip   *string
		res  bool
	}{
		{
			name: "ip missing",
			res:  true,
		},
		{
			name: "not ignored",
			ip:   strPtr("1.2.3.4"),
			res:  false,
		},
		{
			name: "ignored",
			ip:   strPtr("0.0.0.0"),
			res:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := ipIsIgnored(tt.ip)
			require.Equal(t, tt.res, res)
		})
	}
}

func strPtr(v string) *string {
	return &v
}
