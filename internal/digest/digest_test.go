package digest

import (
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/scanner"
)

func makeChange(proto string, port int, dir monitor.Direction) monitor.Change {
	return monitor.Change{
		Port:      scanner.Port{Proto: proto, Number: port},
		Direction: dir,
	}
}

func TestFlushNoOpBeforeInterval(t *testing.T) {
	var buf strings.Builder
	d := New(&buf, 10*time.Minute)
	d.Add([]monitor.Change{makeChange("tcp", 80, monitor.Opened)})
	if d.Flush() {
		t.Fatal("expected no flush before interval")
	}
	if buf.Len() != 0 {
		t.Fatalf("expected no output, got %q", buf.String())
	}
}

func TestFlushWritesAfterInterval(t *testing.T) {
	var buf strings.Builder
	d := New(&buf, time.Millisecond)
	d.Add([]monitor.Change{
		makeChange("tcp", 443, monitor.Opened),
		makeChange("tcp", 22, monitor.Closed),
	})
	time.Sleep(5 * time.Millisecond)
	if !d.Flush() {
		t.Fatal("expected flush after interval")
	}
	out := buf.String()
	if !strings.Contains(out, "Opened: 1") {
		t.Errorf("expected opened count, got: %s", out)
	}
	if !strings.Contains(out, "Closed: 1") {
		t.Errorf("expected closed count, got: %s", out)
	}
}

func TestFlushResetsAccumulator(t *testing.T) {
	var buf strings.Builder
	d := New(&buf, time.Millisecond)
	d.Add([]monitor.Change{makeChange("udp", 53, monitor.Opened)})
	time.Sleep(5 * time.Millisecond)
	d.Flush()
	buf.Reset()
	time.Sleep(5 * time.Millisecond)
	d.Flush()
	if strings.Contains(buf.String(), "udp") {
		t.Error("stale change appeared in second flush")
	}
}

func TestSummaryString(t *testing.T) {
	now := time.Now()
	s := Summary{
		From:   now,
		To:     now.Add(time.Minute),
		Opened: []monitor.Change{makeChange("tcp", 8080, monitor.Opened)},
	}
	out := s.String()
	if !strings.Contains(out, "Port Digest") {
		t.Errorf("missing header in: %s", out)
	}
	if !strings.Contains(out, "+") {
		t.Errorf("missing opened marker in: %s", out)
	}
}

func TestNewDefaultsToStdout(t *testing.T) {
	d := New(nil, time.Minute)
	if d.w == nil {
		t.Error("expected non-nil writer")
	}
}
