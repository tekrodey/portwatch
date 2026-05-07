package escalate_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/escalate"
	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/scanner"
)

func makeChange(port int, dir string) monitor.Change {
	return monitor.Change{
		Port:      scanner.Port{Number: port, Protocol: "tcp"},
		Direction: dir,
	}
}

func TestFirstOccurrenceIsNormal(t *testing.T) {
	e := escalate.New(time.Minute)
	c := makeChange(80, "opened")
	if got := e.Assess(c); got != escalate.LevelNormal {
		t.Fatalf("expected Normal, got %d", got)
	}
}

func TestSecondOccurrenceWithinWindowIsElevated(t *testing.T) {
	e := escalate.New(time.Minute)
	c := makeChange(80, "opened")
	e.Assess(c)
	if got := e.Assess(c); got != escalate.LevelElevated {
		t.Fatalf("expected Elevated, got %d", got)
	}
}

func TestThirdOccurrenceWithinWindowIsCritical(t *testing.T) {
	e := escalate.New(time.Minute)
	c := makeChange(80, "opened")
	e.Assess(c)
	e.Assess(c)
	if got := e.Assess(c); got != escalate.LevelCritical {
		t.Fatalf("expected Critical, got %d", got)
	}
}

func TestOccurrenceAfterWindowResetsToNormal(t *testing.T) {
	now := time.Now()
	clock := func() time.Time { return now }

	e := escalate.New(0) // zero window — immediately expired
	_ = e              // rebuild with injected clock via internal constructor

	// Use the exported path: after window expires count resets.
	e2 := escalate.New(time.Millisecond)
	c := makeChange(443, "closed")
	e2.Assess(c)
	time.Sleep(5 * time.Millisecond)
	if got := e2.Assess(c); got != escalate.LevelNormal {
		t.Fatalf("expected Normal after window expiry, got %d (clock=%v)", got, clock())
	}
}

func TestDistinctDirectionsAreIndependent(t *testing.T) {
	e := escalate.New(time.Minute)
	opened := makeChange(22, "opened")
	closed := makeChange(22, "closed")
	e.Assess(opened)
	e.Assess(opened) // elevated for opened
	if got := e.Assess(closed); got != escalate.LevelNormal {
		t.Fatalf("closed direction should start at Normal, got %d", got)
	}
}

func TestResetClearsCount(t *testing.T) {
	e := escalate.New(time.Minute)
	c := makeChange(8080, "opened")
	e.Assess(c)
	e.Assess(c) // elevated
	e.Reset(c)
	if got := e.Assess(c); got != escalate.LevelNormal {
		t.Fatalf("expected Normal after Reset, got %d", got)
	}
}
