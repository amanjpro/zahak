package cache

type StaticEval struct {
	Hash uint64
	Eval int32
}

type StaticEvalCache struct {
	items []StaticEval
	size  uint32
}

var EmptyStaticEval = StaticEval{0, 0}

var EvalTable StaticEvalCache = NewStaticEvalCache()

func (c *StaticEvalCache) hash(key uint64) uint32 {
	return uint32(key>>32) % c.size
}

func (c *StaticEvalCache) Set(hash uint64, value int32) {
	key := c.hash(hash)
	c.items[key] = StaticEval{hash, value}
}

func (c *StaticEvalCache) Get(hash uint64) (int32, bool) {
	key := c.hash(hash)
	item := c.items[key]
	if item.Hash == hash {
		return item.Eval, true
	}
	return -1, false
}

func NewStaticEvalCache() StaticEvalCache {
	size := 10 * 1024 * 1024 / (64 + 32)
	items := make([]StaticEval, size)
	cache := StaticEvalCache{items, uint32(size)}
	for i := 0; i < int(size); i++ {
		cache.items[i] = EmptyStaticEval
	}
	return cache
}

func ResetEvalCache() {
	EvalTable = NewStaticEvalCache()
}
