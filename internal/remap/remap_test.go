package remap_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/remap"
	"github.com/user/portwatch/internal/scanner"
)

func makeChange(proto string, port int, dir monitor.Direction) monitor.Change {
	return monitor.Change{
		Port:      scanner.Port{Proto: proto, Number: port},
		Direction: dir,
	}
}

func TestNewReturnsDefaultAlias(t *testing.T) {
	r := remap.New()
	got := r.Alias("tcp", 80)
	if got != "tcp:80" {
		t.Fatalf("expected tcp:80, got %s", got)
	}
}

func TestSetAndAlias(t *testing.T) {
	r := remap.New()
	r.Set("tcp", 443, "https")
	if got := r.Alias("tcp", 443); got != "https" {
		t.Fatalf("expected https, got %s", got)
	}
}

func TestAliasProtocolIsDistinct(t *testing.T) {
	r := remap.New()
	r.Set("tcp", 53, "dns-tcp")
	if got := r.Alias("udp", 53); got != "udp:53" {
		t.Fatalf("udp:53 should not inherit tcp alias, got %s", got)
	}
}

func TestApplyAnnotatesChanges(t *testing.T) {
	r := remap.New()
	r.Set("tcp", 80, "http")

	changes := []monitor.Change{
		makeChange("tcp", 80, monitor.Opened),
		makeChange("tcp", 9000, monitor.Opened),
	}

	out := r.Apply(changes)
	if out[0].Label != "http" {
		t.Errorf("expected http, got %s", out[0].Label)
	}
	if out[1].Label != "tcp:9000" {
		t.Errorf("expected tcp:9000, got %s", out[1].Label)
	}
}

func TestLoadFromFile(t *testing.T) {
	rules := map[string]string{
		"tcp:22":  "ssh",
		"tcp:443": "https",
	}
	data, _ := json.Marshal(rules)

	path := filepath.Join(t.TempDir(), "aliases.json")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}

	r, err := remap.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got := r.Alias("tcp", 22); got != "ssh" {
		t.Errorf("expected ssh, got %s", got)
	}
	if got := r.Alias("tcp", 443); got != "https" {
		t.Errorf("expected https, got %s", got)
	}
}

func TestLoadMissingFileReturnsError(t *testing.T) {
	_, err := remap.Load("/nonexistent/aliases.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
