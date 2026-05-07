# presencemap

The `presencemap` package tracks which ports have been observed across scan cycles, recording first-seen and last-seen timestamps along with a running count of sightings.

## Types

### `Map`

A thread-safe store of `Entry` values keyed by an arbitrary string (typically `"proto:port"`).

```go
m := presencemap.New()
m.Touch("tcp:80")          // record a sighting
e, ok := m.Get("tcp:80")  // retrieve metadata
m.Delete("tcp:80")        // remove on port close
```

Persistence is provided via `Save(path)` and `Load(path, clockFn)`.

### `Tracker`

`Tracker` wraps a `Map` and integrates with the `monitor.Change` stream:

```go
tr := presencemap.NewTracker(m)
tr.Apply(changes)     // opened → Touch, closed → Delete
snap := tr.Snapshot() // copy of all current entries
```

## Entry fields

| Field | Description |
|---|---|
| `FirstSeen` | Time the port was first observed |
| `LastSeen` | Time of the most recent sighting |
| `Count` | Total number of times the port was touched |

## Persistence

The map is serialised as JSON and can be reloaded on startup:

```go
m, err := presencemap.Load("/var/lib/portwatch/presence.json", time.Now)
```

A missing file is treated as an empty map (no error).

## Key format

Keys are conventionally formatted as `"proto:port"` (e.g. `"tcp:443"`, `"udp:53"`). The `Key` helper constructs a canonical key from its components:

```go
k := presencemap.Key("tcp", 443) // returns "tcp:443"
m.Touch(k)
```
