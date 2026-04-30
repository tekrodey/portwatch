// Package limiter provides a concurrency limiter that caps the number of
// in-flight operations using a semaphore channel.
package limiter

import (
	"context"
	"errors"
)

// ErrLimitExceeded is returned when Acquire is called with a context that is
// already cancelled or when the semaphore is full and the context expires.
var ErrLimitExceeded = errors.New("limiter: limit exceeded")

// Limiter caps concurrent operations to a fixed number of slots.
type Limiter struct {
	sem chan struct{}
}

// New returns a Limiter that allows at most n concurrent operations.
// It panics if n is less than 1.
func New(n int) *Limiter {
	if n < 1 {
		panic("limiter: n must be >= 1")
	}
	return &Limiter{sem: make(chan struct{}, n)}
}

// Acquire blocks until a slot is available or ctx is cancelled.
// Callers must call Release exactly once after a successful Acquire.
func (l *Limiter) Acquire(ctx context.Context) error {
	select {
	case l.sem <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ErrLimitExceeded
	}
}

// Release frees one slot. It must only be called after a successful Acquire.
func (l *Limiter) Release() {
	<-l.sem
}

// Available returns the number of slots currently free.
func (l *Limiter) Available() int {
	return cap(l.sem) - len(l.sem)
}

// Capacity returns the total number of slots.
func (l *Limiter) Capacity() int {
	return cap(l.sem)
}
