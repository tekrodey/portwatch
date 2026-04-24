# history

The `history` package provides a persistent, bounded log of port change events.

## Types

### `Entry`
Represents a single recorded event:
- `Timestamp` — UTC time the change was detected
- `Port` — port number
- `Proto` — protocol (e.g. `tcp`)
- `Action` — `"opened"` or `"closed"`

### `History`
Thread-safe store backed by a JSON file on disk.

```go
h, err := history.New("/var/lib/portwatch/history.json", 500)
```

- `Add(Entry) error` — append and persist an entry; oldest entries are dropped when `maxSize` is exceeded.
- `All() []Entry` — return a snapshot of all stored entries.

### `Recorder`
Bridges `monitor.Change` slices (produced by the monitor loop) to the `History` store.

```go
rec := history.NewRecorder(h)
rec.Record(changes)
```

## Storage

Entries are stored as a JSON array. The file is rewritten on every `Add` call, keeping the format human-readable and easy to inspect with standard tools.

## Integration

Wire the recorder into the main daemon loop after calling `monitor.Scan`:

```go
changes, _ := mon.Scan()
alert.Notify(changes)
recorder.Record(changes)
```
