package tests

import (
	"fmt"
	"testing"

	"github.com/prasdud/gopher/internal"
)

func TestGetIPv4(t *testing.T) {
	ip, err := internal.GetIPv4()
	if err != nil {
		t.Fatalf("GetIPv4() returned error: %v", err)
	}

	if ip == nil {
		t.Fatal("GetIPv4() returned nil")
	}

	if ip.Address == "" {
		t.Errorf("expected Address to be non-empty, got %q", ip.Address)
	}

	fmt.Printf("IPv4: %s\n", ip.Address)
}
