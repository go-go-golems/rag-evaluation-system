package spec

import "fmt"

// ToWidgetPage lowers a validated PageSpec into the JSON-like WidgetPage shape
// consumed by the existing React app.
func (s PageSpec) ToWidgetPage() JSONObject {
	page := JSONObject{
		"schemaVersion": valueOrDefault(s.SchemaVersion, "0.2.0"),
		"id":            s.ID,
		"title":         s.Title,
	}
	if s.Meta != nil {
		page["meta"] = s.Meta
	}
	if s.Root.Kind != "" {
		page["root"] = s.Root.ToWidgetNode()
	}
	if len(s.Diagnostics) > 0 {
		page["diagnostics"] = validationIssuesToJSON(s.Diagnostics)
	}
	return page
}

// ToWidgetNode lowers a typed node into the current Widget IR node shape.
func (s NodeSpec) ToWidgetNode() JSONObject {
	switch s.Kind {
	case NodeKindText:
		return JSONObject{"kind": "text", "text": s.Text}
	case NodeKindElement:
		node := JSONObject{"kind": "element", "tag": s.Tag}
		if s.Props != nil {
			node["attrs"] = s.Props
		}
		if len(s.Children) > 0 {
			node["children"] = lowerNodes(s.Children)
		}
		return node
	case NodeKindComponent:
		node := JSONObject{"kind": "component", "type": s.Type}
		if s.Props != nil {
			node["props"] = s.Props
		}
		if len(s.Children) > 0 {
			node["children"] = lowerNodes(s.Children)
		}
		if s.Source != nil {
			node["source"] = JSONObject{"module": s.Source.Module, "helper": s.Source.Helper, "path": s.Source.Path}
		}
		if len(s.Diagnostics) > 0 {
			node["diagnostics"] = validationIssuesToJSON(s.Diagnostics)
		}
		return node
	default:
		return JSONObject{"kind": "component", "type": "UnknownWidget", "props": JSONObject{"message": fmt.Sprintf("cannot lower node kind %q", s.Kind)}}
	}
}

// ToNode lowers section intent to the current SectionBlock component.
func (s SectionSpec) ToNode() NodeSpec {
	props := JSONObject{"label": s.Title, "level": s.Level, "rule": true, "density": "flush"}
	if s.AnchorID != "" {
		props["anchorId"] = s.AnchorID
	}
	if s.Caption != "" {
		props["caption"] = s.Caption
	}
	return NodeSpec{Kind: NodeKindComponent, Type: "SectionBlock", Props: props, Children: s.Children}
}

// ToNode lowers a collection intent into a composed Widget IR subtree. The
// output intentionally targets the current runtime components while keeping the
// input side typed.
func (s CollectionSpec) ToNode() NodeSpec {
	keyField := s.keyField()
	children := []NodeSpec{}
	if create := s.createButtonNode(); create.Kind != "" {
		children = append(children, create)
	}
	children = append(children, s.tableNode(keyField))
	if s.Arrangement.Kind == ArrangementKindMasterDetail {
		children = append(children, s.detailNode(keyField))
	}
	return NodeSpec{Kind: NodeKindComponent, Type: "Stack", Props: JSONObject{"gap": "md"}, Children: children}
}

func (s CollectionSpec) tableNode(keyField string) NodeSpec {
	props := JSONObject{
		"rows":      rowsToJSON(s.Rows),
		"getRowKey": keyField,
		"columns":   s.tableColumns(),
	}
	if s.Empty != "" {
		props["emptyMessage"] = s.Empty
	}
	if s.Selection != nil {
		if s.Selection.Value != "" && s.Selection.Value != "__new" {
			props["selectedKey"] = s.Selection.Value
		}
		props["onRowSelect"] = ActionSpec{Kind: ActionKindNavigate, To: fmt.Sprintf("?%s=${row.%s}", s.Selection.Param, keyField)}.ToWidgetAction()
	} else if s.Actions.Open != nil {
		props["onRowSelect"] = s.Actions.Open.ToWidgetAction()
	}
	return NodeSpec{Kind: NodeKindComponent, Type: "DataTable", Props: props}
}

func (s CollectionSpec) tableColumns() []JSONValue {
	columns := []JSONValue{}
	for _, field := range s.Schema.Fields {
		if field.Summary.Elide || field.Semantic == FieldSemanticProse || field.Kind == FieldKindMedia {
			continue
		}
		cell := JSONObject{"field": field.Name, "kind": cellKind(field)}
		if field.Semantic == FieldSemanticKey {
			cell["tone"] = "muted"
		}
		column := JSONObject{"id": field.Name, "header": fieldLabel(field), "cell": cell}
		if field.Layout.Width != "" {
			column["maxWidth"] = field.Layout.Width
		}
		columns = append(columns, column)
	}
	columns = append(columns, s.actionColumns()...)
	return columns
}

func (s CollectionSpec) actionColumns() []JSONValue {
	columns := []JSONValue{}
	if s.Actions.Open != nil {
		columns = append(columns, JSONObject{
			"id": "open", "header": "Open", "maxWidth": "8ch",
			"cell": JSONObject{"kind": "actionButton", "label": "Open", "action": s.Actions.Open.ToWidgetAction()},
		})
	}
	if s.Mode != CollectionModeEdit && s.Mode != CollectionModeManage {
		return columns
	}
	if s.Actions.Reorder != nil {
		columns = append(columns,
			JSONObject{"id": "moveUp", "header": "", "maxWidth": "5ch", "cell": JSONObject{"kind": "actionButton", "label": "↑", "action": actionWithLiteralPayload(*s.Actions.Reorder, "direction", "up").ToWidgetAction()}},
			JSONObject{"id": "moveDown", "header": "", "maxWidth": "5ch", "cell": JSONObject{"kind": "actionButton", "label": "↓", "action": actionWithLiteralPayload(*s.Actions.Reorder, "direction", "down").ToWidgetAction()}},
		)
	}
	if s.Actions.Remove != nil {
		columns = append(columns, JSONObject{
			"id": "delete", "header": "Delete", "maxWidth": "9ch",
			"cell": JSONObject{"kind": "actionButton", "label": "Delete", "action": s.Actions.Remove.ToWidgetAction()},
		})
	}
	return columns
}

func (s CollectionSpec) createButtonNode() NodeSpec {
	if s.Actions.Create == nil || s.Selection == nil || (s.Mode != CollectionModeEdit && s.Mode != CollectionModeManage) {
		return NodeSpec{}
	}
	label := valueOrDefault(s.Actions.Create.Label, "New item")
	button := NodeSpec{Kind: NodeKindComponent, Type: "Button", Props: JSONObject{"action": ActionSpec{Kind: ActionKindNavigate, To: "?" + s.Selection.Param + "=__new"}.ToWidgetAction()}, Children: []NodeSpec{{Kind: NodeKindText, Text: label}}}
	return NodeSpec{Kind: NodeKindComponent, Type: "Inline", Props: JSONObject{"gap": "sm", "justify": "end"}, Children: []NodeSpec{button}}
}

func (s CollectionSpec) detailNode(keyField string) NodeSpec {
	if s.Selection == nil || s.Selection.Value == "" {
		return captionNode("Select a row to open it here.")
	}
	values := JSONObject{}
	title := "New item"
	if s.Selection.Value != "__new" {
		var found bool
		for _, row := range s.Rows {
			if stringifyJSON(row[keyField]) == s.Selection.Value {
				values = row
				found = true
				break
			}
		}
		if !found {
			return captionNode(fmt.Sprintf("No row matches %q.", s.Selection.Value))
		}
		title = "Edit"
		if primary := s.primaryFieldName(); primary != "" {
			if label := stringifyJSON(values[primary]); label != "" {
				title = "Edit: " + label
			}
		}
	}
	form := s.recordFormNode(values, title)
	children := []NodeSpec{form}
	children = append(children, NodeSpec{Kind: NodeKindComponent, Type: "Inline", Props: JSONObject{"gap": "sm"}, Children: []NodeSpec{{Kind: NodeKindComponent, Type: "Button", Props: JSONObject{"action": ActionSpec{Kind: ActionKindNavigate, To: "?" + s.Selection.Param + "="}.ToWidgetAction()}, Children: []NodeSpec{{Kind: NodeKindText, Text: "Close"}}}}})
	return NodeSpec{Kind: NodeKindComponent, Type: "Stack", Props: JSONObject{"gap": "sm"}, Children: children}
}

func (s CollectionSpec) recordFormNode(values JSONObject, title string) NodeSpec {
	props := JSONObject{"title": title}
	if s.Actions.Submit != nil {
		props["formAction"] = s.Actions.Submit.FormAction
		props["method"] = valueOrDefault(s.Actions.Submit.Method, "post")
	}
	return NodeSpec{Kind: NodeKindComponent, Type: "FormPanel", Props: props, Children: s.recordRows(values)}
}

func (s CollectionSpec) recordRows(values JSONObject) []NodeSpec {
	rows := []NodeSpec{}
	grid := []NodeSpec{}
	flush := func() {
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
			rows = append(rows, NodeSpec{Kind: NodeKindComponent, Type: "FieldGrid", Props: JSONObject{"columns": columns}, Children: grid})
		}
		grid = []NodeSpec{}
	}
	for _, field := range s.Schema.Fields {
		row := fieldRowNode(field, values)
		if gridableField(field) {
			grid = append(grid, row)
			continue
		}
		flush()
		rows = append(rows, row)
	}
	flush()
	return rows
}

// ToWidgetAction lowers a typed action to the current ActionSpec JSON shape.
func (s ActionSpec) ToWidgetAction() JSONObject {
	action := JSONObject{"kind": string(s.Kind)}
	if s.Name != "" {
		action["name"] = s.Name
	}
	if s.To != "" {
		action["to"] = s.To
	}
	if s.Event != "" {
		action["event"] = s.Event
	}
	if len(s.Payload.Fields) > 0 {
		action["payload"] = s.Payload.ToJSON()
	}
	if s.Confirm != nil {
		action["confirm"] = s.Confirm.LegacyTemplateString()
	}
	return action
}

// ToJSON lowers payload template fields. Literal values lower to direct JSON;
// path/text descriptors are kept as typed data for Action IR v2 consumers.
func (s PayloadTemplate) ToJSON() JSONObject {
	out := JSONObject{}
	for _, field := range s.Fields {
		switch field.Value.Kind {
		case TemplateValueLiteral:
			out[field.Name] = field.Value.Value
		case TemplateValuePath:
			out[field.Name] = JSONObject{"kind": "path", "path": field.Value.Path}
		case TemplateValueText:
			out[field.Name] = field.Value.Text
		}
	}
	return out
}

// LegacyTemplateString lowers a typed text/path template to the string format
// used by the current frontend interpolation helper.
func (s TemplateSpec) LegacyTemplateString() string {
	out := ""
	for _, part := range s.Parts {
		switch part.Kind {
		case TemplateValueText:
			out += part.Text
		case TemplateValuePath:
			out += "${" + part.Path + "}"
		case TemplateValueLiteral:
			out += stringifyJSON(part.Value)
		}
	}
	return out
}

func fieldRowNode(field FieldSpec, values JSONObject) NodeSpec {
	controlType := "TextInput"
	orientation := "inline"
	controlProps := JSONObject{
		"name":         field.Name,
		"defaultValue": stringifyJSON(values[field.Name]),
		"readOnly":     field.Editor.ReadOnly || field.Semantic == FieldSemanticKey,
	}
	if field.Editor.Placeholder != "" {
		controlProps["placeholder"] = field.Editor.Placeholder
	}
	if field.Validation.Required {
		controlProps["required"] = true
	}
	if field.Validation.MaxLength > 0 {
		controlProps["maxLength"] = field.Validation.MaxLength
	}
	if field.Editor.Control == EditorControlTextarea || field.Semantic == FieldSemanticProse {
		controlType = "TextareaInput"
		orientation = "stacked"
		controlProps["rows"] = intOrDefault(field.Editor.Rows, 4)
	}
	return NodeSpec{Kind: NodeKindComponent, Type: "FormRow", Props: JSONObject{
		"label":       fieldLabel(field),
		"control":     NodeSpec{Kind: NodeKindComponent, Type: controlType, Props: controlProps}.ToWidgetNode(),
		"orientation": orientation,
		"required":    field.Validation.Required,
	}}
}

func captionNode(text string) NodeSpec {
	return NodeSpec{Kind: NodeKindComponent, Type: "Caption", Props: JSONObject{"tone": "muted"}, Children: []NodeSpec{{Kind: NodeKindText, Text: text}}}
}

func lowerNodes(nodes []NodeSpec) []JSONValue {
	out := make([]JSONValue, 0, len(nodes))
	for _, node := range nodes {
		out = append(out, node.ToWidgetNode())
	}
	return out
}

func validationIssuesToJSON(issues []ValidationIssue) []JSONValue {
	out := make([]JSONValue, 0, len(issues))
	for _, issue := range issues {
		out = append(out, JSONObject{"severity": issue.Severity, "code": issue.Code, "path": issue.Path, "message": issue.Message, "hint": issue.Hint})
	}
	return out
}

func rowsToJSON(rows []JSONObject) []JSONValue {
	out := make([]JSONValue, 0, len(rows))
	for _, row := range rows {
		out = append(out, row)
	}
	return out
}

func (s CollectionSpec) keyField() string {
	for _, field := range s.Schema.Fields {
		if field.Semantic == FieldSemanticKey {
			return field.Name
		}
	}
	return "id"
}

func (s CollectionSpec) primaryFieldName() string {
	for _, field := range s.Schema.Fields {
		if field.Semantic == FieldSemanticPrimary {
			return field.Name
		}
	}
	return ""
}

func fieldLabel(field FieldSpec) string {
	if field.Label != "" {
		return field.Label
	}
	return field.Name
}

func cellKind(field FieldSpec) string {
	switch field.Semantic {
	case FieldSemanticKey:
		return "caption"
	case FieldSemanticCount, FieldSemanticSize, FieldSemanticMeasure:
		return "number"
	case FieldSemanticStatus:
		return "status"
	default:
		return "field"
	}
}

func gridableField(field FieldSpec) bool {
	switch field.Semantic {
	case FieldSemanticShort, FieldSemanticKey, FieldSemanticCount, FieldSemanticSize, FieldSemanticMeasure:
		return true
	}
	return field.Editor.Control != EditorControlTextarea && field.Semantic != FieldSemanticProse
}

func actionWithLiteralPayload(action ActionSpec, key string, value JSONValue) ActionSpec {
	fields := append([]PayloadFieldSpec(nil), action.Payload.Fields...)
	fields = append(fields, PayloadFieldSpec{Name: key, Value: TemplateValue{Kind: TemplateValueLiteral, Value: value}})
	action.Payload.Fields = fields
	return action
}

func valueOrDefault(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}

func intOrDefault(value, fallback int) int {
	if value == 0 {
		return fallback
	}
	return value
}

func stringifyJSON(value JSONValue) string {
	if value == nil {
		return ""
	}
	return fmt.Sprint(value)
}
