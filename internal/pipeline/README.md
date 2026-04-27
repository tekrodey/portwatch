# pipeline

The `pipeline` package wires together the individual processing stages that
every detected port-change passes through before reaching an output sink
(notifier, webhook, digest, etc.).

## Concept

A `Pipeline` holds an ordered list of `Stage` functions. Each stage receives a
batch of `monitor.Change` values and returns a (possibly filtered or enriched)
batch. Stages are applied in order; if any stage returns an empty slice the
remaining stages are skipped.

```
raw changes
    │
    ▼
[ FilterStage ]   – drop ignored / non-whitelisted ports
    │
    ▼
[ DedupStage ]    – suppress duplicate events within the dedup window
    │
    ▼
[ DebounceStage ] – hold transient flaps until the debounce window passes
    │
    ▼
[ EnrichStage ]   – attach service names and reverse-DNS labels
    │
    ▼
 output sinks
```

## Usage

```go
p := pipeline.New(
    pipeline.FilterStage(f),
    pipeline.DedupStage(dd),
    pipeline.DebounceStage(db),
    pipeline.EnrichStage(en),
)

result := p.Run(rawChanges)
```

Or use `DefaultStages` to get the standard ordering:

```go
stages := pipeline.DefaultStages(ctx, f, dd, db, en)
p := pipeline.New(stages...)
```

## Custom stages

Any function with the signature `func([]monitor.Change) []monitor.Change` can
be used as a `Stage`, making it straightforward to inject project-specific
transformations or test doubles.
