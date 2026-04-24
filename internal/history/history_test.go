package history_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/history"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "history.json")
}

func TestAddAndAll(t *testing.T) {
	h, err := history.New(tempPath(t), 100)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	e := history.Entry{Timestamp: time.Now(), Port: 8080, Proto: "tcp", Action: "opened"}
	if err := h.Add(e); err != nil {
		t.Fatalf("Add: %v", err)
	}
	entries := h.All()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Port != 8080 {
		t.Errorf("expected port 8080, got %d", entries[0].Port)
	}
}

func TestMaxSizeTruncates(t *testing.T) {
	h, err := history.New(tempPath(t), 3)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	for i := 0; i < 5; i++ {
		_ = h.Add(history.Entry{Timestamp: time.Now(), Port: 1000 + i, Proto: "tcp", Action: "opened"})
	}
	if got := len(h.All()); got != 3 {
		t.Errorf("expected 3 entries after truncation, got %d", got)
	}
}

func TestPersistence(t *testing.T) {
	p := tempPath(t)
	h1, _ := history.New(p, 100)
	_ = h1.Add(history.Entry{Timestamp: time.Now(), Port: 443, Proto: "tcp", Action: "closed"})

	h2, err := history.New(p, 100)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if got := len(h2.All()); got != 1 {
		t.Errorf("expected 1 persisted entry, got %d", got)
	}
}

func TestLoadMissingFileIsOK(t *testing.T) {
	p := filepath.Join(t.TempDir(), "nonexistent.json")
	h, err := history.New(p, 10)
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if got := len(h.All()); got != 0 {
		t.Errorf("expected 0 entries, got %d", got)
	}
}

func TestLoadInvalidJSONReturnsError(t *testing.T) {
	p := tempPath(t)
	_ = os.WriteFile(p, []byte("not-json"), 0o644)
	_, err := history.New(p, 10)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}
