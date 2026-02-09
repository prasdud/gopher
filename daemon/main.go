package main

import (
	"fmt"
	"time"

	"github.com/prasdud/gopher/internal"
	"github.com/prasdud/gopher/internal/collector"
)

// CollectorResult is what each goroutine sends into the shared channel.
type CollectorResult struct {
	Name   string
	Metric collector.Metric
	Err    error
}

func main() {
	// 1. Register all collectors
	collectors := []collector.Collector{
		&collector.UptimeCollector{},
		&collector.RamCollector{},
		&collector.CpuCollector{},
	}

	// 2. One shared channel for all results
	results := make(chan CollectorResult)

	// 3. One goroutine per collector, each with its own ticker
	for _, c := range collectors {
		go func(c collector.Collector) {
			ticker := time.NewTicker(c.Interval())
			for range ticker.C {
				metric, err := c.Collect()
				results <- CollectorResult{
					Name:   c.Name(),
					Metric: metric,
					Err:    err,
				}
			}
		}(c)
	}

	// 4. Aggregator: read from channel forever, print what arrives
	for result := range results {
		if result.Err != nil {
			fmt.Printf("[%s] error: %v\n", result.Name, result.Err)
			continue
		}

		switch m := result.Metric.(type) {
		case *internal.Uptime:
			fmt.Printf("[uptime] %dh %dm %ds\n", m.Hours, m.Minutes, m.Seconds)
		case *internal.RamDetails:
			fmt.Printf("[ram] total: %d GB, free: %d GB, available: %d GB, used: %d GB\n",
				m.TotalRam, m.FreeRam, m.AvailableRam, m.UsedRam)
		case *internal.CpuPercent:
			fmt.Printf("[cpu] usage: %.2f%%\n", m.Percent)
		}
	}
}
