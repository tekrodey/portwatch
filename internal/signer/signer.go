// Package signer provides HMAC-SHA256 request signing for outbound webhook payloads.
package signer

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"strconv"
	"time"
)

// ErrEmptySecret is returned when an empty secret is provided.
var ErrEmptySecret = errors.New("signer: secret must not be empty")

// Signer signs outbound HTTP requests using HMAC-SHA256.
type Signer struct {
	secret []byte
	now    func() time.Time
}

// New returns a Signer using the given secret.
// Returns ErrEmptySecret if secret is empty.
func New(secret string) (*Signer, error) {
	if secret == "" {
		return nil, ErrEmptySecret
	}
	return &Signer{secret: []byte(secret), now: time.Now}, nil
}

func newWithClock(secret string, now func() time.Time) (*Signer, error) {
	s, err := New(secret)
	if err != nil {
		return nil, err
	}
	s.now = now
	return s, nil
}

// Sign adds X-Portwatch-Signature and X-Portwatch-Timestamp headers to the request.
// The signature covers "timestamp.body" using HMAC-SHA256.
func (s *Signer) Sign(req *http.Request, body []byte) {
	ts := strconv.FormatInt(s.now().Unix(), 10)
	sig := s.compute(ts, body)
	req.Header.Set("X-Portwatch-Timestamp", ts)
	req.Header.Set("X-Portwatch-Signature", "sha256="+sig)
}

// Verify returns true if the signature header on req matches the expected
// HMAC for the given body and timestamp header.
func (s *Signer) Verify(req *http.Request, body []byte) bool {
	ts := req.Header.Get("X-Portwatch-Timestamp")
	want := req.Header.Get("X-Portwatch-Signature")
	if ts == "" || want == "" {
		return false
	}
	expected := "sha256=" + s.compute(ts, body)
	return hmac.Equal([]byte(want), []byte(expected))
}

func (s *Signer) compute(ts string, body []byte) string {
	h := hmac.New(sha256.New, s.secret)
	h.Write([]byte(ts))
	h.Write([]byte("."))
	h.Write(body)
	return hex.EncodeToString(h.Sum(nil))
}
