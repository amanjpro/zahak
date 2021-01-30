package main

import (
	"sync"

	"github.com/notnil/chess"
)

type CachedEval struct {
	eval float64
	line []chess.Move
}

type Cache struct {
	mu    sync.Mutex
	items map[uint64]CachedEval
}

var evalCache = Cache{items: make(map[uint64]CachedEval)}

func (c *Cache) Set(key uint64, value CachedEval) {
	// Lock so only one goroutine at a time can access the map c.v.
	evalCache.mu.Lock()
	c.items[key] = value
	evalCache.mu.Unlock()
}

func (c *Cache) Get(key uint64) (CachedEval, bool) {
	// Lock so only one goroutine at a time can access the map c.v.
	evalCache.mu.Lock()
	item, found := c.items[key]
	evalCache.mu.Unlock()
	return item, found
}
