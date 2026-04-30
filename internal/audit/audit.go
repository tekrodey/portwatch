// Package audit provides a structured audit log of all port change events
// with timestamps, severity, and optional annotations.
package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// Entry represents a single audited event.
type Entry struct {
	Timestamp time.Time      `json:"timestamp"`
	Direction string         `json:"direction"` // "opened" | "closed"
	Protocol  string         `json:"protocol"`
	Port      int            `json:"port"`
	Service   string         `json:"service,omitempty"`
	Note      string         `json:"note,omitempty"`
}

// Logger writes structured audit entries to a destination.
type Logger struct {
	out io.Writer
}

// New returns a Logger that writes to w.
// Pass nil to default to os.Stdout.
func New(w io.Writer) *Logger {
	if w == nil {
		w = os.Stdout
	}
	return &Logger{out: w}
}

// Log converts a slice of monitor.Change values into audit entries and
// writes each one as a newline-delimited JSON record.
func (l *Logger) Log(changes []monitor.Change) error {
	for _, c := range changes {
		e := Entry{
			Timestamp: time.Now().UTC(),
			Direction: c.Direction,
			Protocol:  c.Port.Protocol,
			Port:      c.Port.Number,
		}
		line, err := json.Marshal(e)
		if err != nil {
			return fmt.Errorf("audit: marshal: %w", err)
		}
		if _, err := fmt.Fprintf(l.out, "%s\n", line); err != nil {
			return fmt.Errorf("audit: write: %w", err)
		}
	}
	return nil
}
