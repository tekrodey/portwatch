package pipeline_test

import (
	"testing"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/pipeline"
)

// identity is a Stage that returns changes unchanged.
func identity(changes []monitor.Change) []monitor.Change { return changes }

// drop is a Stage that discards all changes.
func drop(_ []monitor.Change) []monitor.Change { return nil }

// double is a Stage that duplicates the slice (useful for counting calls).
func double(changes []monitor.Change) []monitor.Change {
	out := make([]monitor.Change, 0, len(changes)*2)
	out = append(out, changes...)
	out = append(out, changes...)
	return out
}

func makeChange(port int, dir string) monitor.Change {
	return monitor.Change{Port: port, Proto: "tcp", Direction: dir}
}

func TestRunEmptyChanges(t *testing.T) {
	p := pipeline.New(identity)
	got := p.Run(nil)
	if len(got) != 0 {
		t.Fatalf("expected empty, got %v", got)
	}
}

func TestRunAppliesStagesInOrder(t *testing.T) {
	var order []int
	stageA := func(c []monitor.Change) []monitor.Change { order = append(order, 1); return c }
	stageB := func(c []monitor.Change) []monitor.Change { order = append(order, 2); return c }

	p := pipeline.New(stageA, stageB)
	p.Run([]monitor.Change{makeChange(80, "opened")})

	if len(order) != 2 || order[0] != 1 || order[1] != 2 {
		t.Fatalf("unexpected stage order: %v", order)
	}
}

func TestRunShortCircuitsOnEmpty(t *testing.T) {
	calls := 0
	counter := func(c []monitor.Change) []monitor.Change { calls++; return c }

	p := pipeline.New(drop, counter)
	p.Run([]monitor.Change{makeChange(443, "opened")})

	// counter should never be reached because drop empties the slice
	if calls != 0 {
		t.Fatalf("expected 0 calls after drop stage, got %d", calls)
	}
}

func TestRunTransformsChanges(t *testing.T) {
	p := pipeline.New(double)
	input := []monitor.Change{makeChange(22, "opened")}
	got := p.Run(input)
	if len(got) != 2 {
		t.Fatalf("expected 2 changes after double stage, got %d", len(got))
	}
}

func TestRunNoStages(t *testing.T) {
	p := pipeline.New()
	input := []monitor.Change{makeChange(8080, "closed")}
	got := p.Run(input)
	if len(got) != 1 {
		t.Fatalf("expected 1 change with no stages, got %d", len(got))
	}
}
