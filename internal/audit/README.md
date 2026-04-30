# audit

The `audit` package provides a structured, append-only log of every port change
event detected by portwatch.

## Overview

Each `monitor.Change` is serialised as a newline-delimited JSON record and
written to any `io.Writer` (file, stdout, network socket, etc.).

## Entry fields

| Field       | Type   | Description                          |
|-------------|--------|--------------------------------------|
| `timestamp` | string | RFC 3339 UTC time of the event       |
| `direction` | string | `"opened"` or `"closed"`             |
| `protocol`  | string | `"tcp"` or `"udp"`                   |
| `port`      | int    | Port number                          |
| `service`   | string | Optional well-known service name     |
| `note`      | string | Optional free-text annotation        |

## Usage

```go
f, _ := os.OpenFile("audit.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
l := audit.New(f)

// inside your monitor loop:
if err := l.Log(changes); err != nil {
    log.Println("audit error:", err)
}
```

Pass `nil` to `New` to fall back to `os.Stdout`.
