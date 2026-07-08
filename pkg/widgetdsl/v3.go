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
	Sections      []v3SectionSpec
}

type v3SectionSpec struct {
	Title    any
	Caption  string
	AnchorID string
	Tone     string
	Children []any
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
		section := v3SectionSpec{Title: "Content", Children: r.exportChild(value)}
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
		spec.Children = append(spec.Children, r.text(value))
		return obj
	})
	setExport(obj, "view", func(value goja.Value) *goja.Object {
		spec.Children = append(spec.Children, r.exportChild(value)...)
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

func (r *runtime) v3RenderableTitle(value goja.Value) any {
	if value == nil || goja.IsUndefined(value) || goja.IsNull(value) {
		return "Section"
	}
	if _, ok := value.Export().(string); ok {
		return value.String()
	}
	return r.exportRenderable(value)
}

func (r *runtime) v3PageToIR(spec *v3PageSpec) map[string]any {
	children := make([]any, 0, len(spec.Sections))
	for _, section := range spec.Sections {
		children = append(children, r.v3SectionToNode(section))
	}
	out := map[string]any{
		"schemaVersion": spec.SchemaVersion,
		"id":            spec.ID,
		"title":         spec.Title,
		"root":          componentNode("Stack", map[string]any{"gap": "lg"}, children...),
	}
	if len(spec.Meta) > 0 {
		out["meta"] = spec.Meta
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
	return componentNode("SectionBlock", props, spec.Children...)
}

func v3PageValidationIssues(spec *v3PageSpec) []map[string]any {
	issues := []map[string]any{}
	if strings.TrimSpace(spec.ID) == "" {
		issues = append(issues, map[string]any{"severity": "error", "code": "page_id_required", "path": "page.id", "message": "page id is required"})
	}
	if strings.TrimSpace(spec.Title) == "" {
		issues = append(issues, map[string]any{"severity": "error", "code": "page_title_required", "path": "page.title", "message": "page title is required"})
	}
	return issues
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
