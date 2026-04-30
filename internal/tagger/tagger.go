// Package tagger assigns human-readable labels to port changes based on
// configurable rules. Tags can be used downstream by alerting and audit
// pipelines to categorise events (e.g. "known-service", "ephemeral",
// "suspicious").
package tagger

import (
	"sync"

	"github.com/user/portwatch/internal/monitor"
)

// Rule maps a port+protocol pair to a tag label.
type Rule struct {
	Port     int
	Protocol string // "tcp" or "udp"
	Tag      string
}

// Tagger attaches tags to monitor.Change values.
type Tagger struct {
	mu    sync.RWMutex
	rules []Rule
}

// New returns a Tagger pre-loaded with the supplied rules.
func New(rules []Rule) *Tagger {
	r := make([]Rule, len(rules))
	copy(r, rules)
	return &Tagger{rules: r}
}

// AddRule appends a new rule at runtime.
func (t *Tagger) AddRule(r Rule) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.rules = append(t.rules, r)
}

// Tag returns all tags that match the given change. The returned slice is
// nil when no rule matches.
func (t *Tagger) Tag(c monitor.Change) []string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	var tags []string
	for _, r := range t.rules {
		if r.Port == c.Port.Port && r.Protocol == c.Port.Protocol {
			tags = append(tags, r.Tag)
		}
	}
	return tags
}

// TagAll annotates a slice of changes, returning a parallel slice of tag
// sets. Indices correspond 1-to-1 with the input slice.
func (t *Tagger) TagAll(changes []monitor.Change) [][]string {
	out := make([][]string, len(changes))
	for i, c := range changes {
		out[i] = t.Tag(c)
	}
	return out
}
