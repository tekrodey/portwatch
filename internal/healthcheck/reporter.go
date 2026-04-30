package healthcheck

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
)

// Reporter writes health status summaries to a writer.
type Reporter struct {
	w       io.Writer
	checker *Checker
}

// NewReporter returns a Reporter that writes to stdout by default.
func NewReporter(c *Checker) *Reporter {
	return &Reporter{w: os.Stdout, checker: c}
}

// NewReporterWithWriter returns a Reporter that writes to w.
func NewReporterWithWriter(c *Checker, w io.Writer) *Reporter {
	return &Reporter{w: w, checker: c}
}

// PrintText writes a human-readable health summary to the writer.
func (r *Reporter) PrintText() {
	statuses := r.checker.All()
	sort.Slice(statuses, func(i, j int) bool {
		return statuses[i].Name < statuses[j].Name
	})
	overall := "OK"
	if !r.checker.Healthy() {
		overall = "DEGRADED"
	}
	fmt.Fprintf(r.w, "overall: %s\n", overall)
	for _, s := range statuses {
		mark := "✓"
		if !s.Healthy {
			mark = "✗"
		}
		line := fmt.Sprintf("  %s %-20s %s", mark, s.Name, s.LastCheck.Format("15:04:05"))
		if s.Message != "" {
			line += fmt.Sprintf(" — %s", s.Message)
		}
		fmt.Fprintln(r.w, line)
	}
}

// PrintJSON writes a JSON-encoded health report to the writer.
func (r *Reporter) PrintJSON() error {
	type report struct {
		Healthy  bool     `json:"healthy"`
		Statuses []Status `json:"statuses"`
	}
	rep := report{
		Healthy:  r.checker.Healthy(),
		Statuses: r.checker.All(),
	}
	enc := json.NewEncoder(r.w)
	enc.SetIndent("", "  ")
	return enc.Encode(rep)
}
