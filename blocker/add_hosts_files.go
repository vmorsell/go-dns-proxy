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

	nFiles := 0
	nHosts := 0

	log.Print("Adding hosts from files...")
	for i, url := range urls {
		log.Printf("[%d / %d] %s...", i+1, len(urls), url)

		resp, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("get: %w", err)
		}
		defer resp.Body.Close()

		s := bufio.NewScanner(resp.Body)
		n, err := b.addHostsFile(s)
		if err != nil {
			log.Printf("add hosts files: %s", err.Error())
			if n == 0 {
				// Report even partly processed files
				continue
			}
		}
		nFiles++
		nHosts += n

	}

	log.Printf("Blocked %d hosts from %d files in %f s", nHosts, nFiles, time.Since(start).Seconds())
	return nil
}

var errScannerIsNil = fmt.Errorf("scanner is nil")

func (b *blocker) addHostsFile(s *bufio.Scanner) (int, error) {
	n := 0

	if s == nil {
		return n, errScannerIsNil
	}

	for s.Scan() {
		host, err := parseLine(s.Text())
		if err != nil {
			if errors.Is(err, errNoHostInLine) {
				continue
			}
			return n, fmt.Errorf("parse: %w", err)
		}
		b.AddHost(host)
		n++
	}
	if err := s.Err(); err != nil {
		return n, fmt.Errorf("scanner: %w", err)
	}
	return n, nil
}

var errNoHostInLine = fmt.Errorf("no host in line")

func parseLine(line string) (string, error) {
	line = strings.TrimSpace(line)
	re := regexp.MustCompile(`^(?:\s*[0-9\.]{7,}\s+)?([a-zA-Z0-9-_\.]+\.[a-zA-Z][a-zA-Z0-9-]*[a-zA-Z0-9]).*$`)
	res := re.FindStringSubmatch(line)

	if len(res) != 2 {
		return "", errNoHostInLine
	}
	return res[1], nil
}
