package alert

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Alerter writes change notifications to an output destination.
type Alerter struct {
	out    io.Writer
	prefix string
}

// New creates an Alerter that writes to the given writer.
// If w is nil, os.Stdout is used.
func New(w io.Writer) *Alerter {
	if w == nil {
		w = os.Stdout
	}
	return &Alerter{out: w}
}

// Notify writes a formatted alert line for each change in the slice.
func (a *Alerter) Notify(changes []monitor.Change) {
	for _, c := range changes {
		level := levelFor(c)
		timestamp := time.Now().Format(time.RFC3339)
		fmt.Fprintf(a.out, "%s [%s] %s\n", timestamp, level, c.String())
	}
}

// levelFor maps a change type to an alert level.
func levelFor(c monitor.Change) Level {
	switch c.Type {
	case monitor.ChangeOpened:
		return LevelAlert
	case monitor.ChangeClosed:
		return LevelWarn
	default:
		return LevelInfo
	}
}
