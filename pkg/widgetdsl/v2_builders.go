package widgetdsl

import (
	"fmt"
	"strings"

	"github.com/dop251/goja"
	v2spec "github.com/go-go-golems/rag-evaluation-system/pkg/widgetdsl/v2/spec"
)

const v2RefProperty = "__widgetdsl_v2_ref"

type v2Ref struct {
	kind       string
	field      *v2spec.FieldSpec
	schema     *v2spec.SchemaSpec
	collection *v2spec.CollectionSpec
	action     *v2spec.ActionSpec
	selection  *v2spec.SelectionSpec
}

func (r *runtime) installDataV2(exports *goja.Object) {
	setExport(exports, "f", r.v2FieldFactoryObject())
	setExport(exports, "schema", r.v2SchemaCtor)
	setExport(exports, "collection", r.v2CollectionCtor)
	setExport(exports, "selection", r.v2SelectionObject())
	setExport(exports, "action", r.v2ActionObject())
}

func (r *runtime) v2FieldFactoryObject() *goja.Object {
	f := r.vm.NewObject()
	setExport(f, "key", func() *goja.Object {
		return r.v2FieldBuilder(v2spec.FieldSpec{Kind: v2spec.FieldKindString, Semantic: v2spec.FieldSemanticKey, Editor: v2spec.EditorSpec{Control: v2spec.EditorControlText}, Summary: v2spec.SummarySpec{CellKind: "caption"}})
	})
	setExport(f, "primary", func() *goja.Object {
		return r.v2FieldBuilder(v2spec.FieldSpec{Kind: v2spec.FieldKindString, Semantic: v2spec.FieldSemanticPrimary, Editor: v2spec.EditorSpec{Control: v2spec.EditorControlText}, Summary: v2spec.SummarySpec{CellKind: "field"}})
	})
	setExport(f, "short", func() *goja.Object {
		return r.v2FieldBuilder(v2spec.FieldSpec{Kind: v2spec.FieldKindString, Semantic: v2spec.FieldSemanticShort, Editor: v2spec.EditorSpec{Control: v2spec.EditorControlText}, Summary: v2spec.SummarySpec{CellKind: "field"}})
	})
	setExport(f, "prose", func() *goja.Object {
		return r.v2FieldBuilder(v2spec.FieldSpec{Kind: v2spec.FieldKindString, Semantic: v2spec.FieldSemanticProse, Editor: v2spec.EditorSpec{Control: v2spec.EditorControlTextarea, Rows: 4}, Summary: v2spec.SummarySpec{Elide: true}})
	})
	setExport(f, "count", func() *goja.Object {
		return r.v2FieldBuilder(v2spec.FieldSpec{Kind: v2spec.FieldKindNumber, Semantic: v2spec.FieldSemanticCount, Editor: v2spec.EditorSpec{Control: v2spec.EditorControlText}, Summary: v2spec.SummarySpec{CellKind: "number"}})
	})
	setExport(f, "status", func() *goja.Object {
		return r.v2FieldBuilder(v2spec.FieldSpec{Kind: v2spec.FieldKindString, Semantic: v2spec.FieldSemanticStatus, Editor: v2spec.EditorSpec{Control: v2spec.EditorControlText}, Summary: v2spec.SummarySpec{CellKind: "status"}})
	})
	return f
}

func (r *runtime) v2FieldBuilder(field v2spec.FieldSpec) *goja.Object {
	fieldCopy := field
	obj := r.vm.NewObject()
	r.attachV2Ref(obj, &v2Ref{kind: "field", field: &fieldCopy})
	setExport(obj, "label", func(label string) *goja.Object {
		fieldCopy.Label = label
		return obj
	})
	setExport(obj, "width", func(width string) *goja.Object {
		fieldCopy.Layout.Width = width
		return obj
	})
	setExport(obj, "required", func() *goja.Object {
		fieldCopy.Validation.Required = true
		return obj
	})
	setExport(obj, "maxLength", func(limit int) *goja.Object {
		fieldCopy.Validation.MaxLength = limit
		return obj
	})
	setExport(obj, "rows", func(rows int) *goja.Object {
		fieldCopy.Editor.Control = v2spec.EditorControlTextarea
		fieldCopy.Editor.Rows = rows
		return obj
	})
	setExport(obj, "readOnly", func() *goja.Object {
		fieldCopy.Editor.ReadOnly = true
		return obj
	})
	setExport(obj, "build", func() map[string]any {
		return map[string]any{"kind": string(fieldCopy.Kind), "semantic": string(fieldCopy.Semantic), "label": fieldCopy.Label}
	})
	return obj
}

func (r *runtime) v2SchemaCtor(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) == 0 || strings.TrimSpace(call.Arguments[0].String()) == "" {
		panic(r.vm.NewGoError(fmt.Errorf("data.v2.dsl schema(name) requires a non-empty name")))
	}
	schema := &v2spec.SchemaSpec{Name: call.Arguments[0].String()}
	obj := r.vm.NewObject()
	r.attachV2Ref(obj, &v2Ref{kind: "schemaBuilder", schema: schema})
	setExport(obj, "field", func(name string, fieldValue goja.Value) *goja.Object {
		fieldRef := r.mustV2Ref(fieldValue, "field")
		field := *fieldRef.field
		field.Name = name
		schema.Fields = append(schema.Fields, field)
		return obj
	})
	setExport(obj, "build", func() *goja.Object {
		built := *schema
		return r.v2SchemaValue(&built)
	})
	setExport(obj, "validate", func() []map[string]any {
		return validationIssuesForJS(schema.Validate("schema"))
	})
	return obj
}

func (r *runtime) v2SchemaValue(schema *v2spec.SchemaSpec) *goja.Object {
	obj := r.vm.NewObject()
	r.attachV2Ref(obj, &v2Ref{kind: "schema", schema: schema})
	setExport(obj, "validate", func() []map[string]any {
		return validationIssuesForJS(schema.Validate("schema"))
	})
	return obj
}

func (r *runtime) v2CollectionCtor(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) < 2 {
		panic(r.vm.NewGoError(fmt.Errorf("data.v2.dsl collection(name, rows) requires a name and rows")))
	}
	name := strings.TrimSpace(call.Arguments[0].String())
	if name == "" {
		panic(r.vm.NewGoError(fmt.Errorf("data.v2.dsl collection(name, rows) requires a non-empty name")))
	}
	collection := &v2spec.CollectionSpec{
		Name:        name,
		Rows:        v2Rows(call.Arguments[1].Export()),
		Mode:        v2spec.CollectionModeShow,
		Arrangement: v2spec.ArrangementSpec{Kind: v2spec.ArrangementKindTable},
	}
	obj := r.vm.NewObject()
	r.attachV2Ref(obj, &v2Ref{kind: "collectionBuilder", collection: collection})
	setExport(obj, "schema", func(schemaValue goja.Value) *goja.Object {
		schemaRef := r.mustV2Ref(schemaValue, "schema")
		collection.Schema = *schemaRef.schema
		return obj
	})
	setExport(obj, "empty", func(message string) *goja.Object {
		collection.Empty = message
		return obj
	})
	setExport(obj, "select", func(arg goja.Value) *goja.Object {
		if fn, ok := goja.AssertFunction(arg); ok {
			selectionBuilder := r.v2SelectionBuilder()
			ret, err := fn(goja.Undefined(), selectionBuilder)
			if err != nil {
				panic(err)
			}
			if !goja.IsUndefined(ret) && !goja.IsNull(ret) {
				collection.Selection = r.mustV2Ref(ret, "selection").selection
			} else {
				collection.Selection = r.mustV2Ref(selectionBuilder, "selectionBuilder").selection
			}
			return obj
		}
		if arg == nil || goja.IsUndefined(arg) || goja.IsNull(arg) {
			return obj
		}
		collection.Selection = r.mustV2Ref(arg, "selection").selection
		return obj
	})
	setExport(obj, "edit", func(args ...goja.Value) *goja.Object {
		collection.Mode = v2spec.CollectionModeEdit
		if len(args) > 0 && !goja.IsUndefined(args[0]) && !goja.IsNull(args[0]) {
			fn, ok := goja.AssertFunction(args[0])
			if !ok {
				panic(r.vm.NewGoError(fmt.Errorf("data.v2.dsl collection.edit(callback) requires a function when an argument is present")))
			}
			if _, err := fn(goja.Undefined(), r.v2EditorBuilder(collection)); err != nil {
				panic(err)
			}
		}
		return obj
	})
	setExport(obj, "table", func(args ...goja.Value) *goja.Object {
		collection.Arrangement = v2spec.ArrangementSpec{Kind: v2spec.ArrangementKindTable}
		if len(args) > 0 && !goja.IsUndefined(args[0]) && !goja.IsNull(args[0]) {
			fn, ok := goja.AssertFunction(args[0])
			if !ok {
				panic(r.vm.NewGoError(fmt.Errorf("data.v2.dsl collection.table(callback) requires a function when an argument is present")))
			}
			if _, err := fn(goja.Undefined(), r.v2TableBuilder(collection)); err != nil {
				panic(err)
			}
		}
		return obj
	})
	setExport(obj, "masterDetail", func() *goja.Object {
		collection.Arrangement = v2spec.ArrangementSpec{Kind: v2spec.ArrangementKindMasterDetail}
		return obj
	})
	setExport(obj, "validate", func() []map[string]any {
		return validationIssuesForJS(collection.Validate("collection"))
	})
	setExport(obj, "toIR", func() any {
		issues := collection.Validate("collection")
		if v2spec.HasErrors(issues) {
			panic(r.vm.NewGoError(fmt.Errorf("data.v2.dsl collection %q is invalid: %s", collection.Name, firstValidationError(issues))))
		}
		return collection.ToNode().ToWidgetNode()
	})
	return obj
}

func (r *runtime) v2SelectionObject() *goja.Object {
	selection := r.vm.NewObject()
	setExport(selection, "urlParam", func(param string, value goja.Value) *goja.Object {
		return r.v2SelectionValue(&v2spec.SelectionSpec{Kind: v2spec.SelectionKindURLParam, Param: param, Value: stringifyValue(value)})
	})
	return selection
}

func (r *runtime) v2SelectionBuilder() *goja.Object {
	spec := &v2spec.SelectionSpec{}
	obj := r.vm.NewObject()
	r.attachV2Ref(obj, &v2Ref{kind: "selectionBuilder", selection: spec})
	setExport(obj, "urlParam", func(param string, value goja.Value) *goja.Object {
		spec.Kind = v2spec.SelectionKindURLParam
		spec.Param = param
		spec.Value = stringifyValue(value)
		return r.v2SelectionValue(spec)
	})
	return obj
}

func (r *runtime) v2SelectionValue(selection *v2spec.SelectionSpec) *goja.Object {
	obj := r.vm.NewObject()
	r.attachV2Ref(obj, &v2Ref{kind: "selection", selection: selection})
	return obj
}

func (r *runtime) v2EditorBuilder(collection *v2spec.CollectionSpec) *goja.Object {
	obj := r.vm.NewObject()
	setExport(obj, "selectUrl", func(param string, value goja.Value) *goja.Object {
		collection.Selection = &v2spec.SelectionSpec{Kind: v2spec.SelectionKindURLParam, Param: param, Value: stringifyValue(value)}
		return obj
	})
	setExport(obj, "submitPost", func(formAction string) *goja.Object {
		collection.Actions.Submit = &v2spec.SubmitSpec{FormAction: formAction, Method: "post"}
		return obj
	})
	setExport(obj, "create", func(value goja.Value) *goja.Object {
		label := "New item"
		if value != nil && !goja.IsUndefined(value) && !goja.IsNull(value) {
			if obj := value.ToObject(r.vm); obj != nil {
				if labelValue := obj.Get("label"); labelValue != nil && !goja.IsUndefined(labelValue) && !goja.IsNull(labelValue) {
					label = labelValue.String()
				} else if exported, ok := value.Export().(string); ok && exported != "" {
					label = exported
				}
			}
		}
		collection.Actions.Create = &v2spec.CreateActionSpec{Label: label}
		return obj
	})
	setExport(obj, "reorder", func(actionValue goja.Value) *goja.Object {
		collection.Actions.Reorder = r.mustV2Ref(actionValue, "action").action
		return obj
	})
	setExport(obj, "remove", func(actionValue goja.Value) *goja.Object {
		collection.Actions.Remove = r.mustV2Ref(actionValue, "action").action
		return obj
	})
	setExport(obj, "actions", func(callback goja.Value) *goja.Object {
		fn, ok := goja.AssertFunction(callback)
		if !ok {
			panic(r.vm.NewGoError(fmt.Errorf("data.v2.dsl edit.actions(callback) requires a function")))
		}
		if _, err := fn(goja.Undefined(), r.v2CollectionActionsBuilder(collection)); err != nil {
			panic(err)
		}
		return obj
	})
	return obj
}

func (r *runtime) v2CollectionActionsBuilder(collection *v2spec.CollectionSpec) *goja.Object {
	obj := r.vm.NewObject()
	setExport(obj, "reorder", func(actionValue goja.Value) *goja.Object {
		collection.Actions.Reorder = r.mustV2Ref(actionValue, "action").action
		return obj
	})
	setExport(obj, "remove", func(actionValue goja.Value) *goja.Object {
		collection.Actions.Remove = r.mustV2Ref(actionValue, "action").action
		return obj
	})
	return obj
}

func (r *runtime) v2TableBuilder(collection *v2spec.CollectionSpec) *goja.Object {
	obj := r.vm.NewObject()
	setExport(obj, "className", func(className string) *goja.Object {
		collection.Table.ClassName = className
		return obj
	})
	setExport(obj, "rowSelect", func(actionValue goja.Value) *goja.Object {
		actionRef := r.mustV2Ref(actionValue, "action")
		collection.Table.RowSelect = actionRef.action
		return obj
	})
	setExport(obj, "actionColumn", func(id string, header string, label string, actionValue goja.Value, options ...goja.Value) *goja.Object {
		actionRef := r.mustV2Ref(actionValue, "action")
		column := v2spec.TableActionColumnSpec{ID: id, Header: header, Label: label, Action: *actionRef.action}
		if opts := exportOptions(options); opts != nil {
			if maxWidth, ok := opts["maxWidth"].(string); ok {
				column.MaxWidth = maxWidth
			}
		}
		collection.Table.ActionColumns = append(collection.Table.ActionColumns, column)
		return obj
	})
	return obj
}

func (r *runtime) v2ActionObject() *goja.Object {
	action := r.vm.NewObject()
	setExport(action, "navigate", func(to string) *goja.Object {
		return r.v2ActionValue(&v2spec.ActionSpec{Kind: v2spec.ActionKindNavigate, To: to})
	})
	setExport(action, "server", func(name string) *goja.Object {
		return r.v2ActionValue(&v2spec.ActionSpec{Kind: v2spec.ActionKindServer, Name: name})
	})
	return action
}

func (r *runtime) v2ActionValue(action *v2spec.ActionSpec) *goja.Object {
	obj := r.vm.NewObject()
	r.attachV2Ref(obj, &v2Ref{kind: "action", action: action})
	setExport(obj, "confirm", func(text string) *goja.Object {
		action.Confirm = &v2spec.TemplateSpec{Parts: []v2spec.TemplateValue{{Kind: v2spec.TemplateValueText, Text: text}}}
		return obj
	})
	setExport(obj, "payloadPath", func(name string, path string) *goja.Object {
		action.Payload.Fields = append(action.Payload.Fields, v2spec.PayloadFieldSpec{Name: name, Value: v2spec.TemplateValue{Kind: v2spec.TemplateValuePath, Path: path}})
		return obj
	})
	setExport(obj, "payload", func(name string, value goja.Value) *goja.Object {
		action.Payload.Fields = append(action.Payload.Fields, v2spec.PayloadFieldSpec{Name: name, Value: v2spec.TemplateValue{Kind: v2spec.TemplateValueLiteral, Value: value.Export()}})
		return obj
	})
	return obj
}

func (r *runtime) attachV2Ref(obj *goja.Object, ref *v2Ref) {
	if err := obj.DefineDataProperty(v2RefProperty, r.vm.ToValue(ref), goja.FLAG_FALSE, goja.FLAG_FALSE, goja.FLAG_FALSE); err != nil {
		panic(err)
	}
}

func (r *runtime) mustV2Ref(value goja.Value, wantKind string) *v2Ref {
	if value == nil || goja.IsUndefined(value) || goja.IsNull(value) {
		panic(r.vm.NewGoError(fmt.Errorf("expected v2 %s handle, got null/undefined", wantKind)))
	}
	obj := value.ToObject(r.vm)
	refValue := obj.Get(v2RefProperty)
	if refValue == nil || goja.IsUndefined(refValue) || goja.IsNull(refValue) {
		panic(r.vm.NewGoError(fmt.Errorf("expected v2 %s handle", wantKind)))
	}
	ref, ok := refValue.Export().(*v2Ref)
	if !ok || ref == nil {
		panic(r.vm.NewGoError(fmt.Errorf("invalid v2 %s handle", wantKind)))
	}
	if ref.kind != wantKind && (wantKind != "selection" || ref.kind != "selectionBuilder") {
		panic(r.vm.NewGoError(fmt.Errorf("expected v2 %s handle, got %s", wantKind, ref.kind)))
	}
	return ref
}

func v2Rows(value any) []v2spec.JSONObject {
	items := anySlice(value)
	rows := make([]v2spec.JSONObject, 0, len(items))
	for _, item := range items {
		row := v2spec.JSONObject{}
		if m, ok := item.(map[string]any); ok {
			for k, v := range m {
				row[k] = v
			}
		}
		rows = append(rows, row)
	}
	return rows
}

func validationIssuesForJS(issues []v2spec.ValidationIssue) []map[string]any {
	out := make([]map[string]any, 0, len(issues))
	for _, issue := range issues {
		out = append(out, map[string]any{"severity": string(issue.Severity), "code": issue.Code, "path": issue.Path, "message": issue.Message, "hint": issue.Hint})
	}
	return out
}

func firstValidationError(issues []v2spec.ValidationIssue) string {
	for _, issue := range issues {
		if issue.Severity == v2spec.ValidationSeverityError {
			return issue.Code + " at " + issue.Path + ": " + issue.Message
		}
	}
	return "unknown validation error"
}
