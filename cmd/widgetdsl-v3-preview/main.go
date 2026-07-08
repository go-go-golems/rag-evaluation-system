package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/go-go-golems/rag-evaluation-system/pkg/widgetdsl"
)

type example struct {
	ID    string
	Title string
	Path  string
}

func main() {
	addr := flag.String("addr", "127.0.0.1:8097", "HTTP listen address")
	examplesDir := flag.String("examples", "pkg/widgetdsl/testdata/v3/examples", "directory containing widget.dsl v3 example .js files")
	appDir := flag.String("app", "packages/rag-evaluation-site/app-dist", "built rag-evaluation-site app-dist directory")
	flag.Parse()

	examples, err := loadExamples(*examplesDir)
	must(err)
	if len(examples) == 0 {
		fatalf("no examples found in %s", *examplesDir)
	}
	if _, err := os.Stat(filepath.Join(*appDir, "index.html")); err != nil {
		fatalf("app index not found at %s: %v\nRun: cd packages/rag-evaluation-site && pnpm build:app", filepath.Join(*appDir, "index.html"), err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/widget/pages/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var page any
		var err error
		if id == "" || id == "index" {
			page = indexPage(examples)
		} else if ex, ok := examples[id]; ok {
			page, err = renderExample(ex.Path)
		} else {
			http.Error(w, "unknown page", http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pageMap, _ := page.(map[string]any)
		attachNav(pageMap, examples, id)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(pageMap)
	})
	mux.HandleFunc("POST /api/widget/actions/{name}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"toast": "Preview action accepted", "refresh": false})
	})
	mux.Handle("/", spaHandler(os.DirFS(*appDir)))

	log.Printf("widget.dsl v3 preview listening at http://%s/pages/index", *addr)
	log.Printf("examples: %s", strings.Join(exampleIDs(examples), ", "))
	must(http.ListenAndServe(*addr, mux))
}

func loadExamples(dir string) (map[string]example, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*.js"))
	if err != nil {
		return nil, err
	}
	sort.Strings(files)
	out := map[string]example{}
	for _, file := range files {
		base := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))
		out[base] = example{ID: base, Title: titleFromID(base), Path: file}
	}
	return out, nil
}

func exampleIDs(examples map[string]example) []string {
	ids := make([]string, 0, len(examples))
	for id := range examples {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

func titleFromID(id string) string {
	parts := strings.Split(id, "-")
	if len(parts) > 0 && len(parts[0]) == 2 {
		parts = parts[1:]
	}
	for i, part := range parts {
		if part == "cms" {
			parts[i] = "CMS"
			continue
		}
		if part != "" {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return strings.Join(parts, " ")
}

func indexPage(examples map[string]example) map[string]any {
	children := []any{component("Caption", map[string]any{"tone": "muted"}, text("Select an example from the navigation or links below."))}
	for _, id := range exampleIDs(examples) {
		ex := examples[id]
		children = append(children, component("Button", map[string]any{"variant": "ghost", "action": map[string]any{"kind": "navigate", "to": "/pages/" + id}}, text(ex.Title)))
	}
	return map[string]any{
		"schemaVersion": "0.1.0",
		"id":            "index",
		"title":         "widget.dsl v3 examples",
		"root": component("Stack", map[string]any{"gap": "md"},
			component("SectionBlock", map[string]any{"label": "Examples", "level": 1, "rule": true}, children...),
		),
	}
}

func attachNav(page map[string]any, examples map[string]example, active string) {
	if page == nil {
		return
	}
	items := []any{map[string]any{"id": "index", "label": "Index"}}
	for _, id := range exampleIDs(examples) {
		items = append(items, map[string]any{"id": id, "label": examples[id].Title})
	}
	page["meta"] = map[string]any{"shell": "app", "activeNavItemId": active, "navItems": items, "maxWidth": "wide"}
}

func renderExample(path string) (map[string]any, error) {
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
	data, err := json.Marshal(value.Export())
	if err != nil {
		return nil, err
	}
	var out map[string]any
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func spaHandler(files fs.FS) http.Handler {
	server := http.FileServer(http.FS(files))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}
		if f, err := files.Open(path); err == nil {
			_ = f.Close()
			server.ServeHTTP(w, r)
			return
		}
		r.URL.Path = "/index.html"
		server.ServeHTTP(w, r)
	})
}

func component(componentType string, props map[string]any, children ...any) map[string]any {
	out := map[string]any{"kind": "component", "type": componentType}
	if len(props) > 0 {
		out["props"] = props
	}
	if len(children) > 0 {
		out["children"] = children
	}
	return out
}

func text(value string) map[string]any { return map[string]any{"kind": "text", "text": value} }

func must(err error) {
	if err != nil {
		fatalf("%v", err)
	}
}

func fatalf(format string, args ...any) {
	log.Fatalf(format, args...)
}
