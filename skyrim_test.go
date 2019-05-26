package main

import (
	"testing"

	"github.com/gganley/skyrim/internal/skyrim"
)

func TestSingleCycle(t *testing.T) {
	initial := skyrim.SingleCycle(skyrim.SetupCycle())
	for i := 0; i < 100; i++ {
		j := skyrim.SingleCycle(skyrim.SetupCycle())
		for i, _ := range initial {
			if initial[i] != j[i] {
				t.Errorf("Expected: %v, Got: %v", initial, j)
			}
		}
	}
}

func TestFullCycle(t *testing.T) {
	initial := skyrim.FullCycle()
	for i := 0; i < 10; i++ {
		if j := skyrim.FullCycle(); initial != j {
			t.Errorf("Expected: %v, Got: %v", initial, j)
		}
	}
}
