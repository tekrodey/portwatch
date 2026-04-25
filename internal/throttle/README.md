# throttle

The `throttle` package limits how frequently an alert can be emitted for the
same (port, protocol, direction) tuple. This prevents alert storms when a port
flaps rapidly or when the scanner polls faster than events can be acted upon.

## Behaviour

- The first `Allow` call for any key **always passes**.
- Subsequent calls within the configured interval are **suppressed** (return `false`).
- Once the interval has elapsed the next call passes and resets the timer.
- `Reset` removes a key's recorded time so the next `Allow` passes immediately,
  useful when a port definitively closes and you want the next open event to
  fire without waiting.

## Usage

```go
th := throttle.New(30 * time.Second)

if th.Allow(change.Port, change.Proto, change.Direction) {
    notifier.Send([]monitor.Change{change})
}
```

## Testing

Pass a custom `Clock` via `NewWithClock` to control time in tests without
sleeping:

```go
clock, advance := fakeClock(time.Now())
th := throttle.NewWithClock(10*time.Second, clock)
advance(11 * time.Second)
```
