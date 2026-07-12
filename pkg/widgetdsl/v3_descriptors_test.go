package widgetdsl

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/dop251/goja"
)

func TestWidgetV3DescriptorExportsMatchDeclarations(t *testing.T) {
	dts := strings.Join(TypeScriptModule(WidgetV3ModuleName).RawDTS, "\n")
	for _, export := range widgetV3Module.Exports {
		want := "export function " + export.Signature + ";"
		if !strings.Contains(dts, want) {
			t.Fatalf("descriptor export %q missing from DTS", want)
		}
	}
	for _, namespace := range widgetV3Module.Namespaces {
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

func TestWidgetV3DescriptorMatchesDirectRuntimeExports(t *testing.T) {
	vm := goja.New()
	exports := vm.NewObject()
	(&runtime{vm: vm}).installWidgetV3(exports)

	wantRoot := make([]string, 0, len(widgetV3Module.Exports)+len(widgetV3Module.Namespaces))
	for _, export := range widgetV3Module.Exports {
		wantRoot = append(wantRoot, export.Name)
		if got := exports.Get(export.Name); got == nil || goja.IsUndefined(got) {
			t.Fatalf("described root export %q is not installed", export.Name)
		}
	}
	for _, namespace := range widgetV3Module.Namespaces {
		wantRoot = append(wantRoot, namespace.ExportName)
		value := exports.Get(namespace.ExportName)
		if value == nil || goja.IsUndefined(value) || goja.IsNull(value) {
			t.Fatalf("described namespace %q is not installed", namespace.ExportName)
		}
		object := value.ToObject(vm)
		wantMembers := make([]string, 0, len(namespace.Members))
		for _, member := range namespace.Members {
			wantMembers = append(wantMembers, member.Name)
			memberValue := object.Get(member.Name)
			if memberValue == nil || goja.IsUndefined(memberValue) {
				t.Errorf("described member %s.%s is not installed", namespace.ExportName, member.Name)
			}
		}
		assertSameStringSet(t, namespace.ExportName+" members", object.Keys(), wantMembers)
	}
	assertSameStringSet(t, "widget.dsl root exports", exports.Keys(), wantRoot)
}

func TestWidgetV3DescriptorViewsAreDirectMembers(t *testing.T) {
	for _, namespace := range widgetV3Module.Namespaces {
		members := make(map[string]struct{}, len(namespace.Members))
		for _, member := range namespace.Members {
			members[member.Name] = struct{}{}
		}
		for _, view := range namespace.Views {
			if _, ok := members[view.Name]; !ok {
				t.Errorf("descriptor view %s.%s is not listed as a direct namespace member", namespace.ExportName, view.Name)
			}
		}
	}
}

func assertSameStringSet(t *testing.T, label string, got, want []string) {
	t.Helper()
	got = append([]string(nil), got...)
	want = append([]string(nil), want...)
	sort.Strings(got)
	sort.Strings(want)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("%s mismatch\n got: %v\nwant: %v", label, got, want)
	}
}

func TestWidgetV3GeneratedAPIReferenceIncludesDescriptorViews(t *testing.T) {
	md := WidgetV3APIReferenceMarkdown()
	for _, fragment := range []string{
		"# widget.dsl API Reference",
		"## `page`",
		"## `ui` — UINamespace",
		"- `emptyState` (function)",
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

func TestWidgetV3EmbeddedAPIHelpMatchesDescriptorReference(t *testing.T) {
	path := filepath.Join("..", "xgoja", "providers", "widgetsite", "doc", "05-widget-dsl-v3-api-reference.md")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read embedded API help %s: %v", path, err)
	}
	const frontmatterEnd = "---\n\n"
	parts := strings.SplitN(string(data), frontmatterEnd, 2)
	if len(parts) != 2 {
		t.Fatalf("embedded API help %s has no complete frontmatter", path)
	}
	want := strings.TrimPrefix(WidgetV3APIReferenceMarkdown(), "# widget.dsl API Reference\n\n")
	if !strings.HasPrefix(parts[1], want) {
		t.Fatalf("embedded API help descriptor reference is stale; regenerate %s from WidgetV3APIReferenceMarkdown", path)
	}
}
