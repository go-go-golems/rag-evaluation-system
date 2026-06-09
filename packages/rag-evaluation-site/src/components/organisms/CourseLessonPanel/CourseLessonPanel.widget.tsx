import { CourseLessonPanel } from './CourseLessonPanel';
import { defineWidget } from '../../../widgets/registry';
import type { CourseLessonPanelWidgetProps } from '../../../widgets/ir';

export const courseLessonPanelWidget = defineWidget<CourseLessonPanelWidgetProps>({
  type: 'CourseLessonPanel',
  module: 'course.dsl',
  render: (props, _children, ctx) => (
    <CourseLessonPanel
      className={props.className}
      course={props.course}
      activeAgendaItemId={props.activeAgendaItemId}
      onAgendaItemSelect={props.onAgendaItemSelectAction
        ? (agendaItemId) => ctx.dispatchAction(props.onAgendaItemSelectAction!, { agendaItemId, value: agendaItemId, componentType: 'CourseLessonPanel' })
        : undefined}
    />
  ),
});
