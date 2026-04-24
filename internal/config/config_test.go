package config_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/user/portwatch/internal/config"
)

func TestDefaultConfig(t *testing.T) {
	cfg := config.DefaultConfig()
	if cfg.PortRange != "1-1024" {
		t.Errorf("expected port_range '1-1024', got %q", cfg.PortRange)
	}
	if cfg.Interval != 30*time.Second {
		t.Errorf("expected interval 30s, got %v", cfg.Interval)
	}
	if cfg.AlertLevel != "info" {
		t.Errorf("expected alert_level 'info', got %q", cfg.AlertLevel)
	}
	if cfg.LogFile != "" {
		t.Errorf("expected empty log_file, got %q", cfg.LogFile)
	}
}

func TestLoadAndSave(t *testing.T) {
	tmp, err := os.CreateTemp(t.TempDir(), "portwatch-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	tmp.Close()

	orig := &config.Config{
		PortRange:  "80-443",
		Interval:   10 * time.Second,
		AlertLevel: "warn",
		LogFile:    "/var/log/portwatch.log",
	}

	if err := orig.Save(tmp.Name()); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := config.Load(tmp.Name())
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if loaded.PortRange != orig.PortRange {
		t.Errorf("port_range: want %q, got %q", orig.PortRange, loaded.PortRange)
	}
	if loaded.Interval != orig.Interval {
		t.Errorf("interval: want %v, got %v", orig.Interval, loaded.Interval)
	}
	if loaded.AlertLevel != orig.AlertLevel {
		t.Errorf("alert_level: want %q, got %q", orig.AlertLevel, loaded.AlertLevel)
	}
	if loaded.LogFile != orig.LogFile {
		t.Errorf("log_file: want %q, got %q", orig.LogFile, loaded.LogFile)
	}
}

func TestLoadMissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/path/portwatch.json")
	if err == nil {
		t.Fatal("expected error loading missing file, got nil")
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	tmp, err := os.CreateTemp(t.TempDir(), "portwatch-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	tmp.WriteString("{not valid json")
	tmp.Close()

	_, err = config.Load(tmp.Name())
	if err == nil {
		t.Fatal("expected error on invalid JSON, got nil")
	}
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Logf("got error type %T (acceptable): %v", err, err)
	}
}
