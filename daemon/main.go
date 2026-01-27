package main

import (
	"fmt"

	"github.com/prasdud/gopher/internal"
)

func main() {
	uptime, err := internal.GetUptime()
	ramDetails, err := internal.GetRamDetails()
	cpuUsage, err := internal.GetCpuUsage()
	cpuPercent, err := internal.GetCpuPercent()

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Uptime: %d hours, %d minutes, %d seconds\n", uptime.Hours, uptime.Minutes, uptime.Seconds)
	fmt.Printf("Total RAM: %d GB, Free RAM: %d GB, Available RAM: %d GB, Used RAM: %d GB\n", ramDetails.TotalRam, ramDetails.FreeRam, ramDetails.AvailableRam, ramDetails.UsedRam)
	fmt.Printf("CPU Idle: %d, CPU Total: %d\n", cpuUsage.Idle, cpuUsage.Total)
	fmt.Printf("CPU Usage Percent: %.2f%%\n", cpuPercent)
}
