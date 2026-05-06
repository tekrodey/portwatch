package reporter_test

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/reporter"
)

type fakeLister struct {
	ports []monitor.Port
}

func (f *fakeLister) OpenPorts() []monitor.Port { return f.ports }

func TestRunnerWritesOnTick(t *testing.T) {
	var buf bytes.Buffer
	lister := &fakeLister{ports: []monitor.Port{{Proto: "tcp", Number: 22}}}
	rep := reporter.New(&buf, reporter.FormatText)
	runner := reporter.NewRunner(lister, rep, 20*time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 70*time.Millisecond)
	defer cancel()

	_ = runner.Run(ctx)

	out := buf.String()
	if !strings.Contains(out, "tcp/22") {
		t.Errorf("expected port in report output, got: %s", out)
	}
}

func TestRunnerExitsOnContextCancel(t *testing.T) {
	var buf bytes.Buffer
	lister := &fakeLister{}
	rep := reporter.New(&buf, reporter.FormatText)
	runner := reporter.NewRunner(lister, rep, time.Hour)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := runner.Run(ctx)
	if err == nil {
		t.Fatal("expected non-nil error on cancelled context")
	}
}

func TestRunnerDefaultInterval(t *testing.T) {
	// NewRunner with zero interval should not panic and should use default.
	var buf bytes.Buffer
	lister := &fakeLister{}
	rep := reporter.New(&buf, reporter.FormatText)
	runner := reporter.NewRunner(lister, rep, 0)
	if runner == nil {
		t.Fatal("expected non-nil runner")
	}
}
