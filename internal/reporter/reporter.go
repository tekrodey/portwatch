// Package reporter formats and writes periodic port-scan summaries.
package reporter

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// Format controls how the summary is rendered.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Reporter writes a formatted summary of the current port state.
type Reporter struct {
	w      io.Writer
	format Format
}

// New returns a Reporter that writes to w using the given format.
// If w is nil, os.Stdout is used.
func New(w io.Writer, format Format) *Reporter {
	if w == nil {
		w = os.Stdout
	}
	if format == "" {
		format = FormatText
	}
	return &Reporter{w: w, format: format}
}

// Write renders the current set of open ports to the underlying writer.
func (r *Reporter) Write(ports []monitor.Port, ts time.Time) error {
	switch r.format {
	case FormatJSON:
		return r.writeJSON(ports, ts)
	default:
		return r.writeText(ports, ts)
	}
}

func (r *Reporter) writeText(ports []monitor.Port, ts time.Time) error {
	_, err := fmt.Fprintf(r.w, "[%s] open ports (%d):\n",
		ts.Format(time.RFC3339), len(ports))
	if err != nil {
		return err
	}
	for _, p := range ports {
		if _, err := fmt.Fprintf(r.w, "  %s\n", p); err != nil {
			return err
		}
	}
	return nil
}

func (r *Reporter) writeJSON(ports []monitor.Port, ts time.Time) error {
	entries := make([]string, len(ports))
	for i, p := range ports {
		entries[i] = fmt.Sprintf("{\"proto\":%q,\"port\":%d}", p.Proto, p.Number)
	}
	_, err := fmt.Fprintf(r.w,
		"{\"timestamp\":%q,\"count\":%d,\"ports\":[%s]}\n",
		ts.Format(time.RFC3339), len(ports), strings.Join(entries, ","))
	return err
}
