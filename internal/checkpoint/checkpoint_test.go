package checkpoint_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/radovskyb/grpcurl-batch/internal/checkpoint"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "checkpoint.json")
}

func TestNew_MissingFile_ReturnsEmpty(t *testing.T) {
	s, err := checkpoint.New(tempPath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Len() != 0 {
		t.Fatalf("expected 0 entries, got %d", s.Len())
	}
}

func TestRecord_And_Done(t *testing.T) {
	p := tempPath(t)
	s, _ := checkpoint.New(p)

	if err := s.Record("req-1", true); err != nil {
		t.Fatalf("Record: %v", err)
	}
	if !s.Done("req-1") {
		t.Fatal("expected req-1 to be done")
	}
	if s.Done("req-2") {
		t.Fatal("req-2 should not be done")
	}
}

func TestRecord_FailureNotDone(t *testing.T) {
	p := tempPath(t)
	s, _ := checkpoint.New(p)
	_ = s.Record("req-fail", false)
	if s.Done("req-fail") {
		t.Fatal("failed entry should not count as done")
	}
}

func TestPersistence_ReloadsEntries(t *testing.T) {
	p := tempPath(t)
	s1, _ := checkpoint.New(p)
	_ = s1.Record("req-a", true)
	_ = s1.Record("req-b", false)

	s2, err := checkpoint.New(p)
	if err != nil {
		t.Fatalf("reload error: %v", err)
	}
	if !s2.Done("req-a") {
		t.Error("req-a should be done after reload")
	}
	if s2.Done("req-b") {
		t.Error("req-b should not be done after reload")
	}
	if s2.Len() != 2 {
		t.Errorf("expected 2 entries, got %d", s2.Len())
	}
}

func TestNew_CorruptFile_ReturnsError(t *testing.T) {
	p := tempPath(t)
	_ = os.WriteFile(p, []byte("not-json{"), 0o644)
	_, err := checkpoint.New(p)
	if err == nil {
		t.Fatal("expected error for corrupt file")
	}
}

func TestLen_Increments(t *testing.T) {
	s, _ := checkpoint.New(tempPath(t))
	for i := 0; i < 5; i++ {
		_ = s.Record(string(rune('a'+i)), true)
	}
	if s.Len() != 5 {
		t.Errorf("expected 5, got %d", s.Len())
	}
}
