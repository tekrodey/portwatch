package rollup

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/scanner"
)

func makeChange(port int, proto, dir string) monitor.Change {
	return monitor.Change{
		Port:      scanner.Port{Number: port, Proto: proto},
		Direction: dir,
	}
}

func TestFirstOccurrenceIsForwarded(t *testing.T) {
	r := New(DefaultWindow)
	if !r.Add(makeChange(8080, "tcp", "opened")) {
		t.Fatal("expected first occurrence to be forwarded")
	}
}

func TestSecondOccurrenceWithinWindowIsAbsorbed(t *testing.T) {
	r := New(DefaultWindow)
	r.Add(makeChange(8080, "tcp", "opened"))
	if r.Add(makeChange(8080, "tcp", "opened")) {
		t.Fatal("expected duplicate within window to be absorbed")
	}
}

func TestOccurrenceAfterWindowIsForwarded(t *testing.T) {
	now := time.Now()
	clock := func() time.Time { return now }
	r := newWithClock(100*time.Millisecond, clock)

	r.Add(makeChange(8080, "tcp", "opened"))

	// advance clock past window
	now = now.Add(200 * time.Millisecond)

	if !r.Add(makeChange(8080, "tcp", "opened")) {
		t.Fatal("expected occurrence after window to be forwarded")
	}
}

func TestDistinctPortsAreIndependent(t *testing.T) {
	r := New(DefaultWindow)
	r.Add(makeChange(8080, "tcp", "opened"))
	if !r.Add(makeChange(9090, "tcp", "opened")) {
		t.Fatal("expected different port to be forwarded independently")
	}
}

func TestDistinctProtocolsAreIndependent(t *testing.T) {
	r := New(DefaultWindow)
	r.Add(makeChange(53, "tcp", "opened"))
	if !r.Add(makeChange(53, "udp", "opened")) {
		t.Fatal("expected different protocol to be forwarded independently")
	}
}

func TestFlushEmitsHeldChangesAfterWindow(t *testing.T) {
	now := time.Now()
	clock := func() time.Time { return now }
	r := newWithClock(100*time.Millisecond, clock)

	r.Add(makeChange(8080, "tcp", "opened"))
	// second add is absorbed
	r.Add(makeChange(8080, "tcp", "opened"))

	// before window: nothing to flush
	if got := r.Flush(); len(got) != 0 {
		t.Fatalf("expected 0 flushed, got %d", len(got))
	}

	// advance past window
	now = now.Add(200 * time.Millisecond)

	got := r.Flush()
	if len(got) != 1 {
		t.Fatalf("expected 1 flushed, got %d", len(got))
	}
}

func TestFlushRemovesEntryFromBucket(t *testing.T) {
	now := time.Now()
	clock := func() time.Time { return now }
	r := newWithClock(50*time.Millisecond, clock)

	r.Add(makeChange(443, "tcp", "closed"))
	now = now.Add(100 * time.Millisecond)
	r.Flush()

	// after flush the bucket is clear; next add should be forwarded
	if !r.Add(makeChange(443, "tcp", "closed")) {
		t.Fatal("expected forwarding after bucket was cleared by Flush")
	}
}
