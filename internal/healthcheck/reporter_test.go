package healthcheck

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestPrintTextOverallOK(t *testing.T) {
	c := newTestChecker()
	c.Register("scanner")
	c.SetHealthy("scanner", "running")
	var buf bytes.Buffer
	r := NewReporterWithWriter(c, &buf)
	r.PrintText()
	out := buf.String()
	if !strings.Contains(out, "overall: OK") {
		t.Errorf("expected 'overall: OK' in output, got:\n%s", out)
	}
	if !strings.Contains(out, "✓") {
		t.Errorf("expected checkmark in output, got:\n%s", out)
	}
}

func TestPrintTextOverallDegraded(t *testing.T) {
	c := newTestChecker()
	c.Register("scanner")
	// scanner stays unhealthy
	var buf bytes.Buffer
	r := NewReporterWithWriter(c, &buf)
	r.PrintText()
	out := buf.String()
	if !strings.Contains(out, "overall: DEGRADED") {
		t.Errorf("expected 'overall: DEGRADED' in output, got:\n%s", out)
	}
	if !strings.Contains(out, "✗") {
		t.Errorf("expected cross mark in output, got:\n%s", out)
	}
}

func TestPrintJSONStructure(t *testing.T) {
	c := newTestChecker()
	c.Register("monitor")
	c.SetHealthy("monitor", "ok")
	var buf bytes.Buffer
	r := NewReporterWithWriter(c, &buf)
	if err := r.PrintJSON(); err != nil {
		t.Fatalf("PrintJSON returned error: %v", err)
	}
	var out struct {
		Healthy  bool     `json:"healthy"`
		Statuses []Status `json:"statuses"`
	}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}
	if !out.Healthy {
		t.Error("expected healthy=true in JSON output")
	}
	if len(out.Statuses) != 1 {
		t.Errorf("expected 1 status, got %d", len(out.Statuses))
	}
}

func TestNewReporterDefaultsToStdout(t *testing.T) {
	c := New()
	r := NewReporter(c)
	if r.w == nil {
		t.Error("expected non-nil writer")
	}
}
