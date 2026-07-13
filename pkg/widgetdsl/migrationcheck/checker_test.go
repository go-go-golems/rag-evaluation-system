package migrationcheck

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanFileFindsLegacyImportsAndRawComponents(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "page.tsx")
	source := []byte(`import ui from "ui.dsl";
export { thing } from "course.dsl";
const data = require('data.dsl');
const modern = require("widget.dsl");
async function load() { return import("cms.dsl"); }
widget.raw.component("BoardEngine", {});
raw.component("FieldRenderer", {});
`)
	if err := os.WriteFile(path, source, 0o644); err != nil {
		t.Fatal(err)
	}

	findings, err := ScanFile(dir, path)
	if err != nil {
		t.Fatalf("scan file: %v", err)
	}
	want := []Finding{
		{Path: "page.tsx", Line: 1, Kind: "legacy-module-import", Value: "ui.dsl", Text: `import ui from "ui.dsl";`},
		{Path: "page.tsx", Line: 2, Kind: "legacy-module-import", Value: "course.dsl", Text: `export { thing } from "course.dsl";`},
		{Path: "page.tsx", Line: 3, Kind: "legacy-module-import", Value: "data.dsl", Text: `const data = require('data.dsl');`},
		{Path: "page.tsx", Line: 5, Kind: "legacy-module-import", Value: "cms.dsl", Text: `async function load() { return import("cms.dsl"); }`},
		{Path: "page.tsx", Line: 6, Kind: "raw-component-escape-hatch", Value: "raw.component", Text: `widget.raw.component("BoardEngine", {});`},
		{Path: "page.tsx", Line: 7, Kind: "raw-component-escape-hatch", Value: "raw.component", Text: `raw.component("FieldRenderer", {});`},
	}
	if len(findings) != len(want) {
		t.Fatalf("finding count = %d, want %d: %#v", len(findings), len(want), findings)
	}
	for i := range want {
		if findings[i] != want[i] {
			t.Fatalf("finding[%d] = %#v, want %#v", i, findings[i], want[i])
		}
	}
}

func TestScanFileFindsLegacyStructuralShellMetadata(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "page.js")
	source := []byte(`widget.page("Legacy", p => p
  .meta("navItems", items)
  .meta("activeNavItemId", "index")
  .meta("maxWidth", "wide")
  .meta("description", "allowed")
);`)
	if err := os.WriteFile(path, source, 0o644); err != nil {
		t.Fatal(err)
	}

	findings, err := ScanFile(dir, path)
	if err != nil {
		t.Fatalf("scan file: %v", err)
	}
	if len(findings) != 3 {
		t.Fatalf("findings = %#v, want three shell metadata findings", findings)
	}
	for _, finding := range findings {
		if finding.Kind != "legacy-shell-metadata" {
			t.Fatalf("finding = %#v", finding)
		}
	}
}

func TestScanPathsSkipsBuildAndGeneratedDirectories(t *testing.T) {
	dir := t.TempDir()
	goodDir := filepath.Join(dir, "src")
	ignoredDirs := []string{
		filepath.Join(dir, "node_modules"),
		filepath.Join(dir, ".xgoja", "jsverbs"),
		filepath.Join(dir, "internal", "xgojaruntime", "xgoja_embed", "jsverbs"),
	}
	if err := os.MkdirAll(goodDir, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, ignoredDir := range ignoredDirs {
		if err := os.MkdirAll(ignoredDir, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(goodDir, "page.js"), []byte(`const widget = require("widget.dsl");`), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, ignoredDir := range ignoredDirs {
		if err := os.WriteFile(filepath.Join(ignoredDir, "old.js"), []byte(`const ui = require("ui.dsl");`), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	findings, err := ScanPaths(Options{Root: dir, Paths: []string{dir}})
	if err != nil {
		t.Fatalf("scan paths: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings from ignored dirs, got %#v", findings)
	}
}
