package cache

import (
	"context"
	"sync"
	"time"
)

type InMemoryCache[T any] struct {
	mu    sync.RWMutex
	cache map[string]cacheItem[T]
	ttl   time.Duration
}

func NewInMemoryCache[T any](ttl time.Duration) *InMemoryCache[T] {
	return &InMemoryCache[T]{
		cache: make(map[string]cacheItem[T]),
		ttl:   ttl,
	}
}

func (c *InMemoryCache[T]) StartCleanup(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				c.cleanup()
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (c *InMemoryCache[T]) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()

	for id, item := range c.cache {
		if now.Sub(item.lastAccessed) > c.ttl {
			delete(c.cache, id)
		}
	}
}

func (c *InMemoryCache[T]) Set(id string, data T) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[id] = cacheItem[T]{
		data:         data,
		lastAccessed: time.Now(),
	}
}

func (c *InMemoryCache[T]) Get(id string) (T, bool) {
	c.mu.RLock()
	item, ok := c.cache[id]
	c.mu.RUnlock()

	if !ok {
		var zero T
		return zero, false
	}

	if time.Since(item.lastAccessed) > c.ttl {
		c.mu.Lock()
		delete(c.cache, id)
		c.mu.Unlock()

		var zero T
		return zero, false
	}

	// update last accessed (sliding TTL)
	c.mu.Lock()
	item.lastAccessed = time.Now()
	c.cache[id] = item
	c.mu.Unlock()

	return item.data, true
}

type cacheItem[T any] struct {
	data         T
	lastAccessed time.Time
}
