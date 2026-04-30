package healthcheck

import (
	"sync"
	"time"
)

// Status represents the current health of a named component.
type Status struct {
	Name      string
	Healthy   bool
	LastCheck time.Time
	Message   string
}

// Checker tracks the health of registered components.
type Checker struct {
	mu       sync.RWMutex
	statuses map[string]Status
	clock    func() time.Time
}

// New returns a new Checker using the real wall clock.
func New() *Checker {
	return newWithClock(time.Now)
}

func newWithClock(clock func() time.Time) *Checker {
	return &Checker{
		statuses: make(map[string]Status),
		clock:    clock,
	}
}

// Register adds a named component with an initial unhealthy state.
func (c *Checker) Register(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.statuses[name] = Status{Name: name, Healthy: false, LastCheck: c.clock()}
}

// SetHealthy marks a component as healthy with an optional message.
func (c *Checker) SetHealthy(name, message string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.statuses[name] = Status{
		Name:      name,
		Healthy:   true,
		LastCheck: c.clock(),
		Message:   message,
	}
}

// SetUnhealthy marks a component as unhealthy with a reason message.
func (c *Checker) SetUnhealthy(name, message string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.statuses[name] = Status{
		Name:      name,
		Healthy:   false,
		LastCheck: c.clock(),
		Message:   message,
	}
}

// All returns a snapshot of all registered component statuses.
func (c *Checker) All() []Status {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make([]Status, 0, len(c.statuses))
	for _, s := range c.statuses {
		out = append(out, s)
	}
	return out
}

// Healthy returns true if all registered components are healthy.
func (c *Checker) Healthy() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, s := range c.statuses {
		if !s.Healthy {
			return false
		}
	}
	return true
}
