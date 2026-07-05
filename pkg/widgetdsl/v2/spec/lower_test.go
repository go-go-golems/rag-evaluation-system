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
