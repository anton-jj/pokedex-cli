package internal
import ("Time")

type Cache struct {
	items map[string]CacheEntry 
}
type CacheEntry  struct {
	createdAt time.Time
	val []byte
}

func (c *Cache) Add(key string, val []byte) {
	c.items[key] = CacheEntry{time.Now(), val}	
}
func (c *Cache) Get(key string) ([]byte, bool) {
	item, exists := c.items[key]
	if exists {
	return item.val, exists
	}
	return nil, exists
}
