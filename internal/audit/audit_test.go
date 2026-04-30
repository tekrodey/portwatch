package audit_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/audit"
	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/scanner"
)

func makeChange(dir, proto string, port int) monitor.Change {
	return monitor.Change{
		Direction: dir,
		Port:      scanner.Port{Protocol: proto, Number: port},
	}
}

func TestLogWritesJSONLines(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)

	changes := []monitor.Change{
		makeChange("opened", "tcp", 8080),
		makeChange("closed", "udp", 53),
	}

	if err := l.Log(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}

	var e audit.Entry
	if err := json.Unmarshal([]byte(lines[0]), &e); err != nil {
		t.Fatalf("unmarshal line 0: %v", err)
	}
	if e.Direction != "opened" || e.Protocol != "tcp" || e.Port != 8080 {
		t.Errorf("unexpected entry: %+v", e)
	}
}

func TestLogNoOpOnEmpty(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)

	if err := l.Log(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output, got %q", buf.String())
	}
}

func TestLogTimestampIsSet(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)

	_ = l.Log([]monitor.Change{makeChange("opened", "tcp", 443)})

	var e audit.Entry
	if err := json.Unmarshal(buf.Bytes(), &e); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if e.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestNewDefaultsToStdout(t *testing.T) {
	l := audit.New(nil)
	if l == nil {
		t.Fatal("expected non-nil logger")
	}
}
