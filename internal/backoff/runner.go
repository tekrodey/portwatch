package backoff

import (
	"context"
	"time"
)

// Runner executes a fallible operation with exponential back-off until it
// succeeds, the context is cancelled, or the maximum number of attempts is
// exhausted (0 means unlimited).
type Runner struct {
	b           *Backoff
	maxAttempts int
	sleep       func(context.Context, time.Duration) error
}

// NewRunner creates a Runner backed by b.
// maxAttempts == 0 means retry forever until ctx is cancelled.
func NewRunner(b *Backoff, maxAttempts int) *Runner {
	return &Runner{
		b:           b,
		maxAttempts: maxAttempts,
		sleep:       contextSleep,
	}
}

// Run calls fn repeatedly, backing off between failures.
// It returns the first nil error from fn, or the last non-nil error once
// attempts are exhausted, or ctx.Err() if the context is cancelled.
func (r *Runner) Run(ctx context.Context, fn func() error) error {
	var last error
	for attempt := 0; ; attempt++ {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		last = fn()
		if last == nil {
			return nil
		}
		if r.maxAttempts > 0 && attempt+1 >= r.maxAttempts {
			return last
		}
		wait := r.b.Interval(attempt)
		if err := r.sleep(ctx, wait); err != nil {
			return err
		}
	}
}

func contextSleep(ctx context.Context, d time.Duration) error {
	select {
	case <-time.After(d):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
