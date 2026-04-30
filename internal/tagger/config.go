package tagger

import (
	"encoding/json"
	"fmt"
	"os"
)

// FileConfig is the on-disk representation of tagger rules.
type FileConfig struct {
	Rules []FileRule `json:"rules"`
}

// FileRule is a single JSON-serialisable tagging rule.
type FileRule struct {
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	Tag      string `json:"tag"`
}

// LoadConfig reads a JSON file and returns the parsed FileConfig.
// A missing file is treated as an empty config rather than an error.
func LoadConfig(path string) (FileConfig, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return FileConfig{}, nil
	}
	if err != nil {
		return FileConfig{}, fmt.Errorf("tagger: read config: %w", err)
	}
	var cfg FileConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return FileConfig{}, fmt.Errorf("tagger: parse config: %w", err)
	}
	return cfg, nil
}

// NewFromConfig constructs a Tagger from a FileConfig.
func NewFromConfig(cfg FileConfig) *Tagger {
	rules := make([]Rule, len(cfg.Rules))
	for i, r := range cfg.Rules {
		rules[i] = Rule{Port: r.Port, Protocol: r.Protocol, Tag: r.Tag}
	}
	return New(rules)
}
