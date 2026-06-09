import type { RenderableValue } from '../../../widgets/ir';
import { CheckList } from './CheckList';
import { defineWidget } from '../../../widgets/registry';
import type { CheckListWidgetProps } from '../../../widgets/ir';

interface CheckListObjectItem {
  id?: string;
  text: RenderableValue;
}

export const checkListWidget = defineWidget<CheckListWidgetProps>({
  type: 'CheckList',
  module: 'ui.dsl',
  render: (props, _children, ctx) => (
    <CheckList
      className={props.className}
      marker={ctx.renderValue(props.marker)}
      items={props.items.map((item) => (isCheckListObjectItem(item)
        ? { ...item, text: ctx.renderValue(item.text) }
        : ctx.renderValue(item)))}
    />
  ),
});

function isCheckListObjectItem(item: CheckListWidgetProps['items'][number]): item is CheckListObjectItem {
  return typeof item === 'object' && item !== null && !Array.isArray(item) && 'text' in item;
}
