# watchdog

The `watchdog` package provides a lightweight liveness monitor that periodically
calls a user-supplied check function and invokes a configurable handler when the
check exceeds a timeout or returns an error.

## Types

| Type | Description |
|------|-------------|
| `Watchdog` | Core struct; runs a probe on every `Interval` tick |
| `Config` | Interval, Timeout, and OnHang callback |
| `Runner` | Convenience wrapper that starts the loop in a goroutine |

## Usage

```go
cfg := watchdog.Config{
    Interval: 30 * time.Second,
    Timeout:  10 * time.Second,
    OnHang: func(err error) {
        log.Printf("liveness check failed: %v", err)
    },
}

probe := watchdog.ScanProbe(func(ctx context.Context) error {
    return scanner.Scan(ctx)
})

runner := watchdog.NewRunner(cfg, probe)
runner.Start(ctx)
```

## Behaviour

- The first probe fires after one full `Interval`.
- If the probe goroutine does not return within `Timeout`, `OnHang` is called
  with `context.DeadlineExceeded`.
- If the probe returns a non-nil error before the timeout, `OnHang` is called
  with that error.
- A nil `OnHang` is silently replaced with a no-op so callers need not guard
  against it.
