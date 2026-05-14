package window

import (
	"sync"
	"testing"
	"time"
)

func newAt(base time.Time) *Window {
	w := New(time.Second, 10)
	w.now = func() time.Time { return base }
	return w
}

func TestNew_PanicsOnInvalidSize(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for size=0")
		}
	}()
	New(0, 4)
}

func TestNew_PanicsOnInvalidBuckets(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for buckets=0")
		}
	}()
	New(time.Second, 0)
}

func TestCount_EmptyWindow(t *testing.T) {
	w := newAt(time.Now())
	if got := w.Count(); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestAdd_And_Count(t *testing.T) {
	base := time.Now()
	w := newAt(base)
	w.Add(3)
	w.Add(5)
	if got := w.Count(); got != 8 {
		t.Fatalf("expected 8, got %d", got)
	}
}

func TestCount_ExpiresOldEvents(t *testing.T) {
	base := time.Now()
	w := New(time.Second, 10)
	w.now = func() time.Time { return base }
	w.Add(7)

	// Advance time beyond the window.
	w.now = func() time.Time { return base.Add(2 * time.Second) }
	if got := w.Count(); got != 0 {
		t.Fatalf("expected 0 after expiry, got %d", got)
	}
}

func TestReset_ClearsCount(t *testing.T) {
	w := newAt(time.Now())
	w.Add(10)
	w.Reset()
	if got := w.Count(); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestAdd_Concurrent(t *testing.T) {
	w := New(time.Second, 10)
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			w.Add(1)
		}()
	}
	wg.Wait()
	if got := w.Count(); got < 1 {
		t.Fatalf("expected positive count, got %d", got)
	}
}
