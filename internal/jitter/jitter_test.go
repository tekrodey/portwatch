package jitter

import (
	"testing"
	"time"
)

func TestApplyZeroFactorReturnsBase(t *testing.T) {
	j := New(0)
	got := j.Apply(time.Second)
	if got != time.Second {
		t.Fatalf("expected 1s, got %v", got)
	}
}

func TestApplyNegativeBasePassesThrough(t *testing.T) {
	j := New(0.5)
	got := j.Apply(-time.Second)
	if got != -time.Second {
		t.Fatalf("expected -1s, got %v", got)
	}
}

func TestApplyAlwaysGteBase(t *testing.T) {
	j := New(0.5)
	base := 100 * time.Millisecond
	for i := 0; i < 200; i++ {
		got := j.Apply(base)
		if got < base {
			t.Fatalf("result %v is less than base %v", got, base)
		}
	}
}

func TestApplyNeverExceedsBaseTimesFactor(t *testing.T) {
	j := New(0.5)
	base := 100 * time.Millisecond
	max := base + time.Duration(float64(base)*0.5)
	for i := 0; i < 200; i++ {
		got := j.Apply(base)
		if got > max {
			t.Fatalf("result %v exceeds max %v", got, max)
		}
	}
}

func TestFactorClampedAboveOne(t *testing.T) {
	j := New(5.0)
	if j.factor != 1.0 {
		t.Fatalf("expected factor clamped to 1.0, got %v", j.factor)
	}
}

func TestFactorClampedBelowZero(t *testing.T) {
	j := New(-1.0)
	if j.factor != 0.0 {
		t.Fatalf("expected factor clamped to 0.0, got %v", j.factor)
	}
}

func TestDeterministicSourceProducesExpectedResult(t *testing.T) {
	// source always returns 0.5
	j := newWithSource(0.4, func() float64 { return 0.5 })
	base := time.Second
	got := j.Apply(base)
	want := base + time.Duration(float64(base)*0.4*0.5)
	if got != want {
		t.Fatalf("expected %v, got %v", want, got)
	}
}
