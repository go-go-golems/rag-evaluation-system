import { SplitPane } from './SplitPane';
import { defineWidget } from '../../../widgets/registry';
import type { SplitPaneWidgetProps } from '../../../widgets/ir';

export const splitPaneWidget = defineWidget<SplitPaneWidgetProps>({
  type: 'SplitPane',
  module: 'ui.dsl',
  render: (props, _children, ctx) => (
    <SplitPane
      className={props.className}
      left={ctx.renderNode(props.left)}
      right={ctx.renderNode(props.right)}
      ratio={props.ratio}
      divider={props.divider}
      gutter={props.gutter}
    />
  ),
});
