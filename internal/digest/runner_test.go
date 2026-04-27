package digest

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

func TestRunnerFlushesOnTick(t *testing.T) {
	var buf strings.Builder
	d := New(&buf, time.Millisecond)
	ch := make(chan []monitor.Change, 4)
	r := NewRunner(d, ch, 5*time.Millisecond)

	ch <- []monitor.Change{makeChange("tcp", 80, monitor.Opened)}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	done := make(chan struct{})
	go func() { r.Run(ctx); close(done) }()
	<-done

	if !strings.Contains(buf.String(), "Port Digest") {
		t.Errorf("expected digest output, got: %q", buf.String())
	}
}

func TestRunnerFlushesOnCancel(t *testing.T) {
	var buf strings.Builder
	d := New(&buf, time.Millisecond)
	ch := make(chan []monitor.Change, 4)
	r := NewRunner(d, ch, time.Hour) // long tick — flush must come from cancel

	ch <- []monitor.Change{makeChange("tcp", 443, monitor.Opened)}
	time.Sleep(5 * time.Millisecond) // let goroutine consume

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() { r.Run(ctx); close(done) }()

	time.Sleep(10 * time.Millisecond)
	cancel()
	<-done

	if !strings.Contains(buf.String(), "Port Digest") {
		t.Errorf("expected digest output on cancel, got: %q", buf.String())
	}
}

func TestRunnerExitsOnClosedChannel(t *testing.T) {
	var buf strings.Builder
	d := New(&buf, time.Hour)
	ch := make(chan []monitor.Change)
	r := NewRunner(d, ch, time.Hour)

	close(ch)
	done := make(chan struct{})
	go func() { r.Run(context.Background()); close(done) }()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("runner did not exit after channel close")
	}
}
