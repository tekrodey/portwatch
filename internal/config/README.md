# config

Package `config` provides loading, saving, and defaulting of portwatch daemon
configuration.

## Fields

| Field        | JSON key     | Default    | Description                                      |
|-------------|-------------|------------|--------------------------------------------------|
| PortRange   | port_range  | `1-1024`   | Port range to scan (passed to the scanner)       |
| Interval    | interval    | `30s`      | How often to run a scan cycle                    |
| AlertLevel  | alert_level | `info`     | Minimum alert level to emit (`info`/`warn`/`error`) |
| LogFile     | log_file    | `` (stdout)| Optional file path to write alerts               |

## Example `portwatch.json`

```json
{
  "port_range": "1-65535",
  "interval": "60s",
  "alert_level": "warn",
  "log_file": "/var/log/portwatch.log"
}
```

## Usage

```go
cfg, err := config.Load("portwatch.json")
if err != nil {
    // file not found — fall back to defaults
    cfg = config.DefaultConfig()
}
```
