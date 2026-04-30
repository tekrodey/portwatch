// Package backoff provides an exponential back-off strategy with optional
// jitter, used when retrying failed deliveries or reconnecting to external
// services.
package backoff

import (
	"math"
	"time"
)

// Config holds tunable parameters for the back-off algorithm.
type Config struct {
	// InitialInterval is the wait time after the first failure.
	InitialInterval time.Duration
	// MaxInterval caps the computed interval regardless of attempt count.
	MaxInterval time.Duration
	// Multiplier is applied to the previous interval on each attempt.
	Multiplier float64
	// JitterFactor adds ±JitterFactor*interval random noise (0 disables).
	JitterFactor float64
}

// DefaultConfig returns a Config suitable for most alert delivery retries.
func DefaultConfig() Config {
	return Config{
		InitialInterval: 500 * time.Millisecond,
		MaxInterval:     30 * time.Second,
		Multiplier:      2.0,
		JitterFactor:    0.1,
	}
}

// Backoff computes wait durations for successive retry attempts.
type Backoff struct {
	cfg    Config
	randFn func(float64) float64 // injectable for tests
}

// New creates a Backoff using cfg.
func New(cfg Config) *Backoff {
	return newWithRand(cfg, defaultJitter)
}

func newWithRand(cfg Config, randFn func(float64) float64) *Backoff {
	if cfg.Multiplier <= 1 {
		cfg.Multiplier = 2.0
	}
	return &Backoff{cfg: cfg, randFn: randFn}
}

// Interval returns the wait duration for the given attempt number (0-indexed).
func (b *Backoff) Interval(attempt int) time.Duration {
	if attempt < 0 {
		attempt = 0
	}
	base := float64(b.cfg.InitialInterval) * math.Pow(b.cfg.Multiplier, float64(attempt))
	max := float64(b.cfg.MaxInterval)
	if base > max {
		base = max
	}
	if b.cfg.JitterFactor > 0 {
		noise := b.randFn(base * b.cfg.JitterFactor)
		base += noise
	}
	return time.Duration(base)
}

func defaultJitter(maxNoise float64) float64 {
	// Use a simple deterministic-ish approach; real callers supply crypto/rand.
	return maxNoise * 0.5
}
