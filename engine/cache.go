package engine

type CachedEval struct {
	Key  uint64 // 8
	Data uint64 // 8
}

const MOVE_MASK uint64 = 0x000000000FFFFFFF  // move << 0
const EVAL_MASK uint64 = 0x00000FFFF0000000  // eval << 28
const DEPTH_MASK uint64 = 0x0007F00000000000 // depth << 44
const TYPE_MASK uint64 = 0x0038000000000000  // type << 51
const AGE_MASK uint64 = 0xFFC0000000000000   // age << 54

func Pack(hashmove Move, eval int16, depth int8, nodeType NodeType, age uint16) uint64 {
	return uint64(hashmove) |
		uint64(eval)<<28 |
		uint64(depth)<<44 |
		uint64(nodeType)<<51 |
		uint64(age)<<54
}

func Unpack(data uint64) (hashmove Move, eval int16, depth int8, nodeType NodeType, age uint16) {
	hashmove = Move(data & MOVE_MASK)
	eval = int16((data & EVAL_MASK) >> 28)
	depth = int8((data & DEPTH_MASK) >> 44)
	nodeType = NodeType((data & TYPE_MASK) >> 51)
	age = uint16((data & AGE_MASK) >> 54)
	return
}

func (c *CachedEval) Update(key uint64, data uint64) {
	c.Key = key
	c.Data = data
}

type NodeType uint8

const (
	Exact      NodeType = 1 << iota // PV-Node
	UpperBound                      // All-Node
	LowerBound                      // Cut-Node
)

var OldAge = uint16(5)

var CACHE_ENTRY_SIZE = uint32(8 + 8)

type Cache struct {
	items    []CachedEval
	size     uint32
	consumed int
	length   uint64
}

const DEFAULT_CACHE_SIZE = uint32(128)
const MAX_CACHE_SIZE = uint32(24000)

func (c *Cache) Consumed() int {
	return int((float64(c.consumed) / float64(len(c.items))) * 1000)
}

func (c *Cache) index(hash uint64) uint32 {
	return uint32((uint64(uint32(hash)) * c.length) >> 32)

	// return uint32(hash>>32) % uint32(len(c.items))
}

func (c *Cache) Set(hash uint64, hashmove Move, eval int16, depth int8, nodeType NodeType, age uint16) {
	if hashmove == EmptyMove {
		return
	}
	index := c.index(hash)
	entry := c.items[index]

	oldKey := entry.Key
	oldData := entry.Data
	entryHash := (oldKey ^ oldData)

	newData := Pack(hashmove, eval, depth, nodeType, age)
	newKey := hash ^ newData

	if oldData != 0 {
		_, _, entryDepth, entryType, entryAge := Unpack(oldData)
		if hash == entryHash {
			c.items[index].Update(newKey, newData)
			return
		}
		if age-entryAge >= OldAge {
			c.items[index].Update(newKey, newData)
			return
		}
		if entryDepth > depth {
			return
		}
		if entryType == Exact || nodeType != Exact {
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

func (c *Cache) Size() uint32 {
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

func NewCache(megabytes uint32) *Cache {
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
