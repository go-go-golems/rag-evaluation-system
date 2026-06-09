import { ContextTreemap } from './ContextTreemap';
import { defineWidget } from '../../../widgets/registry';
import type { ContextTreemapWidgetProps } from '../../../widgets/ir';

export const contextTreemapWidget = defineWidget<ContextTreemapWidgetProps>({
  type: 'ContextTreemap',
  module: 'context_window.dsl',
  render: (props) => <ContextTreemap className={props.className} snapshot={props.snapshot} selectedPartId={props.selectedPartId} />,
});
