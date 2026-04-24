package config

import (
	"encoding/json"
	"os"
	"time"
)

// Config holds the portwatch daemon configuration.
type Config struct {
	// PortRange defines the range of ports to scan, e.g. "1-1024".
	PortRange string `json:"port_range"`
	// Interval is how often to scan for port changes.
	Interval time.Duration `json:"interval"`
	// AlertLevel controls minimum severity to emit: "info", "warn", or "error".
	AlertLevel string `json:"alert_level"`
	// LogFile is an optional path to write alerts to. Empty means stdout.
	LogFile string `json:"log_file"`
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		PortRange:  "1-1024",
		Interval:   30 * time.Second,
		AlertLevel: "info",
		LogFile:    "",
	}
}

// Load reads a JSON config file from path and returns the parsed Config.
// Fields absent in the file retain their zero values; callers should apply
// defaults before calling Load if desired.
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cfg := DefaultConfig()
	if err := json.NewDecoder(f).Decode(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Save writes the Config to path as indented JSON.
func (c *Config) Save(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(c)
}
