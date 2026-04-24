package notify_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/notify"
	"github.com/user/portwatch/internal/scanner"
)

func makeChange(proto string, port int, kind monitor.ChangeKind) monitor.Change {
	return monitor.Change{
		Port: scanner.Port{Proto: proto, Number: port},
		Kind: kind,
	}
}

func TestSendWritesFormattedLines(t *testing.T) {
	var buf bytes.Buffer
	n := notify.New(&buf, nil)

	changes := []monitor.Change{
		makeChange("tcp", 80, monitor.Opened),
		makeChange("tcp", 443, monitor.Opened),
	}

	if err := n.Send(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "tcp:80") {
		t.Errorf("expected tcp:80 in output, got: %s", out)
	}
	if !strings.Contains(out, "tcp:443") {
		t.Errorf("expected tcp:443 in output, got: %s", out)
	}
}

func TestSendNoOpOnEmpty(t *testing.T) {
	var buf bytes.Buffer
	n := notify.New(&buf, nil)

	if err := n.Send(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output for empty changes, got: %s", buf.String())
	}
}

func TestSendUsesCustomFormatter(t *testing.T) {
	var buf bytes.Buffer
	f := &stubFormatter{msg: "CUSTOM_OUTPUT"}
	n := notify.New(&buf, f)

	changes := []monitor.Change{makeChange("udp", 53, monitor.Closed)}
	if err := n.Send(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "CUSTOM_OUTPUT") {
		t.Errorf("expected custom formatter output, got: %s", buf.String())
	}
}

func TestNewDefaultsToStdout(t *testing.T) {
	// Ensure New does not panic with nil writer and nil formatter.
	n := notify.New(nil, nil)
	if n == nil {
		t.Fatal("expected non-nil Notifier")
	}
}

type stubFormatter struct{ msg string }

func (s *stubFormatter) Format(_ []monitor.Change) string { return s.msg }
