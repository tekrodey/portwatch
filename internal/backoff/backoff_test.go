package backoff

import (
	"testing"
	"time"
)

func zeroJitter(_ float64) float64 { return 0 }

func TestDefaultConfigValues(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.InitialInterval != 500*time.Millisecond {
		t.Fatalf("expected 500ms initial, got %v", cfg.InitialInterval)
	}
	if cfg.MaxInterval != 30*time.Second {
		t.Fatalf("expected 30s max, got %v", cfg.MaxInterval)
	}
	if cfg.Multiplier != 2.0 {
		t.Fatalf("expected multiplier 2.0, got %v", cfg.Multiplier)
	}
}

func TestIntervalGrowsExponentially(t *testing.T) {
	cfg := Config{
		InitialInterval: 100 * time.Millisecond,
		MaxInterval:     10 * time.Second,
		Multiplier:      2.0,
	}
	b := newWithRand(cfg, zeroJitter)

	expected := []time.Duration{
		100 * time.Millisecond,
		200 * time.Millisecond,
		400 * time.Millisecond,
		800 * time.Millisecond,
	}
	for i, want := range expected {
		got := b.Interval(i)
		if got != want {
			t.Errorf("attempt %d: want %v, got %v", i, want, got)
		}
	}
}

func TestIntervalCapsAtMax(t *testing.T) {
	cfg := Config{
		InitialInterval: 1 * time.Second,
		MaxInterval:     3 * time.Second,
		Multiplier:      2.0,
	}
	b := newWithRand(cfg, zeroJitter)

	for _, attempt := range []int{3, 5, 10} {
		got := b.Interval(attempt)
		if got > cfg.MaxInterval {
			t.Errorf("attempt %d: %v exceeds max %v", attempt, got, cfg.MaxInterval)
		}
	}
}

func TestIntervalNegativeAttemptTreatedAsZero(t *testing.T) {
	cfg := Config{
		InitialInterval: 200 * time.Millisecond,
		MaxInterval:     10 * time.Second,
		Multiplier:      2.0,
	}
	b := newWithRand(cfg, zeroJitter)
	if b.Interval(-1) != b.Interval(0) {
		t.Error("negative attempt should equal attempt 0")
	}
}

func TestIntervalMultiplierClampedWhenTooLow(t *testing.T) {
	cfg := Config{
		InitialInterval: 100 * time.Millisecond,
		MaxInterval:     10 * time.Second,
		Multiplier:      0.5, // invalid — should be clamped to 2
	}
	b := newWithRand(cfg, zeroJitter)
	if b.cfg.Multiplier != 2.0 {
		t.Errorf("expected multiplier clamped to 2.0, got %v", b.cfg.Multiplier)
	}
}

func TestJitterIncreasesInterval(t *testing.T) {
	cfg := Config{
		InitialInterval: 1 * time.Second,
		MaxInterval:     30 * time.Second,
		Multiplier:      2.0,
		JitterFactor:    0.2,
	}
	// jitter always returns full noise amount
	b := newWithRand(cfg, func(max float64) float64 { return max })
	base := b.Interval(0)
	if base <= 1*time.Second {
		t.Errorf("expected jitter to increase interval, got %v", base)
	}
}
