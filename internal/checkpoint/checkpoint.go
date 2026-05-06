// Package checkpoint provides persistent progress tracking for batch runs.
// It records which requests have completed so that interrupted runs can be
// resumed without re-executing already-successful calls.
package checkpoint

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Entry records the outcome of a single request.
type Entry struct {
	Key       string    `json:"key"`
	Success   bool      `json:"success"`
	Timestamp time.Time `json:"timestamp"`
}

// Store holds completed entries and can persist them to disk.
type Store struct {
	mu      sync.RWMutex
	path    string
	entries map[string]Entry
}

// New loads an existing checkpoint file from path, or returns an empty Store
// if the file does not exist yet.
func New(path string) (*Store, error) {
	s := &Store{
		path:    path,
		entries: make(map[string]Entry),
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return s, nil
	}
	if err != nil {
		return nil, err
	}
	var list []Entry
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}
	for _, e := range list {
		s.entries[e.Key] = e
	}
	return s, nil
}

// Done reports whether the given key has already been recorded as successful.
func (s *Store) Done(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.entries[key]
	return ok && e.Success
}

// Record marks key with the given success flag and flushes to disk.
func (s *Store) Record(key string, success bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[key] = Entry{Key: key, Success: success, Timestamp: time.Now().UTC()}
	return s.flush()
}

// Len returns the number of recorded entries.
func (s *Store) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.entries)
}

// flush serialises entries to disk. Caller must hold the write lock.
func (s *Store) flush() error {
	list := make([]Entry, 0, len(s.entries))
	for _, e := range s.entries {
		list = append(list, e)
	}
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}
