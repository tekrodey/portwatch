package webhook_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/webhook"
)

func makeChange(port int, proto, dir string) monitor.Change {
	return monitor.Change{
		Port:      scanner.Port{Port: port, Proto: proto},
		Direction: dir,
	}
}

func TestSendPostsJSON(t *testing.T) {
	var got webhook.Payload
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(200)
	}))
	defer srv.Close()

	s := webhook.New(srv.URL, time.Second)
	changes := []monitor.Change{makeChange(8080, "tcp", "opened")}
	if err := s.Send(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(got.Changes))
	}
	if got.Changes[0].Port != 8080 || got.Changes[0].Direction != "opened" {
		t.Errorf("unexpected change entry: %+v", got.Changes[0])
	}
}

func TestSendNoOpOnEmpty(t *testing.T) {
	called := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer srv.Close()

	s := webhook.New(srv.URL, time.Second)
	if err := s.Send(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP call for empty changes")
	}
}

func TestSendReturnsErrorOnBadStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer srv.Close()

	s := webhook.New(srv.URL, time.Second)
	err := s.Send([]monitor.Change{makeChange(22, "tcp", "closed")})
	if err == nil {
		t.Fatal("expected error for 500 response")
	}
}

func TestNewDefaultTimeout(t *testing.T) {
	s := webhook.New("http://example.com", 0)
	if s == nil {
		t.Fatal("expected non-nil sender")
	}
}
