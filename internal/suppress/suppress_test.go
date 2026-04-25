package suppress

import (
	"testing"
	"time"
)

func newFakeSuppress(base time.Time) *Suppressor {
	s := New()
	s.now = func() time.Time { return base }
	return s
}

func TestIsSuppressedReturnsFalseForUnknownKey(t *testing.T) {
	s := newFakeSuppress(time.Now())
	if s.IsSuppressed("tcp:9000") {
		t.Fatal("expected false for unknown key")
	}
}

func TestSuppressAndIsSuppressed(t *testing.T) {
	base := time.Now()
	s := newFakeSuppress(base)
	s.Suppress("tcp:8080", 5*time.Minute, "maintenance")

	if !s.IsSuppressed("tcp:8080") {
		t.Fatal("expected key to be suppressed")
	}
}

func TestSuppressExpires(t *testing.T) {
	base := time.Now()
	s := newFakeSuppress(base)
	s.Suppress("tcp:8080", 1*time.Second, "")

	// advance clock past expiry
	s.now = func() time.Time { return base.Add(2 * time.Second) }

	if s.IsSuppressed("tcp:8080") {
		t.Fatal("expected suppression to have expired")
	}
}

func TestLiftRemovesSuppression(t *testing.T) {
	base := time.Now()
	s := newFakeSuppress(base)
	s.Suppress("udp:53", 10*time.Minute, "dns")
	s.Lift("udp:53")

	if s.IsSuppressed("udp:53") {
		t.Fatal("expected suppression to be lifted")
	}
}

func TestActiveReturnsOnlyLiveSuppression(t *testing.T) {
	base := time.Now()
	s := newFakeSuppress(base)
	s.Suppress("tcp:443", 10*time.Minute, "tls")
	s.Suppress("tcp:80", 1*time.Millisecond, "")

	// advance so tcp:80 expires
	s.now = func() time.Time { return base.Add(1 * time.Second) }

	active := s.Active()
	if len(active) != 1 {
		t.Fatalf("expected 1 active suppression, got %d", len(active))
	}
	if _, ok := active["tcp:443"]; !ok {
		t.Fatal("expected tcp:443 in active suppressions")
	}
}

func TestSuppressOverwritesExistingEntry(t *testing.T) {
	base := time.Now()
	s := newFakeSuppress(base)
	s.Suppress("tcp:22", 1*time.Minute, "first")
	s.Suppress("tcp:22", 10*time.Minute, "second")

	e := s.Active()["tcp:22"]
	if e.Reason != "second" {
		t.Fatalf("expected reason 'second', got %q", e.Reason)
	}
}
