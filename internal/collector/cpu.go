package collector

import (
	"time"

	"github.com/prasdud/gopher/internal"
)

type CpuCollector struct{}

func (c *CpuCollector) Name() string {
	return "cpu"
}

func (c *CpuCollector) Interval() time.Duration {
	return 1 * time.Second
}

func (c *CpuCollector) Collect() (Metric, error) {
	percent, err := internal.GetCpuPercent()
	if err != nil {
		return nil, err
	}
	return &internal.CpuPercent{Percent: percent}, nil
}
