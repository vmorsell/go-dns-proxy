package main

type ResolveStrategy int

const (
	CACHE = iota
	SERVER
)

func (s ResolveStrategy) String() string {
	return [...]string{"CACHE", "SERVER"}[s]
}
