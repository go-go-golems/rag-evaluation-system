package spec

import "testing"

func TestCollectionSpecLowersSimplestTable(t *testing.T) {
	collection := CollectionSpec{
		Name:        "sessions",
		Rows:        []JSONObject{{"sessionId": "s1", "title": "Intro", "turnCount": 12, "status": "ready"}},
		Schema:      sessionSchema(),
		Mode:        CollectionModeShow,
		Arrangement: ArrangementSpec{Kind: ArrangementKindTable},
		Empty:       "No sessions.",
	}
	if issues := collection.Validate("collection"); HasErrors(issues) {
		t.Fatalf("Validate() unexpected errors: %#v", issues)
	}

	node := collection.ToNode().ToWidgetNode()
	if got := node["type"]; got != "Stack" {
		t.Fatalf("root type = %v, want Stack", got)
	}
	table := child(node, 0)
	if got := table["type"]; got != "DataTable" {
		t.Fatalf("table type = %v, want DataTable", got)
	}
	props := table["props"].(JSONObject)
	if got := props["getRowKey"]; got != "sessionId" {
		t.Fatalf("getRowKey = %v, want sessionId", got)
	}
	if _, ok := props["onRowSelect"]; ok {
		t.Fatalf("onRowSelect present for simplest table: %#v", props["onRowSelect"])
	}
	columns := props["columns"].([]JSONValue)
	if len(columns) != 4 {
		t.Fatalf("columns len = %d, want 4", len(columns))
	}
}

func TestCollectionSpecLowersSelectableTable(t *testing.T) {
	collection := CollectionSpec{
		Name:        "sessions",
		Rows:        []JSONObject{{"sessionId": "s1", "title": "Intro"}, {"sessionId": "s2", "title": "Debugging"}},
		Schema:      sessionSchema(),
		Mode:        CollectionModeShow,
		Selection:   &SelectionSpec{Kind: SelectionKindURLParam, Param: "selected", Value: "s2"},
		Arrangement: ArrangementSpec{Kind: ArrangementKindTable},
	}
	if issues := collection.Validate("collection"); HasErrors(issues) {
		t.Fatalf("Validate() unexpected errors: %#v", issues)
	}

	table := child(collection.ToNode().ToWidgetNode(), 0)
	props := table["props"].(JSONObject)
	if got := props["selectedKey"]; got != "s2" {
		t.Fatalf("selectedKey = %v, want s2", got)
	}
	action := props["onRowSelect"].(JSONObject)
	if got := action["kind"]; got != "navigate" {
		t.Fatalf("onRowSelect.kind = %v, want navigate", got)
	}
	if got := action["to"]; got != "?selected=${row.sessionId}" {
		t.Fatalf("onRowSelect.to = %v", got)
	}
}

func TestCollectionSpecLowersMasterDetailEditor(t *testing.T) {
	collection := CollectionSpec{
		Name:   "agenda",
		Rows:   []JSONObject{{"id": "agenda-intro", "number": "14h30", "title": "Intro", "description": "Welcome"}},
		Schema: agendaSchema(),
		Mode:   CollectionModeEdit,
		Selection: &SelectionSpec{
			Kind:  SelectionKindURLParam,
			Param: "agenda",
			Value: "agenda-intro",
		},
		Arrangement: ArrangementSpec{Kind: ArrangementKindMasterDetail},
		Actions: CollectionActions{
			Create: &CreateActionSpec{Label: "New agenda item"},
			Submit: &SubmitSpec{FormAction: "/settings/agenda-item", Method: "post"},
			Remove: &ActionSpec{Kind: ActionKindServer, Name: "admin-delete-agenda-item", Confirm: &TemplateSpec{Parts: []TemplateValue{
				{Kind: TemplateValueText, Text: "Delete “"},
				{Kind: TemplateValuePath, Path: "row.title"},
				{Kind: TemplateValueText, Text: "”?"},
			}}},
		},
	}
	if issues := collection.Validate("collection"); HasErrors(issues) {
		t.Fatalf("Validate() unexpected errors: %#v", issues)
	}

	root := collection.ToNode().ToWidgetNode()
	children := root["children"].([]JSONValue)
	if len(children) != 3 {
		t.Fatalf("children len = %d, want create/table/detail", len(children))
	}
	detail := children[2].(JSONObject)
	form := child(detail, 0)
	if got := form["type"]; got != "FormPanel" {
		t.Fatalf("detail child type = %v, want FormPanel", got)
	}
	props := form["props"].(JSONObject)
	if got := props["title"]; got != "Edit: Intro" {
		t.Fatalf("form title = %v, want Edit: Intro", got)
	}
	if got := props["formAction"]; got != "/settings/agenda-item" {
		t.Fatalf("formAction = %v", got)
	}
}

func TestCollectionSpecLowersNewMasterDetailWithEditableKey(t *testing.T) {
	collection := CollectionSpec{
		Name:        "agenda",
		Rows:        []JSONObject{{"id": "agenda-intro", "title": "Intro"}},
		Schema:      agendaSchema(),
		Mode:        CollectionModeEdit,
		Selection:   &SelectionSpec{Kind: SelectionKindURLParam, Param: "agenda", Value: "__new"},
		Arrangement: ArrangementSpec{Kind: ArrangementKindMasterDetail},
		Actions:     CollectionActions{Submit: &SubmitSpec{FormAction: "/settings/agenda-item", Method: "post"}},
	}
	if issues := collection.Validate("collection"); HasErrors(issues) {
		t.Fatalf("Validate() unexpected errors: %#v", issues)
	}

	root := collection.ToNode().ToWidgetNode()
	detail := child(root, 1)
	form := child(detail, 0)
	fieldGrid := child(form, 0)
	idRow := child(fieldGrid, 0)
	control := idRow["props"].(JSONObject)["control"].(JSONObject)
	controlProps := control["props"].(JSONObject)
	if controlProps["readOnly"] != false {
		t.Fatalf("new item key readOnly = %#v, want false", controlProps["readOnly"])
	}
}

func TestCollectionSpecLowersShapingKeyboardCommandsAndSemanticStyles(t *testing.T) {
	navigate := ActionSpec{Kind: ActionKindNavigate, To: "/jobs", Options: JSONObject{"query": JSONObject{"page": JSONObject{"kind": "accessor", "mode": "context", "path": "page"}}, "preserveQuery": []string{"q"}, "omitEmpty": true}}
	collection := CollectionSpec{
		Name: "jobs", Rows: []JSONObject{{"sessionId": "j1", "title": "Go", "status": "shortlisted"}}, Schema: sessionSchema(), Arrangement: ArrangementSpec{Kind: ArrangementKindTable},
		Shaping: CollectionShapingSpec{Search: &SearchSpec{Name: "q", Value: "go", ResultCount: 47, Submit: &navigate}, Pagination: &PaginationSpec{Page: 2, PageSize: 20, TotalItems: 47, Sizes: []int{20, 50, 100}, OnChange: &navigate}},
		Table:   TableSpec{Keyboard: TableKeyboardSpec{Enabled: true, Mode: "rows", Selection: "manual", VimAliases: true, EnterSelect: true}, Commands: []RowCommandSpec{{ID: "star", Key: "s", Label: "Toggle star", Action: ActionSpec{Kind: ActionKindServer, Name: "job.star"}}}, StyleRules: []SemanticStyleRule{{Field: "status", Equals: "shortlisted", Tone: "success"}}},
	}
	root := collection.ToNode().ToWidgetNode()
	children := root["children"].([]JSONValue)
	if len(children) != 3 {
		t.Fatalf("children len = %d, want search/table/pagination", len(children))
	}
	if children[0].(JSONObject)["type"] != "SearchField" || children[2].(JSONObject)["type"] != "Pagination" {
		t.Fatalf("unexpected shaping order: %#v", children)
	}
	tableProps := children[1].(JSONObject)["props"].(JSONObject)
	if tableProps["keyboard"].(JSONObject)["vimAliases"] != true {
		t.Fatalf("keyboard not lowered: %#v", tableProps["keyboard"])
	}
	if len(tableProps["commands"].([]JSONValue)) != 1 || len(tableProps["styleRules"].([]JSONValue)) != 1 {
		t.Fatalf("commands/styles not lowered: %#v", tableProps)
	}
	pagerAction := children[2].(JSONObject)["props"].(JSONObject)["onPageChangeAction"].(JSONObject)
	if pagerAction["omitEmpty"] != true || pagerAction["query"] == nil {
		t.Fatalf("structured navigation options lost: %#v", pagerAction)
	}
}

func TestCollectionSpecLowersSortableColumns(t *testing.T) {
	collection := CollectionSpec{
		Name:        "jobs",
		Rows:        []JSONObject{{"sessionId": "j1", "title": "Go"}},
		Schema:      sessionSchema(),
		Arrangement: ArrangementSpec{Kind: ArrangementKindTable},
		Table: TableSpec{SortColumns: []TableSortColumnSpec{{
			Field:     "title",
			Direction: "ascending",
			Action:    ActionSpec{Kind: ActionKindNavigate, To: "/jobs?sort=title-desc"},
		}}},
	}
	tableProps := collection.ToNode().ToWidgetNode()["children"].([]JSONValue)[0].(JSONObject)["props"].(JSONObject)
	for _, value := range tableProps["columns"].([]JSONValue) {
		column := value.(JSONObject)
		if column["id"] != "title" {
			continue
		}
		if column["sortable"] != true || column["sortDirection"] != "ascending" || column["onSort"] == nil {
			t.Fatalf("sortable title column = %#v", column)
		}
		return
	}
	t.Fatalf("title column was not lowered: %#v", tableProps["columns"])
}

func TestCollectionSpecValidationRejectsBadArrangement(t *testing.T) {
	collection := CollectionSpec{
		Name:        "sessions",
		Schema:      sessionSchema(),
		Mode:        CollectionModeShow,
		Arrangement: ArrangementSpec{Kind: ArrangementKind("mast-detail")},
	}
	issues := collection.Validate("collection")
	if !HasErrors(issues) {
		t.Fatalf("Validate() had no errors: %#v", issues)
	}
	if issues[0].Code != "collection.arrangement.invalid" {
		t.Fatalf("first issue code = %s", issues[0].Code)
	}
}

func child(node JSONObject, index int) JSONObject {
	return node["children"].([]JSONValue)[index].(JSONObject)
}

func sessionSchema() SchemaSpec {
	return SchemaSpec{Name: "Session", Fields: []FieldSpec{
		{Name: "sessionId", Label: "ID", Kind: FieldKindString, Semantic: FieldSemanticKey},
		{Name: "title", Label: "Title", Kind: FieldKindString, Semantic: FieldSemanticPrimary},
		{Name: "turnCount", Label: "Turns", Kind: FieldKindNumber, Semantic: FieldSemanticCount},
		{Name: "status", Label: "Status", Kind: FieldKindString, Semantic: FieldSemanticStatus},
		{Name: "body", Label: "Body", Kind: FieldKindString, Semantic: FieldSemanticProse},
	}}
}

func agendaSchema() SchemaSpec {
	return SchemaSpec{Name: "AgendaItem", Fields: []FieldSpec{
		{Name: "id", Label: "ID", Kind: FieldKindString, Semantic: FieldSemanticKey},
		{Name: "number", Label: "Time", Kind: FieldKindString, Semantic: FieldSemanticShort},
		{Name: "title", Label: "Title", Kind: FieldKindString, Semantic: FieldSemanticPrimary, Validation: FieldValidation{Required: true, MaxLength: 160}},
		{Name: "description", Label: "Description", Kind: FieldKindString, Semantic: FieldSemanticProse, Editor: EditorSpec{Control: EditorControlTextarea, Rows: 4}},
	}}
}
