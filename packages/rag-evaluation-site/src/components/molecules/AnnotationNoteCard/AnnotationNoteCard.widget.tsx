import { AnnotationNoteCard } from './AnnotationNoteCard';
import { defineWidget } from '../../../widgets/registry';
import type { AnnotationNoteCardWidgetProps } from '../../../widgets/ir';

export const annotationNoteCardWidget = defineWidget<AnnotationNoteCardWidgetProps>({
  type: 'AnnotationNoteCard',
  module: 'context_window.dsl',
  render: (props) => <AnnotationNoteCard className={props.className} annotation={props.annotation} selected={props.selected} index={props.index} />,
});
