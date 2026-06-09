import { TranscriptSessionHeader } from './TranscriptSessionHeader';
import { defineWidget } from '../../../widgets/registry';
import type { TranscriptSessionHeaderWidgetProps } from '../../../widgets/ir';

export const transcriptSessionHeaderWidget = defineWidget<TranscriptSessionHeaderWidgetProps>({
  type: 'TranscriptSessionHeader',
  module: 'context_window.dsl',
  render: (props, _children, ctx) => (
    <TranscriptSessionHeader
      className={props.className}
      title={ctx.renderValue(props.title)}
      subtitle={ctx.renderValue(props.subtitle)}
      messageCount={props.messageCount}
      annotationCount={props.annotationCount}
      tokenTotal={props.tokenTotal}
      rightSlot={props.rightSlot ? ctx.renderNode(props.rightSlot) : undefined}
    />
  ),
});
