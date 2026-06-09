import { TextInput } from './TextInput';
import { defineWidget } from '../../../widgets/registry';
import type { TextInputWidgetProps } from '../../../widgets/ir';

export const textInputWidget = defineWidget<TextInputWidgetProps>({
  type: 'TextInput',
  module: 'ui.dsl',
  render: (props) => (
    <TextInput
      className={props.className}
      name={props.name}
      value={props.value}
      placeholder={props.placeholder}
      type={props.type}
      disabled={props.disabled}
      min={props.min}
      max={props.max}
      readOnly
    />
  ),
});
