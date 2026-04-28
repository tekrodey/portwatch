// Package sampler provides probabilistic sampling for port change events,
// allowing high-volume environments to reduce alert noise by forwarding only
// a fraction of qualifying changes.
package sampler

import (
	"math/rand"
	"sync"

	"github.com/user/portwatch/internal/monitor"
)

// Sampler forwards changes with a given probability in the range (0.0, 1.0].
// A rate of 1.0 forwards every change; 0.5 forwards roughly half.
type Sampler struct {
	mu   sync.Mutex
	rate float64
	rng  *rand.Rand
}

// New returns a Sampler that forwards changes at the given rate.
// rate is clamped to [0.0, 1.0].
func New(rate float64, src rand.Source) *Sampler {
	if rate < 0 {
		rate = 0
	}
	if rate > 1 {
		rate = 1
	}
	if src == nil {
		src = rand.NewSource(42)
	}
	return &Sampler{
		rate: rate,
		rng:  rand.New(src),
	}
}

// Sample filters changes, returning only those selected by the sampling rate.
// When rate is 1.0 the original slice is returned unchanged.
func (s *Sampler) Sample(changes []monitor.Change) []monitor.Change {
	if len(changes) == 0 {
		return changes
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.rate >= 1.0 {
		return changes
	}
	out := changes[:0:0]
	for _, c := range changes {
		if s.rng.Float64() < s.rate {
			out = append(out, c)
		}
	}
	return out
}

// Rate returns the configured sampling rate.
func (s *Sampler) Rate() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.rate
}
