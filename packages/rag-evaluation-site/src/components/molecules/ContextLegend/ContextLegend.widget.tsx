import { ContextLegend } from './ContextLegend';
import { defineWidget } from '../../../widgets/registry';
import type { ContextLegendWidgetProps } from '../../../widgets/ir';

export const contextLegendWidget = defineWidget<ContextLegendWidgetProps>({
  type: 'ContextLegend',
  module: 'context_window.dsl',
  render: (props) => <ContextLegend className={props.className} kinds={props.kinds} mode={props.mode} compact={props.compact} selectedKind={props.selectedKind} />,
});
