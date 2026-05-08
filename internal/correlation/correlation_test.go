package correlation_test

import (
	"context"
	"strings"
	"testing"

	"github.com/sethpollack/grpcurl-batch/internal/correlation"
)

func TestGenerate_ProducesNonEmptyID(t *testing.T) {
	g := correlation.New("")
	id, err := g.Generate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id == "" {
		t.Fatal("expected non-empty ID")
	}
}

func TestGenerate_WithPrefix(t *testing.T) {
	g := correlation.New("req")
	id, err := g.Generate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(id.String(), "req-") {
		t.Fatalf("expected prefix 'req-', got %q", id)
	}
}

func TestGenerate_UniqueIDs(t *testing.T) {
	g := correlation.New("")
	ids := make(map[correlation.ID]struct{}, 100)
	for i := 0; i < 100; i++ {
		id, err := g.Generate()
		if err != nil {
			t.Fatalf("unexpected error on iteration %d: %v", i, err)
		}
		if _, seen := ids[id]; seen {
			t.Fatalf("duplicate ID generated: %q", id)
		}
		ids[id] = struct{}{}
	}
}

func TestWithID_And_FromContext(t *testing.T) {
	ctx := context.Background()
	expected := correlation.ID("test-abc123")
	ctx = correlation.WithID(ctx, expected)

	got, ok := correlation.FromContext(ctx)
	if !ok {
		t.Fatal("expected ID in context, got none")
	}
	if got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}

func TestFromContext_Missing(t *testing.T) {
	_, ok := correlation.FromContext(context.Background())
	if ok {
		t.Fatal("expected no ID in empty context")
	}
}

func TestMustFromContext_Fallback(t *testing.T) {
	id := correlation.MustFromContext(context.Background())
	if id != "unknown" {
		t.Fatalf("expected 'unknown', got %q", id)
	}
}

func TestMustFromContext_Present(t *testing.T) {
	ctx := correlation.WithID(context.Background(), "my-id")
	if got := correlation.MustFromContext(ctx); got != "my-id" {
		t.Fatalf("expected 'my-id', got %q", got)
	}
}
