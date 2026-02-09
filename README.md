# GOpher

<img src="gopher.webp" alt="Gopher" width="600">

Minimal metrics collection and aggregator. In Go. For people who don't want to waste an eternity setting up Grafana dashboards.

## Running

### TUI (local mode)

```sh
go run ./angel
```

Launches a live dashboard in your terminal showing CPU, RAM, and uptime. Metrics update automatically.

| Key | Action |
|---|---|
| `Ctrl+O` | Toggle compact/expanded view |
| `Ctrl+R` | Force rerender |
| `q` / `Ctrl+C` | Quit |

### Daemon (headless)

```sh
go run ./daemon
```

Runs collectors in the background and prints metrics to stdout. Eventually this will write to SQLite for persistence.

### Build

```sh
go build ./angel   # build the TUI binary
go build ./daemon  # build the daemon binary
```