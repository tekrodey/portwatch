package monitor

import (
	"fmt"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// PortState holds the last known state of scanned ports.
type PortState struct {
	OpenPorts map[int]bool
	LastSeen  time.Time
}

// Change represents a detected port state change.
type Change struct {
	Port   int
	Kind   string // "opened" or "closed"
	At     time.Time
}

func (c Change) String() string {
	return fmt.Sprintf("[%s] port %d was %s", c.At.Format(time.RFC3339), c.Port, c.Kind)
}

// Monitor watches for port changes over time.
type Monitor struct {
	scanner  *scanner.Scanner
	previous PortState
	Interval time.Duration
}

// NewMonitor creates a Monitor using the provided Scanner.
func NewMonitor(s *scanner.Scanner, interval time.Duration) *Monitor {
	return &Monitor{
		scanner:  s,
		Interval: interval,
	}
}

// Scan performs a single scan and returns any changes since the last scan.
func (m *Monitor) Scan() ([]Change, error) {
	ports, err := m.scanner.Scan()
	if err != nil {
		return nil, fmt.Errorf("monitor scan: %w", err)
	}

	current := make(map[int]bool, len(ports))
	for _, p := range ports {
		current[p.Port] = true
	}

	now := time.Now()
	var changes []Change

	// Detect newly opened ports.
	for port := range current {
		if !m.previous.OpenPorts[port] {
			changes = append(changes, Change{Port: port, Kind: "opened", At: now})
		}
	}

	// Detect closed ports.
	for port := range m.previous.OpenPorts {
		if !current[port] {
			changes = append(changes, Change{Port: port, Kind: "closed", At: now})
		}
	}

	m.previous = PortState{OpenPorts: current, LastSeen: now}
	return changes, nil
}

// Run starts the monitoring loop, sending changes to the returned channel.
// Close the done channel to stop the loop.
func (m *Monitor) Run(done <-chan struct{}) <-chan Change {
	ch := make(chan Change, 16)
	go func() {
		defer close(ch)
		ticker := time.NewTicker(m.Interval)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				changes, err := m.Scan()
				if err != nil {
					continue
				}
				for _, c := range changes {
					ch <- c
				}
			}
		}
	}()
	return ch
}
