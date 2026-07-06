package widgetdsl

import (
	"strings"
	"testing"
)

func TestDataV2TypeScriptModuleDeclaresTypedFluentBuilders(t *testing.T) {
	mod := TypeScriptModule(DataV2ModuleName)
	if mod == nil {
		t.Fatalf("TypeScriptModule(%q) returned nil", DataV2ModuleName)
	}
	dts := strings.Join(mod.RawDTS, "\n")
	wantFragments := []string{
		"export interface FieldBuilder extends FieldHandle",
		"export interface SchemaBuilder",
		"field(name: string, field: FieldHandle): this;",
		"export interface CollectionBuilder",
		"empty(message: string): this;",
		"select(selection?: SelectionHandle | SelectionCallback | null): this;",
		"edit(callback?: EditCallback | null): this;",
		"table(callback?: TableCallback | null): this;",
		"actionColumn(id: string, header: string, label: string, action: ActionHandle, options?: { maxWidth?: string }): this;",
		"toIR(): WidgetNode;",
		"export interface ActionBuilder extends ActionHandle",
		"payloadPath(name: string, path: string): this;",
		"export const f: FieldFactory;",
		"export function collection(name: string, rows: Record<string, JsonValue>[]): CollectionBuilder;",
		"export const action: ActionFactory;",
	}
	for _, fragment := range wantFragments {
		if !strings.Contains(dts, fragment) {
			t.Fatalf("data.v2.dsl DTS missing %q\n--- DTS ---\n%s", fragment, dts)
		}
	}
}

func TestDataV2TypeScriptModuleDoesNotExposeLegacyOptionBagGrammar(t *testing.T) {
	mod := TypeScriptModule(DataV2ModuleName)
	if mod == nil {
		t.Fatalf("TypeScriptModule(%q) returned nil", DataV2ModuleName)
	}
	dts := strings.Join(mod.RawDTS, "\n")
	for _, forbidden := range []string{
		"export function dataTable",
		"export interface CellSpec",
		"export const cell",
		"export function record(values: Props",
		"export function collection(rows: Props[]",
		"export function formPost",
	} {
		if strings.Contains(dts, forbidden) {
			t.Fatalf("data.v2.dsl DTS contains legacy declaration %q\n--- DTS ---\n%s", forbidden, dts)
		}
	}
}
