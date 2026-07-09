package migrationcheck

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_javascript "github.com/tree-sitter/tree-sitter-javascript/bindings/go"
	tree_sitter_typescript "github.com/tree-sitter/tree-sitter-typescript/bindings/go"
)

var LegacyModules = []string{
	"ui.dsl",
	"data.dsl",
	"data.v2.dsl",
	"context_window.dsl",
	"course.dsl",
	"cms.dsl",
}

var sourceSuffixes = map[string]bool{
	".js": true, ".jsx": true, ".mjs": true, ".cjs": true,
	".ts": true, ".tsx": true, ".mts": true, ".cts": true,
}

var ignoredDirs = map[string]bool{
	".git": true, "node_modules": true, "dist": true, "app-dist": true,
	"storybook-static": true, "coverage": true, ".next": true, ".turbo": true,
}

// Finding is one parser-backed migration finding.
type Finding struct {
	Path  string `json:"path"`
	Line  int    `json:"line"`
	Kind  string `json:"kind"`
	Value string `json:"value"`
	Text  string `json:"text"`
}

// Options controls source discovery.
type Options struct {
	Root  string
	Paths []string
}

// DefaultPaths returns known first-party widget source locations when they exist.
func DefaultPaths(root string) []string {
	candidates := []string{
		filepath.Join(root, "go-go-course", "cmd", "go-go-course", "lib", "pages"),
		filepath.Join(root, "..", "go-go-course", "cmd", "go-go-course", "lib", "pages"),
		filepath.Join(root, "pkg", "widgetdsl", "testdata", "v3", "examples"),
		filepath.Join(root, "examples"),
	}
	ret := []string{}
	seen := map[string]bool{}
	for _, candidate := range candidates {
		abs, err := filepath.Abs(candidate)
		if err != nil {
			continue
		}
		if seen[abs] {
			continue
		}
		if _, err := os.Stat(abs); err == nil {
			seen[abs] = true
			ret = append(ret, abs)
		}
	}
	return ret
}

// ScanPaths discovers source files and scans them using tree-sitter JavaScript/TypeScript grammars.
func ScanPaths(opts Options) ([]Finding, error) {
	root := opts.Root
	if root == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		root = cwd
	}
	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	paths := opts.Paths
	if len(paths) == 0 {
		paths = DefaultPaths(rootAbs)
	}

	files, err := DiscoverSourceFiles(paths)
	if err != nil {
		return nil, err
	}
	findings := []Finding{}
	for _, file := range files {
		fileFindings, scanErr := ScanFile(rootAbs, file)
		if scanErr != nil {
			return nil, scanErr
		}
		findings = append(findings, fileFindings...)
	}
	SortFindings(findings)
	return findings, nil
}

// DiscoverSourceFiles returns JS/TS source files under paths, skipping generated/build directories.
func DiscoverSourceFiles(paths []string) ([]string, error) {
	seen := map[string]bool{}
	files := []string{}
	add := func(path string) error {
		abs, err := filepath.Abs(path)
		if err != nil {
			return err
		}
		if seen[abs] {
			return nil
		}
		seen[abs] = true
		files = append(files, abs)
		return nil
	}
	for _, input := range paths {
		abs, err := filepath.Abs(input)
		if err != nil {
			return nil, err
		}
		info, err := os.Stat(abs)
		if err != nil {
			continue
		}
		if !info.IsDir() {
			if isSourceFile(abs) {
				if err := add(abs); err != nil {
					return nil, err
				}
			}
			continue
		}
		err = filepath.WalkDir(abs, func(path string, d fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if d.IsDir() {
				if ignoredDirs[d.Name()] {
					return filepath.SkipDir
				}
				return nil
			}
			if isSourceFile(path) {
				return add(path)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	sort.Strings(files)
	return files, nil
}

func isSourceFile(path string) bool {
	return sourceSuffixes[strings.ToLower(filepath.Ext(path))]
}

// ScanFile parses one JS/TS source file and returns legacy import/raw escape-hatch findings.
func ScanFile(root, filename string) ([]Finding, error) {
	source, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	parser := tree_sitter.NewParser()
	defer parser.Close()
	if err := parser.SetLanguage(languageForPath(filename)); err != nil {
		return nil, fmt.Errorf("configure parser for %s: %w", filename, err)
	}
	tree := parser.Parse(source, nil)
	if tree == nil {
		return nil, fmt.Errorf("parse %s: parser returned nil tree", filename)
	}
	defer tree.Close()
	rootNode := tree.RootNode()
	if rootNode == nil {
		return nil, fmt.Errorf("parse %s: parser returned nil root", filename)
	}

	displayPath := displayPath(root, filename)
	lines := splitLines(source)
	legacy := map[string]bool{}
	for _, moduleName := range LegacyModules {
		legacy[moduleName] = true
	}

	findings := []Finding{}
	add := func(n *tree_sitter.Node, kind, value string) {
		line := int(n.StartPosition().Row) + 1
		findings = append(findings, Finding{
			Path:  displayPath,
			Line:  line,
			Kind:  kind,
			Value: value,
			Text:  lineText(lines, line),
		})
	}

	walkTree(rootNode, func(n *tree_sitter.Node) {
		switch n.Kind() {
		case "import_statement":
			if specifier, ok := staticSourceField(n, source); ok && legacy[specifier] {
				add(n, "legacy-module-import", specifier)
			}
		case "export_statement":
			if specifier, ok := staticSourceField(n, source); ok && legacy[specifier] {
				add(n, "legacy-module-import", specifier)
			}
		case "call_expression":
			functionNode := n.ChildByFieldName("function")
			if functionNode == nil {
				return
			}
			functionText := strings.TrimSpace(functionNode.Utf8Text(source))
			if functionText == "require" || functionText == "import" {
				args := n.ChildByFieldName("arguments")
				if specifier, ok := firstArgumentStringLiteral(args, source); ok && legacy[specifier] {
					add(n, "legacy-module-import", specifier)
				}
				return
			}
			if functionText == "raw.component" || functionText == "widget.raw.component" {
				add(n, "raw-component-escape-hatch", "raw.component")
			}
		}
	})
	SortFindings(findings)
	return findings, nil
}

func languageForPath(filename string) *tree_sitter.Language {
	switch strings.ToLower(filepath.Ext(filename)) {
	case ".ts", ".mts", ".cts":
		return tree_sitter.NewLanguage(tree_sitter_typescript.LanguageTypescript())
	case ".tsx":
		return tree_sitter.NewLanguage(tree_sitter_typescript.LanguageTSX())
	default:
		return tree_sitter.NewLanguage(tree_sitter_javascript.Language())
	}
}

func staticSourceField(n *tree_sitter.Node, source []byte) (string, bool) {
	if n == nil {
		return "", false
	}
	sourceNode := n.ChildByFieldName("source")
	if sourceNode == nil || sourceNode.Kind() != "string" {
		return "", false
	}
	return unquoteTreeSitterString(sourceNode.Utf8Text(source)), true
}

func firstArgumentStringLiteral(args *tree_sitter.Node, source []byte) (string, bool) {
	if args == nil {
		return "", false
	}
	cursor := args.Walk()
	defer cursor.Close()
	children := args.NamedChildren(cursor)
	if len(children) == 0 {
		return "", false
	}
	first := children[0]
	if first.Kind() != "string" {
		return "", false
	}
	return unquoteTreeSitterString(first.Utf8Text(source)), true
}

func walkTree(n *tree_sitter.Node, visit func(*tree_sitter.Node)) {
	if n == nil {
		return
	}
	visit(n)
	cursor := n.Walk()
	defer cursor.Close()
	for _, child := range n.NamedChildren(cursor) {
		childCopy := child
		walkTree(&childCopy, visit)
	}
}

func unquoteTreeSitterString(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 {
		quote := s[0]
		if (quote == '\'' || quote == '"' || quote == '`') && s[len(s)-1] == quote {
			return s[1 : len(s)-1]
		}
	}
	return s
}

func displayPath(root, filename string) string {
	if root == "" {
		return filepath.ToSlash(filename)
	}
	if rel, err := filepath.Rel(root, filename); err == nil && !strings.HasPrefix(rel, "..") {
		return filepath.ToSlash(rel)
	}
	return filepath.ToSlash(filename)
}

func splitLines(source []byte) []string {
	return strings.Split(string(source), "\n")
}

func lineText(lines []string, line int) string {
	if line <= 0 || line > len(lines) {
		return ""
	}
	return strings.TrimSpace(lines[line-1])
}

func SortFindings(findings []Finding) {
	sort.Slice(findings, func(i, j int) bool {
		if findings[i].Path != findings[j].Path {
			return findings[i].Path < findings[j].Path
		}
		if findings[i].Line != findings[j].Line {
			return findings[i].Line < findings[j].Line
		}
		if findings[i].Kind != findings[j].Kind {
			return findings[i].Kind < findings[j].Kind
		}
		return findings[i].Value < findings[j].Value
	})
}

func FindingsJSON(findings []Finding) ([]byte, error) {
	return json.MarshalIndent(findings, "", "  ")
}
