//go:build amd64
// +build amd64

package engine

import (
	"unsafe"
)

func _prefetch(item unsafe.Pointer)

func (c *Cache) Prefetch(hash uint64) {
	index := c.index(hash)
	p := unsafe.Pointer(&c.items[index])

	_prefetch(p)
}
