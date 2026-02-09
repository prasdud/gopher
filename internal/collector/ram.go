package collector

import (
	"time"

	"github.com/prasdud/gopher/internal"
)

type RamCollector struct{}

func (r *RamCollector) Name() string {
	return "ram"
}

func (r *RamCollector) Interval() time.Duration {
	return 2 * time.Second
}

func (r *RamCollector) Collect() (Metric, error) {
	return internal.GetRamDetails()
}
