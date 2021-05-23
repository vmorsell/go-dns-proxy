package blocker

func (b *blocker) AddIP(ip string) error {
	b.mu.Lock()
	b.ips[ip] = struct{}{}
	b.mu.Unlock()
	return nil
}
