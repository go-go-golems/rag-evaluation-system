package widgetdsl

// Grammar verbs: intent-level authoring on top of the component helpers.
//
// data.dsl gains the data grammar — schema/f (field roles), record (one
// record shown or edited), collection (records shown/edited through an
// arrangement), urlParam (URL-backed selection), formPost (native form
// submits) — and ui.dsl gains the structure grammar (section). Grammar calls
// compile to plain Widget IR built from existing components (SectionBlock,
// FieldGrid, FormPanel, FormRow, DataTable, Stack, Inline, Button, Caption),
// so the renderer needs nothing new. See ticket RAGEVAL-UI-GRAMMAR design-doc
// 02 for the rationale and the module-reorganization story.

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dop251/goja"
)

// fieldRoles drive how a field renders in summaries (table cells), editors
// (form controls), and read-only views. Roles, not types: `number` and
// `duration` are both strings — the role says how they behave.
var fieldRoles = []string{
	"key",     // stable identity; muted in summaries, read-only in editors
	"primary", // the scannable column; required-ish text
	"short",   // one-line text
	"prose",   // multi-line text; elided from summaries
	"count",   // numeric
	"size",    // numeric (bytes-ish)
	"measure", // numeric against an optional limit (bars, treemaps)
	"date",    // date-ish string
	"status",  // small vocabulary rendered as StatusText in summaries
	"tags",    // string list; comma-edited
	"media",   // image/file source; elided from summaries
	"href",    // link target
}

var gridableRoles = map[string]bool{
	"key": true, "short": true, "count": true, "size": true,
	"measure": true, "date": true, "status": true, "href": true,
}

func (r *runtime) installDataGrammar(exports *goja.Object) {
	setExport(exports, "f", r.fieldRoleObject())
	setExport(exports, "schema", r.schemaCtor)
	setExport(exports, "record", r.recordVerb)
	setExport(exports, "collection", r.collectionVerb)
	setExport(exports, "urlParam", func(param string, value goja.Value) map[string]any {
		out := map[string]any{"param": param, "value": ""}
		if value != nil && !goja.IsUndefined(value) && !goja.IsNull(value) {
			out["value"] = strings.TrimSpace(value.String())
		}
		return out
	})
	setExport(exports, "formPost", func(formAction string, options ...goja.Value) map[string]any {
		out := map[string]any{"formAction": formAction, "method": "post"}
		mergeOptions(out, exportOptions(options))
		return out
	})
}

func (r *runtime) fieldRoleObject() *goja.Object {
	f := r.vm.NewObject()
	for _, role := range fieldRoles {
		role := role
		setExport(f, role, func(options ...goja.Value) map[string]any {
			out := map[string]any{"role": role}
			mergeOptions(out, exportOptions(options))
			return out
		})
	}
	return f
}

// schemaCtor preserves field order (goja object keys are insertion-ordered);
// plain map exports would lose it.
func (r *runtime) schemaCtor(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) == 0 || !isPlainObject(call.Arguments[0]) {
		panic(r.vm.NewGoError(fmt.Errorf("data.dsl schema(fields) requires an object of f.* field specs")))
	}
	obj := call.Arguments[0].ToObject(r.vm)
	fields := []any{}
	for _, key := range obj.Keys() {
		spec := exportObject(obj.Get(key))
		if _, ok := spec["role"]; !ok {
			panic(r.vm.NewGoError(fmt.Errorf("data.dsl schema field %q must be built with f.<role>(...)", key)))
		}
		spec["name"] = key
		fields = append(fields, spec)
	}
	return r.vm.ToValue(map[string]any{"__ragSchema": true, "fields": fields})
}

func schemaFields(value any) []map[string]any {
	schema, _ := value.(map[string]any)
	out := []map[string]any{}
	for _, entry := range anySlice(schema["fields"]) {
		if field, ok := entry.(map[string]any); ok {
			out = append(out, field)
		}
	}
	return out
}

func fieldLabel(field map[string]any) string {
	if label := stringFromMap(field, "label", ""); label != "" {
		return label
	}
	name := stringFromMap(field, "name", "")
	if name == "" {
		return ""
	}
	return strings.ToUpper(name[:1]) + name[1:]
}

func anyToString(value any) string {
	switch v := value.(type) {
	case nil:
		return ""
	case string:
		return v
	case float64:
		if v == float64(int64(v)) {
			return strconv.FormatInt(int64(v), 10)
		}
		return strconv.FormatFloat(v, 'f', -1, 64)
	case int64:
		return strconv.FormatInt(v, 10)
	case bool:
		return strconv.FormatBool(v)
	case []any:
		parts := make([]string, 0, len(v))
		for _, item := range v {
			parts = append(parts, anyToString(item))
		}
		return strings.Join(parts, ", ")
	default:
		return fmt.Sprintf("%v", v)
	}
}

// ─── record ─────────────────────────────────────────────────────────

// data.record(values, options) — one record, shown or edited.
// options: schema (required), verb "edit"|"show", arrange "field-grid"|"rows",
// title/subtitle, submit (data.formPost), status/statusMessage/submitLabel/
// resetLabel/footer passthrough to FormPanel.
func (r *runtime) recordVerb(call goja.FunctionCall) goja.Value {
	values := map[string]any{}
	if len(call.Arguments) > 0 && isPlainObject(call.Arguments[0]) {
		values = exportObject(call.Arguments[0])
	}
	options := map[string]any{}
	if len(call.Arguments) > 1 && isPlainObject(call.Arguments[1]) {
		options = exportObject(call.Arguments[1])
	}
	fields := schemaFields(options["schema"])
	if len(fields) == 0 {
		panic(r.vm.NewGoError(fmt.Errorf("data.dsl record(values, {schema}) requires a schema built with data.schema")))
	}
	verb := stringFromMap(options, "verb", "edit")
	if verb == "show" {
		items := []any{}
		for _, field := range fields {
			if field["role"] == "media" {
				continue
			}
			items = append(items, map[string]any{
				"key":   fieldLabel(field),
				"value": anyToString(values[stringFromMap(field, "name", "")]),
			})
		}
		return r.vm.ToValue(componentNode("MetadataGrid", map[string]any{"items": items, "density": "compact"}))
	}

	props := map[string]any{
		"title": stringFromMap(options, "title", "Details"),
	}
	copyIfPresent(props, options, "subtitle")
	copyIfPresent(props, options, "status")
	copyIfPresent(props, options, "statusMessage")
	copyIfPresent(props, options, "submitLabel")
	copyIfPresent(props, options, "resetLabel")
	copyIfPresent(props, options, "footer")
	if submit, ok := options["submit"].(map[string]any); ok {
		copyIfPresent(props, submit, "formAction")
		copyIfPresent(props, submit, "method")
	}
	rows := recordEditorRows(fields, values, stringFromMap(options, "arrange", "field-grid"))
	return r.vm.ToValue(componentNode("FormPanel", props, rows...))
}

func recordEditorRows(fields []map[string]any, values map[string]any, arrange string) []any {
	rows := []any{}
	grid := []any{}
	flushGrid := func() {
		if len(grid) == 0 {
			return
		}
		if len(grid) == 1 {
			rows = append(rows, grid[0])
		} else {
			columns := 2
			if len(grid) >= 3 {
				columns = 3
			}
			rows = append(rows, componentNode("FieldGrid", map[string]any{"columns": columns}, grid...))
		}
		grid = []any{}
	}
	for _, field := range fields {
		row := recordEditorRow(field, values)
		role, _ := field["role"].(string)
		if arrange == "field-grid" && gridableRoles[role] {
			grid = append(grid, row)
			continue
		}
		flushGrid()
		rows = append(rows, row)
	}
	flushGrid()
	return rows
}

func recordEditorRow(field map[string]any, values map[string]any) map[string]any {
	name := stringFromMap(field, "name", "")
	role, _ := field["role"].(string)
	value := anyToString(values[name])

	controlProps := map[string]any{
		"name":         name,
		"defaultValue": value,
		"readOnly":     boolFromMap(field, "readOnly", role == "key" && !boolFromMap(field, "editable", false)),
	}
	copyIfPresent(controlProps, field, "placeholder")
	copyIfPresent(controlProps, field, "maxLength")

	controlType := "TextInput"
	orientation := "inline"
	if role == "prose" {
		controlType = "TextareaInput"
		orientation = "stacked"
		controlProps["rows"] = valueOrDefault(field["rows"], 4)
	}

	rowProps := map[string]any{
		"label":       fieldLabel(field),
		"control":     componentNode(controlType, controlProps),
		"orientation": orientation,
	}
	copyIfPresent(rowProps, field, "hint")
	copyIfPresent(rowProps, field, "required")
	return componentNode("FormRow", rowProps)
}

// ─── collection ─────────────────────────────────────────────────────

// data.collection(rows, options) — records through an arrangement.
// options: schema (required), verb "show"|"edit"|"pick"|"manage",
// arrange "table"|"master-detail", title/caption, select (data.urlParam),
// open (action), submit (data.formPost, for verb edit), reorder (action),
// remove (action with confirm), create (bool or {label}), empty, getRowKey.
func (r *runtime) collectionVerb(call goja.FunctionCall) goja.Value {
	rows := []any{}
	if len(call.Arguments) > 0 {
		rows = anySlice(call.Arguments[0].Export())
	}
	options := map[string]any{}
	if len(call.Arguments) > 1 && isPlainObject(call.Arguments[1]) {
		options = exportObject(call.Arguments[1])
	}
	fields := schemaFields(options["schema"])
	if len(fields) == 0 {
		panic(r.vm.NewGoError(fmt.Errorf("data.dsl collection(rows, {schema}) requires a schema built with data.schema")))
	}
	verb := stringFromMap(options, "verb", "show")
	arrange := stringFromMap(options, "arrange", "table")
	keyField := collectionKeyField(fields, options)
	sel, _ := options["select"].(map[string]any)
	selParam := stringFromMap(sel, "param", "")
	selValue := stringFromMap(sel, "value", "")

	table := collectionTable(rows, fields, options, verb, keyField, selParam, selValue)
	children := []any{}
	if createNode := collectionCreateButton(options, verb, selParam); createNode != nil {
		children = append(children, createNode)
	}
	children = append(children, table)
	if arrange == "master-detail" {
		children = append(children, r.collectionDetail(rows, fields, options, verb, keyField, selParam, selValue))
	}

	body := componentNode("Stack", map[string]any{"gap": "md"}, children...)
	title := stringFromMap(options, "title", "")
	if title == "" {
		return r.vm.ToValue(body)
	}
	sectionProps := map[string]any{"label": title, "level": 2, "rule": true, "density": "flush"}
	if caption := stringFromMap(options, "caption", ""); caption != "" {
		sectionProps["caption"] = caption
	}
	return r.vm.ToValue(componentNode("SectionBlock", sectionProps, body))
}

func collectionKeyField(fields []map[string]any, options map[string]any) string {
	if key := stringFromMap(options, "getRowKey", ""); key != "" {
		return key
	}
	for _, field := range fields {
		if field["role"] == "key" {
			return stringFromMap(field, "name", "id")
		}
	}
	return "id"
}

func collectionTable(rows []any, fields []map[string]any, options map[string]any, verb, keyField, selParam, selValue string) map[string]any {
	columns := []any{}
	for _, field := range fields {
		role, _ := field["role"].(string)
		if role == "prose" || role == "media" {
			continue
		}
		name := stringFromMap(field, "name", "")
		cell := map[string]any{"field": name}
		switch role {
		case "key":
			cell["kind"] = "caption"
			cell["tone"] = "muted"
		case "count", "size", "measure":
			cell["kind"] = "number"
		case "status":
			cell["kind"] = "status"
		default:
			cell["kind"] = "field"
		}
		column := map[string]any{"id": name, "header": fieldLabel(field), "cell": cell}
		if width := stringFromMap(field, "width", ""); width != "" {
			column["maxWidth"] = width
		}
		columns = append(columns, column)
	}
	columns = append(columns, collectionActionColumns(options, verb)...)

	props := map[string]any{
		"rows":      rows,
		"getRowKey": keyField,
		"columns":   columns,
	}
	if empty := stringFromMap(options, "empty", ""); empty != "" {
		props["emptyMessage"] = empty
	}
	if selValue != "" && selValue != "__new" {
		props["selectedKey"] = selValue
	}
	if selParam != "" {
		props["onRowSelect"] = map[string]any{
			"kind": "navigate",
			"to":   fmt.Sprintf("?%s=${row.%s}", selParam, keyField),
		}
	} else if open, ok := normalizeActionSpec(options["open"], nil, nil); ok {
		props["onRowSelect"] = open
	}
	return componentNode("DataTable", props)
}

func collectionActionColumns(options map[string]any, verb string) []any {
	columns := []any{}
	if open, ok := normalizeActionSpec(options["open"], nil, nil); ok {
		columns = append(columns, map[string]any{
			"id": "open", "header": "Open", "maxWidth": "8ch",
			"cell": map[string]any{"kind": "actionButton", "label": "Open", "action": open},
		})
	}
	if verb != "edit" && verb != "manage" {
		return columns
	}
	if reorder, ok := normalizeActionSpec(options["reorder"], nil, nil); ok {
		columns = append(columns,
			map[string]any{"id": "moveUp", "header": "", "maxWidth": "5ch",
				"cell": map[string]any{"kind": "actionButton", "label": "↑", "action": actionWithPayload(reorder, "direction", "up")}},
			map[string]any{"id": "moveDown", "header": "", "maxWidth": "5ch",
				"cell": map[string]any{"kind": "actionButton", "label": "↓", "action": actionWithPayload(reorder, "direction", "down")}},
		)
	}
	if remove, ok := normalizeActionSpec(options["remove"], nil, nil); ok {
		columns = append(columns, map[string]any{
			"id": "delete", "header": "Delete", "maxWidth": "9ch",
			"cell": map[string]any{"kind": "actionButton", "label": "Delete", "action": remove},
		})
	}
	return columns
}

func actionWithPayload(spec map[string]any, key string, value any) map[string]any {
	out := map[string]any{}
	for k, v := range spec {
		out[k] = v
	}
	payload := map[string]any{}
	if existing, ok := out["payload"].(map[string]any); ok {
		for k, v := range existing {
			payload[k] = v
		}
	}
	payload[key] = value
	out["payload"] = payload
	return out
}

func (r *runtime) collectionDetail(rows []any, fields []map[string]any, options map[string]any, verb, keyField, selParam, selValue string) any {
	if selValue == "" {
		return componentNode("Caption", map[string]any{"tone": "muted"},
			map[string]any{"kind": "text", "text": "Select a row to open it here."})
	}
	values := map[string]any{}
	title := "New item"
	if selValue != "__new" {
		found := false
		for _, entry := range rows {
			row, ok := entry.(map[string]any)
			if !ok {
				continue
			}
			if anyToString(row[keyField]) == selValue {
				values = row
				found = true
				break
			}
		}
		if !found {
			return componentNode("Caption", map[string]any{"tone": "muted"},
				map[string]any{"kind": "text", "text": fmt.Sprintf("No row matches %q.", selValue)})
		}
		title = "Edit"
		for _, field := range fields {
			if field["role"] == "primary" {
				if label := anyToString(values[stringFromMap(field, "name", "")]); label != "" {
					title = fmt.Sprintf("Edit: %s", label)
				}
				break
			}
		}
	}

	detailVerb := "show"
	if verb == "edit" || verb == "manage" {
		detailVerb = "edit"
	}
	detailOptions := map[string]any{
		"schema": options["schema"],
		"verb":   detailVerb,
		"title":  stringFromMap(options, "detailTitle", title),
	}
	copyIfPresent(detailOptions, options, "submit")
	copyIfPresent(detailOptions, options, "status")
	copyIfPresent(detailOptions, options, "statusMessage")
	detail := r.recordVerb(goja.FunctionCall{Arguments: []goja.Value{
		r.vm.ToValue(values), r.vm.ToValue(detailOptions),
	}})

	children := []any{detail.Export()}
	if selParam != "" {
		children = append(children, componentNode("Inline", map[string]any{"gap": "sm"},
			componentNode("Button", map[string]any{
				"action": map[string]any{"kind": "navigate", "to": "?" + selParam + "="},
			}, map[string]any{"kind": "text", "text": "Close"}),
		))
	}
	return componentNode("Stack", map[string]any{"gap": "sm"}, children...)
}

func collectionCreateButton(options map[string]any, verb, selParam string) any {
	create := options["create"]
	if create == nil || create == false || selParam == "" {
		return nil
	}
	if verb != "edit" && verb != "manage" {
		return nil
	}
	label := "New item"
	if opts, ok := create.(map[string]any); ok {
		label = stringFromMap(opts, "label", label)
	}
	return componentNode("Inline", map[string]any{"gap": "sm", "justify": "end"},
		componentNode("Button", map[string]any{
			"action": map[string]any{"kind": "navigate", "to": "?" + selParam + "=__new"},
		}, map[string]any{"kind": "text", "text": label}),
	)
}

// ─── section (ui.dsl) ───────────────────────────────────────────────

// ui.section(title, options?, ...children) — flat document sectioning:
// uppercase label + 1px rule, no box. options: level 1|2|3, anchor, caption,
// actions (widget node), rule (default true), density (default "flush"),
// divider.
func (r *runtime) sectionVerb(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) == 0 {
		panic(r.vm.NewGoError(fmt.Errorf("ui.dsl section(title, options?, ...children) requires a title")))
	}
	title := call.Arguments[0]
	var label any
	if _, ok := title.Export().(string); ok {
		label = title.String()
	} else {
		label = r.exportRenderable(title)
	}
	props := map[string]any{
		"label":   label,
		"rule":    true,
		"density": "flush",
		"level":   1,
	}
	rest := call.Arguments[1:]
	if len(rest) > 0 && isPlainObject(rest[0]) && !looksLikeWidgetNodeExport(rest[0]) {
		options := exportObject(rest[0])
		rest = rest[1:]
		copyIfPresent(props, options, "level")
		copyIfPresent(props, options, "caption")
		copyIfPresent(props, options, "actions")
		copyIfPresent(props, options, "divider")
		if value, ok := options["rule"]; ok {
			props["rule"] = value
		}
		if value := stringFromMap(options, "density", ""); value != "" {
			props["density"] = value
		}
		if anchor := stringFromMap(options, "anchor", ""); anchor != "" {
			props["anchorId"] = anchor
		}
	}
	children := r.exportChildren(rest)
	return r.vm.ToValue(componentNode("SectionBlock", props, children...))
}
