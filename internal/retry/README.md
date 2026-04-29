# retry

Provides exponential-backoff retry logic for transient failures such as webhook
delivery or external API calls.

## Usage

```go
r := retry.New(retry.DefaultConfig())
err := r.Do(ctx, func() error {
    return sendWebhook(payload)
})
if errors.Is(err, retry.ErrMaxAttempts) {
    log.Println("webhook delivery failed after all attempts")
}
```

## Config

| Field | Default | Description |
|-------|---------|-------------|
| `MaxAttempts` | `3` | Total number of attempts (including the first). |
| `BaseDelay` | `200ms` | Initial wait between attempts. |
| `MaxDelay` | `5s` | Upper bound on inter-attempt delay. |

Delay doubles after each failure and is capped at `MaxDelay`.

## Cancellation

Passing a cancelled `context.Context` causes `Do` to return `context.Canceled`
immediately, without consuming any remaining attempts.
