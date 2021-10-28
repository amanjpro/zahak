package engine

// import (
// 	"fmt"
// )

type CachedEval struct {
	Key  uint64 // 8
	Data uint64 // 8
}

type NodeType uint8

const (
	Exact      NodeType = 1 << iota // PV-Node
	UpperBound                      // All-Node
	LowerBound                      // Cut-Node
)

type Cache struct {
	items    []CachedEval
	size     int
	consumed int
	length   uint64
}

const OldAge = uint16(5)
const CACHE_ENTRY_SIZE = 8 + 8
const DEFAULT_CACHE_SIZE = 128
const MAX_CACHE_SIZE = 24000

const MOVE_MASK uint64 = 0b1111111111111111111111111111 // move << 0, 28 bits
const EVAL_MASK uint64 = 0b1111111111111111             // eval << 28, 16 bits
const DEPTH_MASK uint64 = 0b1111111                     // depth << 44, 7 bits
const TYPE_MASK uint64 = 0b111                          // type << 51, 3 bits
const AGE_MASK uint64 = 0b1111111111                    // age << 54, 10 bits

func Pack(hashmove Move, eval int16, depth int8, nodeType NodeType, age uint16) uint64 {
	return (uint64(hashmove) & MOVE_MASK) |
		((uint64(eval) & EVAL_MASK) << 28) |
		((uint64(depth) & DEPTH_MASK) << 44) |
		((uint64(nodeType) & TYPE_MASK) << 51) |
		((uint64(age) & AGE_MASK) << 54)
}

func Unpack(data uint64) (hashmove Move, eval int16, depth int8, nodeType NodeType, age uint16) {
	hashmove = Move(data & MOVE_MASK)
	eval = int16((data >> 28) & EVAL_MASK)
	depth = int8((data >> 44) & DEPTH_MASK)
	nodeType = NodeType((data >> 51) & TYPE_MASK)
	age = uint16((data >> 54) & AGE_MASK)
	return
}

func (c *CachedEval) Update(key uint64, data uint64) {
	c.Key = key
	c.Data = data
}

func (c *Cache) Consumed() int {
	return int((float64(c.consumed) / float64(len(c.items))) * 1000)
}

func (c *Cache) index(hash uint64) int {
	return int(hash>>32) % len(c.items)
}

func (c *Cache) Set(hash uint64, hashmove Move, eval int16, depth int8, nodeType NodeType, age uint16) {
	// if hashmove == EmptyMove {
	// 	return
	// }
	index := c.index(hash)

	oldValue := c.items[index]
	oldKey := oldValue.Key
	oldData := oldValue.Data

	newData := Pack(hashmove, eval, depth, nodeType, age)
	// very good for debugging hash issues
	// newHashmove, newEval, newDepth, newNodeType, newAge := Unpack(newData)
	// if hashmove != newHashmove || eval != newEval || depth != newDepth || nodeType != newNodeType || age != newAge {
	// 	panic(fmt.Sprintf(
	// 		"Culprits are: %d %d %d %d %d\nSomehow became: %d %d %d %d %d\n", hashmove, eval, depth, nodeType, age, newHashmove, newEval, newDepth, newNodeType, newAge))
	// }

	newKey := newData ^ hash

	if oldData != 0 {
		_, _, oldDepth, oldType, oldAge := Unpack(oldData)
		if (hash ^ oldData) == oldKey {
			c.items[index].Update(newKey, newData)
			return
		}
		if age-oldAge >= OldAge {
			c.items[index].Update(newKey, newData)
			return
		}
		if oldDepth > depth {
			return
		}
		if oldType == Exact || nodeType != Exact {
			return
		} else if nodeType == Exact {
			c.items[index].Update(newKey, newData)
			return
		}
		c.items[index].Update(newKey, newData)
	} else {
		c.consumed += 1
		c.items[index].Update(newKey, newData)
	}
}

func (c *Cache) Size() int {
	return c.size
}

func (c *Cache) Get(hash uint64) (Move, int16, int8, NodeType, bool) {
	index := c.index(hash)
	value := c.items[index]
	data := value.Data
	key := value.Key
	ok := hash == (key ^ data)
	if ok {
		move, eval, depth, nType, _ := Unpack(data)
		return move, eval, depth, nType, true
	}
	return 0, 0, 0, 0, false
}

func NewCache(megabytes int) *Cache {
	if megabytes > MAX_CACHE_SIZE || megabytes < 1 {
		return nil
	}
	size := int((megabytes * 1024 * 1024) / CACHE_ENTRY_SIZE)
	length := RoundPowerOfTwo(size)

	return &Cache{make([]CachedEval, length), megabytes, 0, uint64(length)}
}

func RoundPowerOfTwo(size int) int {
	var x = 1
	for (x << 1) <= size {
		x <<= 1
	}
	return x
}
