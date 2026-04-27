package webhook

import "time"

// Config holds webhook-specific configuration.
type Config struct {
	// URL is the endpoint to POST change payloads to.
	URL string `json:"url"`

	// Timeout is the per-request HTTP timeout.
	// Defaults to 5 s when zero.
	Timeout time.Duration `json:"timeout_ms"`

	// Enabled controls whether the webhook sender is active.
	Enabled bool `json:"enabled"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Enabled: false,
		Timeout: 5 * time.Second,
	}
}

// NewFromConfig constructs a Sender from a Config.
// Returns nil when the webhook is disabled or has no URL.
func NewFromConfig(cfg Config) *Sender {
	if !cfg.Enabled || cfg.URL == "" {
		return nil
	}
	return New(cfg.URL, cfg.Timeout)
}
