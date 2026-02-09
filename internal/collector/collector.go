package collector

import "time"

// Metric is the common interface all metric results must satisfy.
type Metric interface {
	MetricName() string
}

// Collector is the interface every metric collector implements.
type Collector interface {
	Name() string
	Interval() time.Duration
	Collect() (Metric, error)
}
