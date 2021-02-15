package cache

type QuiescenceEval struct {
	Hash uint64
	Eval int
}

type QuiesceneCache struct {
	items []*QuiescenceEval
	size  uint32
}

var QCache QuiesceneCache

func (c *QuiesceneCache) Set(hash uint64, value *QuiescenceEval) {
	key := uint32(hash>>32) % c.size
	c.items[key] = value
}

func (c *QuiesceneCache) Get(hash uint64) (*QuiescenceEval, bool) {
	key := uint32(hash>>32) % c.size
	item := c.items[key]
	if item != nil && item.Hash == hash {
		return item, true
	}
	return nil, false
}

func NewQCache(megabytes uint32) {
	size := uint32(1000000) // megabytes * 1024 * 1024 / dummySize
	items := make([]*QuiescenceEval, size)
	QCache = QuiesceneCache{items, size} //s, current: 0}
}
