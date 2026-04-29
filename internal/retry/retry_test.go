package retry

import (
	"context"
	"errors"
	"testing"
	"time"
)

var errTemp = errors.New("temporary error")

func noSleep(_ time.Duration) {}

func TestSuccessOnFirstAttempt(t *testing.T) {
	r := newWithSleep(DefaultConfig(), noSleep)
	calls := 0
	err := r.Do(context.Background(), func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestRetriesUpToMaxAttempts(t *testing.T) {
	cfg := Config{MaxAttempts: 3, BaseDelay: time.Millisecond, MaxDelay: time.Millisecond}
	r := newWithSleep(cfg, noSleep)
	calls := 0
	err := r.Do(context.Background(), func() error {
		calls++
		return errTemp
	})
	if !errors.Is(err, ErrMaxAttempts) {
		t.Fatalf("expected ErrMaxAttempts, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestSucceedsOnSecondAttempt(t *testing.T) {
	cfg := Config{MaxAttempts: 3, BaseDelay: time.Millisecond, MaxDelay: time.Millisecond}
	r := newWithSleep(cfg, noSleep)
	calls := 0
	err := r.Do(context.Background(), func() error {
		calls++
		if calls < 2 {
			return errTemp
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if calls != 2 {
		t.Fatalf("expected 2 calls, got %d", calls)
	}
}

func TestCancelledContextStopsRetry(t *testing.T) {
	cfg := Config{MaxAttempts: 5, BaseDelay: time.Millisecond, MaxDelay: time.Millisecond}
	r := newWithSleep(cfg, noSleep)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := r.Do(ctx, func() error { return errTemp })
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestDelayCapAtMaxDelay(t *testing.T) {
	cfg := Config{MaxAttempts: 4, BaseDelay: 100 * time.Millisecond, MaxDelay: 150 * time.Millisecond}
	var delays []time.Duration
	r := newWithSleep(cfg, func(d time.Duration) { delays = append(delays, d) })
	_ = r.Do(context.Background(), func() error { return errTemp })
	for _, d := range delays {
		if d > cfg.MaxDelay {
			t.Fatalf("delay %v exceeded MaxDelay %v", d, cfg.MaxDelay)
		}
	}
}

func TestDefaultConfigValues(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.MaxAttempts != 3 {
		t.Errorf("expected MaxAttempts 3, got %d", cfg.MaxAttempts)
	}
	if cfg.BaseDelay != 200*time.Millisecond {
		t.Errorf("unexpected BaseDelay: %v", cfg.BaseDelay)
	}
	if cfg.MaxDelay != 5*time.Second {
		t.Errorf("unexpected MaxDelay: %v", cfg.MaxDelay)
	}
}
