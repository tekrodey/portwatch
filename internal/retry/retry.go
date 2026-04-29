package retry

import (
	"context"
	"errors"
	"time"
)

// ErrMaxAttempts is returned when all retry attempts are exhausted.
var ErrMaxAttempts = errors.New("retry: max attempts reached")

// Config holds retry policy parameters.
type Config struct {
	MaxAttempts int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
}

// DefaultConfig returns a sensible default retry configuration.
func DefaultConfig() Config {
	return Config{
		MaxAttempts: 3,
		BaseDelay:   200 * time.Millisecond,
		MaxDelay:    5 * time.Second,
	}
}

// Retryer executes operations with exponential backoff.
type Retryer struct {
	cfg   Config
	sleep func(time.Duration)
}

// New returns a Retryer using the provided Config.
func New(cfg Config) *Retryer {
	return &Retryer{cfg: cfg, sleep: time.Sleep}
}

// newWithSleep is used in tests to inject a fake sleep function.
func newWithSleep(cfg Config, sleep func(time.Duration)) *Retryer {
	return &Retryer{cfg: cfg, sleep: sleep}
}

// Do calls fn up to MaxAttempts times, backing off between failures.
// It stops early if ctx is cancelled or fn returns nil.
func (r *Retryer) Do(ctx context.Context, fn func() error) error {
	delay := r.cfg.BaseDelay
	for attempt := 1; attempt <= r.cfg.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := fn(); err == nil {
			return nil
		}
		if attempt == r.cfg.MaxAttempts {
			break
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
		delay *= 2
		if delay > r.cfg.MaxDelay {
			delay = r.cfg.MaxDelay
		}
	}
	return ErrMaxAttempts
}
