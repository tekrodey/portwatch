package batch

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/scanner"
)

func makeChange(port int) monitor.Change {
	return monitor.Change{
		Port:      scanner.Port{Port: port, Proto: "tcp"},
		Direction: monitor.Opened,
	}
}

func TestFlushReturnNilBeforeThreshold(t *testing.T) {
	b := New(3, time.Minute)
	b.Add([]monitor.Change{makeChange(80)})
	if got := b.Flush(); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestFlushReturnsBatchWhenSizeReached(t *testing.T) {
	b := New(2, time.Minute)
	b.Add([]monitor.Change{makeChange(80), makeChange(443)})
	got := b.Flush()
	if len(got) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(got))
	}
}

func TestFlushResetsBufferAfterSizeFlush(t *testing.T) {
	b := New(2, time.Minute)
	b.Add([]monitor.Change{makeChange(80), makeChange(443)})
	b.Flush()
	if got := b.Flush(); got != nil {
		t.Fatalf("expected nil after reset, got %v", got)
	}
}

func TestFlushReturnsBatchWhenIntervalElapsed(t *testing.T) {
	now := time.Now()
	clock := func() time.Time { return now }
	b := newWithClock(100, 50*time.Millisecond, clock)
	b.Add([]monitor.Change{makeChange(8080)})

	// advance clock past interval
	now = now.Add(100 * time.Millisecond)
	got := b.Flush()
	if len(got) != 1 {
		t.Fatalf("expected 1 change after interval, got %d", len(got))
	}
}

func TestForceFlushReturnsAllChanges(t *testing.T) {
	b := New(100, time.Hour)
	b.Add([]monitor.Change{makeChange(22), makeChange(3306)})
	got := b.ForceFlush()
	if len(got) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(got))
	}
}

func TestForceFlushOnEmptyReturnsNil(t *testing.T) {
	b := New(10, time.Minute)
	if got := b.ForceFlush(); got != nil {
		t.Fatalf("expected nil on empty force flush, got %v", got)
	}
}

func TestAddMultipleBatches(t *testing.T) {
	b := New(3, time.Hour)
	b.Add([]monitor.Change{makeChange(80)})
	b.Add([]monitor.Change{makeChange(443), makeChange(8080)})
	got := b.Flush()
	if len(got) != 3 {
		t.Fatalf("expected 3 changes from two adds, got %d", len(got))
	}
}
