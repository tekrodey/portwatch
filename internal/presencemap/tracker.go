package presencemap

import (
	"fmt"

	"github.com/user/portwatch/internal/monitor"
)

// Tracker wraps Map and integrates with the monitor change stream,
// touching entries for opened ports and deleting entries for closed ones.
type Tracker struct {
	m *Map
}

// NewTracker returns a Tracker backed by the provided Map.
func NewTracker(m *Map) *Tracker {
	return &Tracker{m: m}
}

// Apply processes a slice of monitor.Change values, updating the presence map.
func (t *Tracker) Apply(changes []monitor.Change) {
	for _, c := range changes {
		k := key(c)
		switch c.Direction {
		case monitor.Opened:
			t.m.Touch(k)
		case monitor.Closed:
			t.m.Delete(k)
		}
	}
}

// Snapshot returns a copy of all current entries keyed by port string.
func (t *Tracker) Snapshot() map[string]Entry {
	t.m.mu.RLock()
	defer t.m.mu.RUnlock()
	out := make(map[string]Entry, len(t.m.entries))
	for k, e := range t.m.entries {
		out[k] = *e
	}
	return out
}

func key(c monitor.Change) string {
	return fmt.Sprintf("%s:%d", c.Port.Proto, c.Port.Number)
}
