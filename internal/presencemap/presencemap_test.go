package presencemap

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func newFakeClock(start time.Time) func() time.Time {
	t := start
	return func() time.Time {
		v := t
		t = t.Add(time.Second)
		return v
	}
}

func TestTouchSetsFirstAndLastSeen(t *testing.T) {
	m := newWithClock(newFakeClock(epoch))
	m.Touch("tcp:80")
	e, ok := m.Get("tcp:80")
	if !ok {
		t.Fatal("expected entry")
	}
	if !e.FirstSeen.Equal(epoch) {
		t.Errorf("first seen: got %v want %v", e.FirstSeen, epoch)
	}
	if e.Count != 1 {
		t.Errorf("count: got %d want 1", e.Count)
	}
}

func TestTouchUpdatesLastSeen(t *testing.T) {
	m := newWithClock(newFakeClock(epoch))
	m.Touch("tcp:80")
	m.Touch("tcp:80")
	e, _ := m.Get("tcp:80")
	if e.FirstSeen.Equal(e.LastSeen) {
		t.Error("expected last_seen to advance after second touch")
	}
	if e.Count != 2 {
		t.Errorf("count: got %d want 2", e.Count)
	}
}

func TestGetMissingReturnsFalse(t *testing.T) {
	m := New()
	_, ok := m.Get("udp:53")
	if ok {
		t.Error("expected false for unknown key")
	}
}

func TestDeleteRemovesEntry(t *testing.T) {
	m := New()
	m.Touch("tcp:443")
	m.Delete("tcp:443")
	_, ok := m.Get("tcp:443")
	if ok {
		t.Error("expected entry to be removed")
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "presence.json")

	m := newWithClock(newFakeClock(epoch))
	m.Touch("tcp:22")
	m.Touch("tcp:22")

	if err := m.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	m2, err := Load(path, time.Now)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	e, ok := m2.Get("tcp:22")
	if !ok {
		t.Fatal("expected entry after load")
	}
	if e.Count != 2 {
		t.Errorf("count: got %d want 2", e.Count)
	}
}

func TestLoadMissingFileReturnsEmpty(t *testing.T) {
	m, err := Load("/nonexistent/path.json", time.Now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m == nil {
		t.Fatal("expected non-nil map")
	}
}

func TestLoadInvalidJSONReturnsError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not json"), 0o600)
	_, err := Load(path, time.Now)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestSaveProducesValidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "p.json")
	m := New()
	m.Touch("udp:123")
	_ = m.Save(path)
	data, _ := os.ReadFile(path)
	var out map[string]interface{}
	if err := json.Unmarshal(data, &out); err != nil {
		t.Errorf("invalid JSON: %v", err)
	}
}
