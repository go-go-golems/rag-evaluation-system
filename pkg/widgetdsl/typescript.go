package widgetdsl

import (
	"fmt"
	"sort"

	"github.com/go-go-golems/go-go-goja/pkg/tsgen/spec"
)

// TypeScriptModule returns the declaration descriptor for one split Widget DSL
// module. The Widget IR shape is intentionally represented as JSON-like data:
// these helpers are authoring conveniences for serializable React WidgetRenderer
// payloads, and individual component props remain open-ended by design.
func TypeScriptModule(moduleName string) *spec.Module {
	moduleSpec, ok := moduleSpecsByName[moduleName]
	if !ok {
		return nil
	}

	lines := []string{
		"export type JsonPrimitive = string | number | boolean | null;",
		"export type JsonValue = JsonPrimitive | JsonValue[] | { [key: string]: JsonValue };",
		"export type WidgetChild = WidgetNode | string | number | boolean | null | undefined;",
		"export interface WidgetNode { kind: string; [key: string]: any; }",
		"export interface WidgetPage { schemaVersion: string; id: string; title?: string; root?: WidgetNode; sections?: WidgetNode[]; [key: string]: any; }",
		"export interface WidgetAction { kind: string; [key: string]: any; }",
		"export type Props = Record<string, any>;",
		"export function text(value: any): WidgetNode;",
		"export function element(tag: string, attrs?: Props | WidgetChild, ...children: WidgetChild[]): WidgetNode;",
		"export function component(type: string, props?: Props | WidgetChild, ...children: WidgetChild[]): WidgetNode;",
		"export function fragment(...children: WidgetChild[]): WidgetNode[];",
	}
	if moduleSpec.page {
		lines = append(lines, "export function page(options: Props): WidgetPage;")
	}

	helperNames := make([]string, 0, len(moduleSpec.helpers))
	for name := range moduleSpec.helpers {
		helperNames = append(helperNames, name)
	}
	sort.Strings(helperNames)
	for _, name := range helperNames {
		lines = append(lines, fmt.Sprintf("export function %s(props?: Props | WidgetChild, ...children: WidgetChild[]): WidgetNode;", name))
	}

	if moduleSpec.cell {
		lines = append(lines,
			"export interface CellSpec { kind: string; [key: string]: any; }",
			"export const cell: {",
			"field(field: string, options?: Props): CellSpec;",
			"number(field: string, options?: Props): CellSpec;",
			"status(field: string, options?: Props): CellSpec;",
			"caption(field: string, options?: Props): CellSpec;",
			"template(template: string): CellSpec;",
			"link(hrefField: string, labelField: string, options?: Props): CellSpec;",
			"linkButton(hrefField: string, labelField: string, options?: Props): CellSpec;",
			"actionButton(label: any, action: WidgetAction, options?: Props): CellSpec;",
			"constant(value: any): CellSpec;",
			"};",
		)
	}
	if moduleSpec.action {
		lines = append(lines,
			"export const action: {",
			"server(name: string, options?: Props): WidgetAction;",
			"navigate(to: string, options?: Props): WidgetAction;",
			"download(to: string, options?: Props): WidgetAction;",
			"event(name: string, options?: Props): WidgetAction;",
			"copy(value: string, options?: Props): WidgetAction;",
			"};",
		)
	}
	if moduleSpec.name == UIModuleName {
		lines = append(lines,
			"export function section(title: WidgetChild, options?: Props, ...children: WidgetChild[]): WidgetNode;",
		)
	}
	if moduleSpec.name == DataV2ModuleName {
		lines = append(lines, dataV2TypeScriptLines()...)
	}
	if moduleSpec.name == WidgetV3ModuleName {
		lines = widgetV3TypeScriptLines()
	}
	if moduleSpec.name == DataModuleName {
		lines = append(lines,
			"export interface FieldSpec { role: string; [key: string]: any; }",
			"export interface Schema { fields: FieldSpec[]; [key: string]: any; }",
			"export const f: {",
			"key(options?: Props): FieldSpec;",
			"primary(options?: Props): FieldSpec;",
			"short(options?: Props): FieldSpec;",
			"prose(options?: Props): FieldSpec;",
			"count(options?: Props): FieldSpec;",
			"size(options?: Props): FieldSpec;",
			"measure(options?: Props): FieldSpec;",
			"date(options?: Props): FieldSpec;",
			"status(options?: Props): FieldSpec;",
			"tags(options?: Props): FieldSpec;",
			"media(options?: Props): FieldSpec;",
			"href(options?: Props): FieldSpec;",
			"};",
			"export function schema(fields: Record<string, FieldSpec>): Schema;",
			"export function record(values: Props, options: Props): WidgetNode;",
			"export function collection(rows: Props[], options: Props): WidgetNode;",
			"export function urlParam(param: string, value?: any): Props;",
			"export function formPost(formAction: string, options?: Props): Props;",
		)
	}
	if moduleSpec.name == ContextWindowModuleName {
		lines = append(lines,
			"export function contextStyleSwatch(options?: Props): WidgetNode;",
			"export function visualStyle(options: Props): any;",
			"export function legendItem(id: string, label: string, options?: Props): any;",
			"export function styleSet(options: Props): any;",
			"export function paletteStyleSet(options: Props): any;",
			"export function contextSnapshot(options: Props): any;",
			"export function contextPart(id: string, label: string, styleKey: string, tokens: number, options?: Props): any;",
		)
	}
	if len(moduleSpec.recipes) > 0 {
		lines = append(lines, "export const recipes: {")
		recipeNames := append([]string(nil), moduleSpec.recipes...)
		sort.Strings(recipeNames)
		for _, name := range recipeNames {
			lines = append(lines, fmt.Sprintf("%s(options: Props): WidgetNode;", name))
		}
		lines = append(lines, "};")
	}

	return &spec.Module{
		Name:        moduleSpec.name,
		Description: moduleSpec.doc,
		RawDTS:      lines,
	}
}

func widgetV3TypeScriptLines() []string {
	return []string{
		"export type JsonPrimitive = string | number | boolean | null;",
		"export type JsonValue = JsonPrimitive | JsonValue[] | { [key: string]: JsonValue };",
		"export interface WidgetNodeSpec { kind: string; [key: string]: any; }",
		"export interface WidgetPageSpec { schemaVersion?: string; id?: string; title?: string; root?: WidgetNodeSpec; [key: string]: any; }",
		"export interface ActionSpec { kind: string; [key: string]: any; }",
		"export interface BindingSpec { kind: string; [key: string]: any; }",
		"export type WidgetChild = WidgetNodeSpec | string | number | boolean | null | undefined;",
		"export type Fragment<TBuilder> = (builder: TBuilder) => void | TBuilder;",
		"export interface RawNamespace {",
		"text(value: any): WidgetNodeSpec;",
		"element(tag: string, attrs?: Record<string, any> | WidgetChild, ...children: WidgetChild[]): WidgetNodeSpec;",
		"component(type: string, props?: Record<string, any> | WidgetChild, ...children: WidgetChild[]): WidgetNodeSpec;",
		"fragment(...children: WidgetChild[]): WidgetNodeSpec[];",
		"}",
		"export interface ActionNamespace {",
		"server(name: string, options?: Record<string, any>): ActionSpec;",
		"navigate(to: string, options?: Record<string, any>): ActionSpec;",
		"download(to: string, options?: Record<string, any>): ActionSpec;",
		"event(name: string, options?: Record<string, any>): ActionSpec;",
		"copy(value: string, options?: Record<string, any>): ActionSpec;",
		"}",
		"export interface BindingNamespace {",
		"field(path: string): BindingSpec;",
		"path(path: string): BindingSpec;",
		"map(field: string): BindingSpec;",
		"template(template: string): BindingSpec;",
		"context(path: string): BindingSpec;",
		"const(value: JsonValue): BindingSpec;",
		"}",
		"export const raw: RawNamespace;",
		"export const act: ActionNamespace;",
		"export const bind: BindingNamespace;",
		"export const page: Record<string, any>;",
		"export const ui: Record<string, any>;",
		"export const data: Record<string, any>;",
		"export const cms: Record<string, any>;",
		"export const course: Record<string, any>;",
		"export const context: Record<string, any>;",
		"export const schedule: Record<string, any>;",
		"export const time: Record<string, any>;",
		"export const style: Record<string, any>;",
	}
}

func dataV2TypeScriptLines() []string {
	return []string{
		"declare const __widgetdslV2Brand: unique symbol;",
		"export interface ValidationIssue { severity: 'error' | 'warning'; code: string; path: string; message: string; hint?: string; }",
		"export interface FieldDescriptor { kind: string; semantic: string; label?: string; }",
		"export interface FieldHandle { readonly [__widgetdslV2Brand]: 'field'; }",
		"export interface SchemaHandle { readonly [__widgetdslV2Brand]: 'schema'; validate(): ValidationIssue[]; }",
		"export interface SelectionHandle { readonly [__widgetdslV2Brand]: 'selection'; }",
		"export interface ActionHandle { readonly [__widgetdslV2Brand]: 'action'; }",
		"export interface FieldBuilder extends FieldHandle {",
		"label(label: string): this;",
		"width(width: string): this;",
		"required(): this;",
		"maxLength(max: number): this;",
		"rows(rows: number): this;",
		"readOnly(): this;",
		"build(): FieldDescriptor;",
		"}",
		"export interface FieldFactory {",
		"key(): FieldBuilder;",
		"primary(): FieldBuilder;",
		"short(): FieldBuilder;",
		"prose(): FieldBuilder;",
		"count(): FieldBuilder;",
		"status(): FieldBuilder;",
		"}",
		"export interface SchemaBuilder {",
		"field(name: string, field: FieldHandle): this;",
		"build(): SchemaHandle;",
		"validate(): ValidationIssue[];",
		"}",
		"export interface SelectionBuilder {",
		"urlParam(param: string, value?: string | number | boolean | null): SelectionHandle;",
		"}",
		"export type SelectionCallback = (selection: SelectionBuilder) => SelectionHandle | void;",
		"export interface ActionBuilder extends ActionHandle {",
		"confirm(text: string): this;",
		"payloadPath(name: string, path: string): this;",
		"payload(name: string, value: JsonValue): this;",
		"}",
		"export interface ActionFactory {",
		"navigate(to: string): ActionBuilder;",
		"server(name: string): ActionBuilder;",
		"}",
		"export interface TableBuilder {",
		"className(className: string): this;",
		"rowSelect(action: ActionHandle): this;",
		"actionColumn(id: string, header: string, label: string, action: ActionHandle, options?: { maxWidth?: string }): this;",
		"}",
		"export type TableCallback = (table: TableBuilder) => void;",
		"export interface CollectionActionsBuilder {",
		"reorder(action: ActionHandle): this;",
		"remove(action: ActionHandle): this;",
		"}",
		"export interface EditorBuilder {",
		"selectUrl(param: string, value?: string | number | boolean | null): this;",
		"submitPost(formAction: string): this;",
		"create(value?: string | { label?: string }): this;",
		"reorder(action: ActionHandle): this;",
		"remove(action: ActionHandle): this;",
		"actions(callback: (actions: CollectionActionsBuilder) => void): this;",
		"}",
		"export type EditCallback = (editor: EditorBuilder) => void;",
		"export interface CollectionBuilder {",
		"schema(schema: SchemaHandle): this;",
		"empty(message: string): this;",
		"select(selection?: SelectionHandle | SelectionCallback | null): this;",
		"edit(callback?: EditCallback | null): this;",
		"table(callback?: TableCallback | null): this;",
		"masterDetail(): this;",
		"validate(): ValidationIssue[];",
		"toIR(): WidgetNode;",
		"}",
		"export const f: FieldFactory;",
		"export function schema(name: string): SchemaBuilder;",
		"export function collection(name: string, rows: Record<string, JsonValue>[]): CollectionBuilder;",
		"export const selection: { urlParam(param: string, value?: string | number | boolean | null): SelectionHandle; };",
		"export const action: ActionFactory;",
	}
}
