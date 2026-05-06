package reporter_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/reporter"
)

var fixedTime = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func makePorts() []monitor.Port {
	return []monitor.Port{
		{Proto: "tcp", Number: 80},
		{Proto: "tcp", Number: 443},
	}
}

func TestWriteTextContainsPortCount(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf, reporter.FormatText)
	if err := r.Write(makePorts(), fixedTime); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "open ports (2)") {
		t.Errorf("expected port count in output, got: %s", buf.String())
	}
}

func TestWriteTextListsPorts(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf, reporter.FormatText)
	_ = r.Write(makePorts(), fixedTime)
	out := buf.String()
	if !strings.Contains(out, "tcp/80") || !strings.Contains(out, "tcp/443") {
		t.Errorf("expected port lines in output, got: %s", out)
	}
}

func TestWriteJSONIsValidStructure(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf, reporter.FormatJSON)
	_ = r.Write(makePorts(), fixedTime)
	out := buf.String()
	if !strings.HasPrefix(out, "{") || !strings.HasSuffix(strings.TrimSpace(out), "}") {
		t.Errorf("expected JSON object, got: %s", out)
	}
	if !strings.Contains(out, "\"count\":2") {
		t.Errorf("expected count field, got: %s", out)
	}
}

func TestWriteEmptyPorts(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf, reporter.FormatText)
	if err := r.Write(nil, fixedTime); err != nil {
		t.Fatalf("unexpected error on empty ports: %v", err)
	}
	if !strings.Contains(buf.String(), "open ports (0)") {
		t.Errorf("expected zero count, got: %s", buf.String())
	}
}

func TestNewDefaultsToStdout(t *testing.T) {
	r := reporter.New(nil, reporter.FormatText)
	if r == nil {
		t.Fatal("expected non-nil reporter")
	}
}

func TestDefaultFormatIsText(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf, "")
	_ = r.Write(makePorts(), fixedTime)
	if !strings.Contains(buf.String(), "open ports") {
		t.Errorf("expected text format output, got: %s", buf.String())
	}
}
