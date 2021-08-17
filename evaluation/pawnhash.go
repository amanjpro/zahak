package evaluation

import (
	. "github.com/amanjpro/zahak/engine"
)

type PawnEval struct {
	Hash       uint64 // 8
	Middlegame int16  // 2
	Endgame    int16  // 2
}

func (c *PawnEval) Update(hash uint64, mg int16, eg int16) {
	c.Hash = hash
	c.Middlegame = mg
	c.Endgame = eg
}

var PAWN_ENTRY_SIZE = 8 + 2 + 2

type PawnCache struct {
	items          []PawnEval
	size           int
	PawnhashMisses int64
	PawnhashHits   int64
}

const DEFAULT_PAWNHASH_SIZE = 2
const MAX_PAWNHASH_SIZE = 10

func (c *PawnCache) hash(key uint64) uint32 {
	return uint32(key>>32) % uint32(len(c.items))
}

func (c *PawnCache) Set(hash uint64, mg int16, eg int16) {
	key := c.hash(hash)
	c.items[key].Update(hash, mg, eg)
}

func (c *PawnCache) Size() int {
	return c.size
}

func (c *PawnCache) Get(hash uint64) (int16, int16, bool) {
	key := c.hash(hash)
	item := c.items[key]
	if item.Hash == hash {
		c.PawnhashHits += 1
		return item.Middlegame, item.Endgame, true
	}
	c.PawnhashMisses += 1
	return 0, 0, false
}

func NewPawnCache(megabytes int) *PawnCache {
	size := int(megabytes * 1024 * 1024 / PAWN_ENTRY_SIZE)
	return &PawnCache{make([]PawnEval, RoundPowerOfTwo(size)), 1, 0, 0}
}
