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
	Register(reg)
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

func TestDataV2BuilderBuildsSelectableTable(t *testing.T) {
	vm := goja.New()
	reg := require.NewRegistry()
	Register(reg)
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

func TestDataV2BuilderRejectsPresentNonFunctionCallback(t *testing.T) {
	vm := goja.New()
	reg := require.NewRegistry()
	Register(reg)
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
