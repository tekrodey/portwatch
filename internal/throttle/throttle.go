// Package throttle limits how frequently alerts are emitted for a given port
// by tracking the last alert time per (port, protocol, direction) key.
package throttle

import (
	"fmt"
	"sync"
	"time"
)

// Clock allows time to be faked in tests.
type Clock func() time.Time

// Throttle decides whether an alert for a given key should be emitted based
// on a minimum interval between successive alerts.
type Throttle struct {
	mu       sync.Mutex
	last     map[string]time.Time
	interval time.Duration
	now      Clock
}

// New returns a Throttle that suppresses repeated alerts within interval.
func New(interval time.Duration) *Throttle {
	return &Throttle{
		last:     make(map[string]time.Time),
		interval: interval,
		now:      time.Now,
	}
}

// NewWithClock returns a Throttle using a custom clock (useful for testing).
func NewWithClock(interval time.Duration, clock Clock) *Throttle {
	t := New(interval)
	t.now = clock
	return t
}

// Allow returns true if enough time has elapsed since the last alert for key.
// It records the current time as the new last-alert time when it returns true.
func (t *Throttle) Allow(port int, proto, direction string) bool {
	key := fmt.Sprintf("%d/%s/%s", port, proto, direction)
	now := t.now()

	t.mu.Lock()
	defer t.mu.Unlock()

	if last, ok := t.last[key]; ok && now.Sub(last) < t.interval {
		return false
	}
	t.last[key] = now
	return true
}

// Reset clears the recorded time for a specific key, allowing the next call
// to Allow to pass unconditionally.
func (t *Throttle) Reset(port int, proto, direction string) {
	key := fmt.Sprintf("%d/%s/%s", port, proto, direction)
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.last, key)
}

// Len returns the number of tracked keys.
func (t *Throttle) Len() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.last)
}
