package cache

type CachedEval struct {
	Eval  int
	Depth int8
}

type Cache struct {
	items map[uint64]*CachedEval
}

var TranspositionTable Cache

func (c *Cache) Set(key uint64, value *CachedEval) {
	c.items[key] = value
}

func (c *Cache) Get(key uint64) (*CachedEval, bool) {
	item, found := c.items[key]
	if found {
		return item, found
	}
	return nil, false
}

func ResetCache() {
	items := make(map[uint64]*CachedEval, 1000_000)
	TranspositionTable = Cache{items: items} //s, current: 0}
}
