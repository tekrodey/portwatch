package watchdog_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/watchdog"
)

func TestDefaultConfigValues(t *testing.T) {
	cfg := watchdog.DefaultConfig()
	if cfg.Interval != 30*time.Second {
		t.Errorf("expected 30s interval, got %v", cfg.Interval)
	}
	if cfg.Timeout != 10*time.Second {
		t.Errorf("expected 10s timeout, got %v", cfg.Timeout)
	}
}

func TestCheckCalledOnTick(t *testing.T) {
	var calls atomic.Int32
	cfg := watchdog.Config{
		Interval: 20 * time.Millisecond,
		Timeout:  100 * time.Millisecond,
		OnHang:   func(err error) {},
	}
	w := watchdog.New(cfg, func(ctx context.Context) error {
		calls.Add(1)
		return nil
	})

	ctx, cancel := context.WithTimeout(context.Background(), 75*time.Millisecond)
	defer cancel()
	w.Run(ctx)

	if calls.Load() < 2 {
		t.Errorf("expected at least 2 check calls, got %d", calls.Load())
	}
}

func TestOnHangCalledWhenCheckErrors(t *testing.T) {
	var hung atomic.Int32
	checkErr := errors.New("check failed")
	cfg := watchdog.Config{
		Interval: 20 * time.Millisecond,
		Timeout:  100 * time.Millisecond,
		OnHang:   func(err error) { hung.Add(1) },
	}
	w := watchdog.New(cfg, func(ctx context.Context) error {
		return checkErr
	})

	ctx, cancel := context.WithTimeout(context.Background(), 55*time.Millisecond)
	defer cancel()
	w.Run(ctx)

	if hung.Load() == 0 {
		t.Error("expected onHang to be called")
	}
}

func TestOnHangCalledOnTimeout(t *testing.T) {
	var hung atomic.Int32
	cfg := watchdog.Config{
		Interval: 20 * time.Millisecond,
		Timeout:  15 * time.Millisecond,
		OnHang:   func(err error) { hung.Add(1) },
	}
	w := watchdog.New(cfg, func(ctx context.Context) error {
		time.Sleep(50 * time.Millisecond) // always exceeds timeout
		return nil
	})

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Millisecond)
	defer cancel()
	w.Run(ctx)

	if hung.Load() == 0 {
		t.Error("expected onHang to be called on timeout")
	}
}

func TestNilOnHangDoesNotPanic(t *testing.T) {
	cfg := watchdog.Config{
		Interval: 20 * time.Millisecond,
		Timeout:  100 * time.Millisecond,
		OnHang:   nil,
	}
	w := watchdog.New(cfg, func(ctx context.Context) error {
		return errors.New("oops")
	})

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Millisecond)
	defer cancel()

	// Should not panic even though OnHang is nil.
	w.Run(ctx)
}
