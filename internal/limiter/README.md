# limiter

Package `limiter` provides a simple semaphore-based concurrency limiter for
capping the number of in-flight operations inside portwatch.

## Usage

```go
l := limiter.New(4) // at most 4 concurrent operations

if err := l.Acquire(ctx); err != nil {
    // context cancelled or timed out before a slot was free
    return err
}
defer l.Release()

// ... do work ...
```

## API

| Symbol | Description |
|---|---|
| `New(n int) *Limiter` | Create a limiter with `n` slots. Panics if `n < 1`. |
| `Acquire(ctx) error` | Block until a slot is free or ctx expires. |
| `Release()` | Free one slot. Must match every successful `Acquire`. |
| `Available() int` | Number of free slots at this instant. |
| `Capacity() int` | Total number of slots. |

## Notes

- `Acquire` and `Release` must be paired; calling `Release` without a prior
  `Acquire` will block indefinitely.
- Use `context.WithTimeout` or `context.WithDeadline` to bound wait time.
