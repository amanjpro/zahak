package cache

type CachedEval struct {
	Hash  uint64
	Eval  int16
	Depth int8
	Type  NodeType
	Age   uint16
}

type NodeType uint16

const (
	Exact      NodeType = 1 << iota // PV-Node
	UpperBound                      // All-Node
	LowerBound                      // Cut-Node
)

var oldAge = uint16(5)

type Cache struct {
	items []*CachedEval
	size  uint32
}

var TranspositionTable Cache

func (c *Cache) Set(hash uint64, value *CachedEval) {
	key := uint32(hash>>32) % c.size
	oldValue := c.items[key]
	if oldValue != nil {
		if value.Hash == oldValue.Hash {
			c.items[key] = value
		}
		if value.Age-oldValue.Age >= oldAge {
			c.items[key] = value
		}
		if oldValue.Depth > value.Depth {
			return
		}
		if oldValue.Type == Exact || value.Type != Exact {
			return
		} else if value.Type == Exact {
			c.items[key] = value
		}
		c.items[key] = value
	} else {
		c.items[key] = value
	}
}

func (c *Cache) Get(hash uint64) (*CachedEval, bool) {
	key := uint32(hash>>32) % c.size
	item := c.items[key]
	if item != nil && item.Hash == hash {
		return item, true
	}
	return nil, false
}

func NewCache(megabytes uint32) {
	dummySize := uint32(1)
	size := megabytes * 1024 * 1024 / dummySize
	items := make([]*CachedEval, size)
	TranspositionTable = Cache{items, size} //s, current: 0}
}

func ResetCache() {
	items := make([]*CachedEval, 100_000_000)
	TranspositionTable = Cache{items, 100_000_000} //s, current: 0}
}
