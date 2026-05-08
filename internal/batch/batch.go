// Package batch provides utilities for splitting a slice of requests
// into fixed-size chunks for controlled parallel processing.
package batch

import "fmt"

// Item represents a single unit of work identified by a name and payload.
type Item struct {
	Name    string
	Payload string
}

// Batch holds a slice of items forming one processing chunk.
type Batch struct {
	Index int
	Items []Item
}

// Options controls how items are grouped into batches.
type Options struct {
	// Size is the maximum number of items per batch. Must be >= 1.
	Size int
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{Size: 10}
}

// Split divides items into batches of at most opts.Size items each.
// It returns an error if opts.Size is less than 1.
func Split(items []Item, opts Options) ([]Batch, error) {
	if opts.Size < 1 {
		return nil, fmt.Errorf("batch: size must be >= 1, got %d", opts.Size)
	}
	if len(items) == 0 {
		return nil, nil
	}

	var batches []Batch
	for i := 0; i < len(items); i += opts.Size {
		end := i + opts.Size
		if end > len(items) {
			end = len(items)
		}
		batches = append(batches, Batch{
			Index: len(batches),
			Items: items[i:end],
		})
	}
	return batches, nil
}

// Flatten merges all batches back into a single ordered slice of items.
func Flatten(batches []Batch) []Item {
	var out []Item
	for _, b := range batches {
		out = append(out, b.Items...)
	}
	return out
}
