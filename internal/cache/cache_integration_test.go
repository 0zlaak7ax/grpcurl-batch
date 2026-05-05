package cache_test

import (
	"sync"
	"testing"
	"time"

	"github.com/example/grpcurl-batch/internal/cache"
)

// TestConcurrentSetGet verifies that concurrent writers and readers do not
// race (run with -race to validate).
func TestConcurrentSetGet(t *testing.T) {
	c := cache.New(time.Minute)
	const workers = 20

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			key := cache.Key("host:443", "Svc/M", string(rune('a'+id%26)))
			c.Set(key, "result", nil)
			_ = c.Get(key)
		}(i)
	}
	wg.Wait()

	if c.Len() == 0 {
		t.Error("expected at least one cached entry after concurrent writes")
	}
}

// TestPurge_DoesNotRemoveValid ensures Purge leaves live entries intact.
func TestPurge_DoesNotRemoveValid(t *testing.T) {
	c := cache.New(time.Hour)
	key := cache.Key("host:443", "Svc/M", "body")
	c.Set(key, "ok", nil)

	c.Purge()

	if c.Len() != 1 {
		t.Errorf("want 1 entry after purge of non-expired, got %d", c.Len())
	}
	if e := c.Get(key); e == nil {
		t.Error("valid entry was removed by Purge")
	}
}

// TestEntry_IsExpired checks the helper directly.
func TestEntry_IsExpired(t *testing.T) {
	now := time.Now()
	live := &cache.Entry{ExpiresAt: now.Add(time.Minute)}
	if live.IsExpired() {
		t.Error("entry should not be expired")
	}

	dead := &cache.Entry{ExpiresAt: now.Add(-time.Second)}
	if !dead.IsExpired() {
		t.Error("entry should be expired")
	}
}
