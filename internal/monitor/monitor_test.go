package monitor_test

import (
	"net"
	"testing"
	"time"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/scanner"
)

func startListener(t *testing.T) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start listener: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	return port, func() { ln.Close() }
}

func newTestMonitor(t *testing.T, startPort, endPort int) *monitor.Monitor {
	t.Helper()
	s, err := scanner.NewScanner("127.0.0.1", startPort, endPort, 200*time.Millisecond)
	if err != nil {
		t.Fatalf("NewScanner: %v", err)
	}
	return monitor.NewMonitor(s, time.Second)
}

func TestScanNoChangesOnFirstCall(t *testing.T) {
	port, stop := startListener(t)
	defer stop()

	m := newTestMonitor(t, port, port)
	changes, err := m.Scan()
	if err != nil {
		t.Fatalf("Scan error: %v", err)
	}
	// First scan establishes baseline; no prior state means no changes.
	if len(changes) != 0 {
		t.Errorf("expected 0 changes on first scan, got %d", len(changes))
	}
}

func TestScanDetectsOpenedPort(t *testing.T) {
	port, stop := startListener(t)
	defer stop()

	m := newTestMonitor(t, port, port)
	// Baseline scan with port open.
	if _, err := m.Scan(); err != nil {
		t.Fatalf("baseline scan error: %v", err)
	}

	// Close the port so next scan sees it as closed.
	stop()

	changes, err := m.Scan()
	if err != nil {
		t.Fatalf("second scan error: %v", err)
	}
	if len(changes) != 1 || changes[0].Kind != "closed" || changes[0].Port != port {
		t.Errorf("expected 1 closed change for port %d, got %+v", port, changes)
	}
}

func TestChangeString(t *testing.T) {
	c := monitor.Change{
		Port: 8080,
		Kind: "opened",
		At:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
	}
	got := c.String()
	if got == "" {
		t.Error("Change.String() returned empty string")
	}
	for _, want := range []string{"8080", "opened"} {
		if !containsStr(got, want) {
			t.Errorf("Change.String() = %q, want it to contain %q", got, want)
		}
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		}())
}
