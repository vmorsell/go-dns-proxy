package cache

import (
	"strings"

	"github.com/miekg/dns"
)

type ErrNoQuestions struct{}

func (e ErrNoQuestions) Error() string { return "no questions in message" }

func (c *cache) Key(r *dns.Msg) (string, error) {
	if len(r.Question) == 0 {
		return "", ErrNoQuestions{}
	}
	var hosts []string
	for _, rr := range r.Question {
		hosts = append(hosts, rr.Name)
	}
	return strings.Join(hosts, ","), nil
}
