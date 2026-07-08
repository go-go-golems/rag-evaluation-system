package widgetdsl

import (
	"strings"
	"testing"
)

func TestWidgetV3TypeScriptModuleDeclaresRootNamespaces(t *testing.T) {
	mod := TypeScriptModule(WidgetV3ModuleName)
	if mod == nil {
		t.Fatalf("TypeScriptModule(%q) returned nil", WidgetV3ModuleName)
	}
	dts := strings.Join(mod.RawDTS, "\n")
	wantFragments := []string{
		"export interface RawNamespace",
		"component(type: string",
		"export interface ActionNamespace",
		"server(name: string",
		"export interface AccessorSpec",
		"export interface SelectionSpec",
		"export interface ListItemSpec",
		"export interface BindingNamespace",
		"field(path: string): AccessorSpec;",
		"export type Slot<TContext>",
		"export interface SlotHelpers",
		"strong(...children: WidgetChild[]): WidgetNodeSpec;",
		"slot<TContext>(context: TContext, slot?: Slot<TContext> | null, fallback?: Slot<TContext>): this;",
		"export interface FieldSetBuilder",
		"key(name: string, options?: Record<string, any>): this;",
		"export interface CollectionBuilder",
		"table(configure?: Fragment<TableBuilder>): this;",
		"export interface CellNamespace",
		"cycle(field: string, options?: Record<string, any>): Record<string, any>;",
		"export interface MatrixBuilder",
		"valueAt(accessor: AccessorSpec): this;",
		"export interface DataNamespace",
		"fields(nameOrConfigure?: string | Fragment<FieldSetBuilder>, configure?: Fragment<FieldSetBuilder>): FieldSetBuilder;",
		"selection: ((modeOrOptions: 'single' | 'multi' | { mode?: 'single' | 'multi'; keyField?: string; selected?: JsonValue }, options?: { keyField?: string; selected?: JsonValue }) => SelectionSpec) & { urlParam(param: string, value: JsonValue): Record<string, any> };",
		"item(id: string, label: WidgetChild | WidgetChild[], options?: Record<string, any>): ListItemSpec;",
		"export const raw: RawNamespace;",
		"export const act: ActionNamespace;",
		"export const bind: BindingNamespace;",
		"export const data: DataNamespace;",
		"export interface PageBuilder",
		"export interface SectionBuilder",
		"export function page(titleOrOptions: string | Record<string, any>, configure?: Fragment<PageBuilder>): PageBuilder;",
		"export interface UINamespace",
		"callout(options?: Record<string, any> | WidgetChild, ...children: WidgetChild[]): WidgetNodeSpec;",
		"export interface ActionsBuilder",
		"shell(shellSpec: JsonValue): this;",
		"density(value: string): this;",
		"breadcrumb(label: WidgetChild, href?: string): this;",
		"actions(configure: Fragment<ActionsBuilder>): this;",
		"metric(label: WidgetChild, value: WidgetChild, options?: Record<string, any>): this;",
		"metadata(record: Record<string, JsonValue>): this;",
		"export const ui: UINamespace;",
		"export interface CmsNamespace",
		"mediaLibrary(assets: CmsAsset[], configure?: Fragment<CmsMediaLibraryBuilder>): WidgetNodeSpec;",
		"articleQueue(articles: CmsArticleSummary[], configure?: Fragment<CmsArticleQueueBuilder>): WidgetNodeSpec;",
		"markdownEditor(body: string, configure?: Fragment<CmsMarkdownEditorBuilder>): WidgetNodeSpec;",
		"export const cms: CmsNamespace;",
		"export const course: Record<string, any>;",
		"export const context: Record<string, any>;",
		"export const schedule: Record<string, any>;",
		"export const time: Record<string, any>;",
		"export const style: Record<string, any>;",
	}
	for _, fragment := range wantFragments {
		if !strings.Contains(dts, fragment) {
			t.Fatalf("widget.dsl DTS missing %q\n--- DTS ---\n%s", fragment, dts)
		}
	}
}

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
