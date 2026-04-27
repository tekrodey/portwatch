package dedup

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

func makeChange(proto string, port int, dir string) monitor.Change {
	return monitor.Change{
		Protocol:  proto,
		Port:      port,
		Direction: dir,
	}
}

func TestFirstOccurrenceAlwaysPasses(t *testing.T) {
	d := New(5 * time.Second)
	changes := []monitor.Change{
		makeChange("tcp", 8080, "opened"),
		makeChange("tcp", 9090, "opened"),
	}
	out := d.Filter(changes)
	if len(out) != 2 {
		t.Fatalf("expected 2, got %d", len(out))
	}
}

func TestDuplicateWithinWindowSuppressed(t *testing.T) {
	now := time.Now()
	clock := func() time.Time { return now }
	d := newWithClock(10*time.Second, clock)

	c := makeChange("tcp", 8080, "opened")
	out1 := d.Filter([]monitor.Change{c})
	if len(out1) != 1 {
		t.Fatalf("expected 1 on first call, got %d", len(out1))
	}

	out2 := d.Filter([]monitor.Change{c})
	if len(out2) != 0 {
		t.Fatalf("expected 0 on duplicate within window, got %d", len(out2))
	}
}

func TestDuplicateAfterWindowPasses(t *testing.T) {
	now := time.Now()
	clock := func() time.Time { return now }
	d := newWithClock(5*time.Second, clock)

	c := makeChange("tcp", 443, "closed")
	d.Filter([]monitor.Change{c})

	// advance clock past window
	now = now.Add(6 * time.Second)
	out := d.Filter([]monitor.Change{c})
	if len(out) != 1 {
		t.Fatalf("expected 1 after window expired, got %d", len(out))
	}
}

func TestDistinguishesProtocolAndDirection(t *testing.T) {
	d := New(10 * time.Second)
	changes := []monitor.Change{
		makeChange("tcp", 80, "opened"),
		makeChange("udp", 80, "opened"),
		makeChange("tcp", 80, "closed"),
	}
	out := d.Filter(changes)
	if len(out) != 3 {
		t.Fatalf("expected 3 distinct events, got %d", len(out))
	}
}

func TestEmptyInputReturnsNil(t *testing.T) {
	d := New(5 * time.Second)
	out := d.Filter(nil)
	if len(out) != 0 {
		t.Fatalf("expected empty output, got %d", len(out))
	}
}
