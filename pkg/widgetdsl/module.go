package widgetdsl

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/go-go-golems/go-go-goja/modules"
)

const ModuleName = "widget.dsl"

var templatePattern = regexp.MustCompile(`\$\{([^}]+)\}|\$([A-Za-z0-9_.-]+)`)

var componentNames = []string{
	"appShell",
	"appNav",
	"button",
	"caption",
	"codeText",
	"contextKindSwatch",
	"annotationBadge",
	"contextLegend",
	"contextBudgetBar",
	"contextStripDiagram",
	"contextStackDiagram",
	"contextTreemap",
	"contextDiagramPanel",
	"dashboardGrid",
	"dataTable",
	"divider",
	"formRow",
	"inline",
	"metadataGrid",
	"panel",
	"scrollRegion",
	"sectionBlock",
	"selectInput",
	"sidebarShell",
	"splitPane",
	"stack",
	"statusText",
	"tabList",
	"textBlock",
	"textInput",
	"transcriptRoleBadge",
	"transcriptSessionHeader",
	"transcriptMessageCard",
	"annotationNoteCard",
	"annotationRailPanel",
	"transcriptReaderPanel",
	"transcriptWorkspacePanel",
	"anchoredCommentCard",
	"anchoredCommentRail",
}

var componentTypes = map[string]string{
	"appShell":                 "AppShell",
	"appNav":                   "AppNav",
	"button":                   "Button",
	"caption":                  "Caption",
	"codeText":                 "CodeText",
	"contextKindSwatch":        "ContextKindSwatch",
	"annotationBadge":          "AnnotationBadge",
	"contextLegend":            "ContextLegend",
	"contextBudgetBar":         "ContextBudgetBar",
	"contextStripDiagram":      "ContextStripDiagram",
	"contextStackDiagram":      "ContextStackDiagram",
	"contextTreemap":           "ContextTreemap",
	"contextDiagramPanel":      "ContextDiagramPanel",
	"dashboardGrid":            "DashboardGrid",
	"dataTable":                "DataTable",
	"divider":                  "Divider",
	"formRow":                  "FormRow",
	"inline":                   "Inline",
	"metadataGrid":             "MetadataGrid",
	"panel":                    "Panel",
	"scrollRegion":             "ScrollRegion",
	"sectionBlock":             "SectionBlock",
	"selectInput":              "SelectInput",
	"sidebarShell":             "SidebarShell",
	"splitPane":                "SplitPane",
	"stack":                    "Stack",
	"statusText":               "StatusText",
	"tabList":                  "TabList",
	"textBlock":                "Text",
	"textInput":                "TextInput",
	"transcriptRoleBadge":      "TranscriptRoleBadge",
	"transcriptSessionHeader":  "TranscriptSessionHeader",
	"transcriptMessageCard":    "TranscriptMessageCard",
	"annotationNoteCard":       "AnnotationNoteCard",
	"annotationRailPanel":      "AnnotationRailPanel",
	"transcriptReaderPanel":    "TranscriptReaderPanel",
	"transcriptWorkspacePanel": "TranscriptWorkspacePanel",
	"anchoredCommentCard":      "AnchoredCommentCard",
	"anchoredCommentRail":      "AnchoredCommentRail",
}

type module struct{}

var _ modules.NativeModule = (*module)(nil)

func NewLoader() require.ModuleLoader {
	mod := &module{}
	return mod.Loader
}

func Register(reg *require.Registry) {
	if reg == nil {
		return
	}
	loader := NewLoader()
	reg.RegisterNativeModule(ModuleName, loader)
	reg.RegisterNativeModule("rag.dsl", loader)
}

func (module) Name() string { return ModuleName }

func (module) Doc() string {
	return `
widget.dsl is a RAG Evaluation System authoring DSL for producing JSON-compatible
Widget IR. The returned IR is intended to be rendered by the React WidgetRenderer,
not by Go-side HTML generation.

Example:

  const rag = require("widget.dsl");
  exports.page = () => rag.panel({title: "Demo"},
    rag.caption({tone: "success"}, "Rendered by React")
  );
`
}

func (m *module) Loader(vm *goja.Runtime, moduleObj *goja.Object) {
	rt := &runtime{vm: vm}
	exports := moduleObj.Get("exports").(*goja.Object)
	rt.install(exports)
}

func init() {
	modules.Register(&module{})
}

type runtime struct {
	vm *goja.Runtime
}

func (r *runtime) install(exports *goja.Object) {
	setExport(exports, "text", r.text)
	setExport(exports, "element", r.element)
	setExport(exports, "component", r.component)
	setExport(exports, "fragment", r.fragment)
	setExport(exports, "page", r.page)

	for _, name := range componentNames {
		componentType := componentTypes[name]
		setExport(exports, name, r.componentFactory(componentType))
	}

	cell := r.vm.NewObject()
	setExport(cell, "field", func(field string, options ...goja.Value) map[string]any {
		out := map[string]any{"kind": "field", "field": field}
		mergeOptions(out, exportOptions(options))
		return out
	})
	setExport(cell, "number", func(field string, options ...goja.Value) map[string]any {
		out := map[string]any{"kind": "number", "field": field}
		mergeOptions(out, exportOptions(options))
		return out
	})
	setExport(cell, "status", func(field string, options ...goja.Value) map[string]any {
		out := map[string]any{"kind": "status", "field": field}
		mergeOptions(out, exportOptions(options))
		return out
	})
	setExport(cell, "caption", func(field string, options ...goja.Value) map[string]any {
		out := map[string]any{"kind": "caption", "field": field}
		mergeOptions(out, exportOptions(options))
		return out
	})
	setExport(cell, "template", func(template string) map[string]any {
		return map[string]any{"kind": "template", "template": template}
	})
	setExport(cell, "link", func(hrefField string, labelField string, options ...goja.Value) map[string]any {
		out := map[string]any{"kind": "link", "hrefField": hrefField, "labelField": labelField}
		mergeOptions(out, exportOptions(options))
		return out
	})
	setExport(cell, "constant", func(value goja.Value) any {
		return map[string]any{"kind": "constant", "value": r.exportRenderable(value)}
	})
	setExport(exports, "cell", cell)

	action := r.vm.NewObject()
	setExport(action, "server", func(name string, options ...goja.Value) map[string]any {
		out := map[string]any{"kind": "server", "name": name}
		mergeOptions(out, exportOptions(options))
		return out
	})
	setExport(action, "navigate", func(to string, options ...goja.Value) map[string]any {
		out := map[string]any{"kind": "navigate", "to": to}
		mergeOptions(out, exportOptions(options))
		return out
	})
	setExport(action, "event", func(event string, options ...goja.Value) map[string]any {
		out := map[string]any{"kind": "event", "event": event}
		mergeOptions(out, exportOptions(options))
		return out
	})
	setExport(action, "copy", func(value string) map[string]any {
		return map[string]any{"kind": "copy", "value": value}
	})
	setExport(exports, "action", action)

	recipes := r.vm.NewObject()
	setExport(recipes, "metrics", r.metricsRecipe)
	setExport(recipes, "actionToolbar", r.actionToolbarRecipe)
	setExport(recipes, "masterDetailTable", r.masterDetailTableRecipe)
	setExport(exports, "recipes", recipes)
}

func (r *runtime) text(value goja.Value) map[string]any {
	return map[string]any{"kind": "text", "text": stringifyValue(value)}
}

func (r *runtime) element(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) == 0 {
		panic(r.vm.NewGoError(fmt.Errorf("widget.dsl element(tag, attrs?, ...children) requires a tag")))
	}
	tag := strings.TrimSpace(call.Arguments[0].String())
	if tag == "" {
		panic(r.vm.NewGoError(fmt.Errorf("widget.dsl element tag must not be empty")))
	}
	attrs := map[string]any{}
	childStart := 1
	if len(call.Arguments) > 1 && isPlainObject(call.Arguments[1]) && !looksLikeWidgetNodeExport(call.Arguments[1]) {
		attrs = exportObject(call.Arguments[1])
		childStart = 2
	}
	out := map[string]any{"kind": "element", "tag": tag}
	if len(attrs) > 0 {
		out["attrs"] = attrs
	}
	children := r.exportChildren(call.Arguments[childStart:])
	if len(children) > 0 {
		out["children"] = children
	}
	return r.vm.ToValue(out)
}

func (r *runtime) component(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) == 0 {
		panic(r.vm.NewGoError(fmt.Errorf("widget.dsl component(type, props?, ...children) requires a type")))
	}
	componentType := strings.TrimSpace(call.Arguments[0].String())
	if componentType == "" {
		panic(r.vm.NewGoError(fmt.Errorf("widget.dsl component type must not be empty")))
	}
	props, childStart := propsAndChildStart(call.Arguments, 1)
	return r.vm.ToValue(r.buildComponent(componentType, props, call.Arguments[childStart:]))
}

func (r *runtime) componentFactory(componentType string) func(goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		props, childStart := propsAndChildStart(call.Arguments, 0)
		return r.vm.ToValue(r.buildComponent(componentType, props, call.Arguments[childStart:]))
	}
}

func (r *runtime) fragment(call goja.FunctionCall) []any {
	return r.exportChildren(call.Arguments)
}

func (r *runtime) page(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) == 0 || !isPlainObject(call.Arguments[0]) {
		panic(r.vm.NewGoError(fmt.Errorf("widget.dsl page(options) requires an options object")))
	}
	options := exportObject(call.Arguments[0])
	id := stringFromMap(options, "id", "page")
	title := stringFromMap(options, "title", id)
	out := map[string]any{
		"schemaVersion": stringFromMap(options, "schemaVersion", "0.1.0"),
		"id":            id,
		"title":         title,
	}
	if meta, ok := options["meta"].(map[string]any); ok && len(meta) > 0 {
		out["meta"] = meta
	}
	if root, ok := options["root"].(map[string]any); ok && isWidgetNodeExport(root) {
		out["root"] = root
		return r.vm.ToValue(out)
	}
	sections := anySlice(options["sections"])
	children := []any{}
	for _, section := range sections {
		if node, ok := section.(map[string]any); ok && isWidgetNodeExport(node) {
			children = append(children, node)
		}
	}
	out["root"] = map[string]any{
		"kind":     "component",
		"type":     "Stack",
		"props":    map[string]any{"gap": stringFromMap(options, "gap", "lg")},
		"children": children,
	}
	return r.vm.ToValue(out)
}

func (r *runtime) metricsRecipe(call goja.FunctionCall) goja.Value {
	options := firstObject(call.Arguments)
	items := anySlice(options["items"])
	children := []any{}
	for _, raw := range items {
		item, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		label := stringFromMap(item, "label", "Metric")
		status := stringFromMap(item, "status", "ready")
		value := item["value"]
		children = append(children, map[string]any{
			"kind":  "component",
			"type":  "Panel",
			"props": map[string]any{"title": label, "density": "condensed"},
			"children": []any{map[string]any{
				"kind":     "component",
				"type":     "StatusText",
				"props":    map[string]any{"status": status, "icon": boolFromMap(item, "icon", true)},
				"children": []any{map[string]any{"kind": "text", "text": fmt.Sprint(value)}},
			}},
		})
	}
	return r.vm.ToValue(map[string]any{
		"kind":     "component",
		"type":     "DashboardGrid",
		"props":    map[string]any{"recipe": stringFromMap(options, "recipe", metricsRecipeName(len(children)))},
		"children": children,
	})
}

func (r *runtime) actionToolbarRecipe(call goja.FunctionCall) goja.Value {
	options := firstObject(call.Arguments)
	actions := anySlice(options["actions"])
	inlineChildren := []any{}
	for _, raw := range actions {
		item, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		label := stringFromMap(item, "label", stringFromMap(item, "name", "Action"))
		props := map[string]any{"variant": stringFromMap(item, "variant", "secondary")}
		if act, ok := normalizeActionSpec(item["action"], item["name"], item["payload"]); ok {
			props["action"] = act
		}
		inlineChildren = append(inlineChildren, map[string]any{
			"kind":     "component",
			"type":     "Button",
			"props":    props,
			"children": []any{map[string]any{"kind": "text", "text": label}},
		})
	}
	if caption, ok := options["caption"].(string); ok && caption != "" {
		inlineChildren = append(inlineChildren, map[string]any{
			"kind":     "component",
			"type":     "Caption",
			"props":    map[string]any{"tone": stringFromMap(options, "captionTone", "muted")},
			"children": []any{map[string]any{"kind": "text", "text": caption}},
		})
	}
	return r.vm.ToValue(map[string]any{
		"kind":  "component",
		"type":  "Panel",
		"props": map[string]any{"title": stringFromMap(options, "title", "Actions")},
		"children": []any{map[string]any{
			"kind":     "component",
			"type":     "Inline",
			"props":    map[string]any{"gap": stringFromMap(options, "gap", "sm"), "wrap": boolFromMap(options, "wrap", true)},
			"children": inlineChildren,
		}},
	})
}

func (r *runtime) masterDetailTableRecipe(call goja.FunctionCall) goja.Value {
	options := firstObject(call.Arguments)
	rows := anySlice(options["rows"])
	selectedKey := options["selectedKey"]
	tableProps := map[string]any{
		"rows":         rows,
		"getRowKey":    valueOrDefault(options["getRowKey"], "id"),
		"columns":      anySlice(options["columns"]),
		"selectedKey":  selectedKey,
		"emptyMessage": valueOrDefault(options["emptyMessage"], "No rows"),
	}
	if act, ok := normalizeActionSpec(options["onRowSelect"], nil, nil); ok {
		tableProps["onRowSelect"] = act
	}
	tablePanel := map[string]any{
		"kind":     "component",
		"type":     "Panel",
		"props":    map[string]any{"title": stringFromMap(options, "title", "Items")},
		"children": []any{map[string]any{"kind": "component", "type": "DataTable", "props": tableProps}},
	}
	detailNode := r.detailNode(options, selectedRow(rows, selectedKey, tableProps["getRowKey"]))
	return r.vm.ToValue(map[string]any{
		"kind":     "component",
		"type":     "DashboardGrid",
		"props":    map[string]any{"recipe": stringFromMap(options, "recipe", "two-up")},
		"children": []any{tablePanel, detailNode},
	})
}

func (r *runtime) detailNode(options map[string]any, row any) any {
	if detailFn, ok := goja.AssertFunction(r.vm.ToValue(options["detail"])); ok {
		value, err := detailFn(goja.Undefined(), r.vm.ToValue(row))
		if err != nil {
			panic(err)
		}
		if exported, ok := value.Export().(map[string]any); ok && isWidgetNodeExport(exported) {
			return exported
		}
	}
	if fallback, ok := options["detail"].(map[string]any); ok && isWidgetNodeExport(fallback) {
		return fallback
	}
	return map[string]any{
		"kind":     "component",
		"type":     "Panel",
		"props":    map[string]any{"title": stringFromMap(options, "detailTitle", "Details")},
		"children": []any{map[string]any{"kind": "component", "type": "Caption", "props": map[string]any{"tone": "muted"}, "children": []any{map[string]any{"kind": "text", "text": "Select a row to inspect it."}}}},
	}
}

func (r *runtime) buildComponent(componentType string, props map[string]any, childValues []goja.Value) map[string]any {
	out := map[string]any{"kind": "component", "type": componentType}
	if len(props) > 0 {
		out["props"] = props
	}
	children := r.exportChildren(childValues)
	if len(children) > 0 {
		out["children"] = children
	}
	return out
}

func (r *runtime) exportChildren(values []goja.Value) []any {
	out := []any{}
	for _, value := range values {
		out = append(out, r.exportChild(value)...)
	}
	return out
}

func (r *runtime) exportChild(value goja.Value) []any {
	if value == nil || goja.IsUndefined(value) || goja.IsNull(value) {
		return nil
	}
	if isArrayLike(value) {
		obj := value.ToObject(r.vm)
		length := int(obj.Get("length").ToInteger())
		out := []any{}
		for i := 0; i < length; i++ {
			out = append(out, r.exportChild(obj.Get(fmt.Sprintf("%d", i)))...)
		}
		return out
	}
	if isWidgetNode(r.vm, value) {
		return []any{value.Export()}
	}
	return []any{map[string]any{"kind": "text", "text": stringifyValue(value)}}
}

func (r *runtime) exportRenderable(value goja.Value) any {
	children := r.exportChild(value)
	if len(children) == 0 {
		return nil
	}
	if len(children) == 1 {
		return children[0]
	}
	return children
}

func propsAndChildStart(args []goja.Value, index int) (map[string]any, int) {
	if len(args) > index && isPlainObject(args[index]) && !looksLikeWidgetNodeExport(args[index]) {
		return exportObject(args[index]), index + 1
	}
	return map[string]any{}, index
}

func isPlainObject(value goja.Value) bool {
	if value == nil || goja.IsUndefined(value) || goja.IsNull(value) {
		return false
	}
	exported := value.Export()
	_, ok := exported.(map[string]any)
	return ok
}

func isWidgetNode(vm *goja.Runtime, value goja.Value) bool {
	if value == nil || goja.IsUndefined(value) || goja.IsNull(value) {
		return false
	}
	obj := value.ToObject(vm)
	kind := obj.Get("kind")
	if kind == nil || goja.IsUndefined(kind) || goja.IsNull(kind) {
		return false
	}
	s := kind.String()
	return s == "text" || s == "element" || s == "component"
}

func looksLikeWidgetNodeExport(value goja.Value) bool {
	exported, ok := value.Export().(map[string]any)
	if !ok {
		return false
	}
	kind, _ := exported["kind"].(string)
	return kind == "text" || kind == "element" || kind == "component"
}

func isArrayLike(value goja.Value) bool {
	if value == nil || goja.IsUndefined(value) || goja.IsNull(value) {
		return false
	}
	t := value.ExportType()
	if t == nil {
		return false
	}
	k := t.Kind()
	return k == reflect.Slice || k == reflect.Array
}

func exportObject(value goja.Value) map[string]any {
	if value == nil || goja.IsUndefined(value) || goja.IsNull(value) {
		return map[string]any{}
	}
	if exported, ok := value.Export().(map[string]any); ok {
		return exported
	}
	return map[string]any{}
}

func exportOptions(options []goja.Value) map[string]any {
	if len(options) == 0 {
		return map[string]any{}
	}
	return exportObject(options[0])
}

func mergeOptions(out map[string]any, options map[string]any) {
	for k, v := range options {
		out[k] = v
	}
}

func firstObject(args []goja.Value) map[string]any {
	if len(args) == 0 {
		return map[string]any{}
	}
	return exportObject(args[0])
}

func stringFromMap(m map[string]any, key, fallback string) string {
	if value, ok := m[key]; ok {
		if s, ok := value.(string); ok && strings.TrimSpace(s) != "" {
			return s
		}
	}
	return fallback
}

func boolFromMap(m map[string]any, key string, fallback bool) bool {
	if value, ok := m[key]; ok {
		if b, ok := value.(bool); ok {
			return b
		}
	}
	return fallback
}

func anySlice(value any) []any {
	if value == nil {
		return []any{}
	}
	if out, ok := value.([]any); ok {
		return out
	}
	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return []any{}
	}
	out := make([]any, 0, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		out = append(out, rv.Index(i).Interface())
	}
	return out
}

func valueOrDefault(value any, fallback any) any {
	if value == nil {
		return fallback
	}
	return value
}

func isWidgetNodeExport(exported map[string]any) bool {
	kind, _ := exported["kind"].(string)
	return kind == "text" || kind == "element" || kind == "component"
}

func metricsRecipeName(count int) string {
	if count <= 2 {
		return "two-up"
	}
	if count == 3 {
		return "three-up"
	}
	return "four-up"
}

func normalizeActionSpec(action any, name any, payload any) (map[string]any, bool) {
	if spec, ok := action.(map[string]any); ok {
		if kind, _ := spec["kind"].(string); kind != "" {
			return spec, true
		}
	}
	if actionName, ok := action.(string); ok && strings.TrimSpace(actionName) != "" {
		out := map[string]any{"kind": "server", "name": actionName}
		if payload != nil {
			out["payload"] = payload
		}
		return out, true
	}
	if actionName, ok := name.(string); ok && strings.TrimSpace(actionName) != "" {
		out := map[string]any{"kind": "server", "name": actionName}
		if payload != nil {
			out["payload"] = payload
		}
		return out, true
	}
	return nil, false
}

func selectedRow(rows []any, selectedKey any, keySpec any) any {
	if selectedKey == nil {
		return nil
	}
	wanted := stringifyRowKey(selectedKey)
	for _, raw := range rows {
		row, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		if rowKeyForSpec(row, keySpec) == wanted {
			return row
		}
	}
	return nil
}

func rowKeyForSpec(row map[string]any, spec any) string {
	switch typed := spec.(type) {
	case string:
		return stringifyRowKey(valueAtPath(row, typed))
	case map[string]any:
		if field, ok := typed["field"].(string); ok {
			return stringifyRowKey(valueAtPath(row, field))
		}
		if template, ok := typed["template"].(string); ok {
			return renderRowTemplate(template, row)
		}
	}
	return stringifyRowKey(valueAtPath(row, "id"))
}

func stringifyRowKey(value any) string {
	if value == nil {
		return ""
	}
	return fmt.Sprint(value)
}

func valueAtPath(row map[string]any, path string) any {
	if path == "" {
		return nil
	}
	var current any = row
	for _, part := range strings.Split(path, ".") {
		if part == "" {
			continue
		}
		object, ok := current.(map[string]any)
		if !ok {
			return nil
		}
		current = object[part]
	}
	return current
}

func renderRowTemplate(template string, row map[string]any) string {
	return templatePattern.ReplaceAllStringFunc(template, func(match string) string {
		path := strings.TrimPrefix(match, "$")
		path = strings.TrimPrefix(strings.TrimSuffix(path, "}"), "{")
		value := valueAtPath(row, path)
		if value == nil {
			return ""
		}
		return stringifyRowKey(value)
	})
}

func stringifyValue(value goja.Value) string {
	if value == nil || goja.IsUndefined(value) || goja.IsNull(value) {
		return ""
	}
	return value.String()
}

func setExport(obj *goja.Object, name string, value any) {
	if err := obj.Set(name, value); err != nil {
		panic(err)
	}
}
