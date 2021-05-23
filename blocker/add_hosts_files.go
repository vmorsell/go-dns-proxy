package blocker

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func (b *blocker) AddHostsFiles(urls ...string) error {
	start := time.Now()

	for _, url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("get: %w", err)
		}
		defer resp.Body.Close()

		s := bufio.NewScanner(resp.Body)
		if err := b.addHostsFile(s); err != nil {
			return err
		}
	}

	log.Printf("Blocked from %d hostsfiles in %f s", len(urls), time.Since(start).Seconds())
	return nil
}

var errScannerIsNil = fmt.Errorf("scanner is nil")

func (b *blocker) addHostsFile(s *bufio.Scanner) error {
	if s == nil {
		return errScannerIsNil
	}

	for s.Scan() {
		ip, host, err := parseLine(s.Text())
		if err != nil {
			if errors.Is(err, ErrPatternNotFound) {
				continue
			}
			return fmt.Errorf("parse: %w", err)
		}
		if ip != nil {
			b.AddIP(*ip)
		}
		if host != nil {
			b.AddHost(*host)
		}
	}
	if err := s.Err(); err != nil {
		return fmt.Errorf("scanner: %w", err)
	}
	return nil
}

var ErrPatternNotFound = fmt.Errorf("pattern not found")

func parseLine(line string) (*string, *string, error) {
	line = strings.TrimSpace(line)
	re := regexp.MustCompile(`^((?:[0-9]{1,3}\.){3}[0-9]{1,3})\s+([a-z0-9\-_\.]+)`)
	res := re.FindStringSubmatch(line)

	if len(res) != 3 {
		return nil, nil, ErrPatternNotFound
	}

	ip := &res[1]
	host := &res[2]

	if ipIsIgnored(ip) {
		ip = nil
	}

	return ip, host, nil
}

func ipIsIgnored(ip *string) bool {
	if ip == nil {
		return true
	}

	ignore := []string{"0.0.0.0", "127.0.0.1"}
	for _, ii := range ignore {
		if ii == *ip {
			return true
		}
	}
	return false
}
