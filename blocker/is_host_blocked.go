package blocker

func (b *blocker) IsHostBlocked(host string) bool {
	b.mu.RLock()
	_, found := b.hosts[host]
	b.mu.RUnlock()
	return found
}
