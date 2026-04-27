package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// Payload is the JSON body sent to the webhook endpoint.
type Payload struct {
	Timestamp time.Time      `json:"timestamp"`
	Changes   []ChangeEntry  `json:"changes"`
}

// ChangeEntry is a single port-change event in the payload.
type ChangeEntry struct {
	Port      int    `json:"port"`
	Proto     string `json:"proto"`
	Direction string `json:"direction"`
}

// Sender posts change notifications to an HTTP webhook.
type Sender struct {
	client  *http.Client
	url     string
	timeout time.Duration
}

// New returns a Sender that posts to the given URL.
func New(url string, timeout time.Duration) *Sender {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &Sender{
		client:  &http.Client{Timeout: timeout},
		url:     url,
		timeout: timeout,
	}
}

// Send serialises changes and POSTs them to the configured URL.
// It is a no-op when changes is empty.
func (s *Sender) Send(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	entries := make([]ChangeEntry, len(changes))
	for i, c := range changes {
		entries[i] = ChangeEntry{
			Port:      c.Port.Port,
			Proto:     c.Port.Proto,
			Direction: c.Direction,
		}
	}

	p := Payload{Timestamp: time.Now().UTC(), Changes: entries}
	body, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("webhook: marshal: %w", err)
	}

	resp, err := s.client.Post(s.url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: unexpected status %d", resp.StatusCode)
	}
	return nil
}
