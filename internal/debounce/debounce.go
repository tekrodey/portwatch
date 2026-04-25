// Package debounce provides a mechanism to suppress repeated port change
// events within a configurable settling window. This prevents alert storms
// when a service restarts and its port briefly disappears then reappears.
package debounce

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// entry tracks when a change key was first seen.
type entry struct {
	firstSeen time.Time
}

// Debouncer holds pending changes and only releases them once they have
// persisted for at least the configured window duration.
type Debouncer struct {
	mu      sync.Mutex
	window  time.Duration
	pending map[string]entry
	now     func() time.Time
}

// New returns a Debouncer with the given settling window.
func New(window time.Duration) *Debouncer {
	return &Debouncer{
		window:  window,
		pending: make(map[string]entry),
		now:     time.Now,
	}
}

// Filter accepts a slice of changes and returns only those that have been
// continuously pending for at least the settling window. Changes that do not
// yet meet the threshold are retained internally for future calls.
func (d *Debouncer) Filter(changes []monitor.Change) []monitor.Change {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.now()

	// Register any new changes.
	for _, c := range changes {
		k := changeKey(c)
		if _, exists := d.pending[k]; !exists {
			d.pending[k] = entry{firstSeen: now}
		}
	}

	// Build a set of keys present in this batch.
	current := make(map[string]struct{}, len(changes))
	for _, c := range changes {
		current[changeKey(c)] = struct{}{}
	}

	// Evict keys no longer present (the change resolved itself).
	for k := range d.pending {
		if _, ok := current[k]; !ok {
			delete(d.pending, k)
		}
	}

	// Emit changes that have exceeded the window.
	var ready []monitor.Change
	for _, c := range changes {
		k := changeKey(c)
		if e, ok := d.pending[k]; ok && now.Sub(e.firstSeen) >= d.window {
			ready = append(ready, c)
			delete(d.pending, k)
		}
	}
	return ready
}

// changeKey returns a string that uniquely identifies a monitor.Change.
func changeKey(c monitor.Change) string {
	return c.String()
}
