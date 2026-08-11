// Refreshes only the self-referentially excluded Digest field in host-owned
// model manifests. It prints filenames and digest identities, never endpoints,
// credentials, request bodies, or source data. Default mode is preview; pass
// --write only after reviewing the preview.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
)

func main() {
	modelsDir := flag.String("models-dir", "", "host-owned manifests/models directory")
	write := flag.Bool("write", false, "atomically update digest fields")
	flag.Parse()
	if *modelsDir == "" {
		panic("--models-dir is required")
	}
	entries, err := os.ReadDir(*modelsDir)
	if err != nil {
		panic(err)
	}
	var names []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)
	for _, name := range names {
		path := filepath.Join(*modelsDir, name)
		body, err := os.ReadFile(path)
		if err != nil {
			panic(err)
		}
		var manifest ragcontract.ModelManifest
		if err := json.Unmarshal(body, &manifest); err != nil {
			panic(err)
		}
		old := manifest.Digest
		manifest.Digest = ""
		want, err := ragcontract.Digest(manifest)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s old=%s expected=%s changed=%t\n", name, old, want, old != want)
		if !*write || old == want {
			continue
		}
		manifest.Digest = want
		updated, err := json.Marshal(manifest)
		if err != nil {
			panic(err)
		}
		temporary := path + ".tmp"
		if err := os.WriteFile(temporary, append(updated, '\n'), 0o600); err != nil {
			panic(err)
		}
		if err := os.Rename(temporary, path); err != nil {
			panic(err)
		}
	}
}
