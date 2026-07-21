package main

import "testing"

func TestPreparationWorkerCount(t *testing.T) {
	for input, want := range map[int]int{-1: 1, 0: 1, 1: 1, 3: 3} {
		if got := preparationWorkerCount(input); got != want {
			t.Fatalf("preparationWorkerCount(%d) = %d, want %d", input, got, want)
		}
	}
}
