package labelmap_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/labelmap"
)

func writeTempLabels(t *testing.T, entries []labelmap.Entry) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "labels-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(entries); err != nil {
		t.Fatal(err)
	}
	return f.Name()
}

func TestNewIsEmpty(t *testing.T) {
	m := labelmap.New()
	if m.Len() != 0 {
		t.Fatalf("expected 0 entries, got %d", m.Len())
	}
}

func TestSetAndGet(t *testing.T) {
	m := labelmap.New()
	labels := []labelmap.Label{{Key: "env", Value: "prod"}}
	m.Set(443, "tcp", labels)

	got, ok := m.Get(443, "tcp")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if len(got) != 1 || got[0].Key != "env" || got[0].Value != "prod" {
		t.Fatalf("unexpected labels: %+v", got)
	}
}

func TestGetMissingReturnsFalse(t *testing.T) {
	m := labelmap.New()
	_, ok := m.Get(80, "tcp")
	if ok {
		t.Fatal("expected no entry for unregistered port")
	}
}

func TestProtocolIsDistinct(t *testing.T) {
	m := labelmap.New()
	m.Set(53, "tcp", []labelmap.Label{{Key: "role", Value: "dns-tcp"}})
	m.Set(53, "udp", []labelmap.Label{{Key: "role", Value: "dns-udp"}})

	tcp, _ := m.Get(53, "tcp")
	udp, _ := m.Get(53, "udp")
	if tcp[0].Value == udp[0].Value {
		t.Fatal("tcp and udp entries should be distinct")
	}
}

func TestLoadMissingFileReturnsEmpty(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.json")
	m, err := labelmap.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Len() != 0 {
		t.Fatalf("expected empty map, got %d entries", m.Len())
	}
}

func TestLoadPopulatesEntries(t *testing.T) {
	entries := []labelmap.Entry{
		{Port: 22, Protocol: "tcp", Labels: []labelmap.Label{{Key: "service", Value: "ssh"}}},
		{Port: 80, Protocol: "tcp", Labels: []labelmap.Label{{Key: "service", Value: "http"}}},
	}
	path := writeTempLabels(t, entries)

	m, err := labelmap.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Len() != 2 {
		t.Fatalf("expected 2 entries, got %d", m.Len())
	}

	labels, ok := m.Get(22, "tcp")
	if !ok || labels[0].Value != "ssh" {
		t.Fatalf("unexpected labels for port 22: %+v", labels)
	}
}

func TestLoadInvalidJSONReturnsError(t *testing.T) {
	f, _ := os.CreateTemp(t.TempDir(), "bad-*.json")
	f.WriteString("not json")
	f.Close()

	_, err := labelmap.Load(f.Name())
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
