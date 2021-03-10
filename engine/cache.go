package engine

type CachedEval struct {
	Hash  uint64
	Eval  int16
	Depth int8
	Type  NodeType
	Age   uint16
}

type NodeType uint8

const (
	Exact      NodeType = 1 << iota // PV-Node
	UpperBound                      // All-Node
	LowerBound                      // Cut-Node
)

var oldAge = uint16(5)

const CACHE_ENTRY_SIZE = uint32(64+16+8+8+16) / 8

type Cache struct {
	items    []CachedEval
	size     uint32
	consumed int
}

var EmptyEval = CachedEval{0, 0, 0, 0, 0}

func (c *Cache) Consumed() int {
	return int((float64(c.consumed) / float64(len(c.items))) * 1000)
}

var EmptyCache = Cache{nil, 0, 0}
var TranspositionTable Cache = EmptyCache

func (c *Cache) hash(key uint64) uint32 {
	return uint32(key>>32) % c.size
}

func (c *Cache) Set(hash uint64, eval int16, depth int8, nodeType NodeType, age uint16) {
	key := c.hash(hash)
	oldValue := c.items[key]
	if oldValue != EmptyEval {
		if hash == oldValue.Hash {
			c.items[key] = CachedEval{hash, eval, depth, nodeType, age}
			return
		}
		if age-oldValue.Age >= oldAge {
			c.items[key] = CachedEval{hash, eval, depth, nodeType, age}
			return
		}
		if oldValue.Depth > depth {
			return
		}
		if oldValue.Type == Exact || nodeType != Exact {
			return
		} else if nodeType == Exact {
			c.items[key] = CachedEval{hash, eval, depth, nodeType, age}
			return
		}
		c.items[key] = CachedEval{hash, eval, depth, nodeType, age}
	} else {
		c.consumed += 1
		c.items[key] = CachedEval{hash, eval, depth, nodeType, age}
	}
}

func (c *Cache) Get(hash uint64) (int16, int8, NodeType, bool) {
	key := c.hash(hash)
	item := &c.items[key]
	if item.Hash == hash {
		return item.Eval, item.Depth, item.Type, true
	}
	return 0, 0, 0, false
}

func NewCache(megabytes uint32) {
	size := megabytes * 1024 * 1024 / CACHE_ENTRY_SIZE
	items := make([]CachedEval, size)
	TranspositionTable = Cache{items, uint32(size), 0} //s, current: 0}
	for i := 0; i < int(size); i++ {
		TranspositionTable.items[i] = EmptyEval
	}
}

func ResetCache() {
	if TranspositionTable.size != EmptyCache.size {
		TranspositionTable.items = make([]CachedEval, TranspositionTable.size)
		TranspositionTable.consumed = 0
		for i := 0; i < int(TranspositionTable.size); i++ {
			TranspositionTable.items[i] = EmptyEval
		}
	} else {
		NewCache(10)
	}
}
