# webhook

The `webhook` package delivers port-change notifications to an HTTP endpoint
by POSTing a JSON payload whenever the monitor detects opened or closed ports.

## Payload format

```json
{
  "timestamp": "2024-01-15T12:00:00Z",
  "changes": [
    {"port": 8080, "proto": "tcp", "direction": "opened"},
    {"port": 22,   "proto": "tcp", "direction": "closed"}
  ]
}
```

## Usage

```go
cfg := webhook.Config{
    Enabled: true,
    URL:     "https://hooks.example.com/portwatch",
    Timeout: 3 * time.Second,
}

sender := webhook.NewFromConfig(cfg)
if sender != nil {
    if err := sender.Send(changes); err != nil {
        log.Printf("webhook error: %v", err)
    }
}
```

## Configuration

| Field        | Type     | Default | Description                          |
|-------------|----------|---------|--------------------------------------|
| `url`        | string   | `""`    | HTTP endpoint to POST payloads to    |
| `timeout_ms` | duration | `5s`    | Per-request timeout                  |
| `enabled`    | bool     | `false` | Must be `true` to activate sending   |

The sender is a no-op when `changes` is empty, avoiding unnecessary requests.
