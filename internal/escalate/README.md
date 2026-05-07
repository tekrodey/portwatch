# escalate

Tracks how many times the same port change event recurs within a sliding window
and promotes it to a higher severity `Level` so downstream components can react
appropriately.

## Levels

| Level      | Meaning                                      |
|------------|----------------------------------------------|
| `Normal`   | First occurrence within the window           |
| `Elevated` | Second occurrence within the window          |
| `Critical` | Third or subsequent occurrence in the window |

## Usage

```go
esc := escalate.New(5 * time.Minute)

for _, c := range changes {
    level := esc.Assess(c)
    switch level {
    case escalate.LevelCritical:
        // page on-call
    case escalate.LevelElevated:
        // send Slack message
    default:
        // log only
    }
}
```

## Resetting

Call `Reset(change)` once an operator acknowledges an alert to avoid the count
carrying over into the next recurrence window.

```go
esc.Reset(change)
```

## Thread safety

All methods are safe for concurrent use.
