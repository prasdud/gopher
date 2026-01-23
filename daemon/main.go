package main

import (
	"fmt"

	"github.com/prasdud/gopher/internal"
)

func main() {
	uptime, err := internal.GetUptime()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Uptime: %d hours, %d minutes, %d seconds\n", uptime.Hours, uptime.Minutes, uptime.Seconds)
}
