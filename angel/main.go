package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/prasdud/gopher/internal/collector"
	"github.com/prasdud/gopher/tui"
)

func main() {
	// 1. Load theme (fall back to default if file not found)
	theme, err := tui.LoadTheme("styles/default.json")
	if err != nil {
		theme = tui.DefaultTheme()
	}

	// 2. Shared channel for collector results
	results := make(chan collector.CollectorResult)

	// 3. Quit channel to signal goroutines to stop
	quit := make(chan struct{})

	// 4. Register and start collectors
	collectors := []collector.Collector{
		&collector.UptimeCollector{},
		&collector.RamCollector{},
		&collector.CpuCollector{},
		&collector.IPv4Collector{},
	}

	for _, c := range collectors {
		go func(c collector.Collector) {
			// run once at start to grab fresh metrics
			metric, err := c.Collect()
			results <- collector.CollectorResult{
				Name:   c.Name(),
				Metric: metric,
				Err:    err,
			}

			// start the periodic ticker later
			ticker := time.NewTicker(c.Interval())
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					metric, err := c.Collect()
					select {
					case results <- collector.CollectorResult{
						Name:   c.Name(),
						Metric: metric,
						Err:    err,
					}:
					case <-quit:
						return
					}
				case <-quit:
					return
				}
			}
		}(c)
	}

	// 5. Create and run the TUI
	model := tui.NewModel(results, theme)
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// 6. Signal goroutines to stop
	close(quit)
}
