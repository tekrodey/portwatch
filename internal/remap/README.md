# remap

The `remap` package maps port/protocol pairs to human-readable alias strings.
It is used to annotate `monitor.Change` events with a `Label` field before
they reach alerting or reporting stages.

## Usage

```go
// Create an empty remapper and add rules in code.
r := remap.New()
r.Set("tcp", 80, "http")
r.Set("tcp", 443, "https")

fmt.Println(r.Alias("tcp", 80))   // http
fmt.Println(r.Alias("udp", 53))   // udp:53  (no rule → default)
```

```go
// Or load rules from a JSON file.
r, err := remap.Load("/etc/portwatch/aliases.json")
```

### Alias file format

A flat JSON object keyed by `"proto:port"`:

```json
{
  "tcp:22":   "ssh",
  "tcp:80":   "http",
  "tcp:443":  "https",
  "udp:53":   "dns",
  "tcp:5432": "postgres"
}
```

## Pipeline integration

Call `Apply` to annotate a slice of changes in one step:

```go
annotated := remapper.Apply(changes)
```

Each `Change.Label` will be set to the configured alias, or fall back to
`"proto:port"` when no rule matches.
