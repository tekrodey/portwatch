// Package suppress provides a mechanism to temporarily silence alerts
// for specific ports, preventing repeated notifications during known
// maintenance windows or expected port state changes.
package suppress

import (
	"sync"
	"time"
)

// Entry holds suppression metadata for a single port+protocol key.
type Entry struct {
	Expires time.Time
	Reason  string
}

// Suppressor tracks active suppressions keyed by "proto:port".
type Suppressor struct {
	mu      sync.Mutex
	entries map[string]Entry
	now     func() time.Time
}

// New returns a new Suppressor.
func New() *Suppressor {
	return &Suppressor{
		entries: make(map[string]Entry),
		now:     time.Now,
	}
}

// Suppress adds a suppression for the given key (e.g. "tcp:8080") that
// expires after duration d. An empty reason is valid.
func (s *Suppressor) Suppress(key string, d time.Duration, reason string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[key] = Entry{
		Expires: s.now().Add(d),
		Reason:  reason,
	}
}

// IsSuppressed reports whether key is currently suppressed.
func (s *Suppressor) IsSuppressed(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.entries[key]
	if !ok {
		return false
	}
	if s.now().After(e.Expires) {
		delete(s.entries, key)
		return false
	}
	return true
}

// Lift removes a suppression for key immediately, regardless of expiry.
func (s *Suppressor) Lift(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, key)
}

// Active returns a snapshot of all currently active suppressions.
func (s *Suppressor) Active() map[string]Entry {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make(map[string]Entry, len(s.entries))
	now := s.now()
	for k, e := range s.entries {
		if now.Before(e.Expires) {
			out[k] = e
		}
	}
	return out
}
