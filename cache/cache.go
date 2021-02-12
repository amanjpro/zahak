package cache

type CachedEval struct {
	Hash  uint64
	Eval  int
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
	items map[uint32]*CachedEval
}

var TranspositionTable Cache

func (c *Cache) Set(hash uint64, value *CachedEval) {
	key := uint32(hash)
	oldValue, ok := c.items[key]
	if ok {
		if value.Age-oldValue.Age >= oldAge {
			c.items[key] = value
		}
		if oldValue.Type == Exact || value.Type != Exact {
			return
		} else if value.Type == Exact {
			c.items[key] = value
		}
		if oldValue.Depth > value.Depth {
			return
		}
		c.items[key] = value
	} else {
		c.items[key] = value
	}
}

func (c *Cache) Get(hash uint64) (*CachedEval, bool) {
	key := uint32(hash)
	item, found := c.items[key]
	if found && item.Hash == hash {
		return item, found
	}
	return nil, false
}

func ResetCache() {
	items := make(map[uint32]*CachedEval, 100_000_000)
	TranspositionTable = Cache{items: items} //s, current: 0}
}
