package proxy

type ResolveStrategy string

const (
	NONE   ResolveStrategy = "None"
	CACHE  ResolveStrategy = "Cache"
	SERVER ResolveStrategy = "Server"
)

func (s ResolveStrategy) String() string {
	return string(s)
}
