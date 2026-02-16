package tests

import (
	"fmt"
	"testing"

	"github.com/prasdud/gopher/internal"
)

func TestGetRamDetails(t *testing.T) {
	ram, err := internal.GetRamDetails()
	if err != nil {
		t.Fatalf("GetRamDetails() returned error: %v", err)
	}

	if ram == nil {
		t.Fatal("GetRamDetails() returned nil")
	}

	// TotalRam should be positive on any real system
	if ram.TotalRam <= 0 {
		t.Errorf("expected TotalRam > 0, got %f", ram.TotalRam)
	}

	// Used RAM shouldn't exceed total
	if ram.UsedRam > ram.TotalRam {
		t.Errorf("UsedRam (%f) exceeds TotalRam (%f)", ram.UsedRam, ram.TotalRam)
	}

	// Free RAM shouldn't exceed total
	if ram.FreeRam > ram.TotalRam {
		t.Errorf("FreeRam (%f) exceeds TotalRam (%f)", ram.FreeRam, ram.TotalRam)
	}

	// Available RAM shouldn't exceed total
	if ram.AvailableRam > ram.TotalRam {
		t.Errorf("AvailableRam (%f) exceeds TotalRam (%f)", ram.AvailableRam, ram.TotalRam)
	}

	// Print for visual check
	fmt.Printf("Total: %f GB, Used: %f GB, Free: %f GB, Available: %f GB\n",
		ram.TotalRam, ram.UsedRam, ram.FreeRam, ram.AvailableRam)
}

func TestRamDetailsMetricName(t *testing.T) {
	r := &internal.RamDetails{}
	if got := r.MetricName(); got != "ram" {
		t.Errorf("MetricName() = %q, want %q", got, "ram")
	}
}
