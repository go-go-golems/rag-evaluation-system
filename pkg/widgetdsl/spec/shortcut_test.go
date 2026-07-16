package spec

import "testing"

func TestPageShortcutsLowerToPageEnvelope(t *testing.T) {
	page := PageSpec{
		ID: "triage",
		Shortcuts: []PageShortcutSpec{{
			ID:             "accept",
			Key:            "y",
			Modifiers:      []ShortcutModifier{ShortcutModifierControl},
			Label:          "Yes",
			Action:         ActionSpec{Kind: ActionKindServer, Name: "triage.accept"},
			PreventDefault: true,
		}},
	}

	lowered := page.ToWidgetPage()
	shortcuts := lowered["shortcuts"].(JSONObject)
	bindings := shortcuts["bindings"].([]JSONValue)
	binding := bindings[0].(JSONObject)
	if binding["id"] != "accept" || binding["key"] != "y" || binding["label"] != "Yes" {
		t.Fatalf("shortcut identity not lowered: %#v", binding)
	}
	if binding["preventDefault"] != true || binding["allowRepeat"] != false {
		t.Fatalf("shortcut policies not lowered: %#v", binding)
	}
	modifiers := binding["modifiers"].([]JSONValue)
	if len(modifiers) != 1 || modifiers[0] != "Control" {
		t.Fatalf("shortcut modifiers not lowered: %#v", modifiers)
	}
	action := binding["action"].(JSONObject)
	if action["kind"] != "server" || action["name"] != "triage.accept" {
		t.Fatalf("shortcut action not lowered: %#v", action)
	}
}

func TestPageShortcutValidationRejectsDuplicateIdentityAndChord(t *testing.T) {
	page := PageSpec{
		ID: "triage",
		Shortcuts: []PageShortcutSpec{
			{ID: "accept", Key: "y", Label: "Yes", Action: ActionSpec{Kind: ActionKindServer, Name: "triage.accept"}, PreventDefault: true},
			{ID: "accept", Key: "Y", Label: "Again", Action: ActionSpec{Kind: ActionKindServer, Name: "triage.again"}, PreventDefault: true},
		},
	}

	issues := page.Validate()
	assertIssueCode(t, issues, "page.shortcut.id.duplicate")
	assertIssueCode(t, issues, "page.shortcut.chord.duplicate")
	assertIssueCode(t, issues, "page.shortcut.character.unmodified")
}

func TestPageShortcutValidationRejectsInvalidShape(t *testing.T) {
	shortcut := PageShortcutSpec{
		ID:        "",
		Key:       "Control+y",
		Modifiers: []ShortcutModifier{ShortcutModifierControl, ShortcutModifierControl, "Hyper"},
		Label:     "",
		Action:    ActionSpec{},
	}

	issues := shortcut.Validate("shortcuts.bindings[0]")
	for _, code := range []string{
		"page.shortcut.id.required",
		"page.shortcut.key.chord.invalid",
		"page.shortcut.label.required",
		"page.shortcut.modifier.duplicate",
		"page.shortcut.modifier.invalid",
		"action.kind.required",
	} {
		assertIssueCode(t, issues, code)
	}
}

func TestPageShortcutCanonicalChordNormalizesOrderingAndLetterCase(t *testing.T) {
	shortcut := PageShortcutSpec{
		Key:       "Y",
		Modifiers: []ShortcutModifier{ShortcutModifierShift, ShortcutModifierControl},
	}
	if got, want := shortcut.CanonicalChord(), "Control+Shift+y"; got != want {
		t.Fatalf("CanonicalChord() = %q, want %q", got, want)
	}
}

func assertIssueCode(t *testing.T, issues []ValidationIssue, code string) {
	t.Helper()
	for _, issue := range issues {
		if issue.Code == code {
			return
		}
	}
	t.Fatalf("missing issue code %q in %#v", code, issues)
}
