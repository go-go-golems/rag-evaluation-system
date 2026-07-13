package widgetdsl

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
)

func TestDataV2BuilderBuildsSimpleTable(t *testing.T) {
	vm := goja.New()
	reg := require.NewRegistry()
	registerLegacyModulesForTests(reg)
	reg.Enable(vm)

	value, err := vm.RunString(`
		const data = require("data.v2.dsl");
		const schema = data.schema("Session")
			.field("sessionId", data.f.key().label("ID"))
			.field("title", data.f.primary().required().maxLength(120))
			.field("turnCount", data.f.count().label("Turns"))
			.field("status", data.f.status())
			.field("body", data.f.prose().rows(3))
			.build();
		const ir = data.collection("sessions", [{ sessionId: "s1", title: "Intro", turnCount: 12, status: "ready", body: "Long" }])
			.schema(schema)
			.table()
			.toIR();
		JSON.stringify(ir);
	`)
	if err != nil {
		t.Fatalf("build simple table: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal([]byte(value.String()), &got); err != nil {
		t.Fatalf("decode IR: %v", err)
	}
	children := got["children"].([]any)
	table := children[0].(map[string]any)
	props := table["props"].(map[string]any)
	if props["getRowKey"] != "sessionId" {
		t.Fatalf("getRowKey = %#v", props["getRowKey"])
	}
	if _, ok := props["onRowSelect"]; ok {
		t.Fatalf("simple table unexpectedly has onRowSelect: %#v", props["onRowSelect"])
	}
	columns := props["columns"].([]any)
	if len(columns) != 4 {
		t.Fatalf("columns len = %d, want 4", len(columns))
	}
}

func TestDataV2BuilderBuildsTableWithExplicitActionColumns(t *testing.T) {
	vm := goja.New()
	reg := require.NewRegistry()
	registerLegacyModulesForTests(reg)
	reg.Enable(vm)

	value, err := vm.RunString(`
		const data = require("data.v2.dsl");
		const schema = data.schema("Material")
			.field("file", data.f.key().label("File"))
			.field("title", data.f.primary().label("Title"))
			.field("size", data.f.short().label("Size"))
			.build();
		const ir = data.collection("materials", [{ file: "deck.md", title: "Deck", href: "/slides/deck.md", size: "2 KB" }])
			.schema(schema)
			.empty("No files.")
			.table(t => t
				.className("course-material-table")
				.actionColumn("open", "Open", "Open", data.action.navigate("${row.href}"), { maxWidth: "8ch" })
				.actionColumn("delete", "Delete", "Delete", data.action.server("delete-material").confirm("Delete ${row.file}?"), { maxWidth: "9ch" }))
			.toIR();
		JSON.stringify(ir.children[0].props);
	`)
	if err != nil {
		t.Fatalf("build table with action columns: %v", err)
	}
	var props map[string]any
	if err := json.Unmarshal([]byte(value.String()), &props); err != nil {
		t.Fatalf("decode props: %v", err)
	}
	if props["emptyMessage"] != "No files." {
		t.Fatalf("emptyMessage = %#v", props["emptyMessage"])
	}
	if props["className"] != "course-material-table" {
		t.Fatalf("className = %#v", props["className"])
	}
	columns := props["columns"].([]any)
	if len(columns) != 5 {
		t.Fatalf("columns len = %d, want 5", len(columns))
	}
	open := columns[len(columns)-2].(map[string]any)
	if open["id"] != "open" || open["maxWidth"] != "8ch" {
		t.Fatalf("open column = %#v", open)
	}
	cell := open["cell"].(map[string]any)
	action := cell["action"].(map[string]any)
	if action["to"] != "${row.href}" {
		t.Fatalf("open action.to = %#v", action["to"])
	}
}

func TestDataV2BuilderBuildsSelectableTable(t *testing.T) {
	vm := goja.New()
	reg := require.NewRegistry()
	registerLegacyModulesForTests(reg)
	reg.Enable(vm)

	value, err := vm.RunString(`
		const data = require("data.v2.dsl");
		const schema = data.schema("Session")
			.field("sessionId", data.f.key())
			.field("title", data.f.primary())
			.build();
		const ir = data.collection("sessions", [{ sessionId: "s1", title: "Intro" }, { sessionId: "s2", title: "Debugging" }])
			.schema(schema)
			.select(s => s.urlParam("selected", "s2"))
			.table(t => t.rowSelect(data.action.navigate("/pages/sessions?selected=${row.sessionId}")))
			.toIR();
		JSON.stringify(ir.children[0].props);
	`)
	if err != nil {
		t.Fatalf("build selectable table: %v", err)
	}
	var props map[string]any
	if err := json.Unmarshal([]byte(value.String()), &props); err != nil {
		t.Fatalf("decode props: %v", err)
	}
	if props["selectedKey"] != "s2" {
		t.Fatalf("selectedKey = %#v", props["selectedKey"])
	}
	action := props["onRowSelect"].(map[string]any)
	if action["to"] != "/pages/sessions?selected=${row.sessionId}" {
		t.Fatalf("onRowSelect.to = %#v", action["to"])
	}
}

func TestDataV2BuilderBuildsMasterDetailEditor(t *testing.T) {
	vm := goja.New()
	reg := require.NewRegistry()
	registerLegacyModulesForTests(reg)
	reg.Enable(vm)

	value, err := vm.RunString(`
		const data = require("data.v2.dsl");
		const schema = data.schema("AgendaItem")
			.field("id", data.f.key().label("ID"))
			.field("number", data.f.short().label("Time"))
			.field("title", data.f.primary().required())
			.field("description", data.f.prose().rows(4))
			.build();
		const ir = data.collection("agenda", [{ id: "agenda-intro", number: "14h30", title: "Intro", description: "Welcome" }])
			.schema(schema)
			.edit(e => e
				.selectUrl("agenda", "agenda-intro")
				.submitPost("/settings/agenda-item")
				.create({ label: "New agenda item" })
				.actions(a => a
					.reorder(data.action.server("admin-reorder-course-agenda"))
					.remove(data.action.server("admin-delete-agenda-item").confirm("Delete ${row.title}?"))))
			.masterDetail()
			.toIR();
		JSON.stringify(ir);
	`)
	if err != nil {
		t.Fatalf("build master-detail editor: %v", err)
	}
	var root map[string]any
	if err := json.Unmarshal([]byte(value.String()), &root); err != nil {
		t.Fatalf("decode root: %v", err)
	}
	children := root["children"].([]any)
	if len(children) != 3 {
		t.Fatalf("children len = %d, want create/table/detail", len(children))
	}
	table := children[1].(map[string]any)
	tableProps := table["props"].(map[string]any)
	if tableProps["selectedKey"] != "agenda-intro" {
		t.Fatalf("selectedKey = %#v", tableProps["selectedKey"])
	}
	detail := children[2].(map[string]any)
	form := detail["children"].([]any)[0].(map[string]any)
	formProps := form["props"].(map[string]any)
	if formProps["formAction"] != "/settings/agenda-item" {
		t.Fatalf("formAction = %#v", formProps["formAction"])
	}
	if formProps["title"] != "Edit: Intro" {
		t.Fatalf("title = %#v", formProps["title"])
	}
	columns := tableProps["columns"].([]any)
	deleteColumn := columns[len(columns)-1].(map[string]any)
	deleteCell := deleteColumn["cell"].(map[string]any)
	deleteAction := deleteCell["action"].(map[string]any)
	if deleteAction["confirm"] != "Delete ${row.title}?" {
		t.Fatalf("delete confirm = %#v", deleteAction["confirm"])
	}
}

func TestDataV2BuilderRejectsPresentNonFunctionCallback(t *testing.T) {
	vm := goja.New()
	reg := require.NewRegistry()
	registerLegacyModulesForTests(reg)
	reg.Enable(vm)

	_, err := vm.RunString(`
		const data = require("data.v2.dsl");
		const schema = data.schema("Session")
			.field("id", data.f.key())
			.field("title", data.f.primary())
			.build();
		data.collection("sessions", [{ id: "s1", title: "Intro" }])
			.schema(schema)
			.table({ not: "a function" })
			.toIR();
	`)
	if err == nil {
		t.Fatalf("expected non-function callback error")
	}
	if !strings.Contains(err.Error(), "requires a function") {
		t.Fatalf("error = %v", err)
	}
}
