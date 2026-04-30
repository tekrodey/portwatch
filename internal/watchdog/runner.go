package watchdog

import (
	"context"
	"log"
	"time"
)

// Runner wraps a Watchdog with a default check function that verifies a
// user-supplied liveness probe completes within the configured timeout. It
// logs hang events to the standard logger.
type Runner struct {
	dog   *Watchdog
	probe func(ctx context.Context) error
}

// NewRunner creates a Runner using the provided Config and liveness probe.
// If cfg.OnHang is nil it is replaced with a default that logs the error.
func NewRunner(cfg Config, probe func(ctx context.Context) error) *Runner {
	if cfg.OnHang == nil {
		cfg.OnHang = func(err error) {
			log.Printf("[watchdog] liveness probe hung or failed: %v", err)
		}
	}
	return &Runner{
		dog:   New(cfg, probe),
		probe: probe,
	}
}

// Start launches the watchdog loop in a background goroutine and returns
// immediately. The loop runs until ctx is cancelled.
func (r *Runner) Start(ctx context.Context) {
	go r.dog.Run(ctx)
}

// Ping is a convenience liveness probe that simply checks whether the context
// is still live. It can be used as the probe argument when no deeper check is
// needed.
func Ping(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

// ScanProbe returns a probe function that calls scanFn and records the
// elapsed time. It is intended to wrap the scanner's Scan method so the
// watchdog can detect a stalled scan loop.
func ScanProbe(scanFn func(ctx context.Context) error) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		start := time.Now()
		err := scanFn(ctx)
		if err != nil {
			log.Printf("[watchdog] scan probe failed after %v: %v", time.Since(start), err)
		}
		return err
	}
}
