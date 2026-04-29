// Package batch groups a stream of changes into fixed-size or time-bounded
// batches before forwarding them to the next stage.
package batch

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// Batcher accumulates changes and flushes them either when the batch reaches
// MaxSize or when the Interval elapses, whichever comes first.
type Batcher struct {
	mu       sync.Mutex
	buf      []monitor.Change
	maxSize  int
	interval time.Duration
	clock    func() time.Time
	lastFlush time.Time
}

// New returns a Batcher with the given maximum batch size and flush interval.
func New(maxSize int, interval time.Duration) *Batcher {
	return newWithClock(maxSize, interval, time.Now)
}

func newWithClock(maxSize int, interval time.Duration, clock func() time.Time) *Batcher {
	return &Batcher{
		maxSize:  maxSize,
		interval: interval,
		clock:    clock,
		lastFlush: clock(),
	}
}

// Add appends changes to the internal buffer.
func (b *Batcher) Add(changes []monitor.Change) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.buf = append(b.buf, changes...)
}

// Flush returns the accumulated changes if the batch is ready (size or time
// threshold reached) and resets the buffer. Returns nil if not yet ready.
func (b *Batcher) Flush() []monitor.Change {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := b.clock()
	sizeReady := b.maxSize > 0 && len(b.buf) >= b.maxSize
	timeReady := b.interval > 0 && now.Sub(b.lastFlush) >= b.interval

	if !sizeReady && !timeReady {
		return nil
	}

	out := make([]monitor.Change, len(b.buf))
	copy(out, b.buf)
	b.buf = b.buf[:0]
	b.lastFlush = now
	return out
}

// ForceFlush returns all buffered changes regardless of thresholds.
func (b *Batcher) ForceFlush() []monitor.Change {
	b.mu.Lock()
	defer b.mu.Unlock()
	if len(b.buf) == 0 {
		return nil
	}
	out := make([]monitor.Change, len(b.buf))
	copy(out, b.buf)
	b.buf = b.buf[:0]
	b.lastFlush = b.clock()
	return out
}
