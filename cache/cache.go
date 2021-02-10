package cache

type CachedEval struct {
	Eval  int
	Depth int8
}

type Cache struct {
	itemss  [5]map[uint64]*CachedEval
	current int
}

var TranspositionTable Cache

func (c *Cache) Rotate() {
	nextCurrent := (c.current + 1) % len(c.itemss)
	c.itemss[nextCurrent] = make(map[uint64]*CachedEval, 1000_000)
	c.current = nextCurrent
}

func (c *Cache) Set(key uint64, value *CachedEval) {
	c.itemss[c.current][key] = value
}

func (c *Cache) Get(key uint64) (*CachedEval, bool) {
	for i, j := 0, c.current; i <= len(c.itemss); i++ {
		item, found := c.itemss[j%len(c.itemss)][key]
		if found {
			return item, found
		}
		j++
	}
	return nil, false
}

func ResetCache() {
	var itemss [5]map[uint64]*CachedEval
	for i := 0; i < len(itemss); i++ {
		itemss[i] = make(map[uint64]*CachedEval, 1000_000)
	}
	TranspositionTable = Cache{itemss: itemss, current: 0}
}
