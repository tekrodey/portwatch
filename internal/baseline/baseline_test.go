package baseline_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/baseline"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "baseline.json")
}

func TestNewMissingFileIsEmpty(t *testing.T) {
	b, err := baseline.New(tempPath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := len(b.All()); got != 0 {
		t.Fatalf("expected 0 entries, got %d", got)
	}
}

func TestSetAndContains(t *testing.T) {
	b, _ := baseline.New(tempPath(t))

	if err := b.Set("tcp", 8080); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if !b.Contains("tcp", 8080) {
		t.Error("expected tcp:8080 to be in baseline")
	}
	if b.Contains("udp", 8080) {
		t.Error("udp:8080 should not be in baseline")
	}
}

func TestRemove(t *testing.T) {
	b, _ := baseline.New(tempPath(t))
	_ = b.Set("tcp", 443)

	if err := b.Remove("tcp", 443); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if b.Contains("tcp", 443) {
		t.Error("tcp:443 should have been removed")
	}
}

func TestPersistence(t *testing.T) {
	path := tempPath(t)
	b, _ := baseline.New(path)
	_ = b.Set("tcp", 22)
	_ = b.Set("udp", 53)

	// reload from disk
	b2, err := baseline.New(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if !b2.Contains("tcp", 22) {
		t.Error("tcp:22 missing after reload")
	}
	if !b2.Contains("udp", 53) {
		t.Error("udp:53 missing after reload")
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	path := tempPath(t)
	if err := os.WriteFile(path, []byte("not-json"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := baseline.New(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestAllReturnsSnapshot(t *testing.T) {
	b, _ := baseline.New(tempPath(t))
	_ = b.Set("tcp", 80)
	_ = b.Set("tcp", 443)

	entries := b.All()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}
