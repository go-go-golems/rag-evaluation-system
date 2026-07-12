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

func TestWidgetV3DescriptorMatchesNestedNamespaces(t *testing.T) {
	vm := goja.New()
	exports := vm.NewObject()
	(&runtime{vm: vm}).installWidgetV3(exports)

	for _, namespace := range widgetV3Module.NestedNamespaces {
		object := resolveV3DescriptorPath(t, vm, exports, namespace.Path)
		assertSameStringSet(t, namespace.Path+" members", object.Keys(), v3DescriptorMemberNames(namespace.Members))
	}
}

func TestWidgetV3DescriptorMatchesBuilderRuntimeMethods(t *testing.T) {
	vm := goja.New()
	exports := vm.NewObject()
	(&runtime{vm: vm}).installWidgetV3(exports)
	if err := vm.Set("widget", exports); err != nil {
		t.Fatalf("install widget global: %v", err)
	}

	value, err := vm.RunString(`
		const builders = {};
		const capture = (name) => (builder) => { builders[name] = Object.keys(builder); };

		widget.page("Probe", page => {
			capture("PageBuilder")(page);
			page.section("Section", section => {
				capture("SectionBuilder")(section);
				section.actions(capture("ActionsBuilder"));
			});
		});
		capture("FieldSetBuilder")(widget.data.fields());
		const collection = widget.data.collection([]);
		capture("CollectionBuilder")(collection);
		collection.table(capture("TableBuilder"));
		collection.edit(capture("EditorBuilder"));
		widget.data.matrix([], capture("MatrixBuilder"));
		widget.schedule.availabilityPoll({ options: [], responses: [] }, capture("SchedulePollBuilder"));
		widget.time.month([], capture("TimeMonthBuilder"));
		widget.time.week([], capture("TimeWeekBuilder"));
		widget.context.styleSet(capture("ContextStyleSetBuilder"));
		widget.context.diagram({ parts: [] }, capture("ContextDiagramBuilder"));
		widget.context.workspace({ messages: [], annotations: [] }, capture("ContextWorkspaceBuilder"));
		widget.course.shell({}, capture("CourseShellBuilder"));
		widget.course.landing({}, capture("CourseLandingBuilder"));
		widget.course.slideDeck({ slides: [] }, capture("CourseSlideDeckBuilder"));
		widget.course.handouts({ docs: [] }, capture("CourseHandoutsBuilder"));
		widget.course.metadataForm({}, capture("CourseMetadataFormBuilder"));
		widget.course.materialUploads({}, capture("CourseMaterialUploadsBuilder"));
		widget.cms.mediaLibrary([], capture("CmsMediaLibraryBuilder"));
		widget.cms.articleQueue([], capture("CmsArticleQueueBuilder"));
		widget.cms.markdownEditor("", capture("CmsMarkdownEditorBuilder"));
		const crmFields = widget.crm.fields();
		capture("CrmFieldsBuilder")(crmFields);
		const pipeline = widget.crm.pipeline("Probe");
		capture("CrmPipelineBuilder")(pipeline);
		widget.crm.pipelineBoard(pipeline, [], capture("CrmPipelineBoardBuilder"));
		widget.crm.recordFields({}, crmFields, capture("CrmRecordFieldsBuilder"));
		widget.crm.activityFeed([], capture("CrmActivityFeedBuilder"));
		builders;
	`)
	if err != nil {
		t.Fatalf("probe v3 builders: %v", err)
	}
	got := value.Export().(map[string]any)
	for _, builder := range widgetV3Module.Builders {
		rawMethods, ok := got[builder.TypeName].([]any)
		if !ok {
			t.Fatalf("builder probe did not capture %s (got %#v)", builder.TypeName, got[builder.TypeName])
		}
		methods := make([]string, 0, len(rawMethods))
		for _, method := range rawMethods {
			methods = append(methods, method.(string))
		}
		assertSameStringSet(t, builder.TypeName+" methods", methods, builder.Methods)
	}
	if len(got) != len(widgetV3Module.Builders) {
		t.Fatalf("builder descriptor count mismatch: runtime captured %d, descriptors contain %d", len(got), len(widgetV3Module.Builders))
	}
}

func TestWidgetV3ComposableBuilderDeclarationsMatchDescriptors(t *testing.T) {
	dts := strings.Join(TypeScriptModule(WidgetV3ModuleName).RawDTS, "\n")
	for _, builder := range widgetV3Module.Builders {
		want := "export interface " + builder.TypeName + " extends ComposableBuilder<" + builder.TypeName + ">"
		if !strings.Contains(dts, want) {
			t.Errorf("builder descriptor %s is not declared as composable", builder.TypeName)
		}
	}
}

func TestWidgetV3BuilderUseComposesFragments(t *testing.T) {
	vm := goja.New()
	exports := vm.NewObject()
	(&runtime{vm: vm}).installWidgetV3(exports)
	if err := vm.Set("widget", exports); err != nil {
		t.Fatalf("install widget global: %v", err)
	}
	value, err := vm.RunString(`
		let fragmentCalls = 0;
		const compact = table => { fragmentCalls += 1; return table.className("compact"); };
		widget.data.collection([{ id: "one", title: "One" }], collection => collection
			.use(c => { fragmentCalls += 1; return c.empty("None"); })
			.table(table => table.use(compact)))
			.toNode();
		fragmentCalls;
	`)
	if err != nil {
		t.Fatalf("compose builder fragments: %v", err)
	}
	if got := value.ToInteger(); got != 2 {
		t.Fatalf("composed fragment calls = %d, want 2", got)
	}
}

func resolveV3DescriptorPath(t *testing.T, vm *goja.Runtime, root *goja.Object, path string) *goja.Object {
	t.Helper()
	object := root
	for _, segment := range strings.Split(path, ".") {
		value := object.Get(segment)
		if value == nil || goja.IsUndefined(value) || goja.IsNull(value) {
			t.Fatalf("descriptor path %q is missing segment %q", path, segment)
		}
		object = value.ToObject(vm)
	}
	return object
}

func v3DescriptorMemberNames(members []v3MemberDescriptor) []string {
	names := make([]string, 0, len(members))
	for _, member := range members {
		names = append(names, member.Name)
	}
	return names
}

func TestWidgetV3ActionContextDescriptorsAreUniqueAndIdentifyComponents(t *testing.T) {
	seen := map[string]struct{}{}
	for _, context := range widgetV3Module.ActionContexts {
		if _, ok := seen[context.Name]; ok {
			t.Errorf("duplicate action context descriptor %q", context.Name)
		}
		seen[context.Name] = struct{}{}
		if context.Component == "" {
			t.Errorf("action context %q has no component", context.Name)
		}
		foundComponentType := false
		for _, field := range context.Fields {
			if field == "componentType" {
				foundComponentType = true
				break
			}
		}
		if !foundComponentType {
			t.Errorf("action context %q does not document componentType", context.Name)
		}
	}
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
