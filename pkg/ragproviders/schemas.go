package ragproviders

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	jsonschema "github.com/santhosh-tekuri/jsonschema/v6"
)

type FileSchemaRegistry struct {
	mu     sync.RWMutex
	raw    map[string]json.RawMessage
	schema map[string]*jsonschema.Schema
}

func LoadSchemaRegistry(dir string) (*FileSchemaRegistry, error) {
	if dir == "" {
		return nil, fmt.Errorf("RAG_SCHEMA_DIRECTORY_REQUIRED")
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("RAG_SCHEMA_DIRECTORY: %w", err)
	}
	registry := &FileSchemaRegistry{raw: map[string]json.RawMessage{}, schema: map[string]*jsonschema.Schema{}}
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			return nil, fmt.Errorf("RAG_SCHEMA_READ: %w", err)
		}
		var document any
		if err := json.Unmarshal(data, &document); err != nil {
			return nil, fmt.Errorf("RAG_SCHEMA_DECODE: %w", err)
		}
		fileID := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
		location := "mem://rag-schema/" + fileID
		compiler := jsonschema.NewCompiler()
		if err := compiler.AddResource(location, document); err != nil {
			return nil, fmt.Errorf("RAG_SCHEMA_REGISTER: %w", err)
		}
		compiled, err := compiler.Compile(location)
		if err != nil {
			return nil, fmt.Errorf("RAG_SCHEMA_COMPILE: %w", err)
		}
		ids := []string{fileID}
		if schemaID, ok := document.(map[string]any)["$id"].(string); ok && schemaID != "" {
			ids = append(ids, schemaID)
		}
		for _, id := range ids {
			if _, exists := registry.schema[id]; exists {
				return nil, fmt.Errorf("RAG_SCHEMA_DUPLICATE")
			}
			registry.raw[id] = append([]byte(nil), data...)
			registry.schema[id] = compiled
		}
	}
	if len(registry.schema) == 0 {
		return nil, fmt.Errorf("RAG_SCHEMA_EMPTY")
	}
	return registry, nil
}
func (r *FileSchemaRegistry) Validate(name string, document json.RawMessage) error {
	if r == nil {
		return fmt.Errorf("RAG_OUTPUT_SCHEMA_VALIDATOR_UNAVAILABLE")
	}
	r.mu.RLock()
	schema, ok := r.schema[name]
	r.mu.RUnlock()
	if !ok {
		return fmt.Errorf("RAG_OUTPUT_SCHEMA_MISSING")
	}
	var value any
	dec := json.NewDecoder(strings.NewReader(string(document)))
	dec.UseNumber()
	if err := dec.Decode(&value); err != nil {
		return fmt.Errorf("RAG_OUTPUT_SCHEMA_JSON: %w", err)
	}
	var extra any
	if err := dec.Decode(&extra); err != io.EOF {
		if err == nil {
			return fmt.Errorf("RAG_OUTPUT_SCHEMA_JSON_MULTIPLE_VALUES")
		}
		return fmt.Errorf("RAG_OUTPUT_SCHEMA_JSON: %w", err)
	}
	if err := schema.Validate(value); err != nil {
		return fmt.Errorf("RAG_OUTPUT_SCHEMA_INVALID: %w", err)
	}
	return nil
}
func (r *FileSchemaRegistry) Raw(name string) (json.RawMessage, error) {
	if r == nil {
		return nil, fmt.Errorf("RAG_OUTPUT_SCHEMA_REGISTRY_UNAVAILABLE")
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	value, ok := r.raw[name]
	if !ok {
		return nil, fmt.Errorf("RAG_OUTPUT_SCHEMA_MISSING")
	}
	return append([]byte(nil), value...), nil
}
