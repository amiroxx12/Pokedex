package internal

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val []byte
}
type Cache struct {
	interval time.Duration
	entries map[string]cacheEntry
	mu sync.Mutex
}

func NewCache(d time.Duration) *Cache {
	c := &Cache{
		interval: d,
		entries: make(map[string]cacheEntry),
	}
	go c.reapLoop(d)
	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = cacheEntry{
		createdAt : time.Now(),
		val: val,
	}
}
func (c *Cache) Get(key string) ([]byte,bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, exists := c.entries[key]
	if !exists {
		return nil,false
	}
	return entry.val,true
}
func (c *Cache) reapLoop(d time.Duration) {
	ticker := time.NewTicker(d)
	for {
		<- ticker.C
		now := time.Now()
		c.mu.Lock()
		for k,entry := range c.entries {
			if now.Sub(entry.createdAt) > d {
				delete(c.entries,k)
			}
		}
		c.mu.Unlock()
	}
}