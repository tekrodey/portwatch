# ratelimit

Provides per-key rate limiting for port change events, suppressing repeated
alerts within a configurable cooldown window.

## Usage

```go
lim := ratelimit.New(30 * time.Second)

if lim.Allow(change) {
    // forward the change
}
```

## Behaviour

- The first occurrence of a (port, protocol, direction) tuple always passes.
- Subsequent occurrences within the cooldown window are suppressed.
- After the window expires the next occurrence passes and the timer resets.
- Direction (opened / closed) is treated as part of the key so an open and a
  close for the same port are tracked independently.
