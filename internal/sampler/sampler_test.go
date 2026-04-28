package sampler

import (
	"math/rand"
	"testing"

	"github.com/user/portwatch/internal/monitor"
)

func makeChange(port int) monitor.Change {
	return monitor.Change{Port: port, Proto: "tcp", Direction: "opened"}
}

func makeChanges(n int) []monitor.Change {
	changes := make([]monitor.Change, n)
	for i := range changes {
		changes[i] = makeChange(8000 + i)
	}
	return changes
}

func TestRateOneForwardsAll(t *testing.T) {
	s := New(1.0, rand.NewSource(1))
	input := makeChanges(20)
	out := s.Sample(input)
	if len(out) != len(input) {
		t.Fatalf("expected %d changes, got %d", len(input), len(out))
	}
}

func TestRateZeroDropsAll(t *testing.T) {
	s := New(0.0, rand.NewSource(1))
	input := makeChanges(20)
	out := s.Sample(input)
	if len(out) != 0 {
		t.Fatalf("expected 0 changes, got %d", len(out))
	}
}

func TestEmptyInputReturnsEmpty(t *testing.T) {
	s := New(1.0, rand.NewSource(1))
	out := s.Sample(nil)
	if out != nil {
		t.Fatalf("expected nil, got %v", out)
	}
}

func TestHalfRateApproximatesHalf(t *testing.T) {
	s := New(0.5, rand.NewSource(99))
	input := makeChanges(1000)
	out := s.Sample(input)
	if len(out) < 400 || len(out) > 600 {
		t.Fatalf("expected ~500 changes at 0.5 rate, got %d", len(out))
	}
}

func TestRateClampedAboveOne(t *testing.T) {
	s := New(5.0, nil)
	if s.Rate() != 1.0 {
		t.Fatalf("expected rate 1.0, got %f", s.Rate())
	}
}

func TestRateClampedBelowZero(t *testing.T) {
	s := New(-0.5, nil)
	if s.Rate() != 0.0 {
		t.Fatalf("expected rate 0.0, got %f", s.Rate())
	}
}

func TestNilSourceUsesDefault(t *testing.T) {
	s := New(1.0, nil)
	out := s.Sample(makeChanges(5))
	if len(out) != 5 {
		t.Fatalf("expected 5 changes, got %d", len(out))
	}
}
