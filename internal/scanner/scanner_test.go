package scanner

import (
	"fmt"
	"net"
	"testing"
	"time"
)

// startTestListener opens a TCP listener on a random port and returns it.
func startTestListener(t *testing.T) (net.Listener, int) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start test listener: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	return ln, port
}

func TestPortString(t *testing.T) {
	p := Port{Protocol: "tcp", Address: "127.0.0.1", Port: 8080}
	want := "tcp://127.0.0.1:8080"
	if got := p.String(); got != want {
		t.Errorf("Port.String() = %q, want %q", got, want)
	}
}

func TestNewScannerDefaults(t *testing.T) {
	s := NewScanner("127.0.0.1", 1, 1024)
	if s.Timeout != 500*time.Millisecond {
		t.Errorf("expected default timeout 500ms, got %v", s.Timeout)
	}
}

func TestScanDetectsOpenPort(t *testing.T) {
	ln, port := startTestListener(t)
	defer ln.Close()

	s := NewScanner("127.0.0.1", port, port)
	s.Timeout = 200 * time.Millisecond

	ports, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan() error: %v", err)
	}
	if len(ports) != 1 {
		t.Fatalf("expected 1 open port, got %d", len(ports))
	}
	if ports[0].Port != port {
		t.Errorf("expected port %d, got %d", port, ports[0].Port)
	}
}

func TestScanInvalidRange(t *testing.T) {
	cases := []struct{ min, max int }{
		{0, 100},
		{100, 70000},
		{500, 100},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf("%d-%d", tc.min, tc.max), func(t *testing.T) {
			s := NewScanner("127.0.0.1", tc.min, tc.max)
			_, err := s.Scan()
			if err == nil {
				t.Error("expected error for invalid range, got nil")
			}
		})
	}
}
