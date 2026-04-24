# notify

The `notify` package is responsible for delivering port-change events to the
operator. It decouples *formatting* from *delivery* so that both concerns can
be extended or tested independently.

## Types

### `Notifier`

Created with `New(out io.Writer, f Formatter)`. Writes formatted change events
to `out`. Defaults to `os.Stdout` when `out` is `nil`.

```go
n := notify.New(os.Stdout, nil)
n.Send(changes) // writes text lines to stdout
```

### `Formatter`

An interface with a single method:

```go
Format(changes []monitor.Change) string
```

The default implementation, `TextFormatter`, produces RFC3339-timestamped lines:

```
[2024-05-01T12:00:00Z] opened tcp:8080
[2024-05-01T12:00:00Z] closed tcp:22
```

## Extension

Provide your own `Formatter` to emit JSON, send webhooks, or integrate with
external alerting systems without changing the `Notifier` itself.
