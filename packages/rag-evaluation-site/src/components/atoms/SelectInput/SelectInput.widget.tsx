import { SelectInput } from './SelectInput';
import { defineWidget } from '../../../widgets/registry';
import type { SelectInputWidgetProps } from '../../../widgets/ir';

export const selectInputWidget = defineWidget<SelectInputWidgetProps>({
  type: 'SelectInput',
  module: 'ui.dsl',
  render: (props, children, ctx) => (
    <SelectInput className={props.className} name={props.name} value={props.value} disabled={props.disabled}>
      {(props.options ?? []).map((option) => (
        <option key={String(option.value)} value={option.value} disabled={option.disabled}>
          {ctx.renderValue(option.label)}
        </option>
      ))}
      {children}
    </SelectInput>
  ),
});
