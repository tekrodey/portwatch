package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/scanner"
)

func makeChange(t monitor.ChangeType, port int) monitor.Change {
	return monitor.Change{
		Type: t,
		Port: scanner.Port{Number: port, Protocol: "tcp"},
	}
}

func TestNotifyOpened(t *testing.T) {
	var buf bytes.Buffer
	a := alert.New(&buf)
	a.Notify([]monitor.Change{makeChange(monitor.ChangeOpened, 8080)})

	out := buf.String()
	if !strings.Contains(out, "[ALERT]") {
		t.Errorf("expected ALERT level for opened port, got: %s", out)
	}
	if !strings.Contains(out, "8080") {
		t.Errorf("expected port 8080 in output, got: %s", out)
	}
}

func TestNotifyClosed(t *testing.T) {
	var buf bytes.Buffer
	a := alert.New(&buf)
	a.Notify([]monitor.Change{makeChange(monitor.ChangeClosed, 9090)})

	out := buf.String()
	if !strings.Contains(out, "[WARN]") {
		t.Errorf("expected WARN level for closed port, got: %s", out)
	}
}

func TestNotifyEmpty(t *testing.T) {
	var buf bytes.Buffer
	a := alert.New(&buf)
	a.Notify([]monitor.Change{})

	if buf.Len() != 0 {
		t.Errorf("expected no output for empty changes, got: %s", buf.String())
	}
}

func TestNewDefaultsToStdout(t *testing.T) {
	// Just ensure New(nil) does not panic.
	a := alert.New(nil)
	if a == nil {
		t.Error("expected non-nil Alerter")
	}
}
