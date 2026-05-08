package batch_test

import (
	"testing"

	"github.com/grpcurl-batch/internal/batch"
)

func makeItems(n int) []batch.Item {
	items := make([]batch.Item, n)
	for i := range items {
		items[i] = batch.Item{Name: "req", Payload: "{}"}
	}
	return items
}

func TestSplit_EmptyInput(t *testing.T) {
	batches, err := batch.Split(nil, batch.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(batches) != 0 {
		t.Fatalf("expected 0 batches, got %d", len(batches))
	}
}

func TestSplit_InvalidSize(t *testing.T) {
	_, err := batch.Split(makeItems(5), batch.Options{Size: 0})
	if err == nil {
		t.Fatal("expected error for size=0, got nil")
	}
}

func TestSplit_ExactMultiple(t *testing.T) {
	items := makeItems(9)
	batches, err := batch.Split(items, batch.Options{Size: 3})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(batches) != 3 {
		t.Fatalf("expected 3 batches, got %d", len(batches))
	}
	for i, b := range batches {
		if b.Index != i {
			t.Errorf("batch %d: wrong index %d", i, b.Index)
		}
		if len(b.Items) != 3 {
			t.Errorf("batch %d: expected 3 items, got %d", i, len(b.Items))
		}
	}
}

func TestSplit_Remainder(t *testing.T) {
	items := makeItems(7)
	batches, err := batch.Split(items, batch.Options{Size: 3})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(batches) != 3 {
		t.Fatalf("expected 3 batches, got %d", len(batches))
	}
	last := batches[2]
	if len(last.Items) != 1 {
		t.Errorf("last batch: expected 1 item, got %d", len(last.Items))
	}
}

func TestFlatten_RoundTrip(t *testing.T) {
	items := makeItems(11)
	batches, err := batch.Split(items, batch.Options{Size: 4})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := batch.Flatten(batches)
	if len(got) != len(items) {
		t.Fatalf("flatten: expected %d items, got %d", len(items), len(got))
	}
}

func TestFlatten_Empty(t *testing.T) {
	got := batch.Flatten(nil)
	if len(got) != 0 {
		t.Fatalf("expected empty slice, got %d items", len(got))
	}
}
