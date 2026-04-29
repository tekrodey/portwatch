package circuit

import (
	"errors"
	"testing"
	"time"
)

var errFake = errors.New("fake error")

// fakeClock returns a function that returns a fixed time, allowing tests to
// control the passage of time.
func fakeClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestClosedByDefault(t *testing.T) {
	b := New(3, time.Minute)
	if b.State() != StateClosed {
		t.Fatal("expected StateClosed on a new breaker")
	}
}

func TestDoSuccessResetsFailures(t *testing.T) {
	b := New(3, time.Minute)
	_ = b.Do(func() error { return errFake })
	_ = b.Do(func() error { return errFake })
	if err := b.Do(func() error { return nil }); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if b.State() != StateClosed {
		t.Fatal("expected StateClosed after successful call")
	}
}

func TestOpensAfterMaxFailures(t *testing.T) {
	b := newWithClock(3, time.Minute, fakeClock(time.Now()))
	for i := 0; i < 3; i++ {
		_ = b.Do(func() error { return errFake })
	}
	if b.State() != StateOpen {
		t.Fatal("expected StateOpen after max failures")
	}
}

func TestDoReturnsErrOpenWhenOpen(t *testing.T) {
	base := time.Now()
	b := newWithClock(2, time.Minute, fakeClock(base))
	_ = b.Do(func() error { return errFake })
	_ = b.Do(func() error { return errFake })

	err := b.Do(func() error { return nil })
	if !errors.Is(err, ErrOpen) {
		t.Fatalf("expected ErrOpen, got %v", err)
	}
}

func TestResetsAfterWindow(t *testing.T) {
	base := time.Now()
	clock := fakeClock(base)
	b := newWithClock(2, 30*time.Second, clock)
	_ = b.Do(func() error { return errFake })
	_ = b.Do(func() error { return errFake })

	if b.State() != StateOpen {
		t.Fatal("expected open")
	}

	// Advance clock past the reset window.
	b.now = fakeClock(base.Add(31 * time.Second))

	if b.State() != StateClosed {
		t.Fatal("expected StateClosed after reset window")
	}
}

func TestDoAllowsCallAfterReset(t *testing.T) {
	base := time.Now()
	b := newWithClock(2, 10*time.Second, fakeClock(base))
	_ = b.Do(func() error { return errFake })
	_ = b.Do(func() error { return errFake })

	b.now = fakeClock(base.Add(11 * time.Second))

	if err := b.Do(func() error { return nil }); err != nil {
		t.Fatalf("expected successful call after reset, got %v", err)
	}
}
