import { ContextKindSwatch } from './ContextKindSwatch';
import { defineWidget } from '../../../widgets/registry';
import type { ContextKindSwatchWidgetProps } from '../../../widgets/ir';

export const contextKindSwatchWidget = defineWidget<ContextKindSwatchWidgetProps>({
  type: 'ContextKindSwatch',
  module: 'context_window.dsl',
  render: (props) => <ContextKindSwatch className={props.className} kind={props.kind} mode={props.mode} size={props.size} selected={props.selected} />,
});
