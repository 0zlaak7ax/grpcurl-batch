package formatter

import (
	"fmt"
	"io"
)

// Report holds aggregated results from a batch run.
type Report struct {
	Total   int
	Passed  int
	Failed  int
	Results []Result
}

// Collector accumulates Results and produces a final Report.
type Collector struct {
	results []Result
}

// Add appends a Result to the collector.
func (c *Collector) Add(r Result) {
	c.results = append(c.results, r)
}

// Report returns the aggregated Report.
func (c *Collector) Report() Report {
	rep := Report{Results: c.results, Total: len(c.results)}
	for _, r := range c.results {
		if r.Success {
			rep.Passed++
		} else {
			rep.Failed++
		}
	}
	return rep
}

// PrintSummary writes a human-readable summary of the Report to w.
func PrintSummary(w io.Writer, rep Report) error {
	_, err := fmt.Fprintf(w,
		"\n--- Batch Summary ---\nTotal: %d | Passed: %d | Failed: %d\n",
		rep.Total, rep.Passed, rep.Failed)
	return err
}
