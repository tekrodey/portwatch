package backoff

import (
	"context"
	"errors"
	"testing"
	"time"
)

func noSleep(_ context.Context, _ time.Duration) error { return nil }

func newTestRunner(maxAttempts int) *Runner {
	cfg := Config{
		InitialInterval: 1 * time.Millisecond,
		MaxInterval:     10 * time.Millisecond,
		Multiplier:      2.0,
	}
	r := NewRunner(New(cfg), maxAttempts)
	r.sleep = noSleep
	return r
}

func TestRunSucceedsOnFirstAttempt(t *testing.T) {
	r := newTestRunner(3)
	err := r.Run(context.Background(), func() error { return nil })
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunRetriesAndEventuallySucceeds(t *testing.T) {
	r := newTestRunner(5)
	calls := 0
	err := r.Run(context.Background(), func() error {
		calls++
		if calls < 3 {
			return errors.New("not yet")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestRunExhaustsMaxAttempts(t *testing.T) {
	r := newTestRunner(3)
	sentinel := errors.New("always fails")
	calls := 0
	err := r.Run(context.Background(), func() error {
		calls++
		return sentinel
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected exactly 3 calls, got %d", calls)
	}
}

func TestRunRespectsContextCancellation(t *testing.T) {
	cfg := Config{
		InitialInterval: 10 * time.Millisecond,
		MaxInterval:     1 * time.Second,
		Multiplier:      2.0,
	}
	r := NewRunner(New(cfg), 0) // unlimited, real sleep
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately
	err := r.Run(ctx, func() error { return errors.New("fail") })
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}
