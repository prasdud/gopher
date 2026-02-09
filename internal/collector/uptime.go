package collector

import (
	"time"

	"github.com/prasdud/gopher/internal"
)

type UptimeCollector struct{}

func (u *UptimeCollector) Name() string {
	return "uptime"
}

func (u *UptimeCollector) Interval() time.Duration {
	return 5 * time.Second
}

func (u *UptimeCollector) Collect() (Metric, error) {
	return internal.GetUptime()
}
