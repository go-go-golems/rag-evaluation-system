package experiments

import "testing"

func TestCanonicalJSONSortsNestedObjectKeys(t *testing.T) {
	first, err := CanonicalJSON(map[string]any{"z": 1, "a": map[string]any{"b": true, "a": "x"}})
	if err != nil {
		t.Fatal(err)
	}
	second, err := CanonicalJSON(map[string]any{"a": map[string]any{"a": "x", "b": true}, "z": 1})
	if err != nil {
		t.Fatal(err)
	}
	if string(first) != `{"a":{"a":"x","b":true},"z":1}` {
		t.Fatalf("canonical JSON = %s", first)
	}
	if string(first) != string(second) {
		t.Fatalf("key order changed canonical bytes: %s != %s", first, second)
	}
}

func TestFingerprintIsSchemaNamespacedAndArrayOrdered(t *testing.T) {
	first, err := Fingerprint("plan/v1", map[string]any{"items": []string{"a", "b"}})
	if err != nil {
		t.Fatal(err)
	}
	changedSchema, err := Fingerprint("plan/v2", map[string]any{"items": []string{"a", "b"}})
	if err != nil {
		t.Fatal(err)
	}
	changedOrder, err := Fingerprint("plan/v1", map[string]any{"items": []string{"b", "a"}})
	if err != nil {
		t.Fatal(err)
	}
	if first == changedSchema || first == changedOrder {
		t.Fatalf("fingerprint did not distinguish semantic input")
	}
}
