import { MetadataGrid } from './MetadataGrid';
import { defineWidget } from '../../../widgets/registry';
import type { MetadataGridWidgetProps } from '../../../widgets/ir';

export const metadataGridWidget = defineWidget<MetadataGridWidgetProps>({
  type: 'MetadataGrid',
  module: 'ui.dsl',
  render: (props, _children, ctx) => (
    <MetadataGrid
      className={props.className}
      density={props.density}
      items={props.items.map((item) => ({
        key: ctx.renderValue(item.key),
        value: ctx.renderValue(item.value),
        copyValue: item.copyValue,
      }))}
    />
  ),
});
