package snapshot

import (
	"encoding/json"
	"os"
	"time"
)

// Entry represents a single port state captured at a point in time.
type Entry struct {
	Port     int       `json:"port"`
	Protocol string    `json:"protocol"`
	State    string    `json:"state"`
	Captured time.Time `json:"captured"`
}

// Snapshot holds a complete set of observed port entries.
type Snapshot struct {
	Entries   []Entry   `json:"entries"`
	CreatedAt time.Time `json:"created_at"`
}

// New returns an empty Snapshot stamped with the current time.
func New() *Snapshot {
	return &Snapshot{CreatedAt: time.Now()}
}

// Add appends a port entry to the snapshot.
func (s *Snapshot) Add(port int, protocol, state string) {
	s.Entries = append(s.Entries, Entry{
		Port:     port,
		Protocol: protocol,
		State:    state,
		Captured: time.Now(),
	})
}

// Save serialises the snapshot to the given file path as JSON.
func (s *Snapshot) Save(path string) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// Load reads a snapshot from a JSON file. If the file does not exist
// a fresh empty snapshot is returned without an error.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return New(), nil
	}
	if err != nil {
		return nil, err
	}
	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	return &s, nil
}
