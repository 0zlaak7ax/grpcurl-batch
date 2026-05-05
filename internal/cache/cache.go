// Package cache provides a simple in-memory result cache for deduplicating
// identical gRPC requests within a single batch run.
package cache

import (
	"crypto/sha256"
	"fmt"
	"sync"
	"time"
)

// Entry holds a cached result and its expiry time.
type Entry struct {
	Output    string
	Err       error
	CachedAt  time.Time
	ExpiresAt time.Time
}

// IsExpired reports whether the entry has passed its TTL.
func (e *Entry) IsExpired() bool {
	return time.Now().After(e.ExpiresAt)
}

// Cache is a thread-safe in-memory store keyed by request fingerprint.
type Cache struct {
	mu      sync.RWMutex
	entries map[string]*Entry
	ttl     time.Duration
}

// New creates a Cache with the given TTL. A zero TTL disables caching.
func New(ttl time.Duration) *Cache {
	return &Cache{
		entries: make(map[string]*Entry),
		ttl:     ttl,
	}
}

// Key derives a cache key from the target address, method and request body.
func Key(address, method, body string) string {
	h := sha256.New()
	h.Write([]byte(address + "\x00" + method + "\x00" + body))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// Get returns the cached entry for key, or nil if absent or expired.
func (c *Cache) Get(key string) *Entry {
	if c.ttl == 0 {
		return nil
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.entries[key]
	if !ok || e.IsExpired() {
		return nil
	}
	return e
}

// Set stores a result under key.
func (c *Cache) Set(key, output string, err error) {
	if c.ttl == 0 {
		return
	}
	now := time.Now()
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = &Entry{
		Output:    output,
		Err:       err,
		CachedAt:  now,
		ExpiresAt: now.Add(c.ttl),
	}
}

// Len returns the number of entries currently stored (including expired).
func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.entries)
}

// Purge removes all expired entries.
func (c *Cache) Purge() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, e := range c.entries {
		if e.IsExpired() {
			delete(c.entries, k)
		}
	}
}
