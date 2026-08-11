// Refreshes host-owned prompt TemplateDigest and manifest Digest fields from
// the exact adjacent .txt/.md template. It prints only filenames and digests;
// prompt text, endpoints, credentials, and source data are never emitted.
// Default mode is preview; --write atomically updates only derived digests.
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
	promptsDir := flag.String("prompts-dir", "", "host-owned manifests/prompts directory")
	write := flag.Bool("write", false, "atomically update derived digest fields")
	flag.Parse()
	if *promptsDir == "" {
		panic("--prompts-dir is required")
	}
	entries, err := os.ReadDir(*promptsDir)
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
		path := filepath.Join(*promptsDir, name)
		body, err := os.ReadFile(path)
		if err != nil {
			panic(err)
		}
		var manifest ragcontract.PromptManifest
		if err := json.Unmarshal(body, &manifest); err != nil {
			panic(err)
		}
		textPath, err := promptTextPath(*promptsDir, manifest.PromptID)
		if err != nil {
			panic(err)
		}
		text, err := os.ReadFile(textPath)
		if err != nil {
			panic(err)
		}
		oldTemplate, oldManifest := manifest.TemplateDigest, manifest.Digest
		template, err := ragcontract.Digest(string(text))
		if err != nil {
			panic(err)
		}
		manifest.TemplateDigest = template
		manifest.Digest = ""
		want, err := ragcontract.Digest(manifest)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s template_old=%s template_expected=%s manifest_old=%s manifest_expected=%s changed=%t\n", name, oldTemplate, template, oldManifest, want, oldTemplate != template || oldManifest != want)
		if !*write || (oldTemplate == template && oldManifest == want) {
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

func promptTextPath(dir, promptID string) (string, error) {
	for _, extension := range []string{".txt", ".md"} {
		path := filepath.Join(dir, promptID+extension)
		if _, err := os.Stat(path); err == nil {
			return path, nil
		} else if !os.IsNotExist(err) {
			return "", err
		}
	}
	return "", fmt.Errorf("prompt template missing for %q", promptID)
}
