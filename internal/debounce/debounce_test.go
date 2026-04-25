package debounce

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/scanner"
)

func makeChange(proto string, port int, dir monitor.Direction) monitor.Change {
	return monitor.Change{
		Port:      scanner.Port{Proto: proto, Number: port},
		Direction: dir,
	}
}

func TestFilterHoldsChangesBelowWindow(t *testing.T) {
	d := New(5 * time.Second)
	now := time.Now()
	d.now = func() time.Time { return now }

	changes := []monitor.Change{makeChange("tcp", 8080, monitor.Opened)}

	got := d.Filter(changes)
	if len(got) != 0 {
		t.Fatalf("expected 0 ready changes, got %d", len(got))
	}
}

func TestFilterReleasesChangeAfterWindow(t *testing.T) {
	d := New(5 * time.Second)
	now := time.Now()
	d.now = func() time.Time { return now }

	changes := []monitor.Change{makeChange("tcp", 8080, monitor.Opened)}

	// First call — registers the change.
	d.Filter(changes)

	// Advance time past the window.
	d.now = func() time.Time { return now.Add(6 * time.Second) }

	got := d.Filter(changes)
	if len(got) != 1 {
		t.Fatalf("expected 1 ready change, got %d", len(got))
	}
}

func TestFilterEvictsResolvedChange(t *testing.T) {
	d := New(5 * time.Second)
	now := time.Now()
	d.now = func() time.Time { return now }

	changes := []monitor.Change{makeChange("tcp", 9090, monitor.Closed)}
	d.Filter(changes)

	// Next call arrives with no changes — the pending entry should be evicted.
	d.now = func() time.Time { return now.Add(3 * time.Second) }
	got := d.Filter(nil)

	if len(got) != 0 {
		t.Fatalf("expected 0 ready changes after eviction, got %d", len(got))
	}
	if len(d.pending) != 0 {
		t.Fatalf("expected pending map to be empty, got %d entries", len(d.pending))
	}
}

func TestFilterDistinguishesDifferentChanges(t *testing.T) {
	d := New(2 * time.Second)
	now := time.Now()
	d.now = func() time.Time { return now }

	changes := []monitor.Change{
		makeChange("tcp", 80, monitor.Opened),
		makeChange("udp", 53, monitor.Opened),
	}
	d.Filter(changes)

	d.now = func() time.Time { return now.Add(3 * time.Second) }
	got := d.Filter(changes)

	if len(got) != 2 {
		t.Fatalf("expected 2 ready changes, got %d", len(got))
	}
}
