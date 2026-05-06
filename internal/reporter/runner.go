package reporter

import (
	"context"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// PortLister is the interface satisfied by monitor.Monitor for listing
// currently open ports.
type PortLister interface {
	OpenPorts() []monitor.Port
}

// Runner periodically calls a PortLister and writes a summary via Reporter.
type Runner struct {
	lister   PortLister
	rep      *Reporter
	interval time.Duration
}

// NewRunner returns a Runner that queries lister every interval and writes
// summaries via rep.
func NewRunner(lister PortLister, rep *Reporter, interval time.Duration) *Runner {
	if interval <= 0 {
		interval = time.Minute
	}
	return &Runner{lister: lister, rep: rep, interval: interval}
}

// Run blocks, writing a report on every tick, until ctx is cancelled.
func (r *Runner) Run(ctx context.Context) error {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case t := <-ticker.C:
			ports := r.lister.OpenPorts()
			if err := r.rep.Write(ports, t); err != nil {
				return err
			}
		}
	}
}
