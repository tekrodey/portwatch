package signer

import (
	"net/http"
	"testing"
	"time"
)

var fixedNow = func() time.Time {
	return time.Unix(1_700_000_000, 0)
}

func newTestSigner(t *testing.T) *Signer {
	t.Helper()
	s, err := newWithClock("supersecret", fixedNow)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return s
}

func TestNewEmptySecretReturnsError(t *testing.T) {
	_, err := New("")
	if err != ErrEmptySecret {
		t.Fatalf("expected ErrEmptySecret, got %v", err)
	}
}

func TestNewValidSecretReturnsNonNil(t *testing.T) {
	s, err := New("abc")
	if err != nil || s == nil {
		t.Fatalf("expected valid signer, got err=%v signer=%v", err, s)
	}
}

func TestSignAddsHeaders(t *testing.T) {
	s := newTestSigner(t)
	body := []byte(`{"port":8080}`)
	req, _ := http.NewRequest(http.MethodPost, "http://example.com", nil)
	s.Sign(req, body)

	if req.Header.Get("X-Portwatch-Timestamp") == "" {
		t.Error("expected X-Portwatch-Timestamp to be set")
	}
	if req.Header.Get("X-Portwatch-Signature") == "" {
		t.Error("expected X-Portwatch-Signature to be set")
	}
}

func TestVerifyValidSignatureReturnsTrue(t *testing.T) {
	s := newTestSigner(t)
	body := []byte(`{"port":8080}`)
	req, _ := http.NewRequest(http.MethodPost, "http://example.com", nil)
	s.Sign(req, body)

	if !s.Verify(req, body) {
		t.Error("expected Verify to return true for valid signature")
	}
}

func TestVerifyTamperedBodyReturnsFalse(t *testing.T) {
	s := newTestSigner(t)
	body := []byte(`{"port":8080}`)
	req, _ := http.NewRequest(http.MethodPost, "http://example.com", nil)
	s.Sign(req, body)

	if s.Verify(req, []byte(`{"port":9999}`)) {
		t.Error("expected Verify to return false for tampered body")
	}
}

func TestVerifyMissingHeadersReturnsFalse(t *testing.T) {
	s := newTestSigner(t)
	req, _ := http.NewRequest(http.MethodPost, "http://example.com", nil)
	if s.Verify(req, []byte("body")) {
		t.Error("expected Verify to return false when headers are missing")
	}
}

func TestSignatureIsDeterministic(t *testing.T) {
	s := newTestSigner(t)
	body := []byte("payload")

	req1, _ := http.NewRequest(http.MethodPost, "http://example.com", nil)
	req2, _ := http.NewRequest(http.MethodPost, "http://example.com", nil)
	s.Sign(req1, body)
	s.Sign(req2, body)

	if req1.Header.Get("X-Portwatch-Signature") != req2.Header.Get("X-Portwatch-Signature") {
		t.Error("expected identical signatures for same input and clock")
	}
}
