package cache

import (
	"errors"
	jobvacancies "job_vacancies/internal/job_vacancies"
	"sync"
	"time"
)

type InMemoryCache struct {
	mu    sync.RWMutex
	cache map[string]cacheItem //map of id and slice of jobs
	ttl   time.Duration
}

func NewInMemoryCache(ttl time.Duration) *InMemoryCache {
	return &InMemoryCache{
		cache: make(map[string]cacheItem),
		ttl:   ttl,
	}
}

func (c *InMemoryCache) StartCleanup(interval time.Duration) {
	go func() {
		for {
			time.Sleep(interval)
			c.cleanup()
		}
	}()
}

func (c *InMemoryCache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()

	for id, item := range c.cache {
		if now.Sub(item.lastAccessed) > c.ttl {
			delete(c.cache, id)
		}
	}

}

func (c *InMemoryCache) Set(id string, data []jobvacancies.Job) {
	c.mu.Lock()
	defer c.mu.Unlock()
	item := cacheItem{jobs: data, lastAccessed: time.Now()}
	c.cache[id] = item

}

func (c *InMemoryCache) Get(id string) ([]jobvacancies.Job, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	item, ok := c.cache[id]
	if !ok {
		return nil, errors.New("cache id doesnt exist")
	}
	if !item.validateTTL(c.ttl) {
		delete(c.cache, id)
		return nil, errors.New("search expired")
	}
	item.lastAccessed = time.Now()
	c.cache[id] = item
	return c.cache[id].jobs, nil
}

type cacheItem struct {
	jobs         []jobvacancies.Job
	lastAccessed time.Time
}

func (i cacheItem) validateTTL(ttl time.Duration) bool {
	if time.Since(i.lastAccessed) < ttl {
		return true
	}
	return false
}
