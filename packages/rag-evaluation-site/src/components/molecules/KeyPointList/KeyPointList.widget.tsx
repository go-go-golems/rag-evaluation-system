import type { RenderableValue } from '../../../widgets/ir';
import { KeyPointList } from './KeyPointList';
import { defineWidget } from '../../../widgets/registry';
import type { KeyPointListWidgetProps } from '../../../widgets/ir';

interface KeyPointObjectItem {
  id?: string;
  index?: RenderableValue;
  title?: RenderableValue;
  text: RenderableValue;
  meta?: RenderableValue;
}

export const keyPointListWidget = defineWidget<KeyPointListWidgetProps>({
  type: 'KeyPointList',
  module: 'ui.dsl',
  render: (props, _children, ctx) => (
    <KeyPointList
      className={props.className}
      markerTone={props.markerTone}
      items={props.items.map((item) => (isKeyPointObjectItem(item)
        ? {
            ...item,
            index: ctx.renderValue(item.index),
            title: ctx.renderValue(item.title),
            text: ctx.renderValue(item.text),
            meta: ctx.renderValue(item.meta),
          }
        : ctx.renderValue(item)))}
    />
  ),
});

function isKeyPointObjectItem(item: KeyPointListWidgetProps['items'][number]): item is KeyPointObjectItem {
  return typeof item === 'object' && item !== null && !Array.isArray(item) && 'text' in item;
}
