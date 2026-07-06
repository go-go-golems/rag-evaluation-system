import type { ReactNode } from "react";
import type { WidgetActionContext } from "./actions";
import type { ActionSpec, ComponentNode, WidgetNode } from "./ir";

export type WidgetModule = "ui.dsl" | "data.dsl" | "context_window.dsl" | "course.dsl" | "cms.dsl";

export interface RenderContext {
	renderNode(node: WidgetNode): ReactNode;
	renderChildren(children?: WidgetNode[]): ReactNode[];
	renderValue(value: unknown): ReactNode;
	bindAction(
		action: ActionSpec | undefined,
		context: WidgetActionContext,
	): (() => void) | undefined;
	dispatchAction(action: ActionSpec, context: WidgetActionContext): void;
}

export interface WidgetAdapter<P = unknown> {
	type: string;
	module: WidgetModule;
	render(props: P, children: ReactNode[], ctx: RenderContext, node: ComponentNode): ReactNode;
}

export interface WidgetRegistry {
	get(type: string): WidgetAdapter | undefined;
	has(type: string): boolean;
	entries(): WidgetAdapter[];
}

export function defineWidget<P>(adapter: WidgetAdapter<P>): WidgetAdapter<P> {
	return adapter;
}

export function createWidgetRegistry(adapters: readonly WidgetAdapter[]): WidgetRegistry {
	const byType = new Map<string, WidgetAdapter>();
	for (const adapter of adapters) {
		if (byType.has(adapter.type)) {
			throw new Error(`Duplicate widget adapter for type ${adapter.type}`);
		}
		byType.set(adapter.type, adapter);
	}
	return {
		get: (type) => byType.get(type),
		has: (type) => byType.has(type),
		entries: () => Array.from(byType.values()),
	};
}

export function mergeWidgetRegistries(...registries: readonly WidgetRegistry[]): WidgetRegistry {
	return createWidgetRegistry(registries.flatMap((registry) => registry.entries()));
}
