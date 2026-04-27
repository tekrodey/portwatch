package enricher

import (
	"errors"
	"testing"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/scanner"
)

func makeChange(addr, proto string, number int, dir monitor.Direction) monitor.Change {
	return monitor.Change{
		Port: scanner.Port{
			Addr:   addr,
			Proto:  proto,
			Number: number,
		},
		Direction: dir,
	}
}

func TestEnrichReturnsServiceName(t *testing.T) {
	e := newWithLookups(
		func(string) ([]string, error) { return nil, errors.New("no dns") },
		func(_, port string) (string, error) {
			if port == "80" {
				return "http", nil
			}
			return "", errors.New("unknown")
		},
	)

	c := makeChange("127.0.0.1", "tcp", 80, monitor.Opened)
	m := e.Enrich(c)

	if m.ServiceName != "http" {
		t.Errorf("expected service 'http', got %q", m.ServiceName)
	}
}

func TestEnrichReturnsReverseDNS(t *testing.T) {
	e := newWithLookups(
		func(addr string) ([]string, error) {
			if addr == "10.0.0.1" {
				return []string{"myhost.local."}, nil
			}
			return nil, errors.New("no ptr")
		},
		func(string, string) (string, error) { return "", errors.New("unknown") },
	)

	c := makeChange("10.0.0.1", "tcp", 443, monitor.Opened)
	m := e.Enrich(c)

	if m.ReverseDNS != "myhost.local" {
		t.Errorf("expected reverse DNS 'myhost.local', got %q", m.ReverseDNS)
	}
}

func TestEnrichGracefulOnLookupErrors(t *testing.T) {
	e := newWithLookups(
		func(string) ([]string, error) { return nil, errors.New("fail") },
		func(string, string) (string, error) { return "", errors.New("fail") },
	)

	c := makeChange("192.168.1.1", "udp", 53, monitor.Closed)
	m := e.Enrich(c)

	if m.ServiceName != "" || m.ReverseDNS != "" || m.Hostname != "" {
		t.Errorf("expected empty meta on all lookup failures, got %+v", m)
	}
}

func TestNewReturnsNonNil(t *testing.T) {
	e := New()
	if e == nil {
		t.Fatal("expected non-nil Enricher from New()")
	}
}
