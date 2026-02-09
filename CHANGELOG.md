# Changelog

All notable changes to this project will be documented in this file.

Format follows [Keep a Changelog](https://keepachangelog.com/). Versioning follows [Semantic Versioning](https://semver.org/).

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
