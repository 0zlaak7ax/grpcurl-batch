// Package notify provides pluggable notification backends for
// reporting batch run outcomes (e.g. Slack, webhook, log).
package notify

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Summary holds the high-level outcome of a batch run passed to notifiers.
type Summary struct {
	Total   int
	Passed  int
	Failed  int
	Elapsed time.Duration
}

// Notifier is the interface implemented by every notification backend.
type Notifier interface {
	Notify(ctx context.Context, s Summary) error
}

// Webhook sends a JSON payload to an arbitrary HTTP endpoint.
type Webhook struct {
	URL    string
	Client *http.Client
}

// NewWebhook constructs a Webhook notifier with a sensible default client.
func NewWebhook(url string) *Webhook {
	return &Webhook{
		URL: url,
		Client: &http.Client{Timeout: 10 * time.Second},
	}
}

// Notify sends the summary as a JSON body via HTTP POST.
func (w *Webhook) Notify(ctx context.Context, s Summary) error {
	body := fmt.Sprintf(
		`{"total":%d,"passed":%d,"failed":%d,"elapsed_ms":%d}`,
		s.Total, s.Passed, s.Failed, s.Elapsed.Milliseconds(),
	)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, w.URL,
		strings.NewReader(body))
	if err != nil {
		return fmt.Errorf("notify webhook: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := w.Client.Do(req)
	if err != nil {
		return fmt.Errorf("notify webhook: do request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("notify webhook: unexpected status %d", resp.StatusCode)
	}
	return nil
}

// Log writes a human-readable line to any io.Writer (e.g. os.Stderr).
type Log struct {
	Out io.Writer
}

// Notify writes a one-line summary to the configured writer.
func (l *Log) Notify(_ context.Context, s Summary) error {
	_, err := fmt.Fprintf(l.Out,
		"[notify] total=%d passed=%d failed=%d elapsed=%s\n",
		s.Total, s.Passed, s.Failed, s.Elapsed.Round(time.Millisecond))
	return err
}
