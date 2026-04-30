package formatter

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

// Format represents the output format type.
type Format string

const (
	FormatJSON    Format = "json"
	FormatText    Format = "text"
	FormatSummary Format = "summary"
)

// Result holds the outcome of a single gRPC call.
type Result struct {
	Method   string        `json:"method"`
	Success  bool          `json:"success"`
	Attempts int           `json:"attempts"`
	Duration time.Duration `json:"duration_ms"`
	Output   string        `json:"output,omitempty"`
	Error    string        `json:"error,omitempty"`
}

// Formatter writes results to an io.Writer in a given format.
type Formatter struct {
	format Format
	out    io.Writer
}

// New creates a new Formatter.
func New(format Format, out io.Writer) *Formatter {
	return &Formatter{format: format, out: out}
}

// Write outputs a single Result according to the configured format.
func (f *Formatter) Write(r Result) error {
	switch f.format {
	case FormatJSON:
		return f.writeJSON(r)
	case FormatSummary:
		return f.writeSummary(r)
	default:
		return f.writeText(r)
	}
}

func (f *Formatter) writeJSON(r Result) error {
	r.Duration = r.Duration / time.Millisecond
	b, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("formatter: marshal error: %w", err)
	}
	_, err = fmt.Fprintln(f.out, string(b))
	return err
}

func (f *Formatter) writeText(r Result) error {
	status := "OK"
	if !r.Success {
		status = "FAIL"
	}
	_, err := fmt.Fprintf(f.out, "[%s] %s (attempts=%d, duration=%dms)\n",
		status, r.Method, r.Attempts, r.Duration/time.Millisecond)
	if err != nil {
		return err
	}
	if r.Output != "" {
		_, err = fmt.Fprintf(f.out, "  output: %s\n", strings.TrimSpace(r.Output))
	}
	if r.Error != "" {
		_, err = fmt.Fprintf(f.out, "  error:  %s\n", strings.TrimSpace(r.Error))
	}
	return err
}

func (f *Formatter) writeSummary(r Result) error {
	status := "✓"
	if !r.Success {
		status = "✗"
	}
	_, err := fmt.Fprintf(f.out, "%s %s\n", status, r.Method)
	return err
}
