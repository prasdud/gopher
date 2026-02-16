# Changelog

All notable changes to this project will be documented in this file.

Format follows [Keep a Changelog](https://keepachangelog.com/). Versioning follows [Semantic Versioning](https://semver.org/).

## [1.2.0] - 2026-02-16

### Added

- IPv4 metric collection: Fetches public IP from `api.ipify.org` (`internal/ip-v4.go`)
- IPv4 integration in TUI: Displayed in a prominent bold, bordered box next to the GOPHER header
- Immediate collection: Modified `angel/main.go` to run all collectors once immediately on startup instead of waiting for the first ticker interval
- Test suite: Added `tests/ip-v4_test.go` for public IP verification

### Changed

- TUI Layout: Removed redundant IP address line from the main metrics list to keep the interface cleaner

## [1.1.0] - 2026-02-16

### Added

- Swap memory metrics: `TotalSwap`, `UsedSwap`, `AvailableSwap` in `RamDetails` (parsed from `SwapTotal`/`SwapFree` in `/proc/meminfo`)
- Swap usage display in TUI — compact mode shows `SWP` bar, expanded mode adds swap row inside the Memory box
- Test suite in `tests/` directory — `ram_test.go` with tests for RAM values, bounds validation, swap, and `MetricName()`

### Changed

- RAM fields changed from `int64` to `float64` across the codebase (`internal/ram.go`, `tui/tui.go`, `tui/components.go`) for sub-GB precision
- `FormatRAM()` now takes `float64` args and formats with 2 decimal places
- Uptime label in compact view now reads `UPTIME ⏱` instead of just `⏱`

### Fixed

- Better vertical spacing between metric groups in compact TUI view
- Expanded view format strings updated from `%d` to `%.2f` to match the `float64` migration

## [1.0.0] - 2026-02-10

### Added

- Live TUI dashboard using bubbletea + lipgloss (`angel/main.go`, `tui/`)
- Compact mode: inline CPU/RAM progress bars with uptime
- Expanded mode: bordered boxes with detailed RAM stats (toggle with Ctrl+O)
- Responsive layout: side-by-side at >80 cols, stacked when narrow, "Terminal too small" at <40x10
- JSON-based theming system (`tui/theme.go`, `styles/default.json`)
- Reusable TUI components: progress bar, bordered box, ASCII header (`tui/components.go`)
- Clean goroutine shutdown via quit channel in angel

### Changed

- Moved `CollectorResult` from `daemon/main.go` to `internal/collector/collector.go` so both daemon and angel can share it

## [0.5.0] - 2026-02-09

### Changed

- Rewrote `daemon/main.go` from one-shot print to concurrent daemon loop
- Each collector runs in its own goroutine with independent tick intervals (CPU 1s, RAM 2s, uptime 5s)
- Fan-in aggregator reads from a shared channel and prints results as they arrive

### Fixed

- Error handling no longer shadows `err` — each collector handles errors independently

## [0.4.0] - 2026-02-09

### Added

- `Collector` and `Metric` interfaces in `internal/collector/collector.go`
- `UptimeCollector`, `RamCollector`, `CpuCollector` wrappers in `internal/collector/`
- `MetricName()` method on `Uptime`, `RamDetails`, `CpuUsage`, and `CpuPercent` structs
- `CpuPercent` wrapper struct for CPU percentage (needed to satisfy `Metric` interface)

## [0.3.0] - 2026-01-27

### Added

- CPU usage metric via `/proc/stat` parsing (`internal/cpu-usage.go`)
- `GetCpuUsage()` — returns raw idle and total CPU cycles
- `GetCpuPercent()` — delta-based CPU percentage over 200ms sample window

## [0.2.0] - 2026-01-24

### Added

- RAM metrics via `/proc/meminfo` parsing (`internal/ram.go`)
- Reports total, free, available, and used RAM in GB
- Used RAM calculation matches the `free` command (Total - Free - Buffers - Cached - SReclaimable)
- README with usage instructions
- Project mascot image

## [0.1.0] - 2026-01-24

### Added

- System uptime metric via `unix.Sysinfo()` syscall (`internal/uptime.go`)
- Returns hours, minutes, seconds breakdown

## [0.0.1] - 2026-01-23

### Added

- Initial project scaffold
- Go module (`github.com/prasdud/gopher`) with Go 1.24.9
- `daemon/` and `angel/` package structure
- `golang.org/x/sys` dependency for unix syscalls
