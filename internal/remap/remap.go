// Package remap translates port numbers to human-readable aliases.
// Rules are loaded from a simple JSON map of "proto:port" -> "alias".
package remap

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/user/portwatch/internal/monitor"
)

// Remapper maps port/protocol pairs to alias strings.
type Remapper struct {
	mu    sync.RWMutex
	rules map[string]string
}

// New returns an empty Remapper.
func New() *Remapper {
	return &Remapper{rules: make(map[string]string)}
}

// Load reads alias rules from a JSON file.
// The file must be a flat object: {"tcp:80": "http", "tcp:443": "https"}.
func Load(path string) (*Remapper, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("remap: open %s: %w", path, err)
	}
	defer f.Close()

	var rules map[string]string
	if err := json.NewDecoder(f).Decode(&rules); err != nil {
		return nil, fmt.Errorf("remap: decode %s: %w", path, err)
	}
	return &Remapper{rules: rules}, nil
}

// Set adds or updates a single alias rule.
func (r *Remapper) Set(proto string, port int, alias string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rules[key(proto, port)] = alias
}

// Alias returns the alias for the given port/protocol, or the default
// string "proto:port" if no rule is configured.
func (r *Remapper) Alias(proto string, port int) string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if a, ok := r.rules[key(proto, port)]; ok {
		return a
	}
	return fmt.Sprintf("%s:%d", proto, port)
}

// Apply annotates each change in the slice with its alias via the Label field.
func (r *Remapper) Apply(changes []monitor.Change) []monitor.Change {
	out := make([]monitor.Change, len(changes))
	for i, c := range changes {
		c.Label = r.Alias(c.Port.Proto, c.Port.Number)
		out[i] = c
	}
	return out
}

func key(proto string, port int) string {
	return fmt.Sprintf("%s:%d", proto, port)
}
