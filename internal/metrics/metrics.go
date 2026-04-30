// Package metrics provides lightweight in-process counters and gauges
// for tracking portwatch runtime statistics (scans run, alerts sent, etc.).
package metrics

import (
	"fmt"
	"io"
	"os"
	"sort"
	"sync"
	"sync/atomic"
)

// Counter is a monotonically increasing uint64 counter.
type Counter struct {
	value uint64
}

// Inc increments the counter by 1.
func (c *Counter) Inc() { atomic.AddUint64(&c.value, 1) }

// Add increments the counter by n.
func (c *Counter) Add(n uint64) { atomic.AddUint64(&c.value, n) }

// Value returns the current counter value.
func (c *Counter) Value() uint64 { return atomic.LoadUint64(&c.value) }

// Gauge is a signed integer that can go up or down.
type Gauge struct {
	mu    sync.Mutex
	value int64
}

// Set sets the gauge to an absolute value.
func (g *Gauge) Set(v int64) {
	g.mu.Lock()
	g.value = v
	g.mu.Unlock()
}

// Inc increments the gauge by 1.
func (g *Gauge) Inc() {
	g.mu.Lock()
	g.value++
	g.mu.Unlock()
}

// Dec decrements the gauge by 1.
func (g *Gauge) Dec() {
	g.mu.Lock()
	g.value--
	g.mu.Unlock()
}

// Value returns the current gauge value.
func (g *Gauge) Value() int64 {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.value
}

// Registry holds named counters and gauges.
type Registry struct {
	mu       sync.Mutex
	counters map[string]*Counter
	gauges   map[string]*Gauge
}

// New returns an empty Registry.
func New() *Registry {
	return &Registry{
		counters: make(map[string]*Counter),
		gauges:   make(map[string]*Gauge),
	}
}

// Counter returns (creating if necessary) the named counter.
func (r *Registry) Counter(name string) *Counter {
	r.mu.Lock()
	defer r.mu.Unlock()
	if c, ok := r.counters[name]; ok {
		return c
	}
	c := &Counter{}
	r.counters[name] = c
	return c
}

// Gauge returns (creating if necessary) the named gauge.
func (r *Registry) Gauge(name string) *Gauge {
	r.mu.Lock()
	defer r.mu.Unlock()
	if g, ok := r.gauges[name]; ok {
		return g
	}
	g := &Gauge{}
	r.gauges[name] = g
	return g
}

// Print writes all metric values in a simple key=value format to w.
// If w is nil, os.Stdout is used.
func (r *Registry) Print(w io.Writer) {
	if w == nil {
		w = os.Stdout
	}
	r.mu.Lock()
	names := make([]string, 0, len(r.counters)+len(r.gauges))
	for n := range r.counters {
		names = append(names, "counter:"+n)
	}
	for n := range r.gauges {
		names = append(names, "gauge:"+n)
	}
	r.mu.Unlock()

	sort.Strings(names)
	for _, key := range names {
		switch key[:8] {
		case "counter:":
			n := key[8:]
			fmt.Fprintf(w, "counter_%s=%d\n", n, r.Counter(n).Value())
		case "gauge:  ", "gauge:":
			n := key[6:]
			fmt.Fprintf(w, "gauge_%s=%d\n", n, r.Gauge(n).Value())
		}
	}
}
