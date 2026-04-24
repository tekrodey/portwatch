package notify

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// Formatter controls how change events are rendered.
type Formatter interface {
	Format(changes []monitor.Change) string
}

// Notifier sends formatted change events to a destination.
type Notifier struct {
	out       io.Writer
	formatter Formatter
}

// New creates a Notifier writing to out using the given Formatter.
// If out is nil, os.Stdout is used.
func New(out io.Writer, f Formatter) *Notifier {
	if out == nil {
		out = os.Stdout
	}
	if f == nil {
		f = &TextFormatter{}
	}
	return &Notifier{out: out, formatter: f}
}

// Send writes formatted changes to the output. It is a no-op when changes is empty.
func (n *Notifier) Send(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}
	msg := n.formatter.Format(changes)
	_, err := fmt.Fprintln(n.out, msg)
	return err
}

// TextFormatter renders changes as human-readable text lines.
type TextFormatter struct{}

func (t *TextFormatter) Format(changes []monitor.Change) string {
	ts := time.Now().UTC().Format(time.RFC3339)
	lines := make([]string, 0, len(changes))
	for _, c := range changes {
		lines = append(lines, fmt.Sprintf("[%s] %s", ts, c.String()))
	}
	return strings.Join(lines, "\n")
}
