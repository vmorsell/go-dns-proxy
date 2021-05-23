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
		host, err := parseLine(s.Text())
		if err != nil {
			if errors.Is(err, errNoHostInLine) {
				continue
			}
			return fmt.Errorf("parse: %w", err)
		}
		b.AddHost(host)
	}
	if err := s.Err(); err != nil {
		return fmt.Errorf("scanner: %w", err)
	}
	return nil
}

var errNoHostInLine = fmt.Errorf("no host in line")

func parseLine(line string) (string, error) {
	line = strings.TrimSpace(line)
	re := regexp.MustCompile(`^(?:[0-9]{1,3}\.){3}[0-9]{1,3}\s+([a-z0-9\-_\.]+)`)
	res := re.FindStringSubmatch(line)

	if len(res) != 2 {
		return "", errNoHostInLine
	}
	return res[1], nil
}
