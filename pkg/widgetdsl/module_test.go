package widgetdsl

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/go-go-golems/go-go-goja/pkg/engine"
)

func TestSplitModulesExportExpectedHelpersAndOmitCrossDomainHelpers(t *testing.T) {
	vm := goja.New()
	reg := require.NewRegistry()
	Register(reg)
	reg.Enable(vm)

	value, err := vm.RunString(`
		const ui = require("ui.dsl");
		const data = require("data.dsl");
		const contextWindow = require("context_window.dsl");
		const course = require("course.dsl");
		({
			uiPage: typeof ui.page,
			uiPanel: typeof ui.panel,
			uiFormPanel: typeof ui.formPanel,
			uiTextareaInput: typeof ui.textareaInput,
			uiDataTable: typeof ui.dataTable,
			dataTable: typeof data.dataTable,
			dataCellField: typeof data.cell.field,
			dataPage: typeof data.page,
			contextDiagramPanel: typeof contextWindow.contextDiagramPanel,
			contextStyleSwatch: typeof contextWindow.contextStyleSwatch,
			contextVisualStyle: typeof contextWindow.visualStyle,
			contextStyleSet: typeof contextWindow.styleSet,
			contextPart: typeof contextWindow.contextPart,
			contextStudioNavIconFromContext: typeof contextWindow.contextStudioNavIcon,
			courseStudioNavIcon: typeof course.contextStudioNavIcon,
			courseStudioShell: typeof course.courseStudioShell,
		});
	`)
	if err != nil {
		t.Fatalf("require split modules: %v", err)
	}
	got := value.Export().(map[string]any)
	wantFunctions := []string{"uiPage", "uiPanel", "uiFormPanel", "uiTextareaInput", "dataTable", "dataCellField", "contextDiagramPanel", "contextStyleSwatch", "contextVisualStyle", "contextStyleSet", "contextPart", "courseStudioNavIcon", "courseStudioShell"}
	for _, name := range wantFunctions {
		if got[name] != "function" {
			t.Fatalf("%s export = %#v, want function (all: %#v)", name, got[name], got)
		}
	}
	wantUndefined := []string{"uiDataTable", "dataPage", "contextStudioNavIconFromContext"}
	for _, name := range wantUndefined {
		if got[name] != "undefined" {
			t.Fatalf("%s export = %#v, want undefined (all: %#v)", name, got[name], got)
		}
	}
}

func TestOldBucketModulesAreAbsent(t *testing.T) {
	vm := goja.New()
	reg := require.NewRegistry()
	Register(reg)
	reg.Enable(vm)

	value, err := vm.RunString(`
		function canRequire(name) {
			try { require(name); return true; } catch (error) { return false; }
		}
		({ widget: canRequire("widget.dsl"), rag: canRequire("rag.dsl") });
	`)
	if err != nil {
		t.Fatalf("check old modules: %v", err)
	}
	got := value.Export().(map[string]any)
	if got["widget"] != false || got["rag"] != false {
		t.Fatalf("old bucket modules should be absent, got %#v", got)
	}
}

func TestBuildsWidgetIRAcrossSplitModules(t *testing.T) {
	vm := goja.New()
	reg := require.NewRegistry()
	Register(reg)
	reg.Enable(vm)

	value, err := vm.RunString(`
		const ui = require("ui.dsl");
		const data = require("data.dsl");
		const contextWindow = require("context_window.dsl");
		const course = require("course.dsl");
		const styleSet = contextWindow.styleSet({ legend: [contextWindow.legendItem("prompt", "Prompt")], styles: { prompt: contextWindow.visualStyle({ pattern: "checker", fill: "#dde6f2", line: "#4f74a8" }) } });
		const snapshot = { id: "ctx", title: "Window", limit: 1000, parts: [contextWindow.contextPart("p", "Prompt", "prompt", 300)] };
		const slide = { id: "s1", number: "01", title: "Window", view: "budget", snapshotId: "ctx", notes: ["Budget"] };
		const page = ui.page({
			id: "split",
			title: "Split modules",
			sections: [
				ui.formPanel({ title: "Settings", method: "post", formAction: "/settings", submitLabel: "Save" },
					ui.formRow({
						label: "Display name",
						required: true,
						hint: "Shown on uploads",
						control: ui.textInput({ name: "displayName", defaultValue: "Manuel", readOnly: false, maxLength: 40 })
					})
				),
				ui.panel({ title: "Table" }, data.dataTable({
					rows: [{ id: "a", title: "Alpha", status: "done" }],
					getRowKey: "id",
					columns: [
						{ id: "title", header: "Title", cell: data.cell.field("title") },
						{ id: "status", header: "Status", cell: data.cell.status("status", { icon: true }) }
					]
				})),
				contextWindow.contextDiagramPanel({ snapshot, styleSet, initialView: "budget" }),
				course.courseStudioShell({
					sections: [{ id: "course", label: "Course", items: [{ id: "slides", label: "Slides", icon: course.contextStudioNavIcon({ id: "slides" }) }] }],
					activeItemId: "slides",
					title: "Studio"
				}, course.courseSlidePanel({ slide, snapshot, index: 0, total: 1 }))
			]
		});
		JSON.stringify(page);
	`)
	if err != nil {
		t.Fatalf("build split-module page: %v", err)
	}
	var decoded map[string]any
	if err := json.Unmarshal([]byte(value.String()), &decoded); err != nil {
		t.Fatalf("split-module page is not JSON: %v\n%s", err, value.String())
	}
	assertString(t, decoded, "id", "split")
	root := decoded["root"].(map[string]any)
	children := root["children"].([]any)
	assertString(t, children[0].(map[string]any), "type", "FormPanel")
	assertString(t, children[1].(map[string]any), "type", "Panel")
	assertString(t, children[2].(map[string]any), "type", "ContextDiagramPanel")
	assertString(t, children[3].(map[string]any), "type", "CourseStudioShell")
}

func TestSplitModuleRecipesAreJSONSerializable(t *testing.T) {
	vm := goja.New()
	reg := require.NewRegistry()
	Register(reg)
	reg.Enable(vm)

	value, err := vm.RunString(`
		const ui = require("ui.dsl");
		const data = require("data.dsl");
		const contextWindow = require("context_window.dsl");
		const course = require("course.dsl");
		const rows = [{ id: 1, name: "Alpha", status: "running" }];
		const styleSet = contextWindow.styleSet({ legend: [contextWindow.legendItem("prompt", "Prompt")], styles: { prompt: contextWindow.visualStyle({ pattern: "checker", fill: "#dde6f2", line: "#4f74a8" }) } });
		const snapshot = { id: "ctx", title: "Window", limit: 1000, parts: [contextWindow.contextPart("p", "Prompt", "prompt", 300)] };
		const transcript = { title: "Session", messages: [{ id: "m1", role: "user", text: "hello" }], annotations: [] };
		const slide = { id: "s1", number: "01", title: "Window", view: "budget", snapshotId: "ctx", notes: ["Budget"] };
		const bundle = { intro: "Docs", docs: [{ id: "d1", title: "Guide", file: "guide.md", format: "markdown", description: "Guide", body: "# Guide" }] };
		const sections = [{ id: "course", label: "Course", items: [{ id: "slides", label: "Slides" }] }];
		const page = ui.page({ id: "recipes", title: "Recipes", sections: [
			ui.recipes.metrics({ items: [{ label: "Total", value: rows.length, status: "ready" }] }),
			ui.recipes.actionToolbar({ title: "Controls", actions: [{ label: "Add", action: ui.action.server("add-query") }] }),
			data.recipes.masterDetailTable({
				rows,
				columns: [{ id: "name", header: "Name", cell: data.cell.field("name") }],
				selectedKey: 1,
				detail: row => ui.panel({ title: "Selected" }, row.name)
			}),
			contextWindow.recipes.contextDiagram({ snapshot, styleSet, view: "budget" }),
			contextWindow.recipes.annotatedTranscript({ transcript, onAnnotationSelect: contextWindow.action.server("select-annotation") }),
			course.recipes.courseStudio({ sections, activeItemId: "slides", main: course.recipes.courseSlide({ slide, snapshot, index: 0, total: 1 }) }),
			course.recipes.handout({ bundle, selectedDocumentId: "d1", onSelect: course.action.server("select-doc") })
		]});
		JSON.stringify(page);
	`)
	if err != nil {
		t.Fatalf("build split-module recipes: %v", err)
	}
	var decoded map[string]any
	if err := json.Unmarshal([]byte(value.String()), &decoded); err != nil {
		t.Fatalf("split-module recipe page is not JSON: %v\n%s", err, value.String())
	}
	root := decoded["root"].(map[string]any)
	children := root["children"].([]any)
	if len(children) != 7 {
		t.Fatalf("recipe children len = %d, want 7: %#v", len(children), children)
	}
	assertString(t, children[2].(map[string]any), "type", "DashboardGrid")
	assertString(t, children[3].(map[string]any), "type", "ContextDiagramPanel")
	assertString(t, children[5].(map[string]any), "type", "CourseStudioShell")
}

func TestContextWindowStyleSetHelpersBuildExpectedShape(t *testing.T) {
	vm := goja.New()
	reg := require.NewRegistry()
	Register(reg)
	reg.Enable(vm)

	value, err := vm.RunString(`
		const contextWindow = require("context_window.dsl");
		const styleSet = contextWindow.paletteStyleSet({
			palette: "Signal Orange / Cyan",
			entries: [
				{ id: "prompt", label: "Prompt", accent: "b", pattern: "checker" },
				{ id: "evidence", label: "Evidence", accent: "a", pattern: "stipple" },
				{ id: "answer", label: "Answer", accent: "a", pattern: "solid" },
			]
		});
		const snapshot = contextWindow.contextSnapshot({
			id: "ctx",
			title: "Window",
			limit: 1000,
			parts: [contextWindow.contextPart("p", "Prompt", "prompt", 300)]
		});
		JSON.stringify({ styleSet, snapshot });
	`)
	if err != nil {
		t.Fatalf("build style helpers: %v", err)
	}
	var decoded map[string]any
	if err := json.Unmarshal([]byte(value.String()), &decoded); err != nil {
		t.Fatalf("style helper output is not JSON: %v", err)
	}
	styleSet := decoded["styleSet"].(map[string]any)
	styles := styleSet["styles"].(map[string]any)
	if _, ok := styles["prompt"].(map[string]any); !ok {
		t.Fatalf("styleSet.styles.prompt missing: %#v", styleSet)
	}
	snapshot := decoded["snapshot"].(map[string]any)
	parts := snapshot["parts"].([]any)
	part := parts[0].(map[string]any)
	assertString(t, part, "styleKey", "prompt")
	if _, hasKind := part["kind"]; hasKind {
		t.Fatalf("contextPart emitted forbidden kind field: %#v", part)
	}
}

func TestContextDiagramRecipeRequiresStyleSet(t *testing.T) {
	vm := goja.New()
	reg := require.NewRegistry()
	Register(reg)
	reg.Enable(vm)

	_, err := vm.RunString(`
		const contextWindow = require("context_window.dsl");
		const snapshot = { id: "ctx", title: "Window", limit: 1000, parts: [] };
		contextWindow.recipes.contextDiagram({ snapshot, view: "budget" });
	`)
	if err == nil {
		t.Fatalf("contextDiagram recipe without styleSet should fail")
	}
	if !strings.Contains(err.Error(), "requires styleSet") {
		t.Fatalf("error = %v, want useful styleSet message", err)
	}
}

func TestEngineRegistrarRegistersSplitModulesOnly(t *testing.T) {
	factory, err := engine.NewRuntimeFactoryBuilder().WithModules(NewRegistrar()).Build()
	if err != nil {
		t.Fatalf("build runtime factory: %v", err)
	}
	rt, err := factory.NewRuntime()
	if err != nil {
		t.Fatalf("create runtime: %v", err)
	}
	defer func() { _ = rt.Close(context.Background()) }()

	value, err := rt.VM.RunString(`
		function canRequire(name) {
			try { require(name); return true; } catch (error) { return false; }
		}
		const ui = require("ui.dsl");
		const data = require("data.dsl");
		const contextWindow = require("context_window.dsl");
		const course = require("course.dsl");
		({
			uiPanel: typeof ui.panel,
			dataTable: typeof data.dataTable,
			contextDiagramPanel: typeof contextWindow.contextDiagramPanel,
			courseStudioShell: typeof course.courseStudioShell,
			widget: canRequire("widget.dsl"),
			rag: canRequire("rag.dsl"),
		});
	`)
	if err != nil {
		t.Fatalf("require split modules through engine registrar: %v", err)
	}
	got := value.Export().(map[string]any)
	for _, name := range []string{"uiPanel", "dataTable", "contextDiagramPanel", "courseStudioShell"} {
		if got[name] != "function" {
			t.Fatalf("%s export = %#v, want function (all: %#v)", name, got[name], got)
		}
	}
	if got["widget"] != false || got["rag"] != false {
		t.Fatalf("old bucket modules should be absent from engine registrar, got %#v", got)
	}
}

func assertString(t *testing.T, m map[string]any, key, want string) {
	t.Helper()
	if got, _ := m[key].(string); got != want {
		t.Fatalf("%s = %#v, want %q (map=%#v)", key, m[key], want, m)
	}
}

func TestCmsModuleExportsHelpersRecipesAndBoundaries(t *testing.T) {
	vm := goja.New()
	reg := require.NewRegistry()
	Register(reg)
	reg.Enable(vm)

	value, err := vm.RunString(`
		const ui = require("ui.dsl");
		const contextWindow = require("context_window.dsl");
		const cms = require("cms.dsl");
		const library = cms.recipes.mediaLibrary({
			assets: [{ id: "a", kind: "image", title: "A", filename: "a.png", mime: "image/png", size: 10, src: "/course-assets/a.png", tags: [], status: "published", createdAt: "", updatedAt: "" }],
			selectedAssetIds: ["a"],
			onFilesSelected: "admin-upload-course-material",
			onAssetSelect: cms.action.navigate("?asset=$assetId"),
		});
		const list = cms.recipes.articleList({
			articles: [{ id: "x", slug: "x", title: "X", status: "draft", tags: [], updatedAt: "" }],
			onRowAction: { kind: "event", event: "row-action", confirm: "Really?" },
		});
		({
			cmsMediaLibraryPanel: typeof cms.mediaLibraryPanel,
			cmsTag: typeof cms.tag,
			cmsMediaThumb: typeof cms.mediaThumb,
			cmsMarkdownEditor: typeof cms.markdownEditor,
			cmsActionServer: typeof cms.action.server,
			cmsPage: typeof cms.page,
			cmsCell: typeof cms.cell,
			uiMediaLibraryPanel: typeof ui.mediaLibraryPanel,
			contextGroupedStrip: typeof contextWindow.contextGroupedStripDiagram,
			libraryType: library.type,
			libraryUploadKind: library.props.onFilesSelectedAction.kind,
			libraryUploadName: library.props.onFilesSelectedAction.name,
			librarySelectKind: library.props.onAssetSelectAction.kind,
			listType: list.type,
			listRowActionConfirm: list.props.onRowActionAction.confirm,
		});
	`)
	if err != nil {
		t.Fatalf("require cms.dsl: %v", err)
	}
	got := value.Export().(map[string]any)
	wantFunctions := []string{"cmsMediaLibraryPanel", "cmsTag", "cmsMediaThumb", "cmsMarkdownEditor", "cmsActionServer", "contextGroupedStrip"}
	for _, name := range wantFunctions {
		if got[name] != "function" {
			t.Fatalf("%s = %#v, want function (all: %#v)", name, got[name], got)
		}
	}
	wantUndefined := []string{"cmsPage", "cmsCell", "uiMediaLibraryPanel"}
	for _, name := range wantUndefined {
		if got[name] != "undefined" {
			t.Fatalf("%s = %#v, want undefined (all: %#v)", name, got[name], got)
		}
	}
	if got["libraryType"] != "MediaLibraryPanel" {
		t.Fatalf("libraryType = %#v, want MediaLibraryPanel", got["libraryType"])
	}
	if got["libraryUploadKind"] != "server" || got["libraryUploadName"] != "admin-upload-course-material" {
		t.Fatalf("upload action = %#v/%#v, want server/admin-upload-course-material", got["libraryUploadKind"], got["libraryUploadName"])
	}
	if got["librarySelectKind"] != "navigate" {
		t.Fatalf("select action kind = %#v, want navigate", got["librarySelectKind"])
	}
	if got["listType"] != "ArticleListPanel" {
		t.Fatalf("listType = %#v, want ArticleListPanel", got["listType"])
	}
	if got["listRowActionConfirm"] != "Really?" {
		t.Fatalf("confirm passthrough = %#v, want Really?", got["listRowActionConfirm"])
	}
}
