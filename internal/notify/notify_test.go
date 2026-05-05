package notify_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"grpcurl-batch/internal/notify"
)

var exampleSummary = notify.Summary{
	Total:   5,
	Passed:  4,
	Failed:  1,
	Elapsed: 120 * time.Millisecond,
}

func TestLog_Notify(t *testing.T) {
	var buf bytes.Buffer
	l := &notify.Log{Out: &buf}
	if err := l.Notify(context.Background(), exampleSummary); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"total=5", "passed=4", "failed=1"} {
		if !strings.Contains(out, want) {
			t.Errorf("output %q missing %q", out, want)
		}
	}
}

func TestWebhook_Notify_Success(t *testing.T) {
	var received string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var b strings.Builder
		_, _ = b.ReadFrom(r.Body)
		received = b.String()
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	wh := notify.NewWebhook(srv.URL)
	if err := wh.Notify(context.Background(), exampleSummary); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(received, `"failed":1`) {
		t.Errorf("payload missing failed field: %s", received)
	}
}

func TestWebhook_Notify_HTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	wh := notify.NewWebhook(srv.URL)
	if err := wh.Notify(context.Background(), exampleSummary); err == nil {
		t.Fatal("expected error for 500 response")
	}
}

func TestWebhook_Notify_BadURL(t *testing.T) {
	wh := notify.NewWebhook("://bad-url")
	if err := wh.Notify(context.Background(), exampleSummary); err == nil {
		t.Fatal("expected error for malformed URL")
	}
}

func TestWebhook_Notify_ContextCancelled(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	wh := notify.NewWebhook(srv.URL)
	if err := wh.Notify(ctx, exampleSummary); err == nil {
		t.Fatal("expected error for cancelled context")
	}
}
