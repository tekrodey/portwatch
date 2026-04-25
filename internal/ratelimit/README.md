# ratelimit

The `ratelimit` package provides a simple cooldown-based rate limiter for
suppressing repeated port-change alerts within a configurable time window.

## Overview

When a port flaps (opens and closes rapidly) or a scan cycle fires more
frequently than expected, the same alert can fire many times in quick
succession. The `Limiter` tracks the last time each `(port, protocol,
direction)` triple produced an alert and silently drops duplicates that
arrive before the cooldown expires.

## Usage

```go
limiter := ratelimit.New(30 * time.Second)

for _, change := range changes {
    if !limiter.Allow(change.Port, change.Protocol, change.Direction) {
        continue // suppressed — still within cooldown
    }
    notifier.Send([]monitor.Change{change})
}
```

## API

| Symbol | Description |
|---|---|
| `New(cooldown)` | Create a new Limiter with the given cooldown duration |
| `Allow(port, proto, dir)` | Returns `true` if the event should be forwarded |
| `Reset(port, proto, dir)` | Clears the suppression state for one event key |
| `Stats()` | Returns a summary string of tracked suppressions |

## Notes

- Each `(port, protocol, direction)` triple is tracked independently.
- `direction` is either `"opened"` or `"closed"`.
- The limiter is safe for concurrent use.
