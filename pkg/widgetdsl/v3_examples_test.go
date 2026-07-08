package widgetdsl

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
)

func TestWidgetV3GoldenExamplesRenderStableIR(t *testing.T) {
	examples, err := filepath.Glob(filepath.Join("testdata", "v3", "examples", "*.js"))
	if err != nil {
		t.Fatalf("glob examples: %v", err)
	}
	sort.Strings(examples)
	if len(examples) == 0 {
		t.Fatal("expected widget.dsl v3 examples")
	}
	goldenDir := filepath.Join("testdata", "v3", "golden")
	update := os.Getenv("WIDGETDSL_UPDATE_GOLDEN") == "1"
	for _, example := range examples {
		example := example
		t.Run(strings.TrimSuffix(filepath.Base(example), filepath.Ext(example)), func(t *testing.T) {
			page := renderWidgetV3ExampleForTest(t, example)
			data, err := json.MarshalIndent(page, "", "  ")
			if err != nil {
				t.Fatalf("marshal rendered page: %v", err)
			}
			data = append(data, '\n')
			goldenPath := filepath.Join(goldenDir, strings.TrimSuffix(filepath.Base(example), filepath.Ext(example))+".json")
			if update {
				if err := os.WriteFile(goldenPath, data, 0o644); err != nil {
					t.Fatalf("update golden: %v", err)
				}
				return
			}
			want, err := os.ReadFile(goldenPath)
			if err != nil {
				t.Fatalf("read golden %s: %v", goldenPath, err)
			}
			if !jsonEqual(data, want) {
				t.Fatalf("rendered IR differs from golden %s\nset WIDGETDSL_UPDATE_GOLDEN=1 to refresh", goldenPath)
			}
		})
	}
}

func jsonEqual(a []byte, b []byte) bool {
	var left any
	var right any
	if err := json.Unmarshal(a, &left); err != nil {
		return false
	}
	if err := json.Unmarshal(b, &right); err != nil {
		return false
	}
	return reflect.DeepEqual(left, right)
}

func renderWidgetV3ExampleForTest(t *testing.T, path string) any {
	t.Helper()
	source, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read example: %v", err)
	}
	vm := goja.New()
	reg := require.NewRegistry()
	Register(reg)
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
		t.Fatalf("run %s: %v", path, err)
	}
	return normalizeJSONValue(value.Export())
}
