// Package circuit implements a simple circuit-breaker that stops forwarding
// changes to a downstream sink after too many consecutive errors, and
// automatically resets after a configurable recovery window.
package circuit

import (
	"errors"
	"sync"
	"time"
)

// ErrOpen is returned by Send when the circuit is open.
var ErrOpen = errors.New("circuit open: too many consecutive errors")

// State represents the current circuit state.
type State int

const (
	StateClosed State = iota
	StateOpen
)

// Breaker trips after MaxFailures consecutive errors and resets after
// ResetAfter has elapsed.
type Breaker struct {
	mu          sync.Mutex
	failures    int
	maxFailures int
	resetAfter  time.Duration
	openedAt    time.Time
	now         func() time.Time
}

// New returns a Breaker that opens after maxFailures consecutive errors and
// resets after resetAfter.
func New(maxFailures int, resetAfter time.Duration) *Breaker {
	return newWithClock(maxFailures, resetAfter, time.Now)
}

func newWithClock(maxFailures int, resetAfter time.Duration, now func() time.Time) *Breaker {
	return &Breaker{
		maxFailures: maxFailures,
		resetAfter:  resetAfter,
		now:         now,
	}
}

// State returns the current state of the breaker, automatically transitioning
// from Open back to Closed if the reset window has elapsed.
func (b *Breaker) State() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state()
}

// state must be called with b.mu held.
func (b *Breaker) state() State {
	if b.failures >= b.maxFailures {
		if b.now().Sub(b.openedAt) >= b.resetAfter {
			b.failures = 0
			return StateClosed
		}
		return StateOpen
	}
	return StateClosed
}

// Do calls fn if the circuit is closed. If fn returns an error the failure
// counter is incremented; a successful call resets it. Returns ErrOpen when
// the circuit is open.
func (b *Breaker) Do(fn func() error) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.state() == StateOpen {
		return ErrOpen
	}

	err := fn()
	if err != nil {
		b.failures++
		if b.failures >= b.maxFailures {
			b.openedAt = b.now()
		}
		return err
	}

	b.failures = 0
	return nil
}
