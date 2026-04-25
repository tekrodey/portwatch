package ratelimit

import (
	"testing"
	"time"
)

func newFakeLimiter(cooldown time.Duration) (*Limiter, *time.Time) {
	t := time.Now()
	l := New(cooldown)
	l.now = func() time.Time { return t }
	return l, &t
}

func TestAllowFirstCallAlwaysPasses(t *testing.T) {
	l, _ := newFakeLimiter(5 * time.Second)
	if !l.Allow(8080, "tcp", "opened") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllowSuppressesWithinCooldown(t *testing.T) {
	l, clock := newFakeLimiter(10 * time.Second)

	l.Allow(8080, "tcp", "opened") // prime
	*clock = clock.Add(5 * time.Second)

	if l.Allow(8080, "tcp", "opened") {
		t.Fatal("expected call within cooldown to be suppressed")
	}
}

func TestAllowPassesAfterCooldownExpires(t *testing.T) {
	l, clock := newFakeLimiter(10 * time.Second)

	l.Allow(8080, "tcp", "opened")
	*clock = clock.Add(11 * time.Second)

	if !l.Allow(8080, "tcp", "opened") {
		t.Fatal("expected call after cooldown to be allowed")
	}
}

func TestAllowDistinguishesDirection(t *testing.T) {
	l, _ := newFakeLimiter(10 * time.Second)

	l.Allow(8080, "tcp", "opened")

	if !l.Allow(8080, "tcp", "closed") {
		t.Fatal("expected different direction to be allowed independently")
	}
}

func TestAllowDistinguishesProtocol(t *testing.T) {
	l, _ := newFakeLimiter(10 * time.Second)

	l.Allow(8080, "tcp", "opened")

	if !l.Allow(8080, "udp", "opened") {
		t.Fatal("expected different protocol to be allowed independently")
	}
}

func TestResetClearsSuppressionState(t *testing.T) {
	l, clock := newFakeLimiter(30 * time.Second)

	l.Allow(9000, "tcp", "opened")
	*clock = clock.Add(5 * time.Second)
	l.Reset(9000, "tcp", "opened")

	if !l.Allow(9000, "tcp", "opened") {
		t.Fatal("expected allow after reset even within cooldown")
	}
}

func TestStatsReturnsString(t *testing.T) {
	l, _ := newFakeLimiter(5 * time.Second)
	l.Allow(1234, "tcp", "opened")
	s := l.Stats()
	if s == "" {
		t.Fatal("expected non-empty stats string")
	}
}
