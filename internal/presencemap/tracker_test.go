package presencemap

import (
	"testing"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/scanner"
)

func makeChange(proto string, number int, dir monitor.Direction) monitor.Change {
	return monitor.Change{
		Port:      scanner.Port{Proto: proto, Number: number},
		Direction: dir,
	}
}

func TestApplyOpenedAddsEntry(t *testing.T) {
	m := New()
	tr := NewTracker(m)
	tr.Apply([]monitor.Change{
		makeChange("tcp", 8080, monitor.Opened),
	})
	_, ok := m.Get("tcp:8080")
	if !ok {
		t.Error("expected entry for opened port")
	}
}

func TestApplyClosedRemovesEntry(t *testing.T) {
	m := New()
	tr := NewTracker(m)
	tr.Apply([]monitor.Change{makeChange("tcp", 8080, monitor.Opened)})
	tr.Apply([]monitor.Change{makeChange("tcp", 8080, monitor.Closed)})
	_, ok := m.Get("tcp:8080")
	if ok {
		t.Error("expected entry to be removed after closed")
	}
}

func TestApplyEmptyChangesIsNoOp(t *testing.T) {
	m := New()
	tr := NewTracker(m)
	tr.Apply(nil)
	snap := tr.Snapshot()
	if len(snap) != 0 {
		t.Errorf("expected empty snapshot, got %d entries", len(snap))
	}
}

func TestSnapshotReturnsCopy(t *testing.T) {
	m := New()
	tr := NewTracker(m)
	tr.Apply([]monitor.Change{
		makeChange("udp", 53, monitor.Opened),
		makeChange("tcp", 443, monitor.Opened),
	})
	snap := tr.Snapshot()
	if len(snap) != 2 {
		t.Errorf("snapshot length: got %d want 2", len(snap))
	}
	// Mutating the snapshot must not affect the underlying map.
	delete(snap, "udp:53")
	_, ok := m.Get("udp:53")
	if !ok {
		t.Error("underlying map should not be affected by snapshot mutation")
	}
}

func TestKeyFormat(t *testing.T) {
	c := makeChange("tcp", 22, monitor.Opened)
	if got := key(c); got != "tcp:22" {
		t.Errorf("key: got %q want %q", got, "tcp:22")
	}
}
