package spec

import "testing"

func TestCollectionSpecLowersMultiSelectionAndBulkActions(t *testing.T) {
	collection := CollectionSpec{
		Name:        "jobs",
		Rows:        []JSONObject{{"id": "a", "title": "A"}, {"id": "b", "title": "B"}},
		Schema:      SchemaSpec{Name: "jobs", Fields: []FieldSpec{{Name: "id", Semantic: FieldSemanticKey}, {Name: "title", Semantic: FieldSemanticPrimary}}},
		Mode:        CollectionModeShow,
		Arrangement: ArrangementSpec{Kind: ArrangementKindTable},
		Table: TableSpec{MultiSelection: &MultiSelectionSpec{
			SelectedKeys: []string{"a", "b"},
			OnChange:     &ActionSpec{Kind: ActionKindServer, Name: "selection-changed"},
			BulkActions:  []BulkActionSpec{{ID: "archive", Label: "Archive", Danger: true, Action: ActionSpec{Kind: ActionKindServer, Name: "archive-jobs"}}},
		}},
	}
	props := child(collection.ToNode().ToWidgetNode(), 0)["props"].(JSONObject)
	selection := props["multiSelection"].(JSONObject)
	if selection["mode"] != "multi" || len(selection["selectedKeys"].([]string)) != 2 {
		t.Fatalf("multiSelection = %#v", selection)
	}
	actions := props["bulkActions"].([]JSONValue)
	if len(actions) != 1 || actions[0].(JSONObject)["id"] != "archive" {
		t.Fatalf("bulkActions = %#v", actions)
	}
}
