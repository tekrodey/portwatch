// Package enricher attaches process and service metadata to port changes.
package enricher

import (
	"fmt"
	"net"
	"strings"

	"github.com/user/portwatch/internal/monitor"
)

// Meta holds optional metadata associated with an open port.
type Meta struct {
	Hostname string
	ReverseDNS string
	ServiceName string
}

// Enricher annotates monitor.Change values with extra context.
type Enricher struct {
	lookupHost  func(addr string) ([]string, error)
	lookupPort  func(network, port string) (string, error)
}

// New returns an Enricher using the standard net lookup functions.
func New() *Enricher {
	return &Enricher{
		lookupHost: net.LookupAddr,
		lookupPort: net.LookupPort,
	}
}

// newWithLookups returns an Enricher with injectable lookup functions (for testing).
func newWithLookups(
	lookupHost func(string) ([]string, error),
	lookupPort func(string, string) (string, error),
) *Enricher {
	return &Enricher{
		lookupHost:  lookupHost,
		lookupPort:  lookupPort,
	}
}

// Enrich returns a Meta struct populated from available system lookups.
func (e *Enricher) Enrich(c monitor.Change) Meta {
	m := Meta{}

	hosts, err := e.lookupHost("127.0.0.1")
	if err == nil && len(hosts) > 0 {
		m.Hostname = strings.TrimSuffix(hosts[0], ".")
	}

	reverse, err := e.lookupHost(c.Port.Addr)
	if err == nil && len(reverse) > 0 {
		m.ReverseDNS = strings.TrimSuffix(reverse[0], ".")
	}

	proto := strings.ToLower(c.Port.Proto)
	portStr := fmt.Sprintf("%d", c.Port.Number)
	svc, err := e.lookupPort(proto, portStr)
	if err == nil {
		m.ServiceName = svc
	}

	return m
}
