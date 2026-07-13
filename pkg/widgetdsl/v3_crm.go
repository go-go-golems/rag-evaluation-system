package widgetdsl

import (
	"fmt"
	"strings"

	"github.com/dop251/goja"
)

const v3CRMRefProperty = "__widgetDSLCRMRef"

type v3CRMRef struct {
	kind     string
	fields   *v3CRMFieldsSpec
	pipeline *v3CRMPipelineSpec
}

type v3CRMFieldsSpec struct {
	Name   string
	Fields []map[string]any
}

type v3CRMPipelineSpec struct {
	ID     string
	Name   string
	Stages []map[string]any
}

// v3CRMObject exposes CRM authoring helpers. These helpers deliberately emit
// the existing generic CRM Widget IR components rather than introducing a
// second renderer contract.
func (r *runtime) v3CRMObject() *goja.Object {
	crm := r.vm.NewObject()
	setExport(crm, "fields", r.v3CRMFields)
	setExport(crm, "pipeline", r.v3CRMPipeline)
	setExport(crm, "pipelineBoard", r.v3CRMPipelineBoard)
	setExport(crm, "recordFields", r.v3CRMRecordFields)
	setExport(crm, "field", r.v3CRMField)
	setExport(crm, "tasksInbox", r.v3CRMTasksInbox)
	setExport(crm, "stat", r.v3CRMStat)
	setExport(crm, "funnel", r.v3CRMFunnel)
	setExport(crm, "intent", r.v3CRMIntentObject())
	return crm
}

func (r *runtime) v3CRMFields(args ...goja.Value) *goja.Object {
	name := "Fields"
	var cb goja.Value
	if len(args) > 0 {
		if _, ok := goja.AssertFunction(args[0]); ok {
			cb = args[0]
		} else {
			name = args[0].String()
		}
	}
	if len(args) > 1 {
		cb = args[1]
	}
	spec := &v3CRMFieldsSpec{Name: name}
	builder := r.v3CRMFieldsBuilder(spec)
	if cb != nil && !goja.IsUndefined(cb) && !goja.IsNull(cb) {
		r.applyV3BuilderCallback(builder, cb, "crm.fields")
	}
	return builder
}

func (r *runtime) v3CRMFieldsBuilder(spec *v3CRMFieldsSpec) *goja.Object {
	obj := r.newV3Builder("crm.fields")
	r.attachCRMRef(obj, &v3CRMRef{kind: "fields", fields: spec})
	add := func(kind, key string, options ...goja.Value) *goja.Object {
		if strings.TrimSpace(key) == "" {
			panic(r.vm.NewGoError(fmt.Errorf("widget.crm.fields key must not be empty")))
		}
		for _, field := range spec.Fields {
			if field["key"] == key {
				panic(r.vm.NewGoError(fmt.Errorf("widget.crm.fields duplicate field key %q", key)))
			}
		}
		field := exportOptions(options)
		field["key"] = key
		field["type"] = kind
		if _, ok := field["label"]; !ok {
			field["label"] = key
		}
		spec.Fields = append(spec.Fields, field)
		return obj
	}
	for _, kind := range []string{"text", "longtext", "email", "phone", "url", "number", "currency", "percent", "date", "datetime", "boolean", "select", "multiselect", "tags", "user"} {
		fieldKind := kind
		setExport(obj, fieldKind, func(key string, options ...goja.Value) *goja.Object {
			return add(fieldKind, key, options...)
		})
	}
	setExport(obj, "relation", func(key, relatedObject string, options ...goja.Value) *goja.Object {
		opts := exportOptions(options)
		opts["relatedObject"] = relatedObject
		return add("relation", key, r.vm.ToValue(opts))
	})
	setExport(obj, "build", func() map[string]any { return v3CRMFieldsSnapshot(spec) })
	setExport(obj, "validate", func() []map[string]any { return v3CRMFieldsIssues(spec) })
	return obj
}

func (r *runtime) v3CRMPipeline(nameOrOptions goja.Value, cb ...goja.Value) *goja.Object {
	name, id := "Pipeline", "pipeline"
	if isPlainObject(nameOrOptions) {
		opts := exportObject(nameOrOptions)
		name = stringFromMap(opts, "name", name)
		id = stringFromMap(opts, "id", id)
	} else if strings.TrimSpace(nameOrOptions.String()) != "" {
		name = nameOrOptions.String()
		id = slugID(name)
	}
	spec := &v3CRMPipelineSpec{ID: id, Name: name}
	builder := r.v3CRMPipelineBuilder(spec)
	if len(cb) > 0 && !goja.IsUndefined(cb[0]) && !goja.IsNull(cb[0]) {
		r.applyV3BuilderCallback(builder, cb[0], "crm.pipeline")
	}
	return builder
}

func (r *runtime) v3CRMPipelineBuilder(spec *v3CRMPipelineSpec) *goja.Object {
	obj := r.newV3Builder("crm.pipeline")
	r.attachCRMRef(obj, &v3CRMRef{kind: "pipeline", pipeline: spec})
	setExport(obj, "stage", func(id, label string, options ...goja.Value) *goja.Object {
		if strings.TrimSpace(id) == "" || strings.TrimSpace(label) == "" {
			panic(r.vm.NewGoError(fmt.Errorf("widget.crm.pipeline stage id and label must not be empty")))
		}
		for _, stage := range spec.Stages {
			if stage["id"] == id {
				panic(r.vm.NewGoError(fmt.Errorf("widget.crm.pipeline duplicate stage id %q", id)))
			}
		}
		stage := exportOptions(options)
		stage["id"] = id
		stage["name"] = label
		stage["order"] = len(spec.Stages)
		spec.Stages = append(spec.Stages, stage)
		return obj
	})
	setExport(obj, "build", func() map[string]any { return v3CRMPipelineSnapshot(spec) })
	setExport(obj, "validate", func() []map[string]any { return v3CRMPipelineIssues(spec) })
	return obj
}

func (r *runtime) v3CRMPipelineBoard(pipelineValue, dealsValue goja.Value, cb ...goja.Value) map[string]any {
	pipeline := r.crmPipelineFromValue(pipelineValue)
	props := v3CRMBoardProps(pipeline, anySlice(dealsValue.Export()), nil)
	builder := r.v3CRMPipelineBoardBuilder(props, pipeline)
	if len(cb) > 0 && !goja.IsUndefined(cb[0]) && !goja.IsNull(cb[0]) {
		r.applyV3BuilderCallback(builder, cb[0], "crm.pipelineBoard")
	}
	return componentNode("BoardEngine", props)
}

func (r *runtime) v3CRMPipelineBoardBuilder(props map[string]any, pipeline map[string]any) *goja.Object {
	obj := r.newV3Builder("crm.pipelineBoard")
	setExport(obj, "summaries", func(value goja.Value) *goja.Object {
		mergeOptions(props, v3CRMBoardProps(pipeline, anySlice(props["cards"]), anySlice(value.Export())))
		return obj
	})
	setExport(obj, "selected", func(id string) *goja.Object { props["selectedCardId"] = id; return obj })
	setExport(obj, "ariaLabel", func(label string) *goja.Object { props["ariaLabel"] = label; return obj })
	setExport(obj, "onMove", func(action goja.Value) *goja.Object { props["onMoveAction"] = action.Export(); return obj })
	setExport(obj, "onOpen", func(action goja.Value) *goja.Object { props["onCardSelectAction"] = action.Export(); return obj })
	return obj
}

func (r *runtime) v3CRMField(value goja.Value, fieldSpec goja.Value, options ...goja.Value) map[string]any {
	props := exportOptions(options)
	props["value"] = value.Export()
	props["spec"] = fieldSpec.Export()
	return componentNode("FieldRenderer", props)
}

func (r *runtime) v3CRMRecordFields(valuesValue, fieldsValue goja.Value, cb ...goja.Value) map[string]any {
	fields := r.crmFieldsFromValue(fieldsValue)
	props := map[string]any{
		"values":   exportObject(valuesValue),
		"sections": v3CRMFieldSections(fields["fields"]),
		"mode":     "read",
	}
	builder := r.v3CRMRecordFieldsBuilder(props)
	if len(cb) > 0 && !goja.IsUndefined(cb[0]) && !goja.IsNull(cb[0]) {
		r.applyV3BuilderCallback(builder, cb[0], "crm.recordFields")
	}
	return componentNode("RecordFieldList", props)
}

func (r *runtime) v3CRMRecordFieldsBuilder(props map[string]any) *goja.Object {
	obj := r.newV3Builder("crm.recordFields")
	setExport(obj, "mode", func(mode string) *goja.Object { props["mode"] = mode; return obj })
	setExport(obj, "refs", func(refs goja.Value) *goja.Object { props["refs"] = refs.Export(); return obj })
	setExport(obj, "onChange", func(action goja.Value) *goja.Object { props["onFieldChangeAction"] = action.Export(); return obj })
	return obj
}

func (r *runtime) v3CRMTasksInbox(tasks goja.Value, cb ...goja.Value) map[string]any {
	rows := make([]any, 0, len(anySlice(tasks.Export())))
	for _, task := range anySlice(tasks.Export()) {
		t := exportMap(task)
		rows = append(rows, componentNode("Inline", map[string]any{"gap": "sm", "justify": "between"},
			componentNode("Text", map[string]any{"size": "compact"}, map[string]any{"kind": "text", "text": fmt.Sprint(t["title"])}),
			componentNode("Caption", map[string]any{}, map[string]any{"kind": "text", "text": fmt.Sprintf("%v · %v", t["priority"], t["dueISO"])}),
		))
	}
	props := map[string]any{"title": "Tasks", "density": "condensed"}
	if len(cb) > 0 && !goja.IsUndefined(cb[0]) && !goja.IsNull(cb[0]) {
		r.applyV3BuilderCallback(r.v3ActionsBuilder(&[]any{}), cb[0], "crm.tasksInbox")
	}
	return componentNode("Panel", props, componentNode("Stack", map[string]any{"gap": "xs"}, rows...))
}

func (r *runtime) v3CRMStat(label, value goja.Value, options ...goja.Value) map[string]any {
	props := exportOptions(options)
	props["label"] = r.v3Renderable(label)
	props["value"] = r.v3Renderable(value)
	return componentNode("StatTile", props)
}

func (r *runtime) v3CRMFunnel(pipelineValue, summariesValue goja.Value, options ...goja.Value) map[string]any {
	pipeline := r.crmPipelineFromValue(pipelineValue)
	byStage := map[string]map[string]any{}
	for _, summary := range anySlice(summariesValue.Export()) {
		byStage[fmt.Sprint(exportMap(summary)["stageId"])] = exportMap(summary)
	}
	segments := []any{}
	for _, stage := range anySlice(pipeline["stages"]) {
		s := exportMap(stage)
		summary := byStage[fmt.Sprint(s["id"])]
		count := any(0)
		if summary != nil && summary["count"] != nil {
			count = summary["count"]
		}
		segments = append(segments, map[string]any{"value": count, "styleKey": s["colorKey"], "label": s["name"]})
	}
	props := exportOptions(options)
	props["segments"] = segments
	props["styleSet"] = v3CRMStageStyleSet()
	props["showCounts"] = true
	props["size"] = "lg"
	return componentNode("SegmentedBar", props)
}

func (r *runtime) v3CRMIntentObject() *goja.Object {
	intent := r.vm.NewObject()
	setExport(intent, "openDeal", func(id goja.Value) map[string]any {
		return map[string]any{"kind": "navigate", "to": "/pages/opportunity?deal=" + v3URLTemplateValue(id)}
	})
	setExport(intent, "moveDeal", func(id, stage goja.Value) map[string]any {
		return map[string]any{"kind": "server", "name": "crm.deal.move", "payload": map[string]any{"dealId": v3CRMActionValue(id), "toStage": v3CRMActionValue(stage)}}
	})
	setExport(intent, "updateField", func(recordID, key, value goja.Value) map[string]any {
		return map[string]any{"kind": "server", "name": "crm.field.update", "payload": map[string]any{"recordId": v3CRMActionValue(recordID), "key": v3CRMActionValue(key), "value": v3CRMActionValue(value)}}
	})
	setExport(intent, "completeTask", func(id goja.Value) map[string]any {
		return map[string]any{"kind": "server", "name": "crm.task.complete", "payload": map[string]any{"taskId": v3CRMActionValue(id)}}
	})
	return intent
}

// v3CRMActionValue converts the concise JavaScript context placeholder form
// ("${cardId}" or "$cardId") into the typed payload contract that the React
// action dispatcher resolves. Other values remain serializable literals.
func v3CRMActionValue(value goja.Value) any {
	if value == nil || goja.IsUndefined(value) || goja.IsNull(value) {
		return nil
	}
	text, ok := value.Export().(string)
	if !ok {
		return value.Export()
	}
	if strings.HasPrefix(text, "${") && strings.HasSuffix(text, "}") && len(text) > 3 {
		return map[string]any{"kind": "path", "path": text[2 : len(text)-1]}
	}
	if strings.HasPrefix(text, "$") && len(text) > 1 && !strings.ContainsAny(text[1:], " /?&=") {
		return map[string]any{"kind": "path", "path": text[1:]}
	}
	return text
}

func (r *runtime) attachCRMRef(obj *goja.Object, ref *v3CRMRef) {
	if err := obj.DefineDataProperty(v3CRMRefProperty, r.vm.ToValue(ref), goja.FLAG_FALSE, goja.FLAG_FALSE, goja.FLAG_FALSE); err != nil {
		panic(err)
	}
}

func (r *runtime) crmFieldsFromValue(value goja.Value) map[string]any {
	if ref := r.crmRef(value, "fields"); ref != nil {
		return v3CRMFieldsSnapshot(ref.fields)
	}
	return exportObject(value)
}

func (r *runtime) crmPipelineFromValue(value goja.Value) map[string]any {
	if ref := r.crmRef(value, "pipeline"); ref != nil {
		return v3CRMPipelineSnapshot(ref.pipeline)
	}
	return exportObject(value)
}

func (r *runtime) crmRef(value goja.Value, want string) *v3CRMRef {
	if value == nil || goja.IsUndefined(value) || goja.IsNull(value) {
		return nil
	}
	ref, _ := value.ToObject(r.vm).Get(v3CRMRefProperty).Export().(*v3CRMRef)
	if ref == nil || ref.kind != want {
		return nil
	}
	return ref
}

func v3CRMFieldsSnapshot(spec *v3CRMFieldsSpec) map[string]any {
	return map[string]any{"name": spec.Name, "fields": append([]map[string]any(nil), spec.Fields...)}
}
func v3CRMPipelineSnapshot(spec *v3CRMPipelineSpec) map[string]any {
	return map[string]any{"id": spec.ID, "name": spec.Name, "stages": append([]map[string]any(nil), spec.Stages...)}
}
func v3CRMFieldsIssues(spec *v3CRMFieldsSpec) []map[string]any {
	if len(spec.Fields) == 0 {
		return []map[string]any{v3ValidationIssue("crm_fields_empty", "fields", "at least one field is required")}
	}
	return nil
}
func v3CRMPipelineIssues(spec *v3CRMPipelineSpec) []map[string]any {
	if len(spec.Stages) == 0 {
		return []map[string]any{v3ValidationIssue("crm_pipeline_empty", "pipeline.stages", "at least one stage is required")}
	}
	return nil
}

func v3CRMFieldSections(fields any) []any {
	order, groups := []string{}, map[string][]any{}
	for _, value := range anySlice(fields) {
		field := exportMap(value)
		group := fmt.Sprint(field["group"])
		if group == "" || group == "<nil>" {
			group = "Details"
		}
		if _, ok := groups[group]; !ok {
			order = append(order, group)
		}
		groups[group] = append(groups[group], field)
	}
	sections := make([]any, 0, len(order))
	for _, group := range order {
		sections = append(sections, map[string]any{"label": group, "fields": groups[group]})
	}
	return sections
}

func v3CRMBoardProps(pipeline map[string]any, cards, summaries []any) map[string]any {
	byStage := map[string]map[string]any{}
	for _, value := range summaries {
		s := exportMap(value)
		byStage[fmt.Sprint(s["stageId"])] = s
	}
	columns := []any{}
	for _, value := range anySlice(pipeline["stages"]) {
		stage := exportMap(value)
		label := fmt.Sprint(stage["name"])
		if s := byStage[fmt.Sprint(stage["id"])]; s != nil {
			label = fmt.Sprintf("%s · %v · %v", label, s["amountTotal"], s["count"])
		}
		columns = append(columns, map[string]any{"id": stage["id"], "header": label, "accent": stage["colorKey"]})
	}
	return map[string]any{"ariaLabel": pipeline["name"], "columns": columns, "cards": cards, "columnField": "stageId", "getCardId": map[string]any{"field": "id"}, "styleSet": v3CRMStageStyleSet(), "card": map[string]any{"title": map[string]any{"kind": "field", "field": "title"}, "subtitle": map[string]any{"kind": "number", "field": "amount", "format": "integer", "fallback": "—"}, "meta": map[string]any{"kind": "field", "field": "ownerId", "fallback": "unassigned"}, "accentField": "status"}}
}

func v3CRMStageStyleSet() map[string]any {
	return map[string]any{"id": "crm-stages", "styles": map[string]any{
		"lead":        map[string]any{"fill": "color-mix(in srgb, var(--mac-text-dim) 20%, var(--mac-surface))", "line": "var(--mac-text-dim)", "labelColor": "var(--mac-text)"},
		"qualified":   map[string]any{"fill": "color-mix(in srgb, var(--mac-accent) 22%, var(--mac-surface))", "line": "var(--mac-accent)", "labelColor": "var(--mac-text)"},
		"proposal":    map[string]any{"fill": "color-mix(in srgb, var(--mac-amber) 26%, var(--mac-surface))", "line": "var(--mac-amber)", "labelColor": "var(--mac-text)"},
		"negotiation": map[string]any{"fill": "color-mix(in srgb, var(--mac-accent-2) 24%, var(--mac-surface))", "line": "var(--mac-accent-2)", "labelColor": "var(--mac-text)"},
		"won":         map[string]any{"fill": "color-mix(in srgb, var(--mac-green) 26%, var(--mac-surface))", "line": "var(--mac-green)", "labelColor": "var(--mac-text)"},
		"lost":        map[string]any{"fill": "color-mix(in srgb, var(--mac-accent-2) 30%, var(--mac-surface))", "line": "var(--mac-accent-2)", "labelColor": "var(--mac-text)"},
	}}
}

func exportMap(value any) map[string]any {
	if out, ok := value.(map[string]any); ok {
		return out
	}
	return map[string]any{}
}
