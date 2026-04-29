package batch

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

func TestRunnerFlushesOnTick(t *testing.T) {
	b := New(100, time.Hour) // only size/time threshold won't fire; runner drives it
	b.Add([]monitor.Change{makeChange(80)})

	var mu sync.Mutex
	var received []monitor.Change

	ctx, cancel := context.WithCancel(context.Background())
	r := NewRunner(b, 20*time.Millisecond, func(batch []monitor.Change) {
		mu.Lock()
		received = append(received, batch...)
		mu.Unlock()
		cancel()
	})

	// Override batcher to use interval flush by setting interval to 0
	b.interval = 0 // force time-based flush to always be ready

	done := make(chan struct{})
	go func() {
		r.Run(ctx)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("runner did not flush within timeout")
	}

	mu.Lock()
	defer mu.Unlock()
	if len(received) == 0 {
		t.Fatal("expected at least one change to be flushed")
	}
}

func TestRunnerForceFlushesOnCancel(t *testing.T) {
	b := New(1000, time.Hour) // high thresholds so normal Flush returns nil
	b.Add([]monitor.Change{makeChange(443), makeChange(8080)})

	var mu sync.Mutex
	var received []monitor.Change

	ctx, cancel := context.WithCancel(context.Background())
	r := NewRunner(b, time.Hour, func(batch []monitor.Change) {
		mu.Lock()
		received = append(received, batch...)
		mu.Unlock()
	})

	done := make(chan struct{})
	go func() {
		r.Run(ctx)
		close(done)
	}()

	cancel()
	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("runner did not exit after cancel")
	}

	mu.Lock()
	defer mu.Unlock()
	if len(received) != 2 {
		t.Fatalf("expected 2 force-flushed changes, got %d", len(received))
	}
}
