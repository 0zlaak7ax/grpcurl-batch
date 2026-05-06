package backoff_test

import (
	"testing"
	"time"

	"github.com/example/grpcurl-batch/internal/backoff"
)

func TestExponential_Delay_GrowsWithAttempt(t *testing.T) {
	e := backoff.New(100*time.Millisecond, 10*time.Second, false)

	d0 := e.Delay(0)
	d1 := e.Delay(1)
	d2 := e.Delay(2)

	if d0 != 100*time.Millisecond {
		t.Errorf("attempt 0: got %v, want 100ms", d0)
	}
	if d1 != 200*time.Millisecond {
		t.Errorf("attempt 1: got %v, want 200ms", d1)
	}
	if d2 != 400*time.Millisecond {
		t.Errorf("attempt 2: got %v, want 400ms", d2)
	}
}

func TestExponential_Delay_CappedAtMax(t *testing.T) {
	e := backoff.New(1*time.Second, 3*time.Second, false)

	d := e.Delay(10)
	if d != 3*time.Second {
		t.Errorf("expected cap at 3s, got %v", d)
	}
}

func TestExponential_Delay_NegativeAttempt(t *testing.T) {
	e := backoff.New(50*time.Millisecond, 1*time.Second, false)
	d := e.Delay(-1)
	if d != 50*time.Millisecond {
		t.Errorf("negative attempt should behave like 0, got %v", d)
	}
}

func TestExponential_Delay_WithJitter(t *testing.T) {
	e := backoff.New(200*time.Millisecond, 5*time.Second, true)
	base := 200 * time.Millisecond
	max := base + base/4

	for i := 0; i < 20; i++ {
		d := e.Delay(0)
		if d < base || d > max {
			t.Errorf("jitter out of range [%v, %v]: got %v", base, max, d)
		}
	}
}

func TestFixed_Delay_AlwaysSame(t *testing.T) {
	f := &backoff.Fixed{Interval: 500 * time.Millisecond}
	for i := 0; i < 5; i++ {
		if d := f.Delay(i); d != 500*time.Millisecond {
			t.Errorf("attempt %d: got %v, want 500ms", i, d)
		}
	}
}

func TestLinear_Delay_Grows(t *testing.T) {
	l := &backoff.Linear{Base: 100 * time.Millisecond, Max: 500 * time.Millisecond}

	if d := l.Delay(0); d != 100*time.Millisecond {
		t.Errorf("attempt 0: got %v", d)
	}
	if d := l.Delay(2); d != 300*time.Millisecond {
		t.Errorf("attempt 2: got %v", d)
	}
}

func TestLinear_Delay_CappedAtMax(t *testing.T) {
	l := &backoff.Linear{Base: 200 * time.Millisecond, Max: 500 * time.Millisecond}
	if d := l.Delay(5); d != 500*time.Millisecond {
		t.Errorf("expected cap, got %v", d)
	}
}
