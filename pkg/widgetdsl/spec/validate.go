package spec

import (
	"fmt"
	"strings"
)

// HasErrors reports whether a validation result contains at least one error.
func HasErrors(issues []ValidationIssue) bool {
	for _, issue := range issues {
		if issue.Severity == ValidationSeverityError {
			return true
		}
	}
	return false
}

// Validate checks page-level invariants and recursively validates the root node.
func (s PageSpec) Validate() []ValidationIssue {
	issues := []ValidationIssue{}
	if strings.TrimSpace(s.ID) == "" {
		issues = append(issues, errorIssue("page.id.required", "id", "Page id is required.", "Set a stable page id before calling toIR()."))
	}
	if s.Shell != nil {
		issues = append(issues, s.Shell.Validate("shell")...)
	}
	seenShortcutIDs := map[string]int{}
	seenShortcutChords := map[string]int{}
	for i, shortcut := range s.Shortcuts {
		path := fmt.Sprintf("shortcuts.bindings[%d]", i)
		issues = append(issues, shortcut.Validate(path)...)
		id := strings.TrimSpace(shortcut.ID)
		if first, ok := seenShortcutIDs[id]; id != "" && ok {
			issues = append(issues, errorIssue("page.shortcut.id.duplicate", path+".id", fmt.Sprintf("Duplicate shortcut id %q; first used at bindings[%d].", id, first), "Use a unique stable command id."))
		} else if id != "" {
			seenShortcutIDs[id] = i
		}
		chord := shortcut.CanonicalChord()
		if first, ok := seenShortcutChords[chord]; chord != "" && ok {
			issues = append(issues, errorIssue("page.shortcut.chord.duplicate", path+".key", fmt.Sprintf("Duplicate shortcut chord %q; first used at bindings[%d].", chord, first), "Choose a unique key and modifier combination."))
		} else if chord != "" {
			seenShortcutChords[chord] = i
		}
	}
	if s.Root.Kind != "" {
		issues = append(issues, s.Root.Validate("root")...)
	}
	return issues
}

// Validate checks the shell and navigation contract before browser transport.
func (s PageShellSpec) Validate(path string) []ValidationIssue {
	issues := []ValidationIssue{}
	switch s.Kind {
	case PageShellKindNone, PageShellKindRootOwned:
		if s.Navigation != nil {
			issues = append(issues, errorIssue("page.shell.navigation.unexpected", path+".navigation", "Only app shells accept navigation.", "Remove navigation or use shell kind app."))
		}
	case PageShellKindApp:
		if s.Navigation == nil {
			issues = append(issues, errorIssue("page.shell.navigation.required", path+".navigation", "App shell navigation is required.", "Configure top or sidebar navigation."))
		} else {
			issues = append(issues, s.Navigation.Validate(path+".navigation")...)
		}
		issues = append(issues, s.Content.Validate(path+".content")...)
	default:
		issues = append(issues, errorIssue("page.shell.kind.invalid", path+".kind", fmt.Sprintf("Unknown page shell kind %q.", s.Kind), "Use none, app, or root-owned."))
	}
	return issues
}

// Validate checks semantic navigation invariants.
func (s NavigationSpec) Validate(path string) []ValidationIssue {
	issues := []ValidationIssue{}
	if s.Placement != NavigationPlacementTop && s.Placement != NavigationPlacementSidebar {
		issues = append(issues, errorIssue("page.shell.navigation.placement.invalid", path+".placement", fmt.Sprintf("Unknown navigation placement %q.", s.Placement), "Use top or sidebar."))
	}
	if s.SidebarWidth != 0 && (s.SidebarWidth < 160 || s.SidebarWidth > 320) {
		issues = append(issues, errorIssue("page.shell.navigation.width.invalid", path+".sidebarWidth", "Sidebar width must be between 160 and 320 pixels.", "Use a compact, readable navigation rail width."))
	}
	if s.NarrowMode != "" && s.NarrowMode != "stack" {
		issues = append(issues, errorIssue("page.shell.navigation.narrow_mode.invalid", path+".narrowMode", fmt.Sprintf("Unknown narrow navigation mode %q.", s.NarrowMode), "Use stack; drawer and top-scroll are not implemented."))
	}
	seen := map[string]bool{}
	for sectionIndex, section := range s.Sections {
		sectionPath := fmt.Sprintf("%s.sections[%d]", path, sectionIndex)
		if strings.TrimSpace(section.ID) == "" {
			issues = append(issues, errorIssue("page.shell.navigation.section.id.required", sectionPath+".id", "Navigation section id is required.", "Set a stable section id."))
		}
		for itemIndex, item := range section.Items {
			itemPath := fmt.Sprintf("%s.items[%d]", sectionPath, itemIndex)
			if strings.TrimSpace(item.ID) == "" {
				issues = append(issues, errorIssue("page.shell.navigation.item.id.required", itemPath+".id", "Navigation item id is required.", "Set a stable item id."))
			} else if seen[item.ID] {
				issues = append(issues, errorIssue("page.shell.navigation.item.id.duplicate", itemPath+".id", fmt.Sprintf("Duplicate navigation item id %q.", item.ID), "Use unique item ids across all sections."))
			}
			seen[item.ID] = true
			if !item.Disabled && item.Action == nil {
				issues = append(issues, errorIssue("page.shell.navigation.item.action.required", itemPath+".action", "Enabled navigation item action is required.", "Provide a serializable navigation action."))
			}
		}
	}
	return issues
}

// Validate checks content viewport tokens.
func (s ContentViewportSpec) Validate(path string) []ValidationIssue {
	issues := []ValidationIssue{}
	if s.MaxWidth != "" && s.MaxWidth != "none" && s.MaxWidth != "wide" && s.MaxWidth != "content" {
		issues = append(issues, errorIssue("page.shell.content.max_width.invalid", path+".maxWidth", fmt.Sprintf("Unknown content max width %q.", s.MaxWidth), "Use none, wide, or content."))
	}
	if s.Padding != "" && s.Padding != "none" && s.Padding != "md" && s.Padding != "lg" {
		issues = append(issues, errorIssue("page.shell.content.padding.invalid", path+".padding", fmt.Sprintf("Unknown content padding %q.", s.Padding), "Use none, md, or lg."))
	}
	if s.Scroll != "" && s.Scroll != "page" && s.Scroll != "main" {
		issues = append(issues, errorIssue("page.shell.content.scroll.invalid", path+".scroll", fmt.Sprintf("Unknown content scroll mode %q.", s.Scroll), "Use page or main."))
	}
	return issues
}

// CanonicalChord returns the normalized chord used for duplicate validation.
func (s PageShortcutSpec) CanonicalChord() string {
	key := strings.TrimSpace(s.Key)
	if len(key) == 1 && key[0] >= 'A' && key[0] <= 'Z' {
		key = strings.ToLower(key)
	}
	present := map[ShortcutModifier]bool{}
	for _, modifier := range s.Modifiers {
		present[modifier] = true
	}
	parts := []string{}
	for _, modifier := range []ShortcutModifier{ShortcutModifierAlt, ShortcutModifierControl, ShortcutModifierMeta, ShortcutModifierShift} {
		if present[modifier] {
			parts = append(parts, string(modifier))
		}
	}
	if key == "" {
		return ""
	}
	return strings.Join(append(parts, key), "+")
}

// Validate checks a page shortcut before browser transport.
func (s PageShortcutSpec) Validate(path string) []ValidationIssue {
	issues := []ValidationIssue{}
	if strings.TrimSpace(s.ID) == "" {
		issues = append(issues, errorIssue("page.shortcut.id.required", path+".id", "Shortcut id is required.", "Use a stable command id such as accept or skip."))
	}
	key := strings.TrimSpace(s.Key)
	if key == "" {
		issues = append(issues, errorIssue("page.shortcut.key.required", path+".key", "Shortcut key is required.", "Use a KeyboardEvent.key value such as y or Enter."))
	} else if strings.Contains(key, "+") && len(key) > 1 {
		issues = append(issues, errorIssue("page.shortcut.key.chord.invalid", path+".key", "Shortcut key must not contain serialized modifiers.", "Put Control, Alt, Meta, or Shift in modifiers and keep key separate."))
	}
	if strings.TrimSpace(s.Label) == "" {
		issues = append(issues, errorIssue("page.shortcut.label.required", path+".label", "Shortcut label is required.", "Describe the command for shortcut help and accessibility."))
	}
	seenModifiers := map[ShortcutModifier]bool{}
	for i, modifier := range s.Modifiers {
		modifierPath := fmt.Sprintf("%s.modifiers[%d]", path, i)
		if !validShortcutModifier(modifier) {
			issues = append(issues, errorIssue("page.shortcut.modifier.invalid", modifierPath, fmt.Sprintf("Unknown shortcut modifier %q.", modifier), "Use Alt, Control, Meta, or Shift."))
		} else if seenModifiers[modifier] {
			issues = append(issues, errorIssue("page.shortcut.modifier.duplicate", modifierPath, fmt.Sprintf("Duplicate shortcut modifier %q.", modifier), "List each modifier once."))
		}
		seenModifiers[modifier] = true
	}
	if len(key) == 1 && len(s.Modifiers) == 0 {
		issues = append(issues, warningIssue("page.shortcut.character.unmodified", path+".key", "Unmodified character shortcuts require a user-facing disable, remap, or focus-only mechanism.", "Expose the host shortcut preference and visible shortcut help."))
	}
	issues = append(issues, s.Action.Validate(path+".action")...)
	return issues
}

// Validate checks node-level shape invariants.
func (s NodeSpec) Validate(path string) []ValidationIssue {
	issues := []ValidationIssue{}
	switch s.Kind {
	case NodeKindText:
		if s.Text == "" {
			issues = append(issues, warningIssue("node.text.empty", path+".text", "Text node is empty.", "Remove the node if this is intentional."))
		}
	case NodeKindElement:
		if strings.TrimSpace(s.Tag) == "" {
			issues = append(issues, errorIssue("node.element.tag.required", path+".tag", "Element node tag is required.", "Use a valid HTML tag name."))
		}
	case NodeKindComponent:
		if strings.TrimSpace(s.Type) == "" {
			issues = append(issues, errorIssue("node.component.type.required", path+".type", "Component node type is required.", "Use a registered Widget IR component type."))
		}
	case "":
		issues = append(issues, errorIssue("node.kind.required", path+".kind", "Node kind is required.", "Use text, element, or component."))
	default:
		issues = append(issues, errorIssue("node.kind.unknown", path+".kind", fmt.Sprintf("Unknown node kind %q.", s.Kind), "Use text, element, or component."))
	}
	for i, child := range s.Children {
		issues = append(issues, child.Validate(fmt.Sprintf("%s.children[%d]", path, i))...)
	}
	return issues
}

// Validate checks section intent before lowering to SectionBlock.
func (s SectionSpec) Validate(path string) []ValidationIssue {
	issues := []ValidationIssue{}
	if strings.TrimSpace(s.Title) == "" {
		issues = append(issues, errorIssue("section.title.required", path+".title", "Section title is required.", "Give every section a visible title."))
	}
	if s.Level < 1 || s.Level > 3 {
		issues = append(issues, errorIssue("section.level.invalid", path+".level", fmt.Sprintf("Section level %d is invalid.", s.Level), "Use level 1, 2, or 3."))
	}
	for i, child := range s.Children {
		issues = append(issues, child.Validate(fmt.Sprintf("%s.children[%d]", path, i))...)
	}
	return issues
}

// Validate checks schema and ordered field invariants.
func (s SchemaSpec) Validate(path string) []ValidationIssue {
	issues := []ValidationIssue{}
	if strings.TrimSpace(s.Name) == "" {
		issues = append(issues, errorIssue("schema.name.required", path+".name", "Schema name is required.", "Use a stable domain name such as AgendaItem or Session."))
	}
	if len(s.Fields) == 0 {
		issues = append(issues, errorIssue("schema.fields.required", path+".fields", "Schema must contain at least one field.", "Add fields before building collections or records."))
	}
	seen := map[string]int{}
	keyCount := 0
	for i, field := range s.Fields {
		fieldPath := fmt.Sprintf("%s.fields[%d]", path, i)
		issues = append(issues, field.Validate(fieldPath)...)
		name := strings.TrimSpace(field.Name)
		if name == "" {
			continue
		}
		if first, ok := seen[name]; ok {
			issues = append(issues, errorIssue("schema.field.duplicate", fieldPath+".name", fmt.Sprintf("Duplicate field name %q; first used at fields[%d].", name, first), "Use unique field names."))
		} else {
			seen[name] = i
		}
		if field.Semantic == FieldSemanticKey {
			keyCount++
		}
	}
	if keyCount > 1 {
		issues = append(issues, errorIssue("schema.key.multiple", path+".fields", "Schema has more than one key field.", "Choose one stable key field for row identity."))
	}
	return issues
}

// Validate checks field-level invariants.
func (s FieldSpec) Validate(path string) []ValidationIssue {
	issues := []ValidationIssue{}
	if strings.TrimSpace(s.Name) == "" {
		issues = append(issues, errorIssue("field.name.required", path+".name", "Field name is required.", "Use the JSON property name for this field."))
	}
	if !validFieldKind(s.Kind) {
		issues = append(issues, errorIssue("field.kind.invalid", path+".kind", fmt.Sprintf("Unknown field kind %q.", s.Kind), "Use a supported field kind."))
	}
	if !validFieldSemantic(s.Semantic) {
		issues = append(issues, errorIssue("field.semantic.invalid", path+".semantic", fmt.Sprintf("Unknown field semantic %q.", s.Semantic), "Use a supported field semantic."))
	}
	if s.Validation.MinLength > 0 && s.Validation.MaxLength > 0 && s.Validation.MinLength > s.Validation.MaxLength {
		issues = append(issues, errorIssue("field.validation.length_range", path+".validation", "Field minLength is greater than maxLength.", "Make minLength less than or equal to maxLength."))
	}
	return issues
}

// Validate checks collection-level invariants.
func (s CollectionSpec) Validate(path string) []ValidationIssue {
	issues := []ValidationIssue{}
	if strings.TrimSpace(s.Name) == "" {
		issues = append(issues, errorIssue("collection.name.required", path+".name", "Collection name is required.", "Use a stable collection name for diagnostics and refs."))
	}
	issues = append(issues, s.Schema.Validate(path+".schema")...)
	if !validCollectionMode(s.Mode) {
		issues = append(issues, errorIssue("collection.mode.invalid", path+".mode", fmt.Sprintf("Unknown collection mode %q.", s.Mode), "Use show, edit, pick, or manage."))
	}
	if !validArrangementKind(s.Arrangement.Kind) {
		issues = append(issues, errorIssue("collection.arrangement.invalid", path+".arrangement.kind", fmt.Sprintf("Unknown collection arrangement %q.", s.Arrangement.Kind), "Use table or master-detail."))
	}
	if s.Selection != nil {
		issues = append(issues, s.Selection.Validate(path+".selection")...)
	}
	for i, column := range s.Table.ActionColumns {
		issues = append(issues, column.Validate(fmt.Sprintf("%s.table.actionColumns[%d]", path, i))...)
	}
	issues = append(issues, validateOptionalAction(s.Table.RowSelect, path+".table.rowSelect")...)
	issues = append(issues, validateOptionalAction(s.Actions.Open, path+".actions.open")...)
	issues = append(issues, validateOptionalAction(s.Actions.Reorder, path+".actions.reorder")...)
	issues = append(issues, validateOptionalAction(s.Actions.Remove, path+".actions.remove")...)
	if s.Actions.Submit != nil && strings.TrimSpace(s.Actions.Submit.FormAction) == "" {
		issues = append(issues, errorIssue("collection.submit.action.required", path+".actions.submit.formAction", "Submit formAction is required.", "Set the native form POST target."))
	}
	return issues
}

// Validate checks URL-backed selection intent.
func (s SelectionSpec) Validate(path string) []ValidationIssue {
	issues := []ValidationIssue{}
	if s.Kind != SelectionKindURLParam {
		issues = append(issues, errorIssue("selection.kind.invalid", path+".kind", fmt.Sprintf("Unknown selection kind %q.", s.Kind), "Use urlParam."))
	}
	if strings.TrimSpace(s.Param) == "" {
		issues = append(issues, errorIssue("selection.param.required", path+".param", "Selection URL parameter is required.", "Use a stable query parameter such as selected or agenda."))
	}
	return issues
}

// Validate checks action shape before it crosses the browser boundary.
func (s ActionSpec) Validate(path string) []ValidationIssue {
	issues := []ValidationIssue{}
	switch s.Kind {
	case ActionKindNavigate, ActionKindDownload:
		if strings.TrimSpace(s.To) == "" {
			issues = append(issues, errorIssue("action.target.required", path+".to", fmt.Sprintf("%s action target is required.", s.Kind), "Set a URL or URL template."))
		}
	case ActionKindServer:
		if strings.TrimSpace(s.Name) == "" {
			issues = append(issues, errorIssue("action.server.name.required", path+".name", "Server action name is required.", "Use a registered server action name."))
		}
	case ActionKindEvent:
		if strings.TrimSpace(s.Event) == "" {
			issues = append(issues, errorIssue("action.event.required", path+".event", "Event action name is required.", "Set the DOM/custom event name."))
		}
	case ActionKindOpenOverlay:
		if strings.TrimSpace(fmt.Sprint(s.Options["target"])) == "" {
			issues = append(issues, errorIssue("action.overlay.target.required", path+".target", "Overlay target is required.", "Use the id of a declared FormDialog."))
		}
	case ActionKindCloseOverlay, ActionKindCopy:
		// Closing the current overlay and copying from context need no fixed target.
	case "":
		issues = append(issues, errorIssue("action.kind.required", path+".kind", "Action kind is required.", "Use navigate, server, download, event, or copy."))
	default:
		issues = append(issues, errorIssue("action.kind.invalid", path+".kind", fmt.Sprintf("Unknown action kind %q.", s.Kind), "Use navigate, server, download, event, or copy."))
	}
	if s.Confirm != nil {
		issues = append(issues, s.Confirm.Validate(path+".confirm")...)
	}
	issues = append(issues, s.Payload.Validate(path+".payload")...)
	return issues
}

// Validate checks payload template fields.
func (s PayloadTemplate) Validate(path string) []ValidationIssue {
	issues := []ValidationIssue{}
	seen := map[string]int{}
	for i, field := range s.Fields {
		fieldPath := fmt.Sprintf("%s.fields[%d]", path, i)
		name := strings.TrimSpace(field.Name)
		if name == "" {
			issues = append(issues, errorIssue("payload.field.name.required", fieldPath+".name", "Payload field name is required.", "Use the JSON payload property name."))
		} else if first, ok := seen[name]; ok {
			issues = append(issues, errorIssue("payload.field.duplicate", fieldPath+".name", fmt.Sprintf("Duplicate payload field %q; first used at fields[%d].", name, first), "Use unique payload field names."))
		} else {
			seen[name] = i
		}
		issues = append(issues, field.Value.Validate(fieldPath+".value")...)
	}
	return issues
}

// Validate checks a text/path/literal template.
func (s TemplateSpec) Validate(path string) []ValidationIssue {
	issues := []ValidationIssue{}
	if len(s.Parts) == 0 {
		issues = append(issues, errorIssue("template.parts.required", path+".parts", "Template must contain at least one part.", "Add text, path, or literal parts."))
	}
	for i, part := range s.Parts {
		issues = append(issues, part.Validate(fmt.Sprintf("%s.parts[%d]", path, i))...)
	}
	return issues
}

// Validate checks a single template value.
func (s TemplateValue) Validate(path string) []ValidationIssue {
	switch s.Kind {
	case TemplateValueText:
		if s.Text == "" {
			return []ValidationIssue{warningIssue("template.text.empty", path+".text", "Template text part is empty.", "Remove the part if intentional.")}
		}
	case TemplateValuePath:
		if strings.TrimSpace(s.Path) == "" {
			return []ValidationIssue{errorIssue("template.path.required", path+".path", "Template path is required.", "Use a context path such as row.id or row.title.")}
		}
	case TemplateValueLiteral:
		// Any JSONValue, including nil, is valid.
	case "":
		return []ValidationIssue{errorIssue("template.value.kind.required", path+".kind", "Template value kind is required.", "Use text, path, or literal.")}
	default:
		return []ValidationIssue{errorIssue("template.value.kind.invalid", path+".kind", fmt.Sprintf("Unknown template value kind %q.", s.Kind), "Use text, path, or literal.")}
	}
	return nil
}

// Validate checks an explicit row-action table column.
func (s TableActionColumnSpec) Validate(path string) []ValidationIssue {
	issues := []ValidationIssue{}
	if strings.TrimSpace(s.ID) == "" {
		issues = append(issues, errorIssue("table_action.id.required", path+".id", "Table action column id is required.", "Use a stable id such as open, edit, or delete."))
	}
	if strings.TrimSpace(s.Label) == "" {
		issues = append(issues, errorIssue("table_action.label.required", path+".label", "Table action column label is required.", "Set the visible button label."))
	}
	issues = append(issues, s.Action.Validate(path+".action")...)
	return issues
}

func validateOptionalAction(action *ActionSpec, path string) []ValidationIssue {
	if action == nil {
		return nil
	}
	return action.Validate(path)
}

func validShortcutModifier(modifier ShortcutModifier) bool {
	switch modifier {
	case ShortcutModifierAlt, ShortcutModifierControl, ShortcutModifierMeta, ShortcutModifierShift:
		return true
	default:
		return false
	}
}

func validFieldKind(kind FieldKind) bool {
	switch kind {
	case FieldKindString, FieldKindNumber, FieldKindBoolean, FieldKindDate, FieldKindMedia, FieldKindURL:
		return true
	default:
		return false
	}
}

func validFieldSemantic(semantic FieldSemantic) bool {
	switch semantic {
	case FieldSemanticKey, FieldSemanticPrimary, FieldSemanticShort, FieldSemanticProse, FieldSemanticCount, FieldSemanticSize, FieldSemanticMeasure, FieldSemanticStatus, FieldSemanticTags:
		return true
	default:
		return false
	}
}

func validCollectionMode(mode CollectionMode) bool {
	switch mode {
	case CollectionModeShow, CollectionModeEdit, CollectionModePick, CollectionModeManage:
		return true
	default:
		return false
	}
}

func validArrangementKind(kind ArrangementKind) bool {
	switch kind {
	case ArrangementKindTable, ArrangementKindMasterDetail:
		return true
	default:
		return false
	}
}

func errorIssue(code, path, message, hint string) ValidationIssue {
	return ValidationIssue{Severity: ValidationSeverityError, Code: code, Path: path, Message: message, Hint: hint}
}

func warningIssue(code, path, message, hint string) ValidationIssue {
	return ValidationIssue{Severity: ValidationSeverityWarning, Code: code, Path: path, Message: message, Hint: hint}
}
