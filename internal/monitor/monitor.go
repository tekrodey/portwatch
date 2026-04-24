package monitor

import (
	"fmt"

	"github.com/user/portwatch/internal/scanner"
)

// ChangeType describes whether a port was opened or closed.
type ChangeType int

const (
	ChangeOpened ChangeType = iota
	ChangeClosed
)

// Change represents a single port state transition.
type Change struct {
	Type ChangeType
	Port scanner.Port
}

// String returns a human-readable description of the change.
func (c Change) String() string {
	action := "opened"
	if c.Type == ChangeClosed {
		action = "closed"
	}
	return fmt.Sprintf("port %s %s", c.Port.String(), action)
}

// Monitor tracks port state between scans and reports changes.
type Monitor struct {
	scanner  *scanner.Scanner
	previous map[int]bool
}

// NewMonitor creates a Monitor using the provided Scanner.
func NewMonitor(s *scanner.Scanner) *Monitor {
	return &Monitor{
		scanner:  s,
		previous: nil,
	}
}

// Scan performs a port scan and returns any changes since the last call.
// On the first call no changes are reported; the baseline is established.
func (m *Monitor) Scan() ([]Change, error) {
	ports, err := m.scanner.Scan()
	if err != nil {
		return nil, fmt.Errorf("monitor scan: %w", err)
	}

	current := make(map[int]bool, len(ports))
	for _, p := range ports {
		current[p.Number] = true
	}

	// First run — establish baseline.
	if m.previous == nil {
		m.previous = current
		return nil, nil
	}

	var changes []Change

	// Detect newly opened ports.
	for _, p := range ports {
		if !m.previous[p.Number] {
			changes = append(changes, Change{Type: ChangeOpened, Port: p})
		}
	}

	// Detect closed ports.
	for num := range m.previous {
		if !current[num] {
			changes = append(changes, Change{
				Type: ChangeClosed,
				Port: scanner.Port{Number: num, Protocol: "tcp"},
			})
		}
	}

	m.previous = current
	return changes, nil
}
