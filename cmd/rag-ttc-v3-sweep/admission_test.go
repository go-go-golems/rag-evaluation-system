package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
)

func TestGenerationAdmissionPersistsPriorAndHardCeiling(t *testing.T) {
	path := filepath.Join(t.TempDir(), "authority.json")
	a, err := newGenerationAdmission(path, 3, 1)
	if err != nil {
		t.Fatal(err)
	}
	if err := a.Admit(); err != nil {
		t.Fatal(err)
	}
	if err := a.Admit(); err != nil {
		t.Fatal(err)
	}
	if err := a.Admit(); err == nil {
		t.Fatal("admission above hard ceiling succeeded")
	}
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var state generationAuthorityState
	if err := json.Unmarshal(body, &state); err != nil {
		t.Fatal(err)
	}
	if state.SchemaVersion != generationAuthoritySchema || state.Maximum != 3 || state.Prior != 1 || state.Admitted != 3 {
		t.Fatalf("state = %#v", state)
	}
}

func TestGenerationAdmissionSerializesConcurrentRequests(t *testing.T) {
	const maximum = 20
	a, err := newGenerationAdmission(filepath.Join(t.TempDir(), "authority.json"), maximum, 0)
	if err != nil {
		t.Fatal(err)
	}
	var admitted atomic.Int64
	var wg sync.WaitGroup
	for range 40 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if a.Admit() == nil {
				admitted.Add(1)
			}
		}()
	}
	wg.Wait()
	if got := admitted.Load(); got != maximum {
		t.Fatalf("admitted = %d, want %d", got, maximum)
	}
}
