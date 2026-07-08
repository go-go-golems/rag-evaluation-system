package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/go-go-golems/rag-evaluation-system/pkg/widgetdsl"
)

func main() {
	examplesDir := flag.String("examples", "pkg/widgetdsl/testdata/v3/examples", "directory containing widget.dsl v3 example .js files")
	outDir := flag.String("out", "pkg/widgetdsl/testdata/v3/rendered", "directory for rendered Widget IR JSON")
	stdout := flag.Bool("stdout", false, "write each rendered page to stdout instead of files")
	flag.Parse()

	files, err := filepath.Glob(filepath.Join(*examplesDir, "*.js"))
	must(err)
	sort.Strings(files)
	if len(files) == 0 {
		fatalf("no .js examples found in %s", *examplesDir)
	}
	if !*stdout {
		must(os.MkdirAll(*outDir, 0o755))
	}
	for _, file := range files {
		value, err := renderExample(file)
		must(err)
		data, err := json.MarshalIndent(value, "", "  ")
		must(err)
		data = append(data, '\n')
		if *stdout {
			fmt.Printf("--- %s ---\n%s", filepath.Base(file), data)
			continue
		}
		name := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file)) + ".json"
		outPath := filepath.Join(*outDir, name)
		must(os.WriteFile(outPath, data, 0o644))
		fmt.Printf("rendered %s -> %s\n", file, outPath)
	}
}

func renderExample(path string) (any, error) {
	source, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	vm := goja.New()
	reg := require.NewRegistry()
	widgetdsl.Register(reg)
	reg.Enable(vm)
	wrapped := `(function(){
` + string(source) + `
if (typeof page !== "undefined") {
  return page && typeof page.toPage === "function" ? page.toPage() : page;
}
throw new Error("example must define const page");
})()`
	value, err := vm.RunString(wrapped)
	if err != nil {
		return nil, fmt.Errorf("run %s: %w", path, err)
	}
	return value.Export(), nil
}

func must(err error) {
	if err != nil {
		fatalf("%v", err)
	}
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
