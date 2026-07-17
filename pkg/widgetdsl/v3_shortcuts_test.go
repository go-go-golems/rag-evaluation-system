package widgetdsl

import (
	"testing"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
)

func TestWidgetV3PageShortcutsEmitTypedBindings(t *testing.T) {
	vm := goja.New()
	reg := require.NewRegistry()
	registerLegacyModulesForTests(reg)
	reg.Enable(vm)

	value, err := vm.RunString(`
		const widget = require("widget.dsl");
		widget.page("Triage", page => page
			.shortcuts(keys => keys
				.bind("accept", "y", widget.act.server("triage.accept"), { label: "Yes" })
				.bind("save", "s", widget.act.server("triage.save"), {
					label: "Save",
					modifiers: ["Control"],
					preventDefault: false,
					allowRepeat: true,
				})
				.bind("copy", "c", widget.act.copy("job-1"), { label: "Copy job id" }))
			.view(widget.ui.caption("Review the current job")))
			.toPage();
	`)
	if err != nil {
		t.Fatalf("build shortcut page: %v", err)
	}

	page := value.Export().(map[string]any)
	bindings := anySlice(anyMap(page["shortcuts"])["bindings"])
	if len(bindings) != 3 {
		t.Fatalf("bindings = %#v, want three", bindings)
	}
	accept := anyMap(bindings[0])
	if accept["id"] != "accept" || accept["key"] != "y" || accept["label"] != "Yes" || accept["preventDefault"] != true || accept["allowRepeat"] != false {
		t.Fatalf("accept binding = %#v", accept)
	}
	acceptAction := anyMap(accept["action"])
	if acceptAction["kind"] != "server" || acceptAction["name"] != "triage.accept" {
		t.Fatalf("accept action = %#v", acceptAction)
	}
	save := anyMap(bindings[1])
	if save["preventDefault"] != false || save["allowRepeat"] != true {
		t.Fatalf("save policies = %#v", save)
	}
	modifiers := anySlice(save["modifiers"])
	if len(modifiers) != 1 || modifiers[0] != "Control" {
		t.Fatalf("save modifiers = %#v", modifiers)
	}
	copyAction := anyMap(anyMap(bindings[2])["action"])
	if copyAction["kind"] != "copy" || copyAction["value"] != "job-1" {
		t.Fatalf("copy action options lost: %#v", copyAction)
	}
}

func TestWidgetV3PageShortcutsRejectDuplicateChord(t *testing.T) {
	vm := goja.New()
	reg := require.NewRegistry()
	registerLegacyModulesForTests(reg)
	reg.Enable(vm)

	_, err := vm.RunString(`
		const widget = require("widget.dsl");
		widget.page("Triage", page => page
			.shortcuts(keys => keys
				.bind("accept", "y", widget.act.server("triage.accept"), { label: "Yes" })
				.bind("reject", "Y", widget.act.server("triage.reject"), { label: "No" })))
			.toPage();
	`)
	if err == nil {
		t.Fatal("duplicate shortcut chord unexpectedly succeeded")
	}
}
