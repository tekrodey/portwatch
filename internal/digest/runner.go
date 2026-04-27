package digest

import (
	"context"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// Runner drives a Digest on a fixed tick schedule, consuming changes from a channel.
type Runner struct {
	digest   *Digest
	changes  <-chan []monitor.Change
	interval time.Duration
}

// NewRunner creates a Runner that feeds changes into d and flushes on interval.
func NewRunner(d *Digest, changes <-chan []monitor.Change, interval time.Duration) *Runner {
	return &Runner{
		digest:   d,
		changes:  changes,
		interval: interval,
	}
}

// Run blocks until ctx is cancelled, accumulating changes and flushing summaries.
func (r *Runner) Run(ctx context.Context) {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			r.digest.Flush()
			return
		case batch, ok := <-r.changes:
			if !ok {
				return
			}
			r.digest.Add(batch)
		case <-ticker.C:
			r.digest.Flush()
		}
	}
}
