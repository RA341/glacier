package cache

import (
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

// Entry CacheEntry holds cached data with expiration
type Entry[T any] struct {
	Data      T
	ExpiresAt time.Time
}

// TTLCache manages cached data
type TTLCache[T any] struct {
	mu         sync.RWMutex
	cache      map[string]Entry[T]
	ttl        time.Duration
	cancelChan chan struct{}
}

// NewStatsCache creates a new cache with specified TTL
func NewStatsCache[T any](ttl time.Duration) *TTLCache[T] {
	cache := &TTLCache[T]{
		cache:      make(map[string]Entry[T]),
		ttl:        ttl,
		mu:         sync.RWMutex{},
		cancelChan: make(chan struct{}),
	}

	go cache.Cleanup()

	return cache
}

// Get retrieves cached stats if not expired
func (c *TTLCache[T]) Get(key string) (T, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.cache[key]
	if !exists || time.Now().After(entry.ExpiresAt) {
		// Declared var zero T to get appropriate zero value for any type T
		var zero T
		return zero, false
	}

	return entry.Data, true
}

// Set stores stats in cache with TTL
func (c *TTLCache[T]) Set(key string, data T) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[key] = Entry[T]{
		Data:      data,
		ExpiresAt: time.Now().Add(c.ttl),
	}
}

func (c *TTLCache[T]) StartCleaner() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.Cleanup()
		case <-c.cancelChan:
			log.Debug().Msg("stopping cache cleaner")
			return
		}
	}
}

func (c *TTLCache[T]) Close() {
	c.cancelChan <- struct{}{}
}

// Cleanup removes expired entries (call periodically)
func (c *TTLCache[T]) Cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, entry := range c.cache {
		if now.After(entry.ExpiresAt) {
			delete(c.cache, key)
		}
	}
}
