// Package baseline manages the expected (trusted) set of open ports.
// It allows users to mark the current port state as "known good" so that
// future scans only alert on deviations from that baseline.
package baseline

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// Entry represents a single baselined port.
type Entry struct {
	Protocol string    `json:"protocol"`
	Port     int       `json:"port"`
	AddedAt  time.Time `json:"added_at"`
}

// Baseline holds the set of trusted protocol/port pairs.
type Baseline struct {
	mu      sync.RWMutex
	entries map[string]Entry
	path    string
}

// New loads a baseline from path, or returns an empty baseline if the file
// does not exist.
func New(path string) (*Baseline, error) {
	b := &Baseline{
		entries: make(map[string]Entry),
		path:    path,
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return b, nil
	}
	if err != nil {
		return nil, err
	}
	var list []Entry
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}
	for _, e := range list {
		b.entries[entryKey(e.Protocol, e.Port)] = e
	}
	return b, nil
}

// Contains reports whether the given protocol/port pair is in the baseline.
func (b *Baseline) Contains(protocol string, port int) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	_, ok := b.entries[entryKey(protocol, port)]
	return ok
}

// Set adds or updates an entry in the baseline and persists it.
func (b *Baseline) Set(protocol string, port int) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.entries[entryKey(protocol, port)] = Entry{
		Protocol: protocol,
		Port:     port,
		AddedAt:  time.Now().UTC(),
	}
	return b.save()
}

// Remove deletes an entry from the baseline and persists the change.
func (b *Baseline) Remove(protocol string, port int) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.entries, entryKey(protocol, port))
	return b.save()
}

// All returns a copy of all baselined entries.
func (b *Baseline) All() []Entry {
	b.mu.RLock()
	defer b.mu.RUnlock()
	out := make([]Entry, 0, len(b.entries))
	for _, e := range b.entries {
		out = append(out, e)
	}
	return out
}

func (b *Baseline) save() error {
	list := make([]Entry, 0, len(b.entries))
	for _, e := range b.entries {
		list = append(list, e)
	}
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(b.path, data, 0o644)
}

func entryKey(protocol string, port int) string {
	return fmt.Sprintf("%s:%d", protocol, port)
}
