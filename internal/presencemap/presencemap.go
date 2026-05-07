// Package presencemap tracks which ports have been seen across scan cycles
// and provides a simple API to query first-seen and last-seen timestamps.
package presencemap

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Entry holds timing metadata for a single port key.
type Entry struct {
	FirstSeen time.Time `json:"first_seen"`
	LastSeen  time.Time `json:"last_seen"`
	Count     int       `json:"count"`
}

// Map is a thread-safe store of port presence entries.
type Map struct {
	mu      sync.RWMutex
	entries map[string]*Entry
	now     func() time.Time
}

// New returns an empty Map using the real clock.
func New() *Map {
	return newWithClock(time.Now)
}

func newWithClock(now func() time.Time) *Map {
	return &Map{
		entries: make(map[string]*Entry),
		now:     now,
	}
}

// Touch records a sighting of key, updating first/last-seen and incrementing count.
func (m *Map) Touch(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	t := m.now()
	if e, ok := m.entries[key]; ok {
		e.LastSeen = t
		e.Count++
		return
	}
	m.entries[key] = &Entry{FirstSeen: t, LastSeen: t, Count: 1}
}

// Get returns the Entry for key and whether it was found.
func (m *Map) Get(key string) (Entry, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	e, ok := m.entries[key]
	if !ok {
		return Entry{}, false
	}
	return *e, true
}

// Delete removes key from the map.
func (m *Map) Delete(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.entries, key)
}

// Save serialises the map to path as JSON.
func (m *Map) Save(path string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	data, err := json.MarshalIndent(m.entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

// Load reads a JSON file written by Save and populates the map.
func Load(path string, now func() time.Time) (*Map, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return newWithClock(now), nil
		}
		return nil, err
	}
	m := newWithClock(now)
	if err := json.Unmarshal(data, &m.entries); err != nil {
		return nil, err
	}
	return m, nil
}
