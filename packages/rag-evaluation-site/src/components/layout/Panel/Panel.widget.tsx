import { Panel } from './Panel';
import { defineWidget } from '../../../widgets/registry';
import type { PanelWidgetProps } from '../../../widgets/ir';

export const panelWidget = defineWidget<PanelWidgetProps>({
  type: 'Panel',
  module: 'ui.dsl',
  render: (props, children, ctx) => (
    <Panel
      className={props.className}
      title={ctx.renderValue(props.title)}
      actions={ctx.renderValue(props.actions)}
      density={props.density}
      fill={props.fill}
    >
      {children}
    </Panel>
  ),
});
