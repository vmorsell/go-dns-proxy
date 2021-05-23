package blocker

func (b *blocker) BlockedIPsLength() int {
	return len(b.ips)
}

func (b *blocker) BlockedHostsLength() int {
	return len(b.ips)
}
