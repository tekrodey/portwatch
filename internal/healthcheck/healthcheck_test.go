package healthcheck

import (
	"testing"
	"time"
)

var fixedTime = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func newTestChecker() *Checker {
	return newWithClock(func() time.Time { return fixedTime })
}

func TestRegisterAddsUnhealthyStatus(t *testing.T) {
	c := newTestChecker()
	c.Register("scanner")
	all := c.All()
	if len(all) != 1 {
		t.Fatalf("expected 1 status, got %d", len(all))
	}
	if all[0].Healthy {
		t.Error("expected initial status to be unhealthy")
	}
	if all[0].Name != "scanner" {
		t.Errorf("expected name 'scanner', got %q", all[0].Name)
	}
}

func TestSetHealthyMarksComponent(t *testing.T) {
	c := newTestChecker()
	c.Register("monitor")
	c.SetHealthy("monitor", "running")
	all := c.All()
	if !all[0].Healthy {
		t.Error("expected component to be healthy")
	}
	if all[0].Message != "running" {
		t.Errorf("expected message 'running', got %q", all[0].Message)
	}
}

func TestSetUnhealthyMarksComponent(t *testing.T) {
	c := newTestChecker()
	c.Register("monitor")
	c.SetHealthy("monitor", "ok")
	c.SetUnhealthy("monitor", "scan failed")
	all := c.All()
	if all[0].Healthy {
		t.Error("expected component to be unhealthy")
	}
	if all[0].Message != "scan failed" {
		t.Errorf("expected message 'scan failed', got %q", all[0].Message)
	}
}

func TestHealthyReturnsTrueWhenAllHealthy(t *testing.T) {
	c := newTestChecker()
	c.Register("a")
	c.Register("b")
	c.SetHealthy("a", "")
	c.SetHealthy("b", "")
	if !c.Healthy() {
		t.Error("expected overall healthy")
	}
}

func TestHealthyReturnsFalseWhenAnyUnhealthy(t *testing.T) {
	c := newTestChecker()
	c.Register("a")
	c.Register("b")
	c.SetHealthy("a", "")
	// b remains unhealthy
	if c.Healthy() {
		t.Error("expected overall unhealthy")
	}
}

func TestLastCheckTimestampIsSet(t *testing.T) {
	c := newTestChecker()
	c.Register("scanner")
	c.SetHealthy("scanner", "ok")
	all := c.All()
	if !all[0].LastCheck.Equal(fixedTime) {
		t.Errorf("expected LastCheck %v, got %v", fixedTime, all[0].LastCheck)
	}
}
