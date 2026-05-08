// Package labelset provides a key-value label store for tagging gRPC
// requests with arbitrary metadata (e.g. environment, team, service).
// Labels are immutable once built and safe for concurrent reads.
package labelset

import (
	"fmt"
	"sort"
	"strings"
)

// LabelSet holds an immutable set of string key-value labels.
type LabelSet struct {
	labels map[string]string
}

// New creates a LabelSet from the provided map. Keys and values are
// trimmed of whitespace; empty keys are silently dropped.
func New(m map[string]string) *LabelSet {
	copy := make(map[string]string, len(m))
	for k, v := range m {
		k = strings.TrimSpace(k)
		if k == "" {
			continue
		}
		copy[k] = strings.TrimSpace(v)
	}
	return &LabelSet{labels: copy}
}

// Get returns the value for key and whether it was present.
func (ls *LabelSet) Get(key string) (string, bool) {
	v, ok := ls.labels[key]
	return v, ok
}

// Keys returns all label keys in sorted order.
func (ls *LabelSet) Keys() []string {
	keys := make([]string, 0, len(ls.labels))
	for k := range ls.labels {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Len returns the number of labels.
func (ls *LabelSet) Len() int { return len(ls.labels) }

// String returns a deterministic, comma-separated "key=value" representation.
func (ls *LabelSet) String() string {
	keys := ls.Keys()
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", k, ls.labels[k]))
	}
	return strings.Join(parts, ",")
}

// Merge returns a new LabelSet that is the union of ls and other.
// Labels from other overwrite labels from ls on key collision.
func (ls *LabelSet) Merge(other *LabelSet) *LabelSet {
	m := make(map[string]string, ls.Len()+other.Len())
	for k, v := range ls.labels {
		m[k] = v
	}
	for k, v := range other.labels {
		m[k] = v
	}
	return &LabelSet{labels: m}
}

// ToMap returns a shallow copy of the underlying label map.
func (ls *LabelSet) ToMap() map[string]string {
	copy := make(map[string]string, len(ls.labels))
	for k, v := range ls.labels {
		copy[k] = v
	}
	return copy
}
