import { Divider } from './Divider';
import { defineWidget } from '../../../widgets/registry';
import type { DividerWidgetProps } from '../../../widgets/ir';

export const dividerWidget = defineWidget<DividerWidgetProps>({
  type: 'Divider',
  module: 'ui.dsl',
  render: (props) => <Divider className={props.className} orientation={props.orientation} />,
});
