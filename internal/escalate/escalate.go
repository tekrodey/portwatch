// Package escalate promotes a change to a higher severity level when the same
// port event recurs more than a configured number of times within a window.
package escalate

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// Level represents an alert severity.
type Level int

const (
	LevelNormal   Level = iota
	LevelElevated       // recurred at least once within the window
	LevelCritical       // recurred at least twice within the window
)

type entry struct {
	count     int
	windowEnd time.Time
}

// Escalator tracks recurrence counts and returns a severity Level for each
// change. It is safe for concurrent use.
type Escalator struct {
	mu      sync.Mutex
	window  time.Duration
	entries map[string]*entry
	now     func() time.Time
}

// New returns an Escalator with the given recurrence window.
func New(window time.Duration) *Escalator {
	return newWithClock(window, time.Now)
}

func newWithClock(window time.Duration, now func() time.Time) *Escalator {
	return &Escalator{
		window:  window,
		entries: make(map[string]*entry),
		now:     now,
	}
}

// Assess returns the Level for the supplied change and updates internal state.
func (e *Escalator) Assess(c monitor.Change) Level {
	k := key(c)
	now := e.now()

	e.mu.Lock()
	defer e.mu.Unlock()

	ent, ok := e.entries[k]
	if !ok || now.After(ent.windowEnd) {
		e.entries[k] = &entry{count: 1, windowEnd: now.Add(e.window)}
		return LevelNormal
	}

	ent.count++
	switch {
	case ent.count >= 3:
		return LevelCritical
	case ent.count == 2:
		return LevelElevated
	default:
		return LevelNormal
	}
}

// Reset clears the recurrence record for a change, e.g. after an operator
// acknowledges the alert.
func (e *Escalator) Reset(c monitor.Change) {
	e.mu.Lock()
	delete(e.entries, key(c))
	e.mu.Unlock()
}

func key(c monitor.Change) string {
	return c.Port.String() + "|" + c.Direction
}
