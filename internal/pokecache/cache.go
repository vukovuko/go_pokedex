package pokecache

import (
	"sync"
	"time"
)

// cacheEntry holds a cached value along with the time it was created.
type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

// Cache is a thread-safe cache that stores responses as raw bytes.
type Cache struct {
	interval time.Duration
	m        map[string]cacheEntry
	mu       sync.Mutex
}

// NewCache creates a new Cache with a configurable expiration interval.
// It starts a background reapLoop to remove expired entries.
func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		interval: interval,
		m:        make(map[string]cacheEntry),
	}
	go c.reapLoop()
	return c
}

// Add stores a value in the cache for the given key.
func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

// Get retrieves a cached value for the given key.
// It returns the value and true if found and not expired; otherwise, it returns nil and false.
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, exists := c.m[key]
	if !exists {
		return nil, false
	}
	if time.Since(entry.createdAt) > c.interval {
		delete(c.m, key)
		return nil, false
	}
	return entry.val, true
}

// reapLoop runs periodically (every c.interval) to remove expired entries.
func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	for {
		<-ticker.C
		c.mu.Lock()
		for key, entry := range c.m {
			if time.Since(entry.createdAt) > c.interval {
				delete(c.m, key)
			}
		}
		c.mu.Unlock()
	}
}
