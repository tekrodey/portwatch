// Package rollup groups rapid successive changes for the same port into a
// single representative event, preventing alert storms during port flapping.
package rollup

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// DefaultWindow is the time span within which changes to the same port are
// collapsed into one.
const DefaultWindow = 3 * time.Second

// Rollup collapses multiple changes for the same (port, proto, direction)
// tuple that arrive within a configurable window into the most-recent one.
type Rollup struct {
	mu      sync.Mutex
	window  time.Duration
	bucket  map[string]entry
	now     func() time.Time
}

type entry struct {
	change    monitor.Change
	receivedAt time.Time
}

// New returns a Rollup with the given deduplication window.
func New(window time.Duration) *Rollup {
	return newWithClock(window, time.Now)
}

func newWithClock(window time.Duration, now func() time.Time) *Rollup {
	return &Rollup{
		window: window,
		bucket: make(map[string]entry),
		now:    now,
	}
}

// Add stores or updates the change. It returns true when the change should be
// forwarded immediately (first occurrence in the window) and false when it
// was absorbed into an existing bucket entry.
func (r *Rollup) Add(c monitor.Change) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	k := key(c)
	e, exists := r.bucket[k]
	now := r.now()

	if exists && now.Sub(e.receivedAt) < r.window {
		// absorb — update the stored change but suppress forwarding
		r.bucket[k] = entry{change: c, receivedAt: e.receivedAt}
		return false
	}

	r.bucket[k] = entry{change: c, receivedAt: now}
	return true
}

// Flush returns all changes that have been held past the window and removes
// them from the bucket. It is intended to be called on a regular ticker so
// that absorbed updates are eventually emitted.
func (r *Rollup) Flush() []monitor.Change {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := r.now()
	var out []monitor.Change
	for k, e := range r.bucket {
		if now.Sub(e.receivedAt) >= r.window {
			out = append(out, e.change)
			delete(r.bucket, k)
		}
	}
	return out
}

func key(c monitor.Change) string {
	return c.Port.Proto + ":" + c.Port.String() + ":" + c.Direction
}
