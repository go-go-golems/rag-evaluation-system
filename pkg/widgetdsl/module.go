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

const (
	UIModuleName            = "ui.dsl"
	DataModuleName          = "data.dsl"
	DataV2ModuleName        = "data.v2.dsl"
	WidgetV3ModuleName      = "widget.dsl"
	ContextWindowModuleName = "context_window.dsl"
	CourseModuleName        = "course.dsl"
	CmsModuleName           = "cms.dsl"
)

type moduleSpec struct {
	name    string
	doc     string
	helpers map[string]string
	page    bool
	cell    bool
	action  bool
	recipes []string
}

var templatePattern = regexp.MustCompile(`\$\{([^}]+)\}|\$([A-Za-z0-9_.-]+)`)

var uiHelpers = map[string]string{
	"appShell": "AppShell",
	"appNav":   "AppNav",
	"button":   "Button",
	// Generic primitives promoted from domain modules (RAGEVAL-UI-GRAMMAR
	// design-doc 02 §3); the domain modules keep their original names as
	// deprecated aliases.
	"breadcrumbs":     "Breadcrumbs",
	"emptyState":      "EmptyState",
	"fieldGrid":       "FieldGrid",
	"markdownArticle": "MarkdownArticle",
	"meterBar":        "MeterBar",
	"pagination":      "Pagination",
	"richArticle":     "RichArticle",
	"searchField":     "SearchField",
	"tag":             "Tag",
	"tileGrid":        "TileGrid",
	"uploadDropArea":  "ContextUploadDropArea",
	"caption":         "Caption",
	"checkList":       "CheckList",
	"codeText":        "CodeText",
	"dashboardGrid":   "DashboardGrid",
	"divider":         "Divider",
	"figureBlock":     "FigureBlock",
	"formPanel":       "FormPanel",
	"formRow":         "FormRow",
	"inline":          "Inline",
	"keyPointList":    "KeyPointList",
	"keyValueStrip":   "KeyValueStrip",
	"metadataGrid":    "MetadataGrid",
	"panel":           "Panel",
	"personSummary":   "PersonSummary",
	"scrollRegion":    "ScrollRegion",
	"sectionBlock":    "SectionBlock",
	"selectInput":     "SelectInput",
	"sidebarNav":      "SidebarNav",
	"sidebarShell":    "SidebarShell",
	"splitPane":       "SplitPane",
	"stack":           "Stack",
	"statusText":      "StatusText",
	"stepList":        "StepList",
	"tabList":         "TabList",
	"textBlock":       "Text",
	"textInput":       "TextInput",
	"textareaInput":   "TextareaInput",
}

var dataHelpers = map[string]string{
	"dataTable": "DataTable",
}

var contextWindowHelpers = map[string]string{
	"anchoredCommentCard":        "AnchoredCommentCard",
	"anchoredCommentRail":        "AnchoredCommentRail",
	"annotationBadge":            "AnnotationBadge",
	"annotationNoteCard":         "AnnotationNoteCard",
	"annotationRailPanel":        "AnnotationRailPanel",
	"contextBudgetBar":           "ContextBudgetBar",
	"contextDiagramPanel":        "ContextDiagramPanel",
	"contextGroupedStripDiagram": "ContextGroupedStripDiagram",
	"contextStyleSwatch":         "ContextStyleSwatch",
	"contextLegend":              "ContextLegend",
	"contextStackDiagram":        "ContextStackDiagram",
	"contextStripDiagram":        "ContextStripDiagram",
	"contextTreemap":             "ContextTreemap",
	"contextTurnPagerPanel":      "ContextTurnPagerPanel",
	"contextUploadDropArea":      "ContextUploadDropArea",
	"transcriptMessageCard":      "TranscriptMessageCard",
	"transcriptReaderPanel":      "TranscriptReaderPanel",
	"transcriptRoleBadge":        "TranscriptRoleBadge",
	"transcriptSessionHeader":    "TranscriptSessionHeader",
	"transcriptWorkspacePanel":   "TranscriptWorkspacePanel",
}

var cmsHelpers = map[string]string{
	"articleListPanel":   "ArticleListPanel",
	"assetTile":          "AssetTile",
	"breadcrumbs":        "Breadcrumbs",
	"cmsShell":           "CmsShell",
	"contentStatusBadge": "ContentStatusBadge",
	"emptyState":         "EmptyState",
	"markdownEditor":     "MarkdownEditor",
	"mediaLibraryPanel":  "MediaLibraryPanel",
	"mediaThumb":         "MediaThumb",
	"meterBar":           "MeterBar",
	"pagination":         "Pagination",
	"searchField":        "SearchField",
	"tag":                "Tag",
	"tileGrid":           "TileGrid",
}

var courseHelpers = map[string]string{
	"contextStudioNavIcon":   "ContextStudioNavIcon",
	"courseLessonPanel":      "CourseLessonPanel",
	"courseSlidePanel":       "CourseSlidePanel",
	"courseStepNav":          "CourseStepNav",
	"courseStudioShell":      "CourseStudioShell",
	"documentListPanel":      "DocumentListPanel",
	"documentPreviewToolbar": "DocumentPreviewToolbar",
	"handoutDocumentShell":   "HandoutDocumentShell",
	"markdownArticle":        "MarkdownArticle",
	"richArticle":            "RichArticle",
	"slideShell":             "SlideShell",
}

var moduleSpecs = []moduleSpec{
	{
		name:    UIModuleName,
		helpers: uiHelpers,
		page:    true,
		action:  true,
		recipes: []string{"metrics", "actionToolbar"},
		doc:     "ui.dsl provides generic Widget IR page, layout, primitive, and foundation helpers.",
	},
	{
		name:    DataModuleName,
		helpers: dataHelpers,
		cell:    true,
		action:  true,
		recipes: []string{"masterDetailTable"},
		doc:     "data.dsl provides data-display helpers, table cell helpers, and data recipes.",
	},
	{
		name: DataV2ModuleName,
		doc:  "data.v2.dsl provides the hard-cutover typed/fluent collection builder experiment.",
	},
	{
		name: WidgetV3ModuleName,
		doc:  "widget.dsl provides the parallel clean Widget DSL v3 root namespace.",
	},
	{
		name:    ContextWindowModuleName,
		helpers: contextWindowHelpers,
		action:  true,
		recipes: []string{"contextDiagram", "annotatedTranscript"},
		doc:     "context_window.dsl provides context-window, transcript, annotation, and anchored-comment helpers.",
	},
	{
		name:    CourseModuleName,
		helpers: courseHelpers,
		action:  true,
		recipes: []string{"courseStudio", "courseSlide", "handout"},
		doc:     "course.dsl provides course, lesson, slide, handout, and course-studio helpers.",
	},
	{
		name:    CmsModuleName,
		helpers: cmsHelpers,
		action:  true,
		recipes: []string{"mediaLibrary", "articleList"},
		doc:     "cms.dsl provides media, asset, and article-management helpers.",
	},
}

var moduleSpecsByName = func() map[string]moduleSpec {
	out := map[string]moduleSpec{}
	for _, spec := range moduleSpecs {
		out[spec.name] = spec
	}
	return out
}()

type module struct {
	spec moduleSpec
}

var _ modules.NativeModule = (*module)(nil)

func NewLoader(moduleName string) require.ModuleLoader {
	spec, ok := moduleSpecsByName[moduleName]
	if !ok {
		panic(fmt.Sprintf("unknown widget DSL module %q", moduleName))
	}
	mod := &module{spec: spec}
	return mod.Loader
}

func Register(reg *require.Registry) {
	if reg == nil {
		return
	}
	for _, spec := range moduleSpecs {
		reg.RegisterNativeModule(spec.name, NewLoader(spec.name))
	}
}

func (m *module) Name() string { return m.spec.name }

func (m *module) Doc() string { return m.spec.doc }

func (m *module) Loader(vm *goja.Runtime, moduleObj *goja.Object) {
	rt := &runtime{vm: vm}
	exports := moduleObj.Get("exports").(*goja.Object)
	rt.install(exports, m.spec)
}

func init() {
	for _, spec := range moduleSpecs {
		modules.Register(&module{spec: spec})
	}
}

type runtime struct {
	vm *goja.Runtime
}

func (r *runtime) install(exports *goja.Object, spec moduleSpec) {
	if spec.name == WidgetV3ModuleName {
		r.installWidgetV3(exports)
		return
	}
	setExport(exports, "text", r.text)
	setExport(exports, "element", r.element)
	setExport(exports, "component", r.component)
	setExport(exports, "fragment", r.fragment)
	if spec.page {
		setExport(exports, "page", r.page)
	}

	for name, componentType := range spec.helpers {
		setExport(exports, name, r.componentFactory(componentType))
	}

	if spec.cell {
		setExport(exports, "cell", r.cellObject())
	}
	if spec.action {
		setExport(exports, "action", r.actionObject())
	}
	if spec.name == ContextWindowModuleName {
		r.installContextWindowStyleHelpers(exports)
	}
	if spec.name == DataModuleName {
		r.installDataGrammar(exports)
	}
	if spec.name == DataV2ModuleName {
		r.installDataV2(exports)
	}
	if spec.name == UIModuleName {
		setExport(exports, "section", r.sectionVerb)
	}
	if len(spec.recipes) > 0 {
		setExport(exports, "recipes", r.recipesObject(spec.recipes))
	}
}

func (r *runtime) installWidgetV3(exports *goja.Object) {
	setExport(exports, "page", r.v3Page)
	setExport(exports, "raw", r.v3RawObject())
	setExport(exports, "act", r.actionObject())
	setExport(exports, "bind", r.bindingObject())
	setExport(exports, "ui", r.v3UIObject())
	setExport(exports, "data", r.v3DataObject())
	setExport(exports, "cms", r.v3CMSObject())
	setExport(exports, "course", r.v3CourseObject())
	setExport(exports, "context", r.v3ContextObject())
	setExport(exports, "schedule", r.v3ScheduleObject())
	setExport(exports, "time", r.v3TimeObject())
	setExport(exports, "style", r.vm.NewObject())
}

func (r *runtime) bindingObject() *goja.Object {
	bind := r.vm.NewObject()
	setExport(bind, "field", func(path string) map[string]any {
		return v3AccessorSpec("field", "field", path)
	})
	setExport(bind, "path", func(path string) map[string]any {
		return v3AccessorSpec("path", "path", path)
	})
	setExport(bind, "map", func(field string) map[string]any {
		return v3AccessorSpec("map", "mapField", field)
	})
	setExport(bind, "template", func(template string) map[string]any {
		return v3AccessorSpec("template", "template", template)
	})
	setExport(bind, "context", func(path string) map[string]any {
		return v3AccessorSpec("context", "path", path)
	})
	setExport(bind, "const", func(value goja.Value) map[string]any {
		return map[string]any{"kind": "const", "value": value.Export()}
	})
	return bind
}

func (r *runtime) cellObject() *goja.Object {
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
	setExport(cell, "linkButton", func(hrefField string, labelField string, options ...goja.Value) map[string]any {
		out := map[string]any{"kind": "linkButton", "hrefField": hrefField, "labelField": labelField}
		mergeOptions(out, exportOptions(options))
		return out
	})
	setExport(cell, "actionButton", func(label goja.Value, action goja.Value, options ...goja.Value) map[string]any {
		out := map[string]any{"kind": "actionButton", "label": r.exportRenderable(label), "action": exportObject(action)}
		mergeOptions(out, exportOptions(options))
		return out
	})
	setExport(cell, "constant", func(value goja.Value) any {
		return map[string]any{"kind": "constant", "value": r.exportRenderable(value)}
	})
	return cell
}

func (r *runtime) actionObject() *goja.Object {
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
	setExport(action, "download", func(to string, options ...goja.Value) map[string]any {
		out := map[string]any{"kind": "download", "to": to}
		mergeOptions(out, exportOptions(options))
		return out
	})
	setExport(action, "event", func(event string, options ...goja.Value) map[string]any {
		out := map[string]any{"kind": "event", "event": event}
		mergeOptions(out, exportOptions(options))
		return out
	})
	setExport(action, "copy", func(value string, options ...goja.Value) map[string]any {
		out := map[string]any{"kind": "copy", "value": value}
		mergeOptions(out, exportOptions(options))
		return out
	})
	return action
}

func (r *runtime) installContextWindowStyleHelpers(exports *goja.Object) {
	setExport(exports, "visualStyle", func(options goja.Value) map[string]any {
		out := exportObject(options)
		if _, ok := out["pattern"]; !ok {
			out["pattern"] = "none"
		}
		if _, ok := out["fill"]; !ok {
			out["fill"] = "var(--mac-surface)"
		}
		return out
	})
	setExport(exports, "legendItem", func(id string, label string, options ...goja.Value) map[string]any {
		out := map[string]any{"id": id, "label": label}
		mergeOptions(out, exportOptions(options))
		return out
	})
	setExport(exports, "styleSet", func(options goja.Value) map[string]any {
		out := exportObject(options)
		if _, ok := out["legend"]; !ok {
			out["legend"] = []any{}
		}
		if _, ok := out["styles"]; !ok {
			out["styles"] = map[string]any{}
		}
		return out
	})
	setExport(exports, "contextPart", func(id string, label string, styleKey string, tokens int, options ...goja.Value) map[string]any {
		out := map[string]any{"id": id, "label": label, "styleKey": styleKey, "tokens": tokens}
		mergeOptions(out, exportOptions(options))
		return out
	})
	setExport(exports, "contextSnapshot", func(options goja.Value) map[string]any {
		out := exportObject(options)
		if _, ok := out["id"]; !ok {
			out["id"] = "ctx"
		}
		if _, ok := out["title"]; !ok {
			out["title"] = "Context"
		}
		if _, ok := out["limit"]; !ok {
			out["limit"] = 0
		}
		if _, ok := out["parts"]; !ok {
			out["parts"] = []any{}
		}
		return out
	})
	setExport(exports, "paletteStyleSet", func(options goja.Value) map[string]any {
		return buildPaletteStyleSet(exportObject(options))
	})
}

func buildPaletteStyleSet(options map[string]any) map[string]any {
	paletteName := stringFromMap(options, "palette", "Dusty Magenta / Blue")
	colors := paletteColors(paletteName)
	entries := anySlice(options["entries"])
	legend := []any{}
	styles := map[string]any{}
	for _, raw := range entries {
		entry, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		id := stringFromMap(entry, "id", "")
		if id == "" {
			continue
		}
		label := stringFromMap(entry, "label", id)
		accent := colorByAccent(colors, stringFromMap(entry, "accent", "a"))
		pattern := stringFromMap(entry, "pattern", "diagonal")
		solid := pattern == "solid" || boolFromMap(entry, "solid", false)
		fillPct := numberFromMap(entry, "fillPct", 18)
		linePct := numberFromMap(entry, "linePct", 60)
		fill := colorMix(accent, fillPct, colors["paper"])
		line := colorMix(accent, linePct, colors["paper"])
		labelColor := colors["ink"]
		if solid {
			fill = accent
			line = "transparent"
			labelColor = "#ffffff"
		}
		styles[id] = map[string]any{"pattern": pattern, "fill": fill, "line": line, "stroke": colors["ink"], "labelColor": labelColor}
		legendItem := map[string]any{"id": id, "label": label}
		if desc, ok := entry["description"]; ok {
			legendItem["description"] = desc
		}
		if hidden, ok := entry["hidden"]; ok {
			legendItem["hidden"] = hidden
		}
		legend = append(legend, legendItem)
	}
	out := map[string]any{"legend": legend, "styles": styles}
	copyIfPresent(out, options, "id")
	copyIfPresent(out, options, "name")
	copyIfPresent(out, options, "legendSize")
	copyIfPresent(out, options, "swatchSize")
	return out
}

func paletteColors(name string) map[string]string {
	switch name {
	case "Signal Orange / Cyan":
		return map[string]string{"paper": "#F5F1E8", "ink": "#141414", "grid": "#BDB7AA", "shadow": "#7E7A72", "accent_a": "#D86F2A", "accent_b": "#2EA6A6", "accent_c": "#F0B45A"}
	case "Slate / Coral":
		return map[string]string{"paper": "#F2F3EF", "ink": "#111314", "grid": "#B8BDBB", "shadow": "#707878", "accent_a": "#5F7F89", "accent_b": "#C46A55", "accent_c": "#D8E0DF"}
	case "Cobalt / Sand":
		return map[string]string{"paper": "#F3EEDC", "ink": "#111318", "grid": "#BFB59A", "shadow": "#766F5E", "accent_a": "#315D91", "accent_b": "#D2A84A", "accent_c": "#A8B8CC"}
	default:
		return map[string]string{"paper": "#F2EEF2", "ink": "#141214", "grid": "#C0B5C1", "shadow": "#776C7A", "accent_a": "#9C527E", "accent_b": "#4F74A8", "accent_c": "#D6B8CB"}
	}
}

func colorByAccent(colors map[string]string, accent string) string {
	switch accent {
	case "b":
		return colors["accent_b"]
	case "c":
		return colors["accent_c"]
	case "grid":
		return colors["grid"]
	case "shadow":
		return colors["shadow"]
	case "ink":
		return colors["ink"]
	default:
		return colors["accent_a"]
	}
}

func colorMix(color string, pct float64, paper string) string {
	return fmt.Sprintf("color-mix(in srgb, %s %.0f%%, %s)", color, pct, paper)
}

func numberFromMap(m map[string]any, key string, fallback float64) float64 {
	if value, ok := m[key]; ok {
		switch typed := value.(type) {
		case int:
			return float64(typed)
		case int64:
			return float64(typed)
		case float64:
			return typed
		case float32:
			return float64(typed)
		}
	}
	return fallback
}

func (r *runtime) recipesObject(names []string) *goja.Object {
	recipes := r.vm.NewObject()
	for _, name := range names {
		switch name {
		case "metrics":
			setExport(recipes, name, r.metricsRecipe)
		case "actionToolbar":
			setExport(recipes, name, r.actionToolbarRecipe)
		case "masterDetailTable":
			setExport(recipes, name, r.masterDetailTableRecipe)
		case "contextDiagram":
			setExport(recipes, name, r.contextDiagramRecipe)
		case "annotatedTranscript":
			setExport(recipes, name, r.annotatedTranscriptRecipe)
		case "courseStudio":
			setExport(recipes, name, r.courseStudioRecipe)
		case "courseSlide":
			setExport(recipes, name, r.courseSlideRecipe)
		case "handout":
			setExport(recipes, name, r.handoutRecipe)
		case "mediaLibrary":
			setExport(recipes, name, r.mediaLibraryRecipe)
		case "articleList":
			setExport(recipes, name, r.articleListRecipe)
		}
	}
	return recipes
}

func (r *runtime) text(value goja.Value) map[string]any {
	return map[string]any{"kind": "text", "text": stringifyValue(value)}
}

func (r *runtime) element(call goja.FunctionCall) goja.Value {
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
		panic(r.vm.NewGoError(fmt.Errorf("widget DSL component(type, props?, ...children) requires a type")))
	}
	componentType := strings.TrimSpace(call.Arguments[0].String())
	if componentType == "" {
		panic(r.vm.NewGoError(fmt.Errorf("widget DSL component type must not be empty")))
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
		panic(r.vm.NewGoError(fmt.Errorf("ui.dsl page(options) requires an options object")))
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

func (r *runtime) contextDiagramRecipe(call goja.FunctionCall) goja.Value {
	options := firstObject(call.Arguments)
	styleSet := options["styleSet"]
	if styleSet == nil {
		if options["palette"] != nil && options["entries"] != nil {
			styleSet = buildPaletteStyleSet(options)
		} else {
			panic(r.vm.NewGoError(fmt.Errorf("context_window.dsl recipes.contextDiagram requires styleSet or palette+entries")))
		}
	}
	props := map[string]any{
		"snapshot": valueOrDefault(options["snapshot"], map[string]any{"id": "empty", "title": "Context", "limit": 0, "parts": []any{}}),
		"styleSet": styleSet,
	}
	if value := options["view"]; value != nil {
		props["initialView"] = value
	} else if value := options["initialView"]; value != nil {
		props["initialView"] = value
	}
	copyIfPresent(props, options, "selectedPartId")
	return r.vm.ToValue(componentNode("ContextDiagramPanel", props))
}

func (r *runtime) annotatedTranscriptRecipe(call goja.FunctionCall) goja.Value {
	options := firstObject(call.Arguments)
	transcript, _ := options["transcript"].(map[string]any)
	props := map[string]any{
		"title":       valueOrDefault(valueFromMap(transcript, "title"), stringFromMap(options, "title", "Transcript")),
		"subtitle":    valueOrDefault(valueFromMap(transcript, "subtitle"), options["subtitle"]),
		"messages":    anySlice(valueOrDefault(valueFromMap(transcript, "messages"), options["messages"])),
		"annotations": anySlice(valueOrDefault(valueFromMap(transcript, "annotations"), options["annotations"])),
		"showNotes":   boolFromMap(options, "showNotes", true),
	}
	if selected := valueOrDefault(options["selectedAnnotationId"], valueFromMap(transcript, "selectedAnnotationId")); selected != nil {
		props["selectedAnnotationId"] = selected
	}
	if act, ok := normalizeActionSpec(options["onAnnotationSelect"], nil, nil); ok {
		props["onAnnotationSelectAction"] = act
	}
	return r.vm.ToValue(componentNode("TranscriptWorkspacePanel", props))
}

func (r *runtime) courseStudioRecipe(call goja.FunctionCall) goja.Value {
	options := firstObject(call.Arguments)
	props := map[string]any{
		"sections": valueOrDefault(options["sections"], []any{}),
		"title":    stringFromMap(options, "title", "Context Window Engineering"),
	}
	copyIfPresent(props, options, "subtitle")
	copyIfPresent(props, options, "activeItemId")
	copyIfPresent(props, options, "contentPadding")
	if act, ok := normalizeActionSpec(options["onNavigate"], nil, nil); ok {
		props["onNavigateAction"] = act
	}
	children := []any{}
	if main, ok := widgetNodeFromAny(options["main"]); ok {
		children = append(children, main)
	}
	return r.vm.ToValue(componentNode("CourseStudioShell", props, children...))
}

func (r *runtime) courseSlideRecipe(call goja.FunctionCall) goja.Value {
	options := firstObject(call.Arguments)
	props := map[string]any{
		"slide":    valueOrDefault(options["slide"], map[string]any{}),
		"snapshot": valueOrDefault(options["snapshot"], map[string]any{"id": "empty", "title": "Context", "limit": 0, "parts": []any{}}),
	}
	copyIfPresent(props, options, "index")
	copyIfPresent(props, options, "total")
	copyIfPresent(props, options, "visualSide")
	copyIfPresent(props, options, "mode")
	if act, ok := normalizeActionSpec(options["onPrevious"], nil, nil); ok {
		props["onPreviousAction"] = act
	}
	if act, ok := normalizeActionSpec(options["onNext"], nil, nil); ok {
		props["onNextAction"] = act
	}
	if act, ok := normalizeActionSpec(options["onPresent"], nil, nil); ok {
		props["onPresentAction"] = act
	}
	if act, ok := normalizeActionSpec(options["onFullscreen"], nil, nil); ok {
		props["onFullscreenAction"] = act
	}
	return r.vm.ToValue(componentNode("CourseSlidePanel", props))
}

func (r *runtime) handoutRecipe(call goja.FunctionCall) goja.Value {
	options := firstObject(call.Arguments)
	bundle, _ := options["bundle"].(map[string]any)
	props := map[string]any{
		"intro":     valueOrDefault(valueFromMap(bundle, "intro"), stringFromMap(options, "intro", "Handout")),
		"documents": anySlice(valueOrDefault(valueFromMap(bundle, "docs"), options["documents"])),
	}
	copyIfPresent(props, options, "selectedDocumentId")
	copyIfPresent(props, options, "title")
	if act, ok := normalizeActionSpec(options["onSelect"], nil, nil); ok {
		props["onDocumentSelectAction"] = act
	}
	if act, ok := normalizeActionSpec(options["onDownload"], nil, nil); ok {
		props["onDownloadAction"] = act
	}
	if act, ok := normalizeActionSpec(options["onDownloadAll"], nil, nil); ok {
		props["onDownloadAllAction"] = act
	}
	return r.vm.ToValue(componentNode("HandoutDocumentShell", props))
}

func (r *runtime) mediaLibraryRecipe(call goja.FunctionCall) goja.Value {
	options := firstObject(call.Arguments)
	props := map[string]any{
		"assets": anySlice(options["assets"]),
	}
	copyIfPresent(props, options, "selectedAssetIds")
	copyIfPresent(props, options, "selectionMode")
	copyIfPresent(props, options, "query")
	copyIfPresent(props, options, "kindFilter")
	copyIfPresent(props, options, "page")
	copyIfPresent(props, options, "pageCount")
	copyIfPresent(props, options, "uploads")
	copyIfPresent(props, options, "showStatusBadges")
	copyIfPresent(props, options, "emptyMessage")
	copyIfPresent(props, options, "title")
	copyIfPresent(props, options, "minTileWidth")
	actionProps := map[string]string{
		"onAssetSelect":      "onAssetSelectAction",
		"onAssetOpen":        "onAssetOpenAction",
		"onQuerySubmit":      "onQuerySubmitAction",
		"onKindFilterChange": "onKindFilterChangeAction",
		"onPageChange":       "onPageChangeAction",
		"onFilesSelected":    "onFilesSelectedAction",
	}
	for option, prop := range actionProps {
		if act, ok := normalizeActionSpec(options[option], nil, nil); ok {
			props[prop] = act
		}
	}
	return r.vm.ToValue(componentNode("MediaLibraryPanel", props))
}

func (r *runtime) articleListRecipe(call goja.FunctionCall) goja.Value {
	options := firstObject(call.Arguments)
	props := map[string]any{
		"articles": anySlice(options["articles"]),
	}
	copyIfPresent(props, options, "selectedArticleId")
	copyIfPresent(props, options, "statusFilter")
	copyIfPresent(props, options, "query")
	copyIfPresent(props, options, "page")
	copyIfPresent(props, options, "pageCount")
	copyIfPresent(props, options, "emptyMessage")
	copyIfPresent(props, options, "title")
	copyIfPresent(props, options, "maxVisibleTags")
	actionProps := map[string]string{
		"onArticleSelect":      "onArticleSelectAction",
		"onCreate":             "onCreateAction",
		"onRowAction":          "onRowActionAction",
		"onStatusFilterChange": "onStatusFilterChangeAction",
		"onQuerySubmit":        "onQuerySubmitAction",
		"onPageChange":         "onPageChangeAction",
	}
	for option, prop := range actionProps {
		if act, ok := normalizeActionSpec(options[option], nil, nil); ok {
			props[prop] = act
		}
	}
	return r.vm.ToValue(componentNode("ArticleListPanel", props))
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

func valueFromMap(m map[string]any, key string) any {
	if m == nil {
		return nil
	}
	return m[key]
}

func copyIfPresent(dst map[string]any, src map[string]any, key string) {
	if value, ok := src[key]; ok && value != nil {
		dst[key] = value
	}
}

func componentNode(componentType string, props map[string]any, children ...any) map[string]any {
	out := map[string]any{"kind": "component", "type": componentType}
	if len(props) > 0 {
		out["props"] = props
	}
	if len(children) > 0 {
		out["children"] = children
	}
	return out
}

func widgetNodeFromAny(value any) (map[string]any, bool) {
	node, ok := value.(map[string]any)
	return node, ok && isWidgetNodeExport(node)
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
