package blocker

func (b *blocker) AddHost(host string) error {
	b.mu.Lock()
	b.hosts[host] = struct{}{}
	b.mu.Unlock()
	return nil
}
