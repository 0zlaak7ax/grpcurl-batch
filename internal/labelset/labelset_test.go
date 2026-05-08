package labelset_test

import (
	"strings"
	"testing"

	"github.com/isobit/grpcurl-batch/internal/labelset"
)

func TestNew_EmptyMap(t *testing.T) {
	ls := labelset.New(map[string]string{})
	if ls.Len() != 0 {
		t.Fatalf("expected 0 labels, got %d", ls.Len())
	}
}

func TestNew_DropsEmptyKeys(t *testing.T) {
	ls := labelset.New(map[string]string{" ": "value", "": "other"})
	if ls.Len() != 0 {
		t.Fatalf("expected empty set, got %d labels", ls.Len())
	}
}

func TestNew_TrimsWhitespace(t *testing.T) {
	ls := labelset.New(map[string]string{" env ": " prod "})
	v, ok := ls.Get("env")
	if !ok {
		t.Fatal("expected key 'env' to be present after trimming")
	}
	if v != "prod" {
		t.Fatalf("expected value 'prod', got %q", v)
	}
}

func TestGet_Missing(t *testing.T) {
	ls := labelset.New(map[string]string{"a": "1"})
	_, ok := ls.Get("missing")
	if ok {
		t.Fatal("expected missing key to return false")
	}
}

func TestKeys_SortedOrder(t *testing.T) {
	ls := labelset.New(map[string]string{"z": "1", "a": "2", "m": "3"})
	keys := ls.Keys()
	if len(keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(keys))
	}
	if keys[0] != "a" || keys[1] != "m" || keys[2] != "z" {
		t.Fatalf("unexpected key order: %v", keys)
	}
}

func TestString_Deterministic(t *testing.T) {
	ls := labelset.New(map[string]string{"env": "prod", "team": "infra"})
	s := ls.String()
	if !strings.Contains(s, "env=prod") || !strings.Contains(s, "team=infra") {
		t.Fatalf("unexpected string representation: %q", s)
	}
	// Should be identical across multiple calls.
	if ls.String() != s {
		t.Fatal("String() is not deterministic")
	}
}

func TestMerge_OtherWins(t *testing.T) {
	a := labelset.New(map[string]string{"env": "staging", "region": "us-east"})
	b := labelset.New(map[string]string{"env": "prod"})
	merged := a.Merge(b)

	v, _ := merged.Get("env")
	if v != "prod" {
		t.Fatalf("expected 'prod' from b, got %q", v)
	}
	v, _ = merged.Get("region")
	if v != "us-east" {
		t.Fatalf("expected 'us-east' from a, got %q", v)
	}
}

func TestMerge_DoesNotMutateOriginal(t *testing.T) {
	a := labelset.New(map[string]string{"k": "original"})
	b := labelset.New(map[string]string{"k": "overridden"})
	_ = a.Merge(b)
	v, _ := a.Get("k")
	if v != "original" {
		t.Fatalf("Merge mutated the receiver; got %q", v)
	}
}

func TestToMap_IsCopy(t *testing.T) {
	ls := labelset.New(map[string]string{"x": "1"})
	m := ls.ToMap()
	m["x"] = "mutated"
	v, _ := ls.Get("x")
	if v != "1" {
		t.Fatal("ToMap returned a reference; mutation affected the LabelSet")
	}
}
