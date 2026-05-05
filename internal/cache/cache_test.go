package cache_test

import (
	"errors"
	"testing"
	"time"

	"github.com/example/grpcurl-batch/internal/cache"
)

func TestNew_ZeroTTL_Disabled(t *testing.T) {
	c := cache.New(0)
	key := cache.Key("host:443", "pkg.Svc/Method", `{}`)
	c.Set(key, "output", nil)
	if got := c.Get(key); got != nil {
		t.Error("expected nil when TTL is zero")
	}
}

func TestGet_MissingKey(t *testing.T) {
	c := cache.New(time.Minute)
	if got := c.Get("nonexistent"); got != nil {
		t.Error("expected nil for missing key")
	}
}

func TestSet_And_Get(t *testing.T) {
	c := cache.New(time.Minute)
	key := cache.Key("host:443", "pkg.Svc/Method", `{"id":1}`)
	c.Set(key, "hello", nil)

	e := c.Get(key)
	if e == nil {
		t.Fatal("expected cached entry")
	}
	if e.Output != "hello" {
		t.Errorf("output: got %q, want %q", e.Output, "hello")
	}
	if e.Err != nil {
		t.Errorf("unexpected error: %v", e.Err)
	}
}

func TestSet_WithError(t *testing.T) {
	c := cache.New(time.Minute)
	key := cache.Key("host:443", "pkg.Svc/Method", `{}`)
	sentinel := errors.New("rpc failed")
	c.Set(key, "", sentinel)

	e := c.Get(key)
	if e == nil {
		t.Fatal("expected cached entry")
	}
	if !errors.Is(e.Err, sentinel) {
		t.Errorf("err: got %v, want %v", e.Err, sentinel)
	}
}

func TestGet_ExpiredEntry(t *testing.T) {
	c := cache.New(10 * time.Millisecond)
	key := cache.Key("host:443", "pkg.Svc/Method", `{}`)
	c.Set(key, "data", nil)
	time.Sleep(20 * time.Millisecond)
	if got := c.Get(key); got != nil {
		t.Error("expected nil for expired entry")
	}
}

func TestPurge_RemovesExpired(t *testing.T) {
	c := cache.New(10 * time.Millisecond)
	for i := 0; i < 3; i++ {
		key := cache.Key("host:443", "pkg.Svc/Method", string(rune('a'+i)))
		c.Set(key, "v", nil)
	}
	if c.Len() != 3 {
		t.Fatalf("want 3 entries, got %d", c.Len())
	}
	time.Sleep(20 * time.Millisecond)
	c.Purge()
	if c.Len() != 0 {
		t.Errorf("want 0 after purge, got %d", c.Len())
	}
}

func TestKey_Deterministic(t *testing.T) {
	a := cache.Key("host:443", "Svc/M", `{"x":1}`)
	b := cache.Key("host:443", "Svc/M", `{"x":1}`)
	if a != b {
		t.Errorf("keys differ: %s vs %s", a, b)
	}
}

func TestKey_DifferentInputs(t *testing.T) {
	a := cache.Key("host:443", "Svc/M", `{"x":1}`)
	b := cache.Key("host:443", "Svc/M", `{"x":2}`)
	if a == b {
		t.Error("different bodies produced identical keys")
	}
}
