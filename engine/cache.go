package engine

import (
	"sync/atomic"
)

type CachedEval struct {
	LowKey   uint32   // 4
	Gate     int32    // 4
	HashMove Move     // 4
	Eval     int16    // 2
	Age      uint16   // 2
	Depth    int8     // 1
	Type     NodeType // 1
}

func (c *CachedEval) Update(lowKey uint32, hashmove Move, eval int16, depth int8, nodeType NodeType, age uint16) {
	c.LowKey = lowKey
	c.HashMove = hashmove
	c.Eval = eval
	c.Depth = depth
	c.Type = nodeType
	c.Age = age
}

type NodeType uint8

const (
	Exact      NodeType = 1 << iota // PV-Node
	UpperBound                      // All-Node
	LowerBound                      // Cut-Node
)

var oldAge = uint16(5)

var EmptyEval = CachedEval{0, 1, EmptyMove, 0, 0, 0, 0}
var CACHE_ENTRY_SIZE = uint32(4 + 4 + 4 + 2 + 2 + 1 + 1)

type Cache struct {
	items    []CachedEval
	size     uint32
	consumed int
}

const DEFAULT_CACHE_SIZE = uint32(128)
const MAX_CACHE_SIZE = uint32(24000)

func (c *Cache) Consumed() int {
	return int((float64(c.consumed) / float64(len(c.items))) * 1000)
}

func (c *Cache) index(hash uint64) uint32 {
	return uint32(hash) % uint32(len(c.items))
}

func (c *Cache) Set(hash uint64, hashmove Move, eval int16, depth int8, nodeType NodeType, age uint16) {
	index := c.index(hash)
	key := uint32(hash >> 32)
	oldValue := &c.items[index]

	if atomic.CompareAndSwapInt32(&oldValue.Gate, 0, 1) {
		if *oldValue != EmptyEval {
			if key == oldValue.LowKey {
				oldValue.Update(key, hashmove, eval, depth, nodeType, age)
			} else if age-oldValue.Age >= oldAge {
				oldValue.Update(key, hashmove, eval, depth, nodeType, age)
			} else if oldValue.Depth <= depth {
				if oldValue.Type == Exact || nodeType != Exact {
				} else if nodeType == Exact {
					oldValue.Update(key, hashmove, eval, depth, nodeType, age)
				} else {
					oldValue.Update(key, hashmove, eval, depth, nodeType, age)
				}
			}
		} else {
			c.consumed += 1
			oldValue.Update(key, hashmove, eval, depth, nodeType, age)
		}
		atomic.StoreInt32(&oldValue.Gate, 0)
	}
}

func (c *Cache) Size() uint32 {
	return c.size
}

func (c *Cache) Get(hash uint64) (move Move, eval int16, depth int8, nType NodeType, ok bool) {
	index := c.index(hash)
	key := uint32(hash >> 32)
	value := &c.items[index]
	if atomic.CompareAndSwapInt32(&value.Gate, 0, 1) {
		if value.LowKey == key {
			move = value.HashMove
			eval = value.Eval
			depth = value.Depth
			nType = value.Type
			ok = true
		}
		atomic.StoreInt32(&value.Gate, 0)
	}
	return
}

func NewCache(megabytes uint32) *Cache {
	if megabytes > MAX_CACHE_SIZE {
		return nil
	}
	size := int(megabytes * 1024 * 1024 / CACHE_ENTRY_SIZE)
	length := RoundPowerOfTwo(size)

	return &Cache{make([]CachedEval, length), megabytes, 0}
}

func RoundPowerOfTwo(size int) int {
	var x = 1
	for (x << 1) <= size {
		x <<= 1
	}
	return x
}
