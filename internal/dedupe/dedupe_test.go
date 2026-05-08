package dedupe_test

import (
	"sync"
	"testing"

	"github.com/your-org/grpcurl-batch/internal/dedupe"
)

func TestNew_EmptyFilter(t *testing.T) {
	f := dedupe.New()
	if f.Len() != 0 {
		t.Fatalf("expected Len 0, got %d", f.Len())
	}
}

func TestIsDuplicate_FirstCallReturnsFalse(t *testing.T) {
	f := dedupe.New()
	if f.IsDuplicate("req-1") {
		t.Fatal("first call should not be a duplicate")
	}
}

func TestIsDuplicate_SecondCallReturnsTrue(t *testing.T) {
	f := dedupe.New()
	f.IsDuplicate("req-1")
	if !f.IsDuplicate("req-1") {
		t.Fatal("second call with same key should be a duplicate")
	}
}

func TestIsDuplicate_DifferentKeysIndependent(t *testing.T) {
	f := dedupe.New()
	if f.IsDuplicate("a") {
		t.Fatal("key 'a' should not be duplicate on first call")
	}
	if f.IsDuplicate("b") {
		t.Fatal("key 'b' should not be duplicate on first call")
	}
	if f.Len() != 2 {
		t.Fatalf("expected Len 2, got %d", f.Len())
	}
}

func TestReset_ClearsSeenKeys(t *testing.T) {
	f := dedupe.New()
	f.IsDuplicate("req-1")
	f.IsDuplicate("req-2")
	f.Reset()
	if f.Len() != 0 {
		t.Fatalf("expected Len 0 after Reset, got %d", f.Len())
	}
	if f.IsDuplicate("req-1") {
		t.Fatal("key should not be duplicate after Reset")
	}
}

func TestIsDuplicate_Concurrent(t *testing.T) {
	f := dedupe.New()
	const goroutines = 50
	var wg sync.WaitGroup
	duplicates := make([]bool, goroutines)

	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func(idx int) {
			defer wg.Done()
			duplicates[idx] = f.IsDuplicate("shared-key")
		}(i)
	}
	wg.Wait()

	firstCount := 0
	for _, d := range duplicates {
		if !d {
			firstCount++
		}
	}
	if firstCount != 1 {
		t.Fatalf("expected exactly 1 non-duplicate across %d goroutines, got %d", goroutines, firstCount)
	}
}
