package watchdog

import (
	"context"
	"time"
)

// Watchdog periodically calls a check function and invokes a handler if the
// check does not complete within the configured timeout.
type Watchdog struct {
	interval time.Duration
	timeout  time.Duration
	check    func(ctx context.Context) error
	onHang   func(err error)
	clock    func() time.Time
}

// Config holds Watchdog configuration.
type Config struct {
	Interval time.Duration
	Timeout  time.Duration
	OnHang   func(err error)
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Interval: 30 * time.Second,
		Timeout:  10 * time.Second,
		OnHang:   func(err error) {},
	}
}

// New creates a Watchdog that calls check on every interval and invokes
// cfg.OnHang if the check exceeds cfg.Timeout.
func New(cfg Config, check func(ctx context.Context) error) *Watchdog {
	if cfg.OnHang == nil {
		cfg.OnHang = func(err error) {}
	}
	return &Watchdog{
		interval: cfg.Interval,
		timeout:  cfg.Timeout,
		check:    check,
		onHang:   cfg.OnHang,
		clock:    time.Now,
	}
}

// Run starts the watchdog loop. It blocks until ctx is cancelled.
func (w *Watchdog) Run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.probe(ctx)
		}
	}
}

func (w *Watchdog) probe(parent context.Context) {
	ctx, cancel := context.WithTimeout(parent, w.timeout)
	defer cancel()

	done := make(chan error, 1)
	go func() { done <- w.check(ctx) }()

	select {
	case err := <-done:
		if err != nil {
			w.onHang(err)
		}
	case <-ctx.Done():
		w.onHang(ctx.Err())
	}
}
