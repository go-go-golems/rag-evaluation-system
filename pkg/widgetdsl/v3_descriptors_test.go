package widgetdsl

import (
	"strings"
	"testing"
)

func TestWidgetV3DescriptorExportsMatchDeclarations(t *testing.T) {
	dts := strings.Join(TypeScriptModule(WidgetV3ModuleName).RawDTS, "\n")
	for _, namespace := range widgetV3NamespaceDescriptors {
		want := "export const " + namespace.ExportName + ": " + namespace.TypeName + ";"
		if !strings.Contains(dts, want) {
			t.Fatalf("descriptor export %q missing from DTS", want)
		}
		for _, view := range namespace.Views {
			if !strings.Contains(dts, view.Signature+";") {
				t.Fatalf("descriptor view %s.%s signature missing from DTS: %s", namespace.ExportName, view.Name, view.Signature)
			}
		}
	}
}

func TestWidgetV3GeneratedAPIReferenceIncludesDescriptorViews(t *testing.T) {
	md := WidgetV3APIReferenceMarkdown()
	for _, fragment := range []string{
		"# widget.dsl API Reference",
		"## `schedule` — ScheduleNamespace",
		"`availabilityPoll(poll: AvailabilityPoll, configure?: Fragment<SchedulePollBuilder>): WidgetNodeSpec` → `MatrixGrid`",
		"## `time` — TimeNamespace",
		"`week(events: CalendarEvent[], configure?: Fragment<TimeWeekBuilder>): WidgetNodeSpec` → `TimeGrid`",
	} {
		if !strings.Contains(md, fragment) {
			t.Fatalf("API reference missing %q\n--- markdown ---\n%s", fragment, md)
		}
	}
}
