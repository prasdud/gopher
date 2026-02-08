# GOpher

Minimal Linux system metrics collector written in Go. For people who don't want to waste an eternity setting up Prometheus + Grafana.

## Project structure

- `daemon/` — Runs on the server. Collects metrics and writes to SQLite.
- `angel/` — Runs on the user's laptop. TUI that connects to the daemon over SSH.
- `internal/` — Shared metric collection code.
  - `uptime.go` — System uptime via `unix.Sysinfo()` syscall. Returns hours/minutes/seconds.
  - `ram.go` — RAM stats parsed from `/proc/meminfo`. Reports total/free/available/used in GB.
  - `cpu-usage.go` — CPU usage from `/proc/stat`. `GetCpuPercent()` samples twice over 200ms to compute delta-based percentage.

## Build and run

```sh
go run ./daemon   # run the metrics collector
go build ./daemon # build the daemon binary
```

## Architecture

### One binary, three modes

1. **`gopher daemon`** — always-on systemd service on the server. Collects metrics, stores in SQLite (`~/.gopher/`), auto-purges after 7 days. Listens on a unix socket (`~/.gopher/gopher.sock`) for angel connections.
2. **`gopher connect user@host`** — SSHes into the server, talks to the daemon over the unix socket, streams live metrics + queries history, renders TUI locally. Auto-deploys the binary to the server on first connect if not present.
3. **`gopher`** — local mode. TUI for monitoring the machine you're on (no SSH, no persistence, just live).

### Data flow

```
Your laptop                              Remote server
┌────────────────┐                   ┌─────────────────────┐
│ gopher connect │ ──── SSH ──────── │  gopher daemon      │
│   (TUI)        │ ← JSON stream ── │    ↓                 │
│                │ ← history query ─ │  SQLite (~/.gopher/) │
└────────────────┘                   └─────────────────────┘
```

### Concurrency model

Each metric collector runs in its own goroutine with its own tick interval. Collectors push results into channels that a central aggregator consumes (fan-in pattern).

```
goroutine: CPU collector ──→ chan CpuMetric ──┐
goroutine: RAM collector ──→ chan RamMetric ──┤
goroutine: Disk collector ─→ chan DiskMetric ─┼──→ aggregator goroutine ──→ SQLite
goroutine: Net collector ──→ chan NetMetric ──┘
```

Why this model:
- **Different metrics have different sample rates.** CPU needs sub-second sampling. Disk barely changes — every 30s is fine.
- **One slow collector doesn't block others.** If disk I/O stalls, CPU and RAM keep flowing.
- **Adding a new metric = adding a goroutine.** Scales without touching existing code.

The aggregator batches writes to SQLite (flush every second) and runs TTL purge on a `time.Ticker`.

### Angel (TUI) concurrency

```
goroutine: SSH stream reader ──→ chan Metrics ──→ TUI render loop (bubbletea)
goroutine: history fetcher ────→ chan History ──↗
```

TUI render loop never blocks on network I/O. SSH hiccup = "reconnecting" state, not a freeze.

### Collector interface

Each metric collector implements:

```go
type Collector interface {
    Name() string
    Interval() time.Duration
    Collect() (Metric, error)
}
```

Daemon starts each collector in a goroutine, wires it to a channel, aggregator fans them in. New metric = implement the interface and register it.

### Future: multi-server

One goroutine per SSH connection, all feeding into the same TUI channel via `select`:

```
goroutine: server-1 SSH ──→ chan Metrics ──┐
goroutine: server-2 SSH ──→ chan Metrics ──┼──→ TUI
goroutine: server-3 SSH ──→ chan Metrics ──┘
```

## Design decisions

| Decision | Choice | Rationale |
|---|---|---|
| Storage | SQLite | Single file, zero setup, good for time-series queries |
| TTL | Fixed 7 days | Auto-purge, no config needed |
| Transport | SSH + unix socket | No ports to open, no auth tokens, user already has SSH keys |
| Daemon | Always-on systemd service | Metrics collected even when angel isn't connected |
| Install | Auto-deploy via angel | First `gopher connect` copies binary to server if missing |
| TUI | bubbletea + lipgloss | Most popular Go TUI stack |
| Metrics | CPU/RAM/uptime now, eventually everything | Collector interface makes adding metrics trivial |

## Key conventions

- **Linux-only**: All metric collection reads from `/proc` or uses unix syscalls.
- **Module path**: `github.com/prasdud/gopher`
- **Go version**: 1.24.9
- **External dependency**: `golang.org/x/sys` (for unix syscalls in uptime)
- All metric types and functions live in the `internal` package.
- Metric functions return pointer-to-struct + error.
- RAM values are converted from kB to GB using integer division. Used RAM calculation matches the `free` command (Total - Free - Buffers - Cached - SReclaimable).

## Known issues

- `daemon/main.go` shadows the `err` variable — only the last error is checked. Earlier errors from `GetUptime`, `GetRamDetails`, and `GetCpuUsage` are silently ignored.
- RAM reporting uses integer GB which loses precision. Should move to float64 or report in MB.
