// Package pipeline wires together the processing stages that every
// detected change passes through before reaching an output sink.
package pipeline

import (
	"context"

	"github.com/user/portwatch/internal/dedup"
	"github.com/user/portwatch/internal/debounce"
	"github.com/user/portwatch/internal/enricher"
	"github.com/user/portwatch/internal/filter"
	"github.com/user/portwatch/internal/monitor"
)

// Stage is a function that accepts a slice of changes and returns a
// (possibly filtered or transformed) slice.
type Stage func([]monitor.Change) []monitor.Change

// Pipeline applies an ordered sequence of Stages to each batch of changes.
type Pipeline struct {
	stages []Stage
}

// New constructs a Pipeline from the supplied stages. Stages are applied
// in the order they are provided.
func New(stages ...Stage) *Pipeline {
	return &Pipeline{stages: stages}
}

// Run passes changes through every stage in order and returns the result.
func (p *Pipeline) Run(changes []monitor.Change) []monitor.Change {
	for _, s := range p.stages {
		if len(changes) == 0 {
			break
		}
		changes = s(changes)
	}
	return changes
}

// DefaultStages builds the standard processing pipeline used by the daemon.
// ctx is forwarded to any stage that requires cancellation awareness.
func DefaultStages(
	_ context.Context,
	f *filter.Filter,
	dd *dedup.Dedup,
	db *debounce.Filter,
	en *enricher.Enricher,
) []Stage {
	return []Stage{
		FilterStage(f),
		DedupStage(dd),
		DebounceStage(db),
		EnrichStage(en),
	}
}

// FilterStage wraps a filter.Filter as a Stage.
func FilterStage(f *filter.Filter) Stage {
	return func(changes []monitor.Change) []monitor.Change {
		return f.Apply(changes)
	}
}

// DedupStage wraps a dedup.Dedup as a Stage.
func DedupStage(d *dedup.Dedup) Stage {
	return func(changes []monitor.Change) []monitor.Change {
		return d.Apply(changes)
	}
}

// DebounceStage wraps a debounce.Filter as a Stage.
func DebounceStage(db *debounce.Filter) Stage {
	return func(changes []monitor.Change) []monitor.Change {
		return db.Apply(changes)
	}
}

// EnrichStage wraps an enricher.Enricher as a Stage.
func EnrichStage(en *enricher.Enricher) Stage {
	return func(changes []monitor.Change) []monitor.Change {
		return en.Enrich(changes)
	}
}
