package collector

import (
	"time"

	"github.com/prasdud/gopher/internal"
)

type IPv4Collector struct{}

func (u *IPv4Collector) Name() string {
	return "ipv4"
}

func (u *IPv4Collector) Interval() time.Duration {
	return 60 * time.Minute
}

func (u *IPv4Collector) Collect() (Metric, error) {
	return internal.GetIPv4()
}
