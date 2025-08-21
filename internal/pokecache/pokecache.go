package internal

import (
	"time"
	"sync"
)

type Cache struct {
	mu sync.Mutex
	items map[string]CacheEntry 
	interval time.Duration
}
type CacheEntry  struct {
	createdAt time.Time
	val []byte
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	c.items[key] = CacheEntry{time.Now(), val}	
	c.mu.Unlock()
}
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	item, exists := c.items[key]
	if exists {
	return item.val, exists
	}
	return nil, exists
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	for {
		// time.Sleep(c.interval)
		<- ticker.C
		c.mu.Lock()
		for k, v := range c.items {
			if time.Since(v.createdAt) > c.interval {
				delete(c.items, k)
			}
		}
		c.mu.Unlock()
	}
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		items: make(map[string]CacheEntry),
		interval: interval,
	}
	go c.reapLoop() 
	return c
}
