package widgetdsl

import (
	"fmt"
	"strings"

	"github.com/dop251/goja"
	widgetspec "github.com/go-go-golems/rag-evaluation-system/pkg/widgetdsl/spec"
)

const v2RefProperty = "__widgetdsl_v2_ref"

type v2Ref struct {
	kind       string
	field      *widgetspec.FieldSpec
	schema     *widgetspec.SchemaSpec
	collection *widgetspec.CollectionSpec
	action     *widgetspec.ActionSpec
	selection  *widgetspec.SelectionSpec
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
		return r.v2FieldBuilder(widgetspec.FieldSpec{Kind: widgetspec.FieldKindString, Semantic: widgetspec.FieldSemanticKey, Editor: widgetspec.EditorSpec{Control: widgetspec.EditorControlText}, Summary: widgetspec.SummarySpec{CellKind: "caption"}})
	})
	setExport(f, "primary", func() *goja.Object {
		return r.v2FieldBuilder(widgetspec.FieldSpec{Kind: widgetspec.FieldKindString, Semantic: widgetspec.FieldSemanticPrimary, Editor: widgetspec.EditorSpec{Control: widgetspec.EditorControlText}, Summary: widgetspec.SummarySpec{CellKind: "field"}})
	})
	setExport(f, "short", func() *goja.Object {
		return r.v2FieldBuilder(widgetspec.FieldSpec{Kind: widgetspec.FieldKindString, Semantic: widgetspec.FieldSemanticShort, Editor: widgetspec.EditorSpec{Control: widgetspec.EditorControlText}, Summary: widgetspec.SummarySpec{CellKind: "field"}})
	})
	setExport(f, "prose", func() *goja.Object {
		return r.v2FieldBuilder(widgetspec.FieldSpec{Kind: widgetspec.FieldKindString, Semantic: widgetspec.FieldSemanticProse, Editor: widgetspec.EditorSpec{Control: widgetspec.EditorControlTextarea, Rows: 4}, Summary: widgetspec.SummarySpec{Elide: true}})
	})
	setExport(f, "count", func() *goja.Object {
		return r.v2FieldBuilder(widgetspec.FieldSpec{Kind: widgetspec.FieldKindNumber, Semantic: widgetspec.FieldSemanticCount, Editor: widgetspec.EditorSpec{Control: widgetspec.EditorControlText}, Summary: widgetspec.SummarySpec{CellKind: "number"}})
	})
	setExport(f, "status", func() *goja.Object {
		return r.v2FieldBuilder(widgetspec.FieldSpec{Kind: widgetspec.FieldKindString, Semantic: widgetspec.FieldSemanticStatus, Editor: widgetspec.EditorSpec{Control: widgetspec.EditorControlText}, Summary: widgetspec.SummarySpec{CellKind: "status"}})
	})
	return f
}

func (r *runtime) v2FieldBuilder(field widgetspec.FieldSpec) *goja.Object {
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
		fieldCopy.Editor.Control = widgetspec.EditorControlTextarea
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
	schema := &widgetspec.SchemaSpec{Name: call.Arguments[0].String()}
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

func (r *runtime) v2SchemaValue(schema *widgetspec.SchemaSpec) *goja.Object {
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
	collection := &widgetspec.CollectionSpec{
		Name:        name,
		Rows:        v2Rows(call.Arguments[1].Export()),
		Mode:        widgetspec.CollectionModeShow,
		Arrangement: widgetspec.ArrangementSpec{Kind: widgetspec.ArrangementKindTable},
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
		collection.Mode = widgetspec.CollectionModeEdit
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
		collection.Arrangement = widgetspec.ArrangementSpec{Kind: widgetspec.ArrangementKindTable}
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
		collection.Arrangement = widgetspec.ArrangementSpec{Kind: widgetspec.ArrangementKindMasterDetail}
		return obj
	})
	setExport(obj, "validate", func() []map[string]any {
		return validationIssuesForJS(collection.Validate("collection"))
	})
	setExport(obj, "toIR", func() any {
		issues := collection.Validate("collection")
		if widgetspec.HasErrors(issues) {
			panic(r.vm.NewGoError(fmt.Errorf("data.v2.dsl collection %q is invalid: %s", collection.Name, firstValidationError(issues))))
		}
		return collection.ToNode().ToWidgetNode()
	})
	return obj
}

func (r *runtime) v2SelectionObject() *goja.Object {
	selection := r.vm.NewObject()
	setExport(selection, "urlParam", func(param string, value goja.Value) *goja.Object {
		return r.v2SelectionValue(&widgetspec.SelectionSpec{Kind: widgetspec.SelectionKindURLParam, Param: param, Value: stringifyValue(value)})
	})
	return selection
}

func (r *runtime) v2SelectionBuilder() *goja.Object {
	spec := &widgetspec.SelectionSpec{}
	obj := r.vm.NewObject()
	r.attachV2Ref(obj, &v2Ref{kind: "selectionBuilder", selection: spec})
	setExport(obj, "urlParam", func(param string, value goja.Value) *goja.Object {
		spec.Kind = widgetspec.SelectionKindURLParam
		spec.Param = param
		spec.Value = stringifyValue(value)
		return r.v2SelectionValue(spec)
	})
	return obj
}

func (r *runtime) v2SelectionValue(selection *widgetspec.SelectionSpec) *goja.Object {
	obj := r.vm.NewObject()
	r.attachV2Ref(obj, &v2Ref{kind: "selection", selection: selection})
	return obj
}

func (r *runtime) v2EditorBuilder(collection *widgetspec.CollectionSpec) *goja.Object {
	obj := r.vm.NewObject()
	setExport(obj, "selectUrl", func(param string, value goja.Value) *goja.Object {
		collection.Selection = &widgetspec.SelectionSpec{Kind: widgetspec.SelectionKindURLParam, Param: param, Value: stringifyValue(value)}
		return obj
	})
	setExport(obj, "submitPost", func(formAction string) *goja.Object {
		collection.Actions.Submit = &widgetspec.SubmitSpec{FormAction: formAction, Method: "post"}
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
		collection.Actions.Create = &widgetspec.CreateActionSpec{Label: label}
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

func (r *runtime) v2CollectionActionsBuilder(collection *widgetspec.CollectionSpec) *goja.Object {
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

func (r *runtime) v2TableBuilder(collection *widgetspec.CollectionSpec) *goja.Object {
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
		column := widgetspec.TableActionColumnSpec{ID: id, Header: header, Label: label, Action: *actionRef.action}
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
		return r.v2ActionValue(&widgetspec.ActionSpec{Kind: widgetspec.ActionKindNavigate, To: to})
	})
	setExport(action, "server", func(name string) *goja.Object {
		return r.v2ActionValue(&widgetspec.ActionSpec{Kind: widgetspec.ActionKindServer, Name: name})
	})
	return action
}

func (r *runtime) v2ActionValue(action *widgetspec.ActionSpec) *goja.Object {
	obj := r.vm.NewObject()
	r.attachV2Ref(obj, &v2Ref{kind: "action", action: action})
	setExport(obj, "confirm", func(text string) *goja.Object {
		action.Confirm = &widgetspec.TemplateSpec{Parts: []widgetspec.TemplateValue{{Kind: widgetspec.TemplateValueText, Text: text}}}
		return obj
	})
	setExport(obj, "payloadPath", func(name string, path string) *goja.Object {
		action.Payload.Fields = append(action.Payload.Fields, widgetspec.PayloadFieldSpec{Name: name, Value: widgetspec.TemplateValue{Kind: widgetspec.TemplateValuePath, Path: path}})
		return obj
	})
	setExport(obj, "payload", func(name string, value goja.Value) *goja.Object {
		action.Payload.Fields = append(action.Payload.Fields, widgetspec.PayloadFieldSpec{Name: name, Value: widgetspec.TemplateValue{Kind: widgetspec.TemplateValueLiteral, Value: value.Export()}})
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

func v2Rows(value any) []widgetspec.JSONObject {
	items := anySlice(value)
	rows := make([]widgetspec.JSONObject, 0, len(items))
	for _, item := range items {
		row := widgetspec.JSONObject{}
		if m, ok := item.(map[string]any); ok {
			for k, v := range m {
				row[k] = v
			}
		}
		rows = append(rows, row)
	}
	return rows
}

func validationIssuesForJS(issues []widgetspec.ValidationIssue) []map[string]any {
	out := make([]map[string]any, 0, len(issues))
	for _, issue := range issues {
		out = append(out, map[string]any{"severity": string(issue.Severity), "code": issue.Code, "path": issue.Path, "message": issue.Message, "hint": issue.Hint})
	}
	return out
}

func firstValidationError(issues []widgetspec.ValidationIssue) string {
	for _, issue := range issues {
		if issue.Severity == widgetspec.ValidationSeverityError {
			return issue.Code + " at " + issue.Path + ": " + issue.Message
		}
	}
	return "unknown validation error"
}
