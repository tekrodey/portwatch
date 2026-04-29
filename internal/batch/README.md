# batch

The `batch` package groups a stream of `monitor.Change` events into fixed-size
or time-bounded batches before forwarding them downstream.

## Types

### `Batcher`

Accumulates changes and flushes them when either threshold is met:

| Threshold  | Description                                      |
|------------|--------------------------------------------------|
| `MaxSize`  | Flush when the buffer reaches this many changes. |
| `Interval` | Flush when this duration has elapsed.            |

```go
b := batch.New(50, 5*time.Second)
b.Add(changes)
if ready := b.Flush(); ready != nil {
    // process batch
}
```

`ForceFlush` drains the buffer unconditionally — useful at shutdown.

### `Runner`

Drives a `Batcher` in a background goroutine, calling a handler for each
non-empty batch on every tick. On context cancellation a final `ForceFlush`
ensures no changes are dropped.

```go
r := batch.NewRunner(b, time.Second, func(batch []monitor.Change) {
    notify.Send(batch)
})
r.Run(ctx)
```

## Notes

- `Batcher` is safe for concurrent use.
- Setting `MaxSize` to `0` disables size-based flushing.
- Setting `Interval` to `0` disables time-based flushing (not recommended for
  the `Runner` as it would never flush on its own).
