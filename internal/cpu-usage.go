package internal

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	idleFieldIndex   = 3
	iowaitFieldIndex = 4
)

type CpuUsage struct {
	Idle  uint64
	Total uint64
}

func (c *CpuUsage) MetricName() string {
	return "cpu"
}

type CpuPercent struct {
	Percent float64
}

func (c *CpuPercent) MetricName() string {
	return "cpu"
}

func GetCpuUsage() (*CpuUsage, error) {
	f, err := os.Open("/proc/stat")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var idle uint64
	var total uint64

	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}
		return nil, errors.New("empty /proc/stat")
	}

	fields := strings.Fields(scanner.Text())
	if len(fields) < 5 {
		return nil, errors.New("invalid /proc/stat format")
	}
	fields = fields[1:] // skip "cpu" prefix

	for i, v := range fields {
		val, e := strconv.ParseUint(v, 10, 64)
		if e != nil {
			return nil, e
		}
		total += val
		if i == idleFieldIndex || i == iowaitFieldIndex {
			idle += val
		}
	}
	return &CpuUsage{
		Idle:  idle,
		Total: total,
	}, nil
}

// GetCpuPercent returns current CPU usage as a percentage (0-100).
// It takes two samples 200ms apart to calculate the delta.
func GetCpuPercent() (float64, error) {
	first, err := GetCpuUsage()
	if err != nil {
		return 0, err
	}

	time.Sleep(200 * time.Millisecond)

	second, err := GetCpuUsage()
	if err != nil {
		return 0, err
	}

	idleDelta := second.Idle - first.Idle
	totalDelta := second.Total - first.Total

	if totalDelta == 0 {
		return 0, nil
	}

	cpuPercent := (1.0 - float64(idleDelta)/float64(totalDelta)) * 100
	return cpuPercent, nil
}
