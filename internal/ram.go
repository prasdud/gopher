package internal

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type RamDetails struct {
	TotalRam     int64
	FreeRam      int64
	AvailableRam int64
	UsedRam      int64
}

func (r *RamDetails) MetricName() string {
	return "ram"
}

func GetRamDetails() (*RamDetails, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	memInfo := make(map[string]int64)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		key := strings.TrimSuffix(fields[0], ":")
		value, err := strconv.ParseInt(fields[1], 10, 64)
		if err != nil {
			continue
		}
		memInfo[key] = value // values are in kB
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	totalKB := memInfo["MemTotal"]
	freeKB := memInfo["MemFree"]
	availableKB := memInfo["MemAvailable"]
	buffersKB := memInfo["Buffers"]
	cachedKB := memInfo["Cached"]
	sReclaimableKB := memInfo["SReclaimable"]

	// Used = Total - Free - Buffers - Cached - SReclaimable (matches `free` command)
	usedKB := totalKB - freeKB - buffersKB - cachedKB - sReclaimableKB

	return &RamDetails{
		TotalRam:     totalKB / 1024 / 1024,     // GB
		FreeRam:      freeKB / 1024 / 1024,      // GB
		AvailableRam: availableKB / 1024 / 1024, // GB
		UsedRam:      usedKB / 1024 / 1024,      // GB
	}, nil
}
