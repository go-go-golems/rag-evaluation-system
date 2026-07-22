package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

const generationAuthoritySchema = "rag-ttc-v3-generation-authority/v1"

type generationAuthorityState struct {
	SchemaVersion string `json:"schemaVersion"`
	Maximum       int    `json:"maximumGenerationRequests"`
	Prior         int    `json:"priorGenerationRequests"`
	Admitted      int    `json:"admittedGenerationRequests"`
}

type generationAdmission struct {
	mu    sync.Mutex
	path  string
	state generationAuthorityState
}

func newGenerationAdmission(path string, maximum, prior int) (*generationAdmission, error) {
	if path == "" || maximum < 1 || prior < 0 || prior > maximum {
		return nil, fmt.Errorf("RAG_SWEEP_GENERATION_AUTHORITY_INVALID")
	}
	a := &generationAdmission{path: path, state: generationAuthorityState{SchemaVersion: generationAuthoritySchema, Maximum: maximum, Prior: prior, Admitted: prior}}
	if err := a.persist(a.state); err != nil {
		return nil, err
	}
	return a, nil
}

func (a *generationAdmission) Admit() error {
	if a == nil {
		return fmt.Errorf("RAG_SWEEP_GENERATION_AUTHORITY_INVALID")
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.state.Admitted >= a.state.Maximum {
		return fmt.Errorf("RAG_SWEEP_GENERATION_REQUEST_CEILING")
	}
	next := a.state
	next.Admitted++
	if err := a.persist(next); err != nil {
		return err
	}
	a.state = next
	return nil
}

func (a *generationAdmission) persist(state generationAuthorityState) error {
	body, err := json.Marshal(state)
	if err != nil {
		return err
	}
	body = append(body, '\n')
	dir := filepath.Dir(a.path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(dir, ".generation-authority-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer func() { _ = os.Remove(tmpName) }()
	if err := tmp.Chmod(0o600); err != nil {
		_ = tmp.Close()
		return err
	}
	if _, err := tmp.Write(body); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Rename(tmpName, a.path); err != nil {
		return err
	}
	directory, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer func() { _ = directory.Close() }()
	return directory.Sync()
}
