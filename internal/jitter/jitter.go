// Package jitter adds randomised delay to a base duration to avoid
// thundering-herd problems when multiple goroutines wake at the same time.
package jitter

import (
	"math/rand"
	"sync"
	"time"
)

// Source is a function that returns a pseudo-random float64 in [0, 1).
type Source func() float64

// Jitter applies a random fraction of the base duration on top of the base,
// scaling by the provided factor (0 ≤ factor ≤ 1).
//
// Example: base=1s, factor=0.5 → result is in [1s, 1.5s).
type Jitter struct {
	mu     sync.Mutex
	source Source
	factor float64
}

// New returns a Jitter that uses the global math/rand source and the given
// factor. factor is clamped to [0, 1].
func New(factor float64) *Jitter {
	return newWithSource(factor, rand.Float64)
}

func newWithSource(factor float64, src Source) *Jitter {
	if factor < 0 {
		factor = 0
	}
	if factor > 1 {
		factor = 1
	}
	return &Jitter{source: src, factor: factor}
}

// Apply returns base + a random fraction of (base * factor).
// If base is zero or negative the value is returned unchanged.
func (j *Jitter) Apply(base time.Duration) time.Duration {
	if base <= 0 || j.factor == 0 {
		return base
	}
	j.mu.Lock()
	r := j.source()
	j.mu.Unlock()
	extra := time.Duration(float64(base) * j.factor * r)
	return base + extra
}
