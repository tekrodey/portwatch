// Package labelmap provides a registry for attaching static key-value labels
// to port/protocol pairs, enabling richer context in alerts and audit logs.
package labelmap

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// Label is a single key-value annotation.
type Label struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Entry maps a port/protocol pair to a set of labels.
type Entry struct {
	Port     int     `json:"port"`
	Protocol string  `json:"protocol"`
	Labels   []Label `json:"labels"`
}

type entryKey struct {
	port     int
	protocol string
}

// Map holds the label registry.
type Map struct {
	mu      sync.RWMutex
	entries map[entryKey][]Label
}

// New returns an empty Map.
func New() *Map {
	return &Map{entries: make(map[entryKey][]Label)}
}

// Load reads a JSON file containing a list of Entry values and populates the Map.
func Load(path string) (*Map, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return New(), nil
		}
		return nil, fmt.Errorf("labelmap: open %s: %w", path, err)
	}
	defer f.Close()

	var entries []Entry
	if err := json.NewDecoder(f).Decode(&entries); err != nil {
		return nil, fmt.Errorf("labelmap: decode %s: %w", path, err)
	}

	m := New()
	for _, e := range entries {
		m.Set(e.Port, e.Protocol, e.Labels)
	}
	return m, nil
}

// Set registers labels for the given port/protocol pair, replacing any existing entry.
func (m *Map) Set(port int, protocol string, labels []Label) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries[entryKey{port, protocol}] = labels
}

// Get returns the labels for the given port/protocol pair and whether an entry exists.
func (m *Map) Get(port int, protocol string) ([]Label, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	labels, ok := m.entries[entryKey{port, protocol}]
	return labels, ok
}

// Len returns the number of registered entries.
func (m *Map) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.entries)
}
