package history

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Entry records a port change event at a point in time.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Port      int       `json:"port"`
	Proto     string    `json:"proto"`
	Action    string    `json:"action"` // "opened" | "closed"
}

// History stores a bounded, persistent log of port change events.
type History struct {
	mu      sync.Mutex
	entries []Entry
	maxSize int
	path    string
}

// New creates a History that persists to path and keeps at most maxSize entries.
func New(path string, maxSize int) (*History, error) {
	h := &History{path: path, maxSize: maxSize}
	if err := h.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return h, nil
}

// Add appends a new entry and persists the log.
func (h *History) Add(e Entry) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.entries = append(h.entries, e)
	if len(h.entries) > h.maxSize {
		h.entries = h.entries[len(h.entries)-h.maxSize:]
	}
	return h.save()
}

// All returns a copy of all stored entries.
func (h *History) All() []Entry {
	h.mu.Lock()
	defer h.mu.Unlock()
	out := make([]Entry, len(h.entries))
	copy(out, h.entries)
	return out
}

func (h *History) load() error {
	data, err := os.ReadFile(h.path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &h.entries)
}

func (h *History) save() error {
	data, err := json.MarshalIndent(h.entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(h.path, data, 0o644)
}
