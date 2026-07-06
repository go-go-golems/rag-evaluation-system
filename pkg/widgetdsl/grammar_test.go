package widgetdsl

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
)

func grammarVM(t *testing.T) *goja.Runtime {
	t.Helper()
	vm := goja.New()
	reg := require.NewRegistry()
	Register(reg)
	reg.Enable(vm)
	return vm
}

func runJSON(t *testing.T, vm *goja.Runtime, src string) map[string]any {
	t.Helper()
	value, err := vm.RunString(src)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	raw, err := json.Marshal(value.Export())
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var out map[string]any
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	return out
}

const agendaSchemaJS = `
	const data = require("data.dsl");
	const schema = data.schema({
		id: data.f.key({ hint: "Stable anchor." }),
		number: data.f.short({ label: "Time", width: "6ch" }),
		duration: data.f.short({ width: "8ch" }),
		title: data.f.primary({ required: true, maxLength: 160 }),
		description: data.f.prose({ rows: 4, maxLength: 800 }),
	});
`

func TestSchemaPreservesFieldOrderAndRoles(t *testing.T) {
	vm := grammarVM(t)
	got := runJSON(t, vm, agendaSchemaJS+`schema`)
	fields, _ := got["fields"].([]any)
	if len(fields) != 5 {
		t.Fatalf("fields = %d, want 5 (%#v)", len(fields), got)
	}
	wantOrder := []string{"id", "number", "duration", "title", "description"}
	wantRoles := []string{"key", "short", "short", "primary", "prose"}
	for i, entry := range fields {
		field := entry.(map[string]any)
		if field["name"] != wantOrder[i] {
			t.Fatalf("field %d name = %v, want %s", i, field["name"], wantOrder[i])
		}
		if field["role"] != wantRoles[i] {
			t.Fatalf("field %d role = %v, want %s", i, field["role"], wantRoles[i])
		}
	}
}

func TestRecordEditCompilesToFormPanelWithFieldGridAndProse(t *testing.T) {
	vm := grammarVM(t)
	got := runJSON(t, vm, agendaSchemaJS+`
		data.record({ id: "agenda-break", number: "16h05", duration: "20 min", title: "Pause", description: "Debrief." }, {
			schema,
			title: "Edit agenda item",
			submit: data.formPost("/settings/agenda-item"),
			status: "idle",
		})
	`)
	if got["type"] != "FormPanel" {
		t.Fatalf("type = %v, want FormPanel", got["type"])
	}
	props := got["props"].(map[string]any)
	if props["formAction"] != "/settings/agenda-item" || props["method"] != "post" {
		t.Fatalf("submit not applied: %#v", props)
	}
	children := got["children"].([]any)
	// id+number+duration are gridable and consecutive → one FieldGrid; then
	// title (primary) and description (prose) as standalone rows.
	first := children[0].(map[string]any)
	if first["type"] != "FieldGrid" {
		t.Fatalf("first child = %v, want FieldGrid (%#v)", first["type"], children)
	}
	gridRows := first["children"].([]any)
	if len(gridRows) != 3 {
		t.Fatalf("grid rows = %d, want 3", len(gridRows))
	}
	idRow := gridRows[0].(map[string]any)["props"].(map[string]any)
	idControl := idRow["control"].(map[string]any)["props"].(map[string]any)
	if idControl["readOnly"] != true {
		t.Fatalf("key field should default readOnly: %#v", idControl)
	}
	last := children[len(children)-1].(map[string]any)
	lastProps := last["props"].(map[string]any)
	control := lastProps["control"].(map[string]any)
	if control["type"] != "TextareaInput" || lastProps["orientation"] != "stacked" {
		t.Fatalf("prose row should be stacked textarea: %#v", last)
	}
	if control["props"].(map[string]any)["defaultValue"] != "Debrief." {
		t.Fatalf("prose defaultValue missing: %#v", control)
	}
}

func TestCollectionTableDerivesColumnsAndSelection(t *testing.T) {
	vm := grammarVM(t)
	got := runJSON(t, vm, agendaSchemaJS+`
		data.collection([
			{ id: "a", number: "14h30", duration: "15 min", title: "Intro", description: "Long prose." },
			{ id: "b", number: "14h45", duration: "35 min", title: "Demo", description: "More prose." },
		], {
			schema,
			title: "Agenda",
			verb: "edit",
			arrange: "master-detail",
			select: data.urlParam("agenda", "b"),
			submit: data.formPost("/settings/agenda-item"),
			reorder: "admin-reorder-agenda",
			remove: { kind: "server", name: "admin-delete-agenda", confirm: "Delete ${row.title}?" },
			create: true,
			empty: "No items.",
		})
	`)
	if got["type"] != "SectionBlock" {
		t.Fatalf("titled collection should wrap in SectionBlock, got %v", got["type"])
	}
	section := got["props"].(map[string]any)
	if section["label"] != "Agenda" || section["rule"] != true {
		t.Fatalf("section props: %#v", section)
	}
	stack := got["children"].([]any)[0].(map[string]any)
	parts := stack["children"].([]any)
	// create button, table, detail
	if len(parts) != 3 {
		t.Fatalf("stack parts = %d, want 3 (create, table, detail)", len(parts))
	}
	table := parts[1].(map[string]any)
	if table["type"] != "DataTable" {
		t.Fatalf("second part = %v, want DataTable", table["type"])
	}
	tableProps := table["props"].(map[string]any)
	if tableProps["selectedKey"] != "b" || tableProps["emptyMessage"] != "No items." {
		t.Fatalf("table props: %#v", tableProps)
	}
	rowSelect := tableProps["onRowSelect"].(map[string]any)
	if rowSelect["to"] != "?agenda=${row.id}" {
		t.Fatalf("onRowSelect = %#v", rowSelect)
	}
	columns := tableProps["columns"].([]any)
	ids := []string{}
	for _, entry := range columns {
		ids = append(ids, entry.(map[string]any)["id"].(string))
	}
	joined := strings.Join(ids, ",")
	if strings.Contains(joined, "description") {
		t.Fatalf("prose column should be elided: %s", joined)
	}
	if !strings.Contains(joined, "moveUp") || !strings.Contains(joined, "moveDown") || !strings.Contains(joined, "delete") {
		t.Fatalf("action columns missing: %s", joined)
	}
	for _, entry := range columns {
		column := entry.(map[string]any)
		if column["id"] == "delete" {
			action := column["cell"].(map[string]any)["action"].(map[string]any)
			if action["confirm"] != "Delete ${row.title}?" {
				t.Fatalf("remove confirm lost: %#v", action)
			}
		}
		if column["id"] == "moveUp" {
			action := column["cell"].(map[string]any)["action"].(map[string]any)
			payload := action["payload"].(map[string]any)
			if action["name"] != "admin-reorder-agenda" || payload["direction"] != "up" {
				t.Fatalf("reorder action: %#v", action)
			}
		}
	}
	detailStack := parts[2].(map[string]any)
	detail := detailStack["children"].([]any)[0].(map[string]any)
	if detail["type"] != "FormPanel" {
		t.Fatalf("detail = %v, want FormPanel", detail["type"])
	}
	if !strings.Contains(detail["props"].(map[string]any)["title"].(string), "Demo") {
		t.Fatalf("detail title should name the primary field: %#v", detail["props"])
	}
}

func TestCollectionWithoutSelectionShowsPlaceholderAndNewRecord(t *testing.T) {
	vm := grammarVM(t)
	got := runJSON(t, vm, agendaSchemaJS+`
		data.collection([], {
			schema, verb: "edit", arrange: "master-detail",
			select: data.urlParam("agenda", "__new"),
			submit: data.formPost("/settings/agenda-item"),
		})
	`)
	parts := got["children"].([]any)
	detailStack := parts[len(parts)-1].(map[string]any)
	detail := detailStack["children"].([]any)[0].(map[string]any)
	title := detail["props"].(map[string]any)["title"]
	if title != "New item" {
		t.Fatalf("__new detail title = %v, want New item", title)
	}
}

func TestSectionCompilesToSectionBlockWithDefaults(t *testing.T) {
	vm := grammarVM(t)
	got := runJSON(t, vm, `
		const ui = require("ui.dsl");
		ui.section("Media library", { level: 2, anchor: "media", caption: "Files under course/media." },
			ui.textBlock({}, "body"))
	`)
	if got["type"] != "SectionBlock" {
		t.Fatalf("type = %v", got["type"])
	}
	props := got["props"].(map[string]any)
	if props["label"] != "Media library" || props["rule"] != true || props["density"] != "flush" {
		t.Fatalf("defaults: %#v", props)
	}
	if props["level"] != float64(2) || props["anchorId"] != "media" {
		t.Fatalf("options: %#v", props)
	}
	if len(got["children"].([]any)) != 1 {
		t.Fatalf("children lost")
	}
}

func TestPromotedUIHelpersAndAliasesCoexist(t *testing.T) {
	vm := grammarVM(t)
	got := runJSON(t, vm, `
		const ui = require("ui.dsl");
		const cms = require("cms.dsl");
		const cw = require("context_window.dsl");
		({
			uiTag: typeof ui.tag, uiPagination: typeof ui.pagination,
			uiTileGrid: typeof ui.tileGrid, uiUploadDropArea: typeof ui.uploadDropArea,
			uiMarkdownArticle: typeof ui.markdownArticle, uiSection: typeof ui.section,
			cmsTag: typeof cms.tag, cwUpload: typeof cw.contextUploadDropArea,
			node: ui.uploadDropArea({ title: "Drop" }),
		})
	`)
	for _, name := range []string{"uiTag", "uiPagination", "uiTileGrid", "uiUploadDropArea", "uiMarkdownArticle", "uiSection", "cmsTag", "cwUpload"} {
		if got[name] != "function" {
			t.Fatalf("%s = %v, want function", name, got[name])
		}
	}
	node := got["node"].(map[string]any)
	if node["type"] != "ContextUploadDropArea" {
		t.Fatalf("uploadDropArea should map to ContextUploadDropArea, got %v", node["type"])
	}
}
