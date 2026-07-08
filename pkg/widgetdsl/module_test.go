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

func TestWidgetV3ModuleExportsRootNamespacesAndKeepsOldModulesAvailable(t *testing.T) {
	vm := goja.New()
	reg := require.NewRegistry()
	Register(reg)
	reg.Enable(vm)

	value, err := vm.RunString(`
		const widget = require("widget.dsl");
		const ui = require("ui.dsl");
		const node = widget.raw.component("Panel", { title: "V3" }, widget.raw.text("Hello"));
		({
			legacyPanel: typeof ui.panel,
			rawComponent: typeof widget.raw.component,
			rawText: typeof widget.raw.text,
			actServer: typeof widget.act.server,
			bindField: typeof widget.bind.field,
			pageFn: typeof widget.page,
			uiNamespace: typeof widget.ui,
			dataNamespace: typeof widget.data,
			cmsNamespace: typeof widget.cms,
			courseNamespace: typeof widget.course,
			contextNamespace: typeof widget.context,
			scheduleNamespace: typeof widget.schedule,
			timeNamespace: typeof widget.time,
			styleNamespace: typeof widget.style,
			node,
			binding: widget.bind.field("title"),
			action: widget.act.server("save", { payload: { id: 1 } }),
			dataSelection: typeof widget.data.selection,
			dataItem: typeof widget.data.item,
		});
	`)
	if err != nil {
		t.Fatalf("require widget.dsl: %v", err)
	}
	got := value.Export().(map[string]any)
	wantFunctions := []string{"legacyPanel", "rawComponent", "rawText", "actServer", "bindField", "pageFn", "dataSelection", "dataItem"}
	for _, name := range wantFunctions {
		if got[name] != "function" {
			t.Fatalf("%s export = %#v, want function (all: %#v)", name, got[name], got)
		}
	}
	wantObjects := []string{"uiNamespace", "dataNamespace", "cmsNamespace", "courseNamespace", "contextNamespace", "scheduleNamespace", "timeNamespace", "styleNamespace"}
	for _, name := range wantObjects {
		if got[name] != "object" {
			t.Fatalf("%s export = %#v, want object (all: %#v)", name, got[name], got)
		}
	}
	node := got["node"].(map[string]any)
	if node["kind"] != "component" || node["type"] != "Panel" {
		t.Fatalf("raw.component emitted %#v", node)
	}
	binding := got["binding"].(map[string]any)
	if binding["kind"] != "accessor" || binding["mode"] != "field" || binding["field"] != "title" {
		t.Fatalf("bind.field emitted %#v", binding)
	}
	action := got["action"].(map[string]any)
	if action["kind"] != "server" || action["name"] != "save" {
		t.Fatalf("act.server emitted %#v", action)
	}
}

func TestWidgetV3PageBuilderEmitsPageIR(t *testing.T) {
	vm := goja.New()
	reg := require.NewRegistry()
	Register(reg)
	reg.Enable(vm)

	value, err := vm.RunString(`
		const widget = require("widget.dsl");
		const quiet = s => s.tone("quiet");
		const builder = widget.page("Hello V3", p => p
			.id("hello-v3")
			.meta("source", "test")
			.use(p => p.title("Hello Widget V3"))
			.section("Intro", s => s
				.use(quiet)
				.caption("Builder callbacks lower to serializable Widget IR.")
				.anchor("intro")
				.text("Hello")
				.view(widget.raw.component("Caption", { tone: "muted" }, "World"))));
		({ issues: builder.validate(), page: builder.toPage() });
	`)
	if err != nil {
		t.Fatalf("build widget.dsl page: %v", err)
	}
	got := value.Export().(map[string]any)
	if issueCount := exportedSliceLen(got["issues"]); issueCount != 0 {
		t.Fatalf("unexpected validation issues: %#v", got["issues"])
	}
	page := got["page"].(map[string]any)
	if page["id"] != "hello-v3" || page["title"] != "Hello Widget V3" {
		t.Fatalf("unexpected page identity: %#v", page)
	}
	meta := page["meta"].(map[string]any)
	if meta["source"] != "test" {
		t.Fatalf("unexpected page meta: %#v", meta)
	}
	root := page["root"].(map[string]any)
	if root["type"] != "Stack" {
		t.Fatalf("unexpected root: %#v", root)
	}
	children := root["children"].([]any)
	if len(children) != 1 {
		t.Fatalf("root children = %#v, want one section", children)
	}
	section := children[0].(map[string]any)
	if section["type"] != "SectionBlock" {
		t.Fatalf("unexpected section: %#v", section)
	}
	props := section["props"].(map[string]any)
	if props["label"] != "Intro" || props["caption"] != "Builder callbacks lower to serializable Widget IR." || props["anchorId"] != "intro" || props["tone"] != "quiet" {
		t.Fatalf("unexpected section props: %#v", props)
	}
	sectionChildren := section["children"].([]any)
	if len(sectionChildren) != 2 {
		t.Fatalf("section children = %#v, want text + caption", sectionChildren)
	}
}

func TestWidgetV3SlotsAndChildNormalization(t *testing.T) {
	vm := goja.New()
	reg := require.NewRegistry()
	Register(reg)
	reg.Enable(vm)

	value, err := vm.RunString(`
		const widget = require("widget.dsl");
		const fallback = (ctx, h) => h.caption("fallback: " + ctx.label);
		const page = widget.page("Slots", p => p.section("Slot section", s => s
			.view(["A", null, undefined, false, widget.raw.fragment("B", false)])
			.slot({ label: "primary" }, (ctx, h) => h.stack(
				{ gap: "sm" },
				h.strong(ctx.label),
				false,
				h.badge("ready", { tone: "success" })
			), fallback)
			.slot({ label: "fallback" }, null, fallback)
		)).toPage();
		page;
	`)
	if err != nil {
		t.Fatalf("build widget.dsl page with slots: %v", err)
	}
	page := value.Export().(map[string]any)
	root := page["root"].(map[string]any)
	section := anySlice(root["children"])[0].(map[string]any)
	children := anySlice(section["children"])
	if len(children) != 4 {
		t.Fatalf("section children = %#v, want A, B, slot stack, fallback caption", children)
	}
	first := children[0].(map[string]any)
	if first["kind"] != "text" || first["text"] != "A" {
		t.Fatalf("first normalized child = %#v, want text A", first)
	}
	second := children[1].(map[string]any)
	if second["kind"] != "text" || second["text"] != "B" {
		t.Fatalf("second normalized child = %#v, want text B", second)
	}
	stack := children[2].(map[string]any)
	if stack["type"] != "Stack" {
		t.Fatalf("slot stack child = %#v, want Stack", stack)
	}
	stackChildren := anySlice(stack["children"])
	if len(stackChildren) != 2 {
		t.Fatalf("slot stack children = %#v, want strong + badge", stackChildren)
	}
	fallbackCaption := children[3].(map[string]any)
	if fallbackCaption["type"] != "Caption" {
		t.Fatalf("fallback child = %#v, want Caption", fallbackCaption)
	}
}

func TestWidgetV3AccessorsSelectionsItemsActionsAndValidation(t *testing.T) {
	vm := goja.New()
	reg := require.NewRegistry()
	Register(reg)
	reg.Enable(vm)

	value, err := vm.RunString(`
		const widget = require("widget.dsl");
		const invalid = widget.page("Invalid", p => p.section("Broken", s => s.view({ kind: "component", props: {} })));
		({
			issues: invalid.validate(),
			field: widget.bind.field("title"),
			path: widget.bind.path("author.name"),
			map: widget.bind.map("cells"),
			template: widget.bind.template("${first} ${last}"),
			context: widget.bind.context("row.id"),
			constant: widget.bind.const(7),
			selection: widget.data.selection({ mode: "multi", keyField: "id", selected: ["a", "b"] }),
			singleSelection: widget.data.selection("single", { keyField: "id", selected: "a" }),
			item: widget.data.item("home", "Home", { href: "/", badge: "New", action: widget.act.navigate("/") }),
			action: widget.act.server("save", {
				confirm: "Save row?",
				payload: { id: widget.bind.context("row.id"), fixed: widget.bind.const(7) },
			}),
		});
	`)
	if err != nil {
		t.Fatalf("build widget.dsl core specs: %v", err)
	}
	got := value.Export().(map[string]any)
	issues := anySlice(got["issues"])
	if len(issues) != 1 {
		t.Fatalf("validation issues = %#v, want one component type issue", issues)
	}
	issue := issues[0].(map[string]any)
	if issue["code"] != "component_type_required" {
		t.Fatalf("validation issue = %#v, want component_type_required", issue)
	}
	field := got["field"].(map[string]any)
	if field["kind"] != "accessor" || field["mode"] != "field" || field["field"] != "title" {
		t.Fatalf("field accessor = %#v", field)
	}
	path := got["path"].(map[string]any)
	if path["kind"] != "accessor" || path["mode"] != "path" || path["path"] != "author.name" {
		t.Fatalf("path accessor = %#v", path)
	}
	mapAccessor := got["map"].(map[string]any)
	if mapAccessor["kind"] != "accessor" || mapAccessor["mode"] != "map" || mapAccessor["mapField"] != "cells" {
		t.Fatalf("map accessor = %#v", mapAccessor)
	}
	constant := got["constant"].(map[string]any)
	if constant["kind"] != "const" || constant["value"] != int64(7) {
		t.Fatalf("const binding = %#v", constant)
	}
	selection := got["selection"].(map[string]any)
	if selection["kind"] != "selection" || selection["mode"] != "multi" || selection["keyField"] != "id" {
		t.Fatalf("selection = %#v", selection)
	}
	singleSelection := got["singleSelection"].(map[string]any)
	if singleSelection["kind"] != "selection" || singleSelection["mode"] != "single" || singleSelection["selected"] != "a" {
		t.Fatalf("single selection = %#v", singleSelection)
	}
	item := got["item"].(map[string]any)
	if item["kind"] != "listItem" || item["id"] != "home" || item["href"] != "/" {
		t.Fatalf("list item = %#v", item)
	}
	action := got["action"].(map[string]any)
	payload := action["payload"].(map[string]any)
	payloadID := payload["id"].(map[string]any)
	if action["kind"] != "server" || action["confirm"] != "Save row?" || payloadID["kind"] != "accessor" || payloadID["mode"] != "context" {
		t.Fatalf("action with payload bindings = %#v", action)
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
		({ rag: canRequire("rag.dsl") });
	`)
	if err != nil {
		t.Fatalf("check old modules: %v", err)
	}
	got := value.Export().(map[string]any)
	if got["rag"] != false {
		t.Fatalf("old bucket module rag.dsl should be absent, got %#v", got)
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
	if got["widget"] != true || got["rag"] != false {
		t.Fatalf("widget.dsl should be present and rag.dsl should be absent from engine registrar, got %#v", got)
	}
}

func exportedSliceLen(value any) int {
	switch v := value.(type) {
	case []any:
		return len(v)
	case []map[string]any:
		return len(v)
	default:
		return -1
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
