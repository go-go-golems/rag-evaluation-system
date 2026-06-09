import type { ReactNode } from 'react';
import type { ActionSpec, ComponentNode, WidgetNode } from './ir';
import type { WidgetActionContext } from './actions';

export type WidgetModule = 'ui.dsl' | 'data.dsl' | 'context_window.dsl' | 'course.dsl';

export interface RenderContext {
  renderNode(node: WidgetNode): ReactNode;
  renderChildren(children?: WidgetNode[]): ReactNode[];
  renderValue(value: unknown): ReactNode;
  bindAction(action: ActionSpec | undefined, context: WidgetActionContext): (() => void) | undefined;
  dispatchAction(action: ActionSpec, context: WidgetActionContext): void;
}

export interface WidgetAdapter<P = unknown> {
  type: string;
  module: WidgetModule;
  render(props: P, children: ReactNode[], ctx: RenderContext, node: ComponentNode): ReactNode;
}

export function defineWidget<P>(adapter: WidgetAdapter<P>): WidgetAdapter<P> {
  return adapter;
}
