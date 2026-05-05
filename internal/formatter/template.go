package formatter

import (
	"bytes"
	"fmt"
	"io"
	"text/template"
	"time"
)

// TemplateData holds the data passed to a custom output template.
type TemplateData struct {
	Method   string
	Address  string
	Success  bool
	Output   string
	Error    string
	Attempts int
	Duration time.Duration
}

// TemplateFormatter renders results using a user-supplied Go text/template string.
type TemplateFormatter struct {
	tmpl *template.Template
	out  io.Writer
}

// NewTemplateFormatter parses tmplStr and returns a TemplateFormatter that
// writes rendered output to out. Returns an error if the template is invalid.
func NewTemplateFormatter(out io.Writer, tmplStr string) (*TemplateFormatter, error) {
	funcMap := template.FuncMap{
		"durationMs": func(d time.Duration) int64 { return d.Milliseconds() },
		"checkmark": func(ok bool) string {
			if ok {
				return "✓"
			}
			return "✗"
		},
	}
	tmpl, err := template.New("result").Funcs(funcMap).Parse(tmplStr)
	if err != nil {
		return nil, fmt.Errorf("formatter: invalid template: %w", err)
	}
	return &TemplateFormatter{tmpl: tmpl, out: out}, nil
}

// Write renders a single Result using the template and writes it to the output writer.
func (tf *TemplateFormatter) Write(r Result) error {
	d := TemplateData{
		Method:   r.Method,
		Address:  r.Address,
		Success:  r.Success,
		Output:   r.Output,
		Error:    r.Error,
		Attempts: r.Attempts,
		Duration: r.Duration,
	}
	var buf bytes.Buffer
	if err := tf.tmpl.Execute(&buf, d); err != nil {
		return fmt.Errorf("formatter: template execution failed: %w", err)
	}
	_, err := fmt.Fprintln(tf.out, buf.String())
	return err
}
