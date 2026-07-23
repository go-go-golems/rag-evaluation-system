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
	if s.Shell != nil {
		page["shell"] = s.Shell.ToJSON()
	}
	if len(s.Shortcuts) > 0 {
		page["shortcuts"] = JSONObject{"bindings": pageShortcutsToJSON(s.Shortcuts)}
	}
	if s.Root.Kind != "" {
		page["root"] = s.Root.ToWidgetNode()
	}
	if len(s.Diagnostics) > 0 {
		page["diagnostics"] = validationIssuesToJSON(s.Diagnostics)
	}
	return page
}

// ToJSON lowers a typed shell into the browser transport shape.
func (s PageShellSpec) ToJSON() JSONObject {
	out := JSONObject{"kind": string(s.Kind)}
	if s.Navigation != nil {
		out["navigation"] = s.Navigation.ToJSON()
	}
	content := s.Content.ToJSON()
	if len(content) > 0 {
		out["content"] = content
	}
	return out
}

func (s NavigationSpec) ToJSON() JSONObject {
	out := JSONObject{
		"placement": string(s.Placement),
		"sections":  navigationSectionsToJSON(s.Sections),
	}
	if s.Brand != nil {
		out["brand"] = s.Brand
	}
	if s.AriaLabel != "" {
		out["ariaLabel"] = s.AriaLabel
	}
	if s.ActiveItem != "" {
		out["activeItemId"] = s.ActiveItem
	}
	if s.SidebarWidth != 0 {
		out["sidebarWidth"] = s.SidebarWidth
	}
	if s.NarrowMode != "" {
		out["narrowMode"] = s.NarrowMode
	}
	return out
}

func navigationSectionsToJSON(sections []NavigationSectionSpec) []any {
	out := make([]any, 0, len(sections))
	for _, section := range sections {
		items := make([]any, 0, len(section.Items))
		for _, item := range section.Items {
			entry := JSONObject{"id": item.ID, "label": item.Label}
			if item.Icon != nil {
				entry["icon"] = item.Icon
			}
			if item.Badge != nil {
				entry["badge"] = item.Badge
			}
			if item.Disabled {
				entry["disabled"] = true
			}
			if item.Action != nil {
				entry["action"] = item.Action
			}
			items = append(items, entry)
		}
		out = append(out, JSONObject{"id": section.ID, "label": section.Label, "items": items})
	}
	return out
}

func (s ContentViewportSpec) ToJSON() JSONObject {
	out := JSONObject{}
	if s.MaxWidth != "" {
		out["maxWidth"] = s.MaxWidth
	}
	if s.Padding != "" {
		out["padding"] = s.Padding
	}
	if s.Scroll != "" {
		out["scroll"] = s.Scroll
	}
	return out
}

func pageShortcutsToJSON(shortcuts []PageShortcutSpec) []JSONValue {
	out := make([]JSONValue, 0, len(shortcuts))
	for _, shortcut := range shortcuts {
		modifiers := make([]JSONValue, 0, len(shortcut.Modifiers))
		for _, modifier := range shortcut.Modifiers {
			modifiers = append(modifiers, string(modifier))
		}
		out = append(out, JSONObject{
			"id":             shortcut.ID,
			"key":            shortcut.Key,
			"modifiers":      modifiers,
			"label":          shortcut.Label,
			"action":         shortcut.Action.ToWidgetAction(),
			"preventDefault": shortcut.PreventDefault,
			"allowRepeat":    shortcut.AllowRepeat,
		})
	}
	return out
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
	if search := s.searchNode(); search.Kind != "" {
		children = append(children, search)
	}
	if create := s.createButtonNode(); create.Kind != "" {
		children = append(children, create)
	}
	children = append(children, s.tableNode(keyField))
	if pager := s.paginationNode(); pager.Kind != "" {
		children = append(children, pager)
	}
	if s.Arrangement.Kind == ArrangementKindMasterDetail {
		children = append(children, s.detailNode(keyField))
	}
	return NodeSpec{Kind: NodeKindComponent, Type: "Stack", Props: JSONObject{"gap": "md"}, Children: children}
}

func (s CollectionSpec) searchNode() NodeSpec {
	search := s.Shaping.Search
	if search == nil {
		return NodeSpec{}
	}
	props := JSONObject{"name": valueOrDefault(search.Name, "q"), "defaultValue": search.Value, "placeholder": search.Placeholder}
	if search.Submit != nil {
		props["onSubmitAction"] = search.Submit.ToWidgetAction()
	}
	if search.Clear != nil {
		props["onClearAction"] = search.Clear.ToWidgetAction()
	}
	if search.ResultCount >= 0 {
		props["resultCount"] = search.ResultCount
	}
	return NodeSpec{Kind: NodeKindComponent, Type: "SearchField", Props: props}
}

func (s CollectionSpec) paginationNode() NodeSpec {
	pager := s.Shaping.Pagination
	if pager == nil {
		return NodeSpec{}
	}
	pageSize := pager.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	pageCount := 1
	if pager.TotalItems > 0 {
		pageCount = (pager.TotalItems + pageSize - 1) / pageSize
	}
	props := JSONObject{"page": intOrDefault(pager.Page, 1), "pageCount": pageCount, "pageSize": pageSize, "totalItems": pager.TotalItems, "pageSizes": pager.Sizes}
	if pager.OnChange != nil {
		props["onPageChangeAction"] = pager.OnChange.ToWidgetAction()
		props["onPageSizeChangeAction"] = pager.OnChange.ToWidgetAction()
	}
	return NodeSpec{Kind: NodeKindComponent, Type: "Pagination", Props: props}
}

func (s CollectionSpec) tableNode(keyField string) NodeSpec {
	props := JSONObject{
		"rows":      rowsToJSON(s.Rows),
		"getRowKey": keyField,
		"columns":   s.tableColumns(),
	}
	if s.Table.ClassName != "" {
		props["className"] = s.Table.ClassName
	}
	if s.Empty != "" {
		props["emptyMessage"] = s.Empty
	}
	if s.Table.Keyboard.Enabled {
		props["keyboard"] = JSONObject{"mode": valueOrDefault(s.Table.Keyboard.Mode, "rows"), "selection": valueOrDefault(s.Table.Keyboard.Selection, "manual"), "vimAliases": s.Table.Keyboard.VimAliases, "enterSelect": s.Table.Keyboard.EnterSelect}
	}
	if s.Table.MultiSelection != nil {
		selection := JSONObject{"mode": "multi", "selectedKeys": s.Table.MultiSelection.SelectedKeys}
		if s.Table.MultiSelection.OnChange != nil {
			selection["onChange"] = s.Table.MultiSelection.OnChange.ToWidgetAction()
		}
		props["multiSelection"] = selection
		if len(s.Table.MultiSelection.BulkActions) > 0 {
			actions := make([]JSONValue, 0, len(s.Table.MultiSelection.BulkActions))
			for _, action := range s.Table.MultiSelection.BulkActions {
				actions = append(actions, JSONObject{"id": action.ID, "label": action.Label, "danger": action.Danger, "disabled": action.Disabled, "action": action.Action.ToWidgetAction()})
			}
			props["bulkActions"] = actions
		}
	}
	if len(s.Table.Commands) > 0 {
		commands := make([]JSONValue, 0, len(s.Table.Commands))
		for _, command := range s.Table.Commands {
			commands = append(commands, JSONObject{"id": command.ID, "key": command.Key, "label": command.Label, "danger": command.Danger, "action": command.Action.ToWidgetAction()})
		}
		props["commands"] = commands
	}
	if len(s.Table.StyleRules) > 0 {
		rules := make([]JSONValue, 0, len(s.Table.StyleRules))
		for _, rule := range s.Table.StyleRules {
			rules = append(rules, JSONObject{"field": rule.Field, "equals": rule.Equals, "tone": rule.Tone})
		}
		props["styleRules"] = rules
	}
	if s.Selection != nil {
		if s.Selection.Value != "" && s.Selection.Value != "__new" {
			props["selectedKey"] = s.Selection.Value
		}
		if s.Table.RowSelect != nil {
			props["onRowSelect"] = s.Table.RowSelect.ToWidgetAction()
		} else {
			props["onRowSelect"] = ActionSpec{Kind: ActionKindNavigate, To: fmt.Sprintf("?%s=${row.%s}", s.Selection.Param, keyField)}.ToWidgetAction()
		}
	} else if s.Table.RowSelect != nil {
		props["onRowSelect"] = s.Table.RowSelect.ToWidgetAction()
	} else if s.Actions.Open != nil {
		props["onRowSelect"] = s.Actions.Open.ToWidgetAction()
	}
	return NodeSpec{Kind: NodeKindComponent, Type: "DataTable", Props: props}
}

func (s CollectionSpec) tableColumns() []JSONValue {
	columns := []JSONValue{}
	sortColumns := map[string]TableSortColumnSpec{}
	for _, sortColumn := range s.Table.SortColumns {
		sortColumns[sortColumn.Field] = sortColumn
	}
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
		if sortColumn, ok := sortColumns[field.Name]; ok {
			column["sortable"] = true
			column["onSort"] = sortColumn.Action.ToWidgetAction()
			if sortColumn.Direction != "" {
				column["sortDirection"] = sortColumn.Direction
			}
		}
		columns = append(columns, column)
	}
	columns = append(columns, s.actionColumns()...)
	return columns
}

func (s CollectionSpec) actionColumns() []JSONValue {
	columns := []JSONValue{}
	for _, actionColumn := range s.Table.ActionColumns {
		column := JSONObject{
			"id":     actionColumn.ID,
			"header": actionColumn.Header,
			"cell":   JSONObject{"kind": "actionButton", "label": actionColumn.Label, "action": actionColumn.Action.ToWidgetAction()},
		}
		if actionColumn.MaxWidth != "" {
			column["maxWidth"] = actionColumn.MaxWidth
		}
		columns = append(columns, column)
	}
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
	form := s.recordFormNode(values, title, s.Selection.Value == "__new")
	children := []NodeSpec{form}
	children = append(children, NodeSpec{Kind: NodeKindComponent, Type: "Inline", Props: JSONObject{"gap": "sm"}, Children: []NodeSpec{{Kind: NodeKindComponent, Type: "Button", Props: JSONObject{"action": ActionSpec{Kind: ActionKindNavigate, To: "?" + s.Selection.Param + "="}.ToWidgetAction()}, Children: []NodeSpec{{Kind: NodeKindText, Text: "Close"}}}}})
	return NodeSpec{Kind: NodeKindComponent, Type: "Stack", Props: JSONObject{"gap": "sm"}, Children: children}
}

func (s CollectionSpec) recordFormNode(values JSONObject, title string, newItem bool) NodeSpec {
	props := JSONObject{"title": title}
	if s.Actions.Submit != nil {
		props["formAction"] = s.Actions.Submit.FormAction
		props["method"] = valueOrDefault(s.Actions.Submit.Method, "post")
	}
	return NodeSpec{Kind: NodeKindComponent, Type: "FormPanel", Props: props, Children: s.recordRows(values, newItem)}
}

func (s CollectionSpec) recordRows(values JSONObject, newItem bool) []NodeSpec {
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
		row := fieldRowNode(field, values, newItem)
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
	for key, value := range s.Options {
		action[key] = value
	}
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

func fieldRowNode(field FieldSpec, values JSONObject, newItem bool) NodeSpec {
	controlType := "TextInput"
	orientation := "inline"
	controlProps := JSONObject{
		"name":         field.Name,
		"defaultValue": stringifyJSON(values[field.Name]),
		"readOnly":     field.Editor.ReadOnly || (field.Semantic == FieldSemanticKey && !newItem),
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
	case FieldSemanticPrimary, FieldSemanticShort, FieldSemanticProse, FieldSemanticTags:
		return "field"
	default:
		return "field"
	}
}

func gridableField(field FieldSpec) bool {
	switch field.Semantic {
	case FieldSemanticShort, FieldSemanticKey, FieldSemanticCount, FieldSemanticSize, FieldSemanticMeasure:
		return true
	case FieldSemanticProse:
		return false
	case FieldSemanticPrimary, FieldSemanticStatus, FieldSemanticTags:
		return field.Editor.Control != EditorControlTextarea
	default:
		return field.Editor.Control != EditorControlTextarea
	}
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
