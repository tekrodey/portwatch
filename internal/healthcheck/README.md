# healthcheck

The `healthcheck` package provides a lightweight component health registry for portwatch.

## Overview

`Checker` tracks the health of named internal components (e.g. `scanner`, `monitor`). Each component can be marked healthy or unhealthy at runtime. `Reporter` formats and writes health summaries as plain text or JSON.

## Usage

```go
c := healthcheck.New()
c.Register("scanner")
c.Register("monitor")

// Update health from your component:
c.SetHealthy("scanner", "last scan ok")
c.SetUnhealthy("monitor", "dial timeout")

// Check overall status:
if !c.Healthy() {
    log.Println("system degraded")
}

// Print a text summary:
r := healthcheck.NewReporter(c)
r.PrintText()

// Or emit JSON:
r.PrintJSON()
```

## Output

### Text
```
overall: DEGRADED
  ✓ scanner             12:00:01 — last scan ok
  ✗ monitor             12:00:02 — dial timeout
```

### JSON
```json
{
  "healthy": false,
  "statuses": [
    {"Name":"scanner","Healthy":true,"LastCheck":"...","Message":"last scan ok"},
    {"Name":"monitor","Healthy":false,"LastCheck":"...","Message":"dial timeout"}
  ]
}
```
