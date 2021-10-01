//go:build !amd64
// +build !amd64

package engine

func (c *Cache) Prefetch(hash uint64) {}
