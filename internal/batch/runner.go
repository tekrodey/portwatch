package batch

import (
	"context"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// Runner periodically polls the Batcher and forwards flushed batches to a
// handler function. It also performs a final ForceFlush when the context is
// cancelled so no changes are silently dropped.
type Runner struct {
	batcher  *Batcher
	tick     time.Duration
	handler  func([]monitor.Change)
}

// NewRunner creates a Runner that ticks every interval and calls handler with
// each non-empty batch.
func NewRunner(b *Batcher, interval time.Duration, handler func([]monitor.Change)) *Runner {
	return &Runner{
		batcher: b,
		tick:    interval,
		handler: handler,
	}
}

// Run blocks until ctx is cancelled, flushing the batcher on every tick.
func (r *Runner) Run(ctx context.Context) {
	ticker := time.NewTicker(r.tick)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if batch := r.batcher.Flush(); len(batch) > 0 {
				r.handler(batch)
			}
		case <-ctx.Done():
			if batch := r.batcher.ForceFlush(); len(batch) > 0 {
				r.handler(batch)
			}
			return
		}
	}
}
