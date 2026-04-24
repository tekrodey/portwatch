package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/snapshot"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "snap.json")
}

func TestNewIsEmpty(t *testing.T) {
	s := snapshot.New()
	if len(s.Entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(s.Entries))
	}
	if s.CreatedAt.IsZero() {
		t.Fatal("expected non-zero CreatedAt")
	}
}

func TestAddEntry(t *testing.T) {
	s := snapshot.New()
	s.Add(8080, "tcp", "open")
	s.Add(53, "udp", "open")

	if len(s.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(s.Entries))
	}
	if s.Entries[0].Port != 8080 || s.Entries[0].Protocol != "tcp" {
		t.Errorf("unexpected first entry: %+v", s.Entries[0])
	}
}

func TestSaveAndLoad(t *testing.T) {
	path := tempPath(t)
	s := snapshot.New()
	s.Add(443, "tcp", "open")

	if err := s.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(loaded.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(loaded.Entries))
	}
	if loaded.Entries[0].Port != 443 {
		t.Errorf("expected port 443, got %d", loaded.Entries[0].Port)
	}
}

func TestLoadMissingFileReturnsEmpty(t *testing.T) {
	s, err := snapshot.Load("/nonexistent/path/snap.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(s.Entries) != 0 {
		t.Fatalf("expected empty snapshot, got %d entries", len(s.Entries))
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	path := tempPath(t)
	if err := os.WriteFile(path, []byte("not-json{"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := snapshot.Load(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}
