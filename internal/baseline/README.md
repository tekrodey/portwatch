# baseline

The `baseline` package manages the **trusted set** of open ports for portwatch.

A baseline represents a "known-good" snapshot of ports that are expected to be
open. When a scan detects a port that is already in the baseline, no alert is
emitted. Only ports that deviate from the baseline (newly opened or
unexpectedly closed) trigger notifications.

## Usage

```go
// Load (or create) a baseline stored at a given path.
b, err := baseline.New("/var/lib/portwatch/baseline.json")
if err != nil {
    log.Fatal(err)
}

// Mark a port as trusted.
_ = b.Set("tcp", 443)

// Check before alerting.
if !b.Contains("tcp", port) {
    alert.Notify(...)
}

// Remove a port that should no longer be trusted.
_ = b.Remove("tcp", 8080)

// Enumerate all trusted entries.
for _, e := range b.All() {
    fmt.Printf("%s:%d (added %s)\n", e.Protocol, e.Port, e.AddedAt.Format(time.RFC3339))
}
```

## Persistence

The baseline is stored as a JSON file. Every call to `Set` or `Remove`
immediately writes the updated list to disk so that the trusted set survives
process restarts.
