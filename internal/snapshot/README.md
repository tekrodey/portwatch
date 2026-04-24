# snapshot

The `snapshot` package provides a lightweight mechanism for capturing and
persisting the current state of observed ports to disk.

## Concepts

- **Entry** — a single port observation: port number, protocol, state, and
  capture timestamp.
- **Snapshot** — an ordered collection of `Entry` values plus a top-level
  `created_at` timestamp representing when the scan was taken.

## Usage

```go
// Create a new snapshot and populate it from scan results.
s := snapshot.New()
s.Add(8080, "tcp", "open")
s.Add(22,   "tcp", "open")

// Persist to disk.
if err := s.Save("/var/lib/portwatch/last.json"); err != nil {
    log.Fatal(err)
}

// Reload on the next run (returns empty snapshot if file is absent).
prev, err := snapshot.Load("/var/lib/portwatch/last.json")
if err != nil {
    log.Fatal(err)
}
```

## File format

Snapshots are stored as indented JSON so they are human-readable and can be
inspected with standard tooling such as `jq`.
