package metrics

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// Print writes a human-readable summary of s to w.
func Print(w io.Writer, s Snapshot) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	rows := []struct {
		label string
		value string
	}{
		{"Total requests:", fmt.Sprintf("%d", s.Total)},
		{"Successful:", fmt.Sprintf("%d", s.Success)},
		{"Failed:", fmt.Sprintf("%d", s.Failure)},
		{"Total duration:", s.TotalDur.Round(time.Millisecond).String()},
		{"Avg latency:", s.AvgLatency().Round(time.Millisecond).String()},
	}
	for _, r := range rows {
		if _, err := fmt.Fprintf(tw, "%s\t%s\n", r.label, r.value); err != nil {
			return err
		}
	}
	return tw.Flush()
}
