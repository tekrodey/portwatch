# digest

The `digest` package produces periodic summary reports of port change activity.
Instead of alerting on every individual change, a digest accumulates events over
a configurable window and emits a single human-readable report at the end.

## Types

### `Summary`
Holds aggregated `Opened` and `Closed` change slices for a time window.
`Summary.String()` renders a formatted report suitable for logging or email.

### `Digest`
Accumulates `monitor.Change` values via `Add` and writes a `Summary` to an
`io.Writer` when `Flush` is called and the interval has elapsed.

```go
d := digest.New(os.Stdout, 5*time.Minute)
d.Add(changes)
d.Flush() // no-op if interval has not elapsed
```

### `Runner`
Drives a `Digest` from a `<-chan []monitor.Change`, flushing on a ticker and
on context cancellation.

```go
ch := make(chan []monitor.Change)
d := digest.New(os.Stdout, 5*time.Minute)
r := digest.NewRunner(d, ch, 5*time.Minute)
go r.Run(ctx)
```

## Notes

- `New(nil, interval)` falls back to `os.Stdout`.
- `Flush` resets the accumulator and advances the window start time.
- The `Runner` performs a final `Flush` when the context is cancelled so no
  events are silently dropped at shutdown.
