package backoff_test

import (
	"testing"
	"time"

	"github.com/example/grpcurl-batch/internal/backoff"
)

// TestExponential_SequenceIsMonotonic verifies that delays grow monotonically
// (without jitter) across a realistic number of attempts.
func TestExponential_SequenceIsMonotonic(t *testing.T) {
	e := backoff.New(10*time.Millisecond, 30*time.Second, false)

	prev := time.Duration(0)
	for i := 0; i < 15; i++ {
		d := e.Delay(i)
		if d < prev {
			t.Errorf("attempt %d: delay %v is less than previous %v", i, d, prev)
		}
		if d > 30*time.Second {
			t.Errorf("attempt %d: delay %v exceeds max", i, d)
		}
		prev = d
	}
}

// TestJitter_NeverExceedsMaxPlusQuarter ensures jitter stays within bounds.
func TestJitter_NeverExceedsMaxPlusQuarter(t *testing.T) {
	max := 2 * time.Second
	e := backoff.New(100*time.Millisecond, max, true)

	// At a high attempt the base delay equals max; jitter adds at most max/4.
	upper := max + max/4
	for i := 0; i < 100; i++ {
		d := e.Delay(20)
		if d > upper {
			t.Errorf("jitter exceeded upper bound %v: got %v", upper, d)
		}
	}
}

// TestLinear_NoMaxUnbounded confirms Linear grows unbounded when Max is zero.
func TestLinear_NoMaxUnbounded(t *testing.T) {
	l := &backoff.Linear{Base: 100 * time.Millisecond}
	expected := 1000 * time.Millisecond
	if d := l.Delay(9); d != expected {
		t.Errorf("expected %v, got %v", expected, d)
	}
}
