// Package ratelimit provides a simple token-bucket rate limiter to suppress
// repeated alerts for the same port change within a cooldown window.
package ratelimit

import (
	"fmt"
	"sync"
	"time"
)

// key identifies a unique port+protocol+direction combination.
type key struct {
	port     int
	protocol string
	direction string // "opened" | "closed"
}

// Limiter tracks the last alert time per port event and suppresses duplicates
// that arrive within the configured cooldown period.
type Limiter struct {
	mu       sync.Mutex
	cooldown time.Duration
	last     map[key]time.Time
	now      func() time.Time // injectable for testing
}

// New creates a Limiter with the given cooldown duration.
// Calls with the same port/protocol/direction within cooldown are suppressed.
func New(cooldown time.Duration) *Limiter {
	return &Limiter{
		cooldown: cooldown,
		last:     make(map[key]time.Time),
		now:      time.Now,
	}
}

// Allow returns true if the event for the given port, protocol and direction
// should be forwarded, and false if it is within the cooldown window.
func (l *Limiter) Allow(port int, protocol, direction string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	k := key{port: port, protocol: protocol, direction: direction}
	now := l.now()

	if t, ok := l.last[k]; ok && now.Sub(t) < l.cooldown {
		return false
	}

	l.last[k] = now
	return true
}

// Reset clears the suppression state for a specific port/protocol/direction.
func (l *Limiter) Reset(port int, protocol, direction string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.last, key{port: port, protocol: protocol, direction: direction})
}

// Stats returns a human-readable summary of currently tracked suppressions.
func (l *Limiter) Stats() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return fmt.Sprintf("ratelimit: tracking %d suppressed events (cooldown %s)", len(l.last), l.cooldown)
}
