import { ContextStackDiagram } from './ContextStackDiagram';
import { defineWidget } from '../../../widgets/registry';
import type { ContextStackDiagramWidgetProps } from '../../../widgets/ir';

export const contextStackDiagramWidget = defineWidget<ContextStackDiagramWidgetProps>({
  type: 'ContextStackDiagram',
  module: 'context_window.dsl',
  render: (props) => <ContextStackDiagram className={props.className} snapshot={props.snapshot} selectedPartId={props.selectedPartId} />,
});
