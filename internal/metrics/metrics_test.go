package metrics_test

import (
	"sync"
	"testing"
	"time"

	"github.com/example/grpcurl-batch/internal/metrics"
)

func TestNew_EmptySnapshot(t *testing.T) {
	c := metrics.New()
	s := c.Snapshot()
	if s.Total != 0 || s.Success != 0 || s.Failure != 0 {
		t.Fatalf("expected zero snapshot, got %+v", s)
	}
	if s.AvgLatency() != 0 {
		t.Fatalf("expected zero avg latency")
	}
}

func TestRecord_Success(t *testing.T) {
	c := metrics.New()
	c.Record(true, 10*time.Millisecond)
	s := c.Snapshot()
	if s.Total != 1 || s.Success != 1 || s.Failure != 0 {
		t.Fatalf("unexpected snapshot: %+v", s)
	}
}

func TestRecord_Failure(t *testing.T) {
	c := metrics.New()
	c.Record(false, 5*time.Millisecond)
	s := c.Snapshot()
	if s.Total != 1 || s.Success != 0 || s.Failure != 1 {
		t.Fatalf("unexpected snapshot: %+v", s)
	}
}

func TestAvgLatency(t *testing.T) {
	c := metrics.New()
	c.Record(true, 20*time.Millisecond)
	c.Record(false, 40*time.Millisecond)
	s := c.Snapshot()
	want := 30 * time.Millisecond
	if s.AvgLatency() != want {
		t.Fatalf("avg latency: got %v want %v", s.AvgLatency(), want)
	}
}

func TestRecord_Concurrent(t *testing.T) {
	c := metrics.New()
	const n = 200
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			c.Record(i%2 == 0, time.Millisecond)
		}(i)
	}
	wg.Wait()
	s := c.Snapshot()
	if s.Total != n {
		t.Fatalf("total: got %d want %d", s.Total, n)
	}
	if s.Success+s.Failure != n {
		t.Fatalf("success+failure mismatch: %d+%d != %d", s.Success, s.Failure, n)
	}
}
