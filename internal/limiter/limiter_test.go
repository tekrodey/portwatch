package limiter_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/user/portwatch/internal/limiter"
)

func TestNewPanicsOnZero(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for n=0")
		}
	}()
	limiter.New(0)
}

func TestAcquireAndRelease(t *testing.T) {
	l := limiter.New(2)
	ctx := context.Background()

	if err := l.Acquire(ctx); err != nil {
		t.Fatalf("first Acquire: %v", err)
	}
	if got := l.Available(); got != 1 {
		t.Fatalf("Available after one Acquire: want 1, got %d", got)
	}
	l.Release()
	if got := l.Available(); got != 2 {
		t.Fatalf("Available after Release: want 2, got %d", got)
	}
}

func TestCapacity(t *testing.T) {
	l := limiter.New(5)
	if got := l.Capacity(); got != 5 {
		t.Fatalf("Capacity: want 5, got %d", got)
	}
}

func TestAcquireBlocksWhenFull(t *testing.T) {
	l := limiter.New(1)
	ctx := context.Background()

	if err := l.Acquire(ctx); err != nil {
		t.Fatalf("first Acquire: %v", err)
	}

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := l.Acquire(ctxTimeout)
	if err == nil {
		t.Fatal("expected ErrLimitExceeded when limiter is full")
	}
	if err != limiter.ErrLimitExceeded {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestConcurrentAcquire(t *testing.T) {
	const (
		cap     = 3
		workers = 9
	)
	l := limiter.New(cap)
	ctx := context.Background()

	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		maxSeen int
		active  int
	)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := l.Acquire(ctx); err != nil {
				return
			}
			defer l.Release()
			mu.Lock()
			active++
			if active > maxSeen {
				maxSeen = active
			}
			mu.Unlock()
			time.Sleep(10 * time.Millisecond)
			mu.Lock()
			active--
			mu.Unlock()
		}()
	}
	wg.Wait()

	if maxSeen > cap {
		t.Fatalf("concurrency exceeded cap: max active=%d, cap=%d", maxSeen, cap)
	}
}
