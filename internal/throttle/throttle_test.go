package throttle

import (
	"testing"
	"time"
)

// fakeClock returns a function that returns a fixed time, with an optional
// advance helper.
func fakeClock(start time.Time) (Clock, func(time.Duration)) {
	current := start
	clock := func() time.Time { return current }
	advance := func(d time.Duration) { current = current.Add(d) }
	return clock, advance
}

func TestAllowFirstCallAlwaysPasses(t *testing.T) {
	clock, _ := fakeClock(time.Now())
	th := NewWithClock(5*time.Second, clock)

	if !th.Allow(8080, "tcp", "opened") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllowSuppressesWithinInterval(t *testing.T) {
	clock, advance := fakeClock(time.Now())
	th := NewWithClock(10*time.Second, clock)

	th.Allow(443, "tcp", "opened")
	advance(5 * time.Second)

	if th.Allow(443, "tcp", "opened") {
		t.Fatal("expected call within interval to be suppressed")
	}
}

func TestAllowPassesAfterIntervalExpires(t *testing.T) {
	clock, advance := fakeClock(time.Now())
	th := NewWithClock(10*time.Second, clock)

	th.Allow(443, "tcp", "opened")
	advance(11 * time.Second)

	if !th.Allow(443, "tcp", "opened") {
		t.Fatal("expected call after interval to be allowed")
	}
}

func TestAllowDistinguishesByDirection(t *testing.T) {
	clock, _ := fakeClock(time.Now())
	th := NewWithClock(10*time.Second, clock)

	th.Allow(80, "tcp", "opened")

	// Different direction should be independent.
	if !th.Allow(80, "tcp", "closed") {
		t.Fatal("expected different direction to be allowed")
	}
}

func TestAllowDistinguishesByProtocol(t *testing.T) {
	clock, _ := fakeClock(time.Now())
	th := NewWithClock(10*time.Second, clock)

	th.Allow(53, "tcp", "opened")

	if !th.Allow(53, "udp", "opened") {
		t.Fatal("expected different protocol to be allowed")
	}
}

func TestResetAllowsImmediateRetrigger(t *testing.T) {
	clock, _ := fakeClock(time.Now())
	th := NewWithClock(30*time.Second, clock)

	th.Allow(22, "tcp", "opened")
	th.Reset(22, "tcp", "opened")

	if !th.Allow(22, "tcp", "opened") {
		t.Fatal("expected allow after reset")
	}
}

func TestLenTracksKeys(t *testing.T) {
	clock, _ := fakeClock(time.Now())
	th := NewWithClock(30*time.Second, clock)

	if th.Len() != 0 {
		t.Fatalf("expected 0 keys, got %d", th.Len())
	}

	th.Allow(80, "tcp", "opened")
	th.Allow(443, "tcp", "opened")

	if th.Len() != 2 {
		t.Fatalf("expected 2 keys, got %d", th.Len())
	}
}
