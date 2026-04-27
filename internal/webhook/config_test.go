package webhook_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/webhook"
)

func TestDefaultConfigDisabled(t *testing.T) {
	cfg := webhook.DefaultConfig()
	if cfg.Enabled {
		t.Error("expected webhook disabled by default")
	}
	if cfg.Timeout != 5*time.Second {
		t.Errorf("expected 5s timeout, got %v", cfg.Timeout)
	}
}

func TestNewFromConfigDisabled(t *testing.T) {
	cfg := webhook.Config{Enabled: false, URL: "http://example.com"}
	if s := webhook.NewFromConfig(cfg); s != nil {
		t.Error("expected nil sender when disabled")
	}
}

func TestNewFromConfigNoURL(t *testing.T) {
	cfg := webhook.Config{Enabled: true, URL: ""}
	if s := webhook.NewFromConfig(cfg); s != nil {
		t.Error("expected nil sender when URL is empty")
	}
}

func TestNewFromConfigEnabled(t *testing.T) {
	cfg := webhook.Config{Enabled: true, URL: "http://example.com", Timeout: 2 * time.Second}
	s := webhook.NewFromConfig(cfg)
	if s == nil {
		t.Fatal("expected non-nil sender for valid config")
	}
}
