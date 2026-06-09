import { ContextStripDiagram } from './ContextStripDiagram';
import { defineWidget } from '../../../widgets/registry';
import type { ContextStripDiagramWidgetProps } from '../../../widgets/ir';

export const contextStripDiagramWidget = defineWidget<ContextStripDiagramWidgetProps>({
  type: 'ContextStripDiagram',
  module: 'context_window.dsl',
  render: (props) => <ContextStripDiagram className={props.className} snapshot={props.snapshot} mode={props.mode} selectedPartId={props.selectedPartId} showLabels={props.showLabels} />,
});
