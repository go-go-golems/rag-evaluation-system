import { FormRow } from './FormRow';
import { defineWidget } from '../../../widgets/registry';
import type { FormRowWidgetProps } from '../../../widgets/ir';

export const formRowWidget = defineWidget<FormRowWidgetProps>({
  type: 'FormRow',
  module: 'ui.dsl',
  render: (props, _children, ctx) => (
    <FormRow
      className={props.className}
      label={ctx.renderValue(props.label)}
      control={ctx.renderNode(props.control)}
      hint={ctx.renderValue(props.hint)}
    />
  ),
});
