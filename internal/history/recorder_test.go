package history_test

import (
	"testing"

	"github.com/user/portwatch/internal/history"
	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/scanner"
)

func TestRecorderOpened(t *testing.T) {
	h, _ := history.New(tempPath(t), 100)
	r := history.NewRecorder(h)

	changes := []monitor.Change{
		{Port: scanner.Port{Port: 9090, Proto: "tcp"}, Closed: false},
	}
	if err := r.Record(changes); err != nil {
		t.Fatalf("Record: %v", err)
	}

	entries := h.All()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Action != "opened" {
		t.Errorf("expected action 'opened', got %q", entries[0].Action)
	}
	if entries[0].Port != 9090 {
		t.Errorf("expected port 9090, got %d", entries[0].Port)
	}
}

func TestRecorderClosed(t *testing.T) {
	h, _ := history.New(tempPath(t), 100)
	r := history.NewRecorder(h)

	changes := []monitor.Change{
		{Port: scanner.Port{Port: 22, Proto: "tcp"}, Closed: true},
	}
	_ = r.Record(changes)

	entries := h.All()
	if entries[0].Action != "closed" {
		t.Errorf("expected action 'closed', got %q", entries[0].Action)
	}
}

func TestRecorderEmpty(t *testing.T) {
	h, _ := history.New(tempPath(t), 100)
	r := history.NewRecorder(h)
	if err := r.Record(nil); err != nil {
		t.Fatalf("unexpected error on empty changes: %v", err)
	}
	if got := len(h.All()); got != 0 {
		t.Errorf("expected 0 entries, got %d", got)
	}
}
