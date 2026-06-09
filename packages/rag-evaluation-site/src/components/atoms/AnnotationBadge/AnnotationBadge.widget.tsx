import { AnnotationBadge } from './AnnotationBadge';
import { defineWidget } from '../../../widgets/registry';
import type { AnnotationBadgeWidgetProps } from '../../../widgets/ir';

export const annotationBadgeWidget = defineWidget<AnnotationBadgeWidgetProps>({
  type: 'AnnotationBadge',
  module: 'context_window.dsl',
  render: (props) => <AnnotationBadge className={props.className} kind={props.kind} label={props.label} selected={props.selected} />,
});
