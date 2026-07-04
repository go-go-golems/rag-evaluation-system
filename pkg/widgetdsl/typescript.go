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
