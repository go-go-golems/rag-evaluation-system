package widgetdsl

import (
	"fmt"
	"testing"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
)

func TestV3CRMFunnelDefaultsMissingStageSummaryToZero(t *testing.T) {
	vm := goja.New()
	reg := require.NewRegistry()
	Register(reg)
	reg.Enable(vm)

	value, err := vm.RunString(`
		const widget = require("widget.dsl");
		const pipeline = widget.crm.pipeline("Sales", (p) =>
			p.stage("lead", "Lead", { colorKey: "lead" })
			 .stage("won", "Won", { colorKey: "won" })
		);
		widget.crm.funnel(pipeline, [{ stageId: "lead", count: 3 }]);
	`)
	if err != nil {
		t.Fatalf("render CRM funnel: %v", err)
	}

	node := value.Export().(map[string]any)
	props := node["props"].(map[string]any)
	segments := props["segments"].([]any)
	if got := fmt.Sprint(segments[0].(map[string]any)["value"]); got != "3" {
		t.Errorf("lead segment value = %s, want 3", got)
	}
	if got := fmt.Sprint(segments[1].(map[string]any)["value"]); got != "0" {
		t.Errorf("missing-stage segment value = %s, want 0", got)
	}
}
