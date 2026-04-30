package tagger_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/tagger"
)

func writeTempConfig(t *testing.T, cfg tagger.FileConfig) string {
	t.Helper()
	data, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	p := filepath.Join(t.TempDir(), "tagger.json")
	if err := os.WriteFile(p, data, 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}
	return p
}

func TestLoadConfigMissingFileIsEmpty(t *testing.T) {
	cfg, err := tagger.LoadConfig("/nonexistent/path/tagger.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Rules) != 0 {
		t.Fatalf("expected empty rules, got %d", len(cfg.Rules))
	}
}

func TestLoadConfigInvalidJSON(t *testing.T) {
	p := filepath.Join(t.TempDir(), "bad.json")
	if err := os.WriteFile(p, []byte("not-json"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := tagger.LoadConfig(p); err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestLoadConfigRoundTrip(t *testing.T) {
	want := tagger.FileConfig{
		Rules: []tagger.FileRule{
			{Port: 22, Protocol: "tcp", Tag: "ssh"},
			{Port: 53, Protocol: "udp", Tag: "dns"},
		},
	}
	p := writeTempConfig(t, want)
	got, err := tagger.LoadConfig(p)
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}
	if len(got.Rules) != len(want.Rules) {
		t.Fatalf("expected %d rules, got %d", len(want.Rules), len(got.Rules))
	}
	for i, r := range got.Rules {
		if r != want.Rules[i] {
			t.Errorf("rule %d: got %+v, want %+v", i, r, want.Rules[i])
		}
	}
}

func TestNewFromConfigBuildsWorkingTagger(t *testing.T) {
	cfg := tagger.FileConfig{
		Rules: []tagger.FileRule{
			{Port: 443, Protocol: "tcp", Tag: "https"},
		},
	}
	tg := tagger.NewFromConfig(cfg)
	tags := tg.Tag(makeChange(443, "tcp"))
	if len(tags) != 1 || tags[0] != "https" {
		t.Fatalf("expected [https], got %v", tags)
	}
}
