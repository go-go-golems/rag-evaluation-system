package widgetdsl

import (
	"testing"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
)

func TestWidgetV3FieldElideOptionOmitsTableColumn(t *testing.T) {
	vm := goja.New()
	reg := require.NewRegistry()
	registerLegacyModulesForTests(reg)
	reg.Enable(vm)

	value, err := vm.RunString(`
		const widget = require("widget.dsl");
		const fields = widget.data.fields(fields => fields
			.key("id", { label: "ID", elide: true })
			.primary("title", { label: "Title" }));
		widget.data.collection("rows", [{ id: "one", title: "One" }], collection => collection
			.schema(fields.build())
			.table()).toNode();
	`)
	if err != nil {
		t.Fatalf("build elided key collection: %v", err)
	}

	root := anyMap(value.Export())
	table := anyMap(anySlice(root["children"])[0])
	columns := anySlice(anyMap(table["props"])["columns"])
	if len(columns) != 1 || anyMap(columns[0])["id"] != "title" {
		t.Fatalf("columns = %#v, want only title", columns)
	}
	if anyMap(table["props"])["getRowKey"] != "id" {
		t.Fatalf("elided key no longer provides row identity: %#v", table)
	}
}
