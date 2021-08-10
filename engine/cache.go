package engine

type CachedEval struct {
	Key  uint64 // 8
	Data uint64 // 8
}

const MOVE_MASK uint64 = 0x000000000FFFFFFF  // move << 0
const EVAL_MASK uint64 = 0x00000FFFF0000000  // eval << 28
const DEPTH_MASK uint64 = 0x0007f00000000000 // depth << 44
const TYPE_MASK uint64 = 0x0038000000000000  // type << 51
const AGE_MASK uint64 = 0xffc0000000000000   // age << 54

func Pack(hashmove Move, eval int16, depth int8, nodeType NodeType, age uint16) uint64 {
	var data uint64
	data = uint64(hashmove) |
		uint64(eval)<<28 |
		uint64(depth)<<44 |
		uint64(nodeType)<<51 |
		uint64(age)<<54
	return data
}

func Unpack(data uint64) (hashmove Move, eval int16, depth int8, nodeType NodeType, age uint16) {
	hashmove = Move(data & MOVE_MASK)
	eval = int16((data & EVAL_MASK) >> 28)
	depth = int8((data & DEPTH_MASK) >> 44)
	nodeType = NodeType((data & TYPE_MASK) >> 51)
	age = uint16((data & AGE_MASK) >> 54)
	return
}

func (c *Cache) Update(index int, key uint64, data uint64) {
	c.items[index].Key = key
	c.items[index].Data = data
}

type NodeType uint8

const (
	Exact      NodeType = 1 << iota // PV-Node
	UpperBound                      // All-Node
	LowerBound                      // Cut-Node
)

var OldAge = uint16(5)

var EmptyEval = CachedEval{0, 0}
var CACHE_ENTRY_SIZE = uint32(8 + 8)

type Cache struct {
	items    []CachedEval
	size     uint32
	consumed int
	count    int
}

const DEFAULT_CACHE_SIZE = uint32(128)
const MAX_CACHE_SIZE = uint32(24000)

func (c *Cache) Consumed() int {
	return int((float64(c.consumed) / float64(len(c.items))) * 1000)
}

func (c *Cache) index(hash uint64) int {
	// return uint64(uint64(uint32(hash))*c.count) >> 32
	return int(hash>>32) % c.count
}

func (c *Cache) Set(hash uint64, hashmove Move, eval int16, depth int8, nodeType NodeType, age uint16) {
	index := c.index(hash)
	oldValue := c.items[index]
	data := Pack(hashmove, eval, depth, nodeType, age)
	entryHash := (oldValue.Key ^ oldValue.Data)
	key := hash ^ data

	if oldValue != EmptyEval {
		_, _, entryDepth, entryType, entryAge := Unpack(oldValue.Data)
		if hash == entryHash {
			c.Update(index, key, data)
		} else if age-entryAge >= OldAge {
			c.Update(index, key, data)
		} else if entryDepth <= depth {
			if entryType == Exact || nodeType != Exact {
				return
			} else if nodeType == Exact {
				c.Update(index, key, data)
			} else {
				c.Update(index, key, data)
			}
		}
	} else {
		c.consumed += 1
		c.Update(index, key, data)
	}
}

func (c *Cache) Size() uint32 {
	return c.size
}

func (c *Cache) Get(hash uint64) (move Move, eval int16, depth int8, nType NodeType, ok bool) {
	index := c.index(hash)
	value := c.items[index]
	ok = hash^value.Data == value.Key
	if ok {
		move, eval, depth, nType, _ = Unpack(value.Data)
	}
	return
}

func NewCache(megabytes uint32) *Cache {
	if megabytes > MAX_CACHE_SIZE {
		return nil
	}
	size := int((megabytes * 1024 * 1024) / CACHE_ENTRY_SIZE)
	length := RoundPowerOfTwo(size)

	return &Cache{make([]CachedEval, length), megabytes, 0, length}
}

func RoundPowerOfTwo(size int) int {
	var x = 1
	for (x << 1) <= size {
		x <<= 1
	}
	return x
}
