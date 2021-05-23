package blocker

func (b *blocker) IsIPBlocked(ip string) bool {
	b.mu.RLock()
	_, found := b.ips[ip]
	b.mu.RUnlock()
	return found
}
