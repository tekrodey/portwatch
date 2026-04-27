// Package dedup provides a filter that suppresses duplicate change events
// within a configurable time window, ensuring each unique port/protocol/direction
// combination is only reported once until the window expires.
package dedup

import (
	"fmt"
	"sync"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// entry tracks when a change was last seen.
type entry struct {
	seenAt time.Time
}

// Dedup filters out duplicate change events within a sliding time window.
type Dedup struct {
	mu     sync.Mutex
	window time.Duration
	seen   map[string]entry
	now    func() time.Time
}

// New returns a Dedup that suppresses repeated changes within window.
func New(window time.Duration) *Dedup {
	return &Dedup{
		window: window,
		seen:   make(map[string]entry),
		now:    time.Now,
	}
}

// newWithClock returns a Dedup with an injectable clock for testing.
func newWithClock(window time.Duration, clock func() time.Time) *Dedup {
	return &Dedup{
		window: window,
		seen:   make(map[string]entry),
		now:    clock,
	}
}

// Filter returns only changes that have not been seen within the window.
// Changes that pass through are recorded; stale entries are evicted.
func (d *Dedup) Filter(changes []monitor.Change) []monitor.Change {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.now()
	d.evict(now)

	var out []monitor.Change
	for _, c := range changes {
		k := key(c)
		if _, exists := d.seen[k]; exists {
			continue
		}
		d.seen[k] = entry{seenAt: now}
		out = append(out, c)
	}
	return out
}

// evict removes entries whose window has expired. Must be called with mu held.
func (d *Dedup) evict(now time.Time) {
	for k, e := range d.seen {
		if now.Sub(e.seenAt) >= d.window {
			delete(d.seen, k)
		}
	}
}

// key builds a unique string for a change event.
func key(c monitor.Change) string {
	return fmt.Sprintf("%s:%d:%s", c.Protocol, c.Port, c.Direction)
}
