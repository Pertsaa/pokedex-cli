package pokecache

import (
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	entries  map[string]cacheEntry
	mu       sync.RWMutex
	duration time.Duration
}

type cacheEntry struct {
	val       []byte
	createdAt time.Time
}

func New(cacheDuration time.Duration) *Cache {
	fmt.Println(cacheDuration)
	c := &Cache{
		entries:  make(map[string]cacheEntry),
		duration: cacheDuration,
	}
	go c.reapLoop()
	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = cacheEntry{
		val:       val,
		createdAt: time.Now(),
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.entries[key]
	if !ok {
		return nil, false
	}

	return entry.val, true
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.duration)
	defer ticker.Stop()
	for {
		<-ticker.C
		for key, entry := range c.entries {
			if time.Since(entry.createdAt) > c.duration {
				c.mu.Lock()
				delete(c.entries, key)
				c.mu.Unlock()
			}
		}
	}
}
