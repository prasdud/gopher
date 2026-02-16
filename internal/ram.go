package internal

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type RamDetails struct {
	TotalRam      float64
	FreeRam       float64
	AvailableRam  float64
	UsedRam       float64
	TotalSwap     float64
	UsedSwap      float64
	AvailableSwap float64
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
	swapTotalKB := memInfo["SwapTotal"]
	swapFreeKB := memInfo["SwapFree"]

	// Used = Total - Free - Buffers - Cached - SReclaimable (matches `free` command)
	usedKB := totalKB - freeKB - buffersKB - cachedKB - sReclaimableKB
	usedSwapKB := swapTotalKB - swapFreeKB

	return &RamDetails{
		TotalRam:      float64(totalKB) / 1024 / 1024,     // GB
		FreeRam:       float64(freeKB) / 1024 / 1024,      // GB
		AvailableRam:  float64(availableKB) / 1024 / 1024, // GB
		UsedRam:       float64(usedKB) / 1024 / 1024,      // GB
		TotalSwap:     float64(swapTotalKB) / 1024 / 1024, // GB
		UsedSwap:      float64(usedSwapKB) / 1024 / 1024,  // GB
		AvailableSwap: float64(swapFreeKB) / 1024 / 1024,  // GB
	}, nil
}
