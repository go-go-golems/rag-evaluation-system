import { SlideShell } from './SlideShell';
import { defineWidget } from '../../../widgets/registry';
import type { SlideShellWidgetProps } from '../../../widgets/ir';

export const slideShellWidget = defineWidget<SlideShellWidgetProps>({
  type: 'SlideShell',
  module: 'course.dsl',
  render: (props, _children, ctx) => (
    <SlideShell
      className={props.className}
      as={props.as}
      eyebrow={ctx.renderValue(props.eyebrow)}
      counter={ctx.renderValue(props.counter)}
      title={ctx.renderValue(props.title)}
      subtitle={ctx.renderValue(props.subtitle)}
      primary={ctx.renderNode(props.primary)}
      secondary={props.secondary ? ctx.renderNode(props.secondary) : undefined}
      primarySide={props.primarySide}
      ratio={props.ratio}
      divider={props.divider}
      footer={props.footer ? ctx.renderNode(props.footer) : undefined}
    />
  ),
});
