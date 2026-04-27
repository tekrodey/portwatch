// Package digest produces periodic summary reports of port change activity.
package digest

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// Summary holds aggregated change counts over a reporting window.
type Summary struct {
	From    time.Time
	To      time.Time
	Opened  []monitor.Change
	Closed  []monitor.Change
}

// String returns a human-readable digest report.
func (s Summary) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "=== Port Digest [%s – %s] ===\n",
		s.From.Format(time.RFC3339), s.To.Format(time.RFC3339))
	fmt.Fprintf(&b, "  Opened: %d  Closed: %d\n", len(s.Opened), len(s.Closed))
	for _, c := range s.Opened {
		fmt.Fprintf(&b, "    + %s\n", c)
	}
	for _, c := range s.Closed {
		fmt.Fprintf(&b, "    - %s\n", c)
	}
	return b.String()
}

// Digest accumulates changes and flushes periodic summaries.
type Digest struct {
	w        io.Writer
	interval time.Duration
	changes  []monitor.Change
	windowStart time.Time
	now      func() time.Time
}

// New creates a Digest that writes summaries to w every interval.
func New(w io.Writer, interval time.Duration) *Digest {
	if w == nil {
		w = os.Stdout
	}
	return &Digest{
		w:           w,
		interval:    interval,
		windowStart: time.Now(),
		now:         time.Now,
	}
}

// Add records a change for inclusion in the next summary.
func (d *Digest) Add(changes []monitor.Change) {
	d.changes = append(d.changes, changes...)
}

// Flush writes the current summary if the interval has elapsed.
// Returns true if a summary was written.
func (d *Digest) Flush() bool {
	if d.now().Sub(d.windowStart) < d.interval {
		return false
	}
	s := d.build()
	fmt.Fprint(d.w, s.String())
	d.changes = nil
	d.windowStart = d.now()
	return true
}

func (d *Digest) build() Summary {
	s := Summary{From: d.windowStart, To: d.now()}
	for _, c := range d.changes {
		if c.Direction == monitor.Opened {
			s.Opened = append(s.Opened, c)
		} else {
			s.Closed = append(s.Closed, c)
		}
	}
	return s
}
