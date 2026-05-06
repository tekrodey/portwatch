# signer

The `signer` package provides HMAC-SHA256 signing and verification for outbound webhook requests sent by portwatch.

## Overview

Each outbound HTTP request is signed by adding two headers:

| Header | Description |
|---|---|
| `X-Portwatch-Timestamp` | Unix timestamp (seconds) at signing time |
| `X-Portwatch-Signature` | `sha256=<hex>` HMAC of `"<timestamp>.<body>"` |

The signature binds the body to a specific point in time, protecting against replay attacks when the receiver validates the timestamp.

## Usage

```go
s, err := signer.New(os.Getenv("PORTWATCH_WEBHOOK_SECRET"))
if err != nil {
    log.Fatal(err)
}

bodyBytes, _ := json.Marshal(payload)
req, _ := http.NewRequest(http.MethodPost, webhookURL, bytes.NewReader(bodyBytes))
s.Sign(req, bodyBytes)
```

## Verification (receiver side)

```go
bodyBytes, _ := io.ReadAll(r.Body)
if !s.Verify(r, bodyBytes) {
    http.Error(w, "invalid signature", http.StatusUnauthorized)
    return
}
```

## Notes

- An empty secret returns `ErrEmptySecret` at construction time.
- Signature comparison uses `hmac.Equal` to prevent timing attacks.
