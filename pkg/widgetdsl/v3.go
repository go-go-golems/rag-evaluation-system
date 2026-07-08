package widgetdsl

import (
	"fmt"
	"strings"

	"github.com/dop251/goja"
)

type v3PageSpec struct {
	SchemaVersion string
	ID            string
	Title         string
	Meta          map[string]any
	Shell         any
	Density       string
	Breadcrumbs   []map[string]any
	Sections      []v3SectionSpec
}

type v3SectionSpec struct {
	Title    any
	Caption  string
	AnchorID string
	Tone     string
	Actions  []any
	Children []v3NodeSpec
}

type v3NodeSpec struct {
	Kind   string
	IR     map[string]any
	Source *v3SourceSpan
}

type v3SourceSpan struct {
	File   string
	Line   int
	Column int
}

type v3SlotSpec struct {
	Function goja.Value
	Fallback goja.Value
}

type v3SelectionSpec struct {
	Mode     string
	KeyField string
	Selected any
}

type v3ListItemSpec struct {
	ID       string
	Label    any
	Icon     any
	Badge    any
	Disabled bool
	Extra    map[string]any
}

func (r *runtime) v3Page(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) == 0 {
		panic(r.vm.NewGoError(fmt.Errorf("widget.dsl page(titleOrOptions, configure?) requires a title string or options object")))
	}
	spec := &v3PageSpec{SchemaVersion: "0.1.0", ID: "page", Title: "Page", Meta: map[string]any{}}
	first := call.Arguments[0]
	if isPlainObject(first) && !looksLikeWidgetNodeExport(first) {
		options := exportObject(first)
		spec.ID = stringFromMap(options, "id", spec.ID)
		spec.Title = stringFromMap(options, "title", spec.Title)
		spec.SchemaVersion = stringFromMap(options, "schemaVersion", spec.SchemaVersion)
		if meta, ok := options["meta"].(map[string]any); ok {
			spec.Meta = meta
		}
	} else {
		title := strings.TrimSpace(first.String())
		if title == "" {
			panic(r.vm.NewGoError(fmt.Errorf("widget.dsl page title must not be empty")))
		}
		spec.Title = title
		spec.ID = slugID(title)
	}
	builder := r.v3PageBuilder(spec)
	if len(call.Arguments) > 1 {
		r.applyV3BuilderCallback(builder, call.Arguments[1], "page")
	}
	return builder
}

func (r *runtime) v3UIObject() *goja.Object {
	ui := r.vm.NewObject()
	setExport(ui, "callout", r.v3ComponentFactory("Panel", map[string]any{"tone": "callout"}))
	setExport(ui, "stack", r.v3ComponentFactory("Stack", nil))
	setExport(ui, "inline", r.v3ComponentFactory("Inline", nil))
	setExport(ui, "card", r.v3ComponentFactory("Panel", nil))
	setExport(ui, "button", r.v3UIButton)
	setExport(ui, "caption", r.v3ComponentFactory("Caption", nil))
	setExport(ui, "badge", r.v3ComponentFactory("Tag", nil))
	setExport(ui, "metadata", r.v3UIMetadata)
	setExport(ui, "form", r.v3ComponentFactory("FormPanel", nil))
	return ui
}

func (r *runtime) v3DataObject() *goja.Object {
	data := r.vm.NewObject()
	setExport(data, "selection", r.v3Selection)
	setExport(data, "item", r.v3ListItem)
	return data
}

func (r *runtime) v3ComponentFactory(componentType string, defaults map[string]any) func(goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		props, childStart := propsAndChildStart(call.Arguments, 0)
		if len(defaults) > 0 {
			merged := map[string]any{}
			for key, value := range defaults {
				merged[key] = value
			}
			for key, value := range props {
				merged[key] = value
			}
			props = merged
		}
		return r.vm.ToValue(r.v3BuildComponent(componentType, props, call.Arguments[childStart:]))
	}
}

func (r *runtime) v3UIButton(label goja.Value, action goja.Value, options ...goja.Value) map[string]any {
	props := exportOptions(options)
	if action != nil && !goja.IsUndefined(action) && !goja.IsNull(action) {
		props["action"] = action.Export()
	}
	return componentNode("Button", props, r.v3NodeSpecsToIR(r.v3ExportChild(label))...)
}

func (r *runtime) v3UIMetadata(record goja.Value, options ...goja.Value) map[string]any {
	props := exportOptions(options)
	props["items"] = v3MetadataItems(exportObject(record))
	return componentNode("MetadataGrid", props)
}

func (r *runtime) v3ActionsBuilder(actions *[]any) *goja.Object {
	obj := r.vm.NewObject()
	setExport(obj, "add", func(label goja.Value, action goja.Value, options ...goja.Value) *goja.Object {
		item := exportOptions(options)
		item["label"] = r.v3Renderable(label)
		if action != nil && !goja.IsUndefined(action) && !goja.IsNull(action) {
			item["action"] = action.Export()
		}
		*actions = append(*actions, item)
		return obj
	})
	setExport(obj, "button", obj.Get("add"))
	return obj
}

func (r *runtime) v3Selection(modeOrOptions goja.Value, options ...goja.Value) map[string]any {
	spec := v3SelectionSpec{Mode: "single"}
	if isPlainObject(modeOrOptions) {
		opts := exportObject(modeOrOptions)
		spec.Mode = stringFromMap(opts, "mode", spec.Mode)
		spec.KeyField = stringFromMap(opts, "keyField", spec.KeyField)
		spec.Selected = opts["selected"]
	} else if modeOrOptions != nil && !goja.IsUndefined(modeOrOptions) && !goja.IsNull(modeOrOptions) {
		spec.Mode = strings.TrimSpace(modeOrOptions.String())
		opts := exportOptions(options)
		spec.KeyField = stringFromMap(opts, "keyField", spec.KeyField)
		spec.Selected = opts["selected"]
	}
	if spec.Mode != "single" && spec.Mode != "multi" {
		panic(r.vm.NewGoError(fmt.Errorf("widget.dsl data.selection mode must be single or multi")))
	}
	out := map[string]any{"kind": "selection", "mode": spec.Mode}
	if spec.KeyField != "" {
		out["keyField"] = spec.KeyField
	}
	if spec.Selected != nil {
		out["selected"] = spec.Selected
	}
	return out
}

func (r *runtime) v3ListItem(id string, label goja.Value, options ...goja.Value) map[string]any {
	if strings.TrimSpace(id) == "" {
		panic(r.vm.NewGoError(fmt.Errorf("widget.dsl data.item id must not be empty")))
	}
	spec := v3ListItemSpec{ID: id, Label: r.v3Renderable(label), Extra: exportOptions(options)}
	out := map[string]any{"kind": "listItem", "id": spec.ID, "label": spec.Label}
	for key, value := range spec.Extra {
		out[key] = value
	}
	if spec.Icon != nil {
		out["icon"] = spec.Icon
	}
	if spec.Badge != nil {
		out["badge"] = spec.Badge
	}
	if spec.Disabled {
		out["disabled"] = true
	}
	return out
}

func (r *runtime) v3PageBuilder(spec *v3PageSpec) *goja.Object {
	obj := r.vm.NewObject()
	setExport(obj, "id", func(id string) *goja.Object {
		if strings.TrimSpace(id) != "" {
			spec.ID = id
		}
		return obj
	})
	setExport(obj, "title", func(title string) *goja.Object {
		if strings.TrimSpace(title) != "" {
			spec.Title = title
		}
		return obj
	})
	setExport(obj, "meta", func(key string, value goja.Value) *goja.Object {
		if spec.Meta == nil {
			spec.Meta = map[string]any{}
		}
		spec.Meta[key] = value.Export()
		return obj
	})
	setExport(obj, "shell", func(shell goja.Value) *goja.Object {
		spec.Shell = shell.Export()
		return obj
	})
	setExport(obj, "density", func(density string) *goja.Object {
		spec.Density = density
		return obj
	})
	setExport(obj, "breadcrumb", func(label goja.Value, href ...string) *goja.Object {
		item := map[string]any{"label": r.v3Renderable(label)}
		if len(href) > 0 && strings.TrimSpace(href[0]) != "" {
			item["href"] = href[0]
		}
		spec.Breadcrumbs = append(spec.Breadcrumbs, item)
		return obj
	})
	setExport(obj, "use", func(fragment goja.Value) *goja.Object {
		r.applyV3BuilderCallback(obj, fragment, "page.use")
		return obj
	})
	setExport(obj, "section", func(title goja.Value, cb ...goja.Value) *goja.Object {
		section := v3SectionSpec{Title: r.v3RenderableTitle(title)}
		sectionBuilder := r.v3SectionBuilder(&section)
		if len(cb) > 0 {
			r.applyV3BuilderCallback(sectionBuilder, cb[0], "section")
		}
		spec.Sections = append(spec.Sections, section)
		return obj
	})
	setExport(obj, "view", func(value goja.Value) *goja.Object {
		section := v3SectionSpec{Title: "Content", Children: r.v3ExportChild(value)}
		spec.Sections = append(spec.Sections, section)
		return obj
	})
	setExport(obj, "validate", func() []map[string]any {
		return v3PageValidationIssues(spec)
	})
	setExport(obj, "toPage", func() map[string]any {
		issues := v3PageValidationIssues(spec)
		if len(issues) > 0 {
			panic(r.vm.NewGoError(fmt.Errorf("widget.dsl page is invalid: %s", issues[0]["message"])))
		}
		return r.v3PageToIR(spec)
	})
	return obj
}

func (r *runtime) v3SectionBuilder(spec *v3SectionSpec) *goja.Object {
	obj := r.vm.NewObject()
	setExport(obj, "caption", func(caption string) *goja.Object {
		spec.Caption = caption
		return obj
	})
	setExport(obj, "anchor", func(anchor string) *goja.Object {
		spec.AnchorID = anchor
		return obj
	})
	setExport(obj, "tone", func(tone string) *goja.Object {
		spec.Tone = tone
		return obj
	})
	setExport(obj, "use", func(fragment goja.Value) *goja.Object {
		r.applyV3BuilderCallback(obj, fragment, "section.use")
		return obj
	})
	setExport(obj, "text", func(value goja.Value) *goja.Object {
		spec.Children = append(spec.Children, r.v3TextNode(value))
		return obj
	})
	setExport(obj, "view", func(value goja.Value) *goja.Object {
		spec.Children = append(spec.Children, r.v3ExportChild(value)...)
		return obj
	})
	setExport(obj, "slot", func(context goja.Value, slot goja.Value, fallback ...goja.Value) *goja.Object {
		var fallbackSlot goja.Value
		if len(fallback) > 0 {
			fallbackSlot = fallback[0]
		}
		nodes := r.callV3Slot(v3SlotSpec{Function: slot, Fallback: fallbackSlot}, context.Export())
		spec.Children = append(spec.Children, nodes...)
		return obj
	})
	setExport(obj, "actions", func(cb goja.Value) *goja.Object {
		actions := r.v3ActionsBuilder(&spec.Actions)
		r.applyV3BuilderCallback(actions, cb, "section.actions")
		return obj
	})
	setExport(obj, "metric", func(label goja.Value, value goja.Value, options ...goja.Value) *goja.Object {
		props := exportOptions(options)
		props["label"] = r.v3Renderable(label)
		props["value"] = r.v3Renderable(value)
		spec.Children = append(spec.Children, v3NodeSpecFromIR(componentNode("KeyValueStrip", map[string]any{"items": []any{props}})))
		return obj
	})
	setExport(obj, "metadata", func(record goja.Value) *goja.Object {
		spec.Children = append(spec.Children, v3NodeSpecFromIR(r.v3MetadataNode(exportObject(record))))
		return obj
	})
	return obj
}

func (r *runtime) applyV3BuilderCallback(builder *goja.Object, cb goja.Value, name string) {
	if cb == nil || goja.IsUndefined(cb) || goja.IsNull(cb) {
		return
	}
	fn, ok := goja.AssertFunction(cb)
	if !ok {
		panic(r.vm.NewGoError(fmt.Errorf("widget.dsl %s callback must be a function", name)))
	}
	if _, err := fn(goja.Undefined(), builder); err != nil {
		panic(err)
	}
}

func (r *runtime) callV3Slot(slot v3SlotSpec, ctx any) []v3NodeSpec {
	return r.callV3SlotFunction(slot.Function, ctx, func(any) []v3NodeSpec {
		return r.callV3SlotFunction(slot.Fallback, ctx, nil)
	})
}

func (r *runtime) callV3SlotFunction(slot goja.Value, ctx any, fallback func(any) []v3NodeSpec) []v3NodeSpec {
	if slot == nil || goja.IsUndefined(slot) || goja.IsNull(slot) {
		if fallback == nil {
			return nil
		}
		return fallback(ctx)
	}
	fn, ok := goja.AssertFunction(slot)
	if !ok {
		panic(r.vm.NewGoError(fmt.Errorf("widget.dsl slot must be a function")))
	}
	value, err := fn(goja.Undefined(), r.vm.ToValue(ctx), r.v3SlotHelpers())
	if err != nil {
		panic(err)
	}
	if isV3EmptySlotResult(value) && fallback != nil {
		return fallback(ctx)
	}
	return r.v3ExportChild(value)
}

func (r *runtime) v3SlotHelpers() *goja.Object {
	h := r.vm.NewObject()
	setExport(h, "text", func(value goja.Value) map[string]any {
		return r.v3TextNode(value).toIR()
	})
	setExport(h, "caption", func(value goja.Value, options ...goja.Value) map[string]any {
		props := exportOptions(options)
		return componentNode("Caption", props, r.v3NodeSpecsToIR(r.v3ExportChild(value))...)
	})
	setExport(h, "strong", func(call goja.FunctionCall) goja.Value {
		return r.vm.ToValue(r.v3BuildElement("strong", map[string]any{}, call.Arguments))
	})
	setExport(h, "stack", func(call goja.FunctionCall) goja.Value {
		props, childStart := propsAndChildStart(call.Arguments, 0)
		return r.vm.ToValue(r.v3BuildComponent("Stack", props, call.Arguments[childStart:]))
	})
	setExport(h, "inline", func(call goja.FunctionCall) goja.Value {
		props, childStart := propsAndChildStart(call.Arguments, 0)
		return r.vm.ToValue(r.v3BuildComponent("Inline", props, call.Arguments[childStart:]))
	})
	setExport(h, "card", func(call goja.FunctionCall) goja.Value {
		props, childStart := propsAndChildStart(call.Arguments, 0)
		return r.vm.ToValue(r.v3BuildComponent("Panel", props, call.Arguments[childStart:]))
	})
	setExport(h, "button", func(label goja.Value, action goja.Value, options ...goja.Value) map[string]any {
		props := exportOptions(options)
		if action != nil && !goja.IsUndefined(action) && !goja.IsNull(action) {
			props["action"] = action.Export()
		}
		return componentNode("Button", props, r.v3NodeSpecsToIR(r.v3ExportChild(label))...)
	})
	setExport(h, "badge", func(value goja.Value, options ...goja.Value) map[string]any {
		props := exportOptions(options)
		return componentNode("Tag", props, r.v3NodeSpecsToIR(r.v3ExportChild(value))...)
	})
	setExport(h, "raw", r.v3RawObject())
	return h
}

func (r *runtime) v3RawObject() *goja.Object {
	raw := r.vm.NewObject()
	setExport(raw, "text", func(value goja.Value) map[string]any {
		return r.v3TextNode(value).toIR()
	})
	setExport(raw, "element", r.v3Element)
	setExport(raw, "component", r.v3Component)
	setExport(raw, "fragment", r.v3Fragment)
	return raw
}

func (r *runtime) v3Element(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) == 0 {
		panic(r.vm.NewGoError(fmt.Errorf("widget DSL element(tag, attrs?, ...children) requires a tag")))
	}
	tag := strings.TrimSpace(call.Arguments[0].String())
	if tag == "" {
		panic(r.vm.NewGoError(fmt.Errorf("widget DSL element tag must not be empty")))
	}
	attrs := map[string]any{}
	childStart := 1
	if len(call.Arguments) > 1 && isPlainObject(call.Arguments[1]) && !looksLikeWidgetNodeExport(call.Arguments[1]) {
		attrs = exportObject(call.Arguments[1])
		childStart = 2
	}
	return r.vm.ToValue(r.v3BuildElement(tag, attrs, call.Arguments[childStart:]))
}

func (r *runtime) v3Component(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) == 0 {
		panic(r.vm.NewGoError(fmt.Errorf("widget DSL component(type, props?, ...children) requires a type")))
	}
	componentType := strings.TrimSpace(call.Arguments[0].String())
	if componentType == "" {
		panic(r.vm.NewGoError(fmt.Errorf("widget DSL component type must not be empty")))
	}
	props, childStart := propsAndChildStart(call.Arguments, 1)
	return r.vm.ToValue(r.v3BuildComponent(componentType, props, call.Arguments[childStart:]))
}

func (r *runtime) v3Fragment(call goja.FunctionCall) goja.Value {
	return r.vm.ToValue(r.v3NodeSpecsToIR(r.v3ExportChildren(call.Arguments)))
}

func (r *runtime) v3RenderableTitle(value goja.Value) any {
	if value == nil || goja.IsUndefined(value) || goja.IsNull(value) {
		return "Section"
	}
	if exported, ok := value.Export().(bool); ok && !exported {
		return "Section"
	}
	if _, ok := value.Export().(string); ok {
		return value.String()
	}
	return r.v3Renderable(value)
}

func (r *runtime) v3Renderable(value goja.Value) any {
	nodes := r.v3ExportChild(value)
	if len(nodes) == 0 {
		return nil
	}
	if len(nodes) == 1 {
		return nodes[0].toIR()
	}
	return r.v3NodeSpecsToIR(nodes)
}

func (r *runtime) v3PageToIR(spec *v3PageSpec) map[string]any {
	children := make([]any, 0, len(spec.Sections)+1)
	if len(spec.Breadcrumbs) > 0 {
		children = append(children, componentNode("Breadcrumbs", map[string]any{"items": spec.Breadcrumbs}))
	}
	for _, section := range spec.Sections {
		children = append(children, r.v3SectionToNode(section))
	}
	rootProps := map[string]any{"gap": "lg"}
	if spec.Density != "" {
		rootProps["density"] = spec.Density
	}
	out := map[string]any{
		"schemaVersion": spec.SchemaVersion,
		"id":            spec.ID,
		"title":         spec.Title,
		"root":          componentNode("Stack", rootProps, children...),
	}
	if len(spec.Meta) > 0 {
		out["meta"] = spec.Meta
	}
	if spec.Shell != nil {
		out["shell"] = spec.Shell
	}
	return out
}

func (r *runtime) v3SectionToNode(spec v3SectionSpec) map[string]any {
	props := map[string]any{"label": spec.Title, "level": 1, "rule": true, "density": "flush"}
	if spec.Caption != "" {
		props["caption"] = spec.Caption
	}
	if spec.AnchorID != "" {
		props["anchorId"] = spec.AnchorID
	}
	if spec.Tone != "" {
		props["tone"] = spec.Tone
	}
	if len(spec.Actions) > 0 {
		props["actions"] = spec.Actions
	}
	return componentNode("SectionBlock", props, r.v3NodeSpecsToIR(spec.Children)...)
}

func (r *runtime) v3ExportChildren(values []goja.Value) []v3NodeSpec {
	out := []v3NodeSpec{}
	for _, value := range values {
		out = append(out, r.v3ExportChild(value)...)
	}
	return out
}

func (r *runtime) v3ExportChild(value goja.Value) []v3NodeSpec {
	if value == nil || goja.IsUndefined(value) || goja.IsNull(value) {
		return nil
	}
	if exported, ok := value.Export().(bool); ok && !exported {
		return nil
	}
	if isArrayLike(value) {
		obj := value.ToObject(r.vm)
		length := int(obj.Get("length").ToInteger())
		out := []v3NodeSpec{}
		for i := 0; i < length; i++ {
			out = append(out, r.v3ExportChild(obj.Get(fmt.Sprintf("%d", i)))...)
		}
		return out
	}
	if isWidgetNode(r.vm, value) {
		return []v3NodeSpec{v3NodeSpecFromIR(value.Export().(map[string]any))}
	}
	return []v3NodeSpec{r.v3TextNode(value)}
}

func (r *runtime) v3TextNode(value goja.Value) v3NodeSpec {
	return v3NodeSpecFromIR(map[string]any{"kind": "text", "text": stringifyValue(value)})
}

func (r *runtime) v3BuildElement(tag string, attrs map[string]any, childValues []goja.Value) map[string]any {
	out := map[string]any{"kind": "element", "tag": tag}
	if len(attrs) > 0 {
		out["attrs"] = attrs
	}
	children := r.v3NodeSpecsToIR(r.v3ExportChildren(childValues))
	if len(children) > 0 {
		out["children"] = children
	}
	return out
}

func (r *runtime) v3BuildComponent(componentType string, props map[string]any, childValues []goja.Value) map[string]any {
	out := map[string]any{"kind": "component", "type": componentType}
	if len(props) > 0 {
		out["props"] = props
	}
	children := r.v3NodeSpecsToIR(r.v3ExportChildren(childValues))
	if len(children) > 0 {
		out["children"] = children
	}
	return out
}

func (r *runtime) v3NodeSpecsToIR(nodes []v3NodeSpec) []any {
	out := make([]any, 0, len(nodes))
	for _, node := range nodes {
		out = append(out, node.toIR())
	}
	return out
}

func v3NodeSpecFromIR(ir map[string]any) v3NodeSpec {
	return v3NodeSpec{Kind: stringFromMap(ir, "kind", ""), IR: ir}
}

func (n v3NodeSpec) toIR() map[string]any {
	out := map[string]any{}
	for k, v := range n.IR {
		out[k] = v
	}
	if n.Source != nil {
		out["source"] = map[string]any{"file": n.Source.File, "line": n.Source.Line, "column": n.Source.Column}
	}
	return out
}

func isV3EmptySlotResult(value goja.Value) bool {
	if value == nil || goja.IsUndefined(value) || goja.IsNull(value) {
		return true
	}
	if exported, ok := value.Export().(bool); ok && !exported {
		return true
	}
	return false
}

func v3PageValidationIssues(spec *v3PageSpec) []map[string]any {
	issues := []map[string]any{}
	if strings.TrimSpace(spec.ID) == "" {
		issues = append(issues, v3ValidationIssue("page_id_required", "page.id", "page id is required"))
	}
	if strings.TrimSpace(spec.Title) == "" {
		issues = append(issues, v3ValidationIssue("page_title_required", "page.title", "page title is required"))
	}
	for sectionIndex, section := range spec.Sections {
		sectionPath := fmt.Sprintf("page.sections[%d]", sectionIndex)
		if section.Title == nil {
			issues = append(issues, v3ValidationIssue("section_title_required", sectionPath+".title", "section title is required"))
		}
		for childIndex, child := range section.Children {
			issues = append(issues, v3NodeValidationIssues(child, fmt.Sprintf("%s.children[%d]", sectionPath, childIndex))...)
		}
	}
	return issues
}

func v3NodeValidationIssues(node v3NodeSpec, path string) []map[string]any {
	issues := []map[string]any{}
	switch node.Kind {
	case "text":
		if _, ok := node.IR["text"]; !ok {
			issues = append(issues, v3ValidationIssue("text_value_required", path+".text", "text node requires a text value"))
		}
	case "element":
		if strings.TrimSpace(stringFromMap(node.IR, "tag", "")) == "" {
			issues = append(issues, v3ValidationIssue("element_tag_required", path+".tag", "element node requires a tag"))
		}
	case "component":
		if strings.TrimSpace(stringFromMap(node.IR, "type", "")) == "" {
			issues = append(issues, v3ValidationIssue("component_type_required", path+".type", "component node requires a type"))
		}
	default:
		issues = append(issues, v3ValidationIssue("node_kind_invalid", path+".kind", "node kind must be text, element, or component"))
	}
	for childIndex, child := range anySlice(node.IR["children"]) {
		childPath := fmt.Sprintf("%s.children[%d]", path, childIndex)
		childNode, ok := widgetNodeFromAny(child)
		if !ok {
			issues = append(issues, v3ValidationIssue("node_child_invalid", childPath, "node child must be a widget node"))
			continue
		}
		issues = append(issues, v3NodeValidationIssues(v3NodeSpecFromIR(childNode), childPath)...)
	}
	return issues
}

func (r *runtime) v3MetadataNode(record map[string]any) map[string]any {
	return componentNode("MetadataGrid", map[string]any{"items": v3MetadataItems(record)})
}

func v3MetadataItems(record map[string]any) []any {
	items := make([]any, 0, len(record))
	for key, value := range record {
		items = append(items, map[string]any{"key": key, "label": key, "value": value})
	}
	return items
}

func v3ValidationIssue(code string, path string, message string) map[string]any {
	return map[string]any{"severity": "error", "code": code, "path": path, "message": message}
}

func v3AccessorSpec(mode string, valueKey string, value string) map[string]any {
	out := map[string]any{"kind": "accessor", "mode": mode}
	if strings.TrimSpace(value) != "" {
		out[valueKey] = value
	}
	return out
}

func slugID(s string) string {
	lower := strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	lastDash := false
	for _, r := range lower {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			lastDash = false
			continue
		}
		if !lastDash {
			b.WriteByte('-')
			lastDash = true
		}
	}
	out := strings.Trim(b.String(), "-")
	if out == "" {
		return "page"
	}
	return out
}
