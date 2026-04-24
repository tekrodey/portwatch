# portwatch

A lightweight CLI daemon that monitors open ports and alerts on unexpected changes.

---

## Installation

```bash
go install github.com/yourusername/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/portwatch.git && cd portwatch && go build -o portwatch .
```

---

## Usage

Start the daemon with a default scan interval of 30 seconds:

```bash
portwatch start
```

Specify a custom interval and log file:

```bash
portwatch start --interval 60 --log /var/log/portwatch.log
```

Define a baseline of expected ports to suppress known alerts:

```bash
portwatch start --allow 22,80,443
```

When an unexpected port is detected, `portwatch` outputs an alert:

```
[ALERT] New open port detected: 4444 (TCP) at 2024-01-15 03:22:11
[ALERT] Port closed unexpectedly: 443 (TCP) at 2024-01-15 03:25:44
```

### Commands

| Command | Description |
|---|---|
| `start` | Start the monitoring daemon |
| `snapshot` | Print current open ports and exit |
| `version` | Print version information |

---

## License

[MIT](LICENSE)