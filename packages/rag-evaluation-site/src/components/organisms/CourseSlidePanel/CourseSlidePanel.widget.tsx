import { CourseSlidePanel } from './CourseSlidePanel';
import { defineWidget } from '../../../widgets/registry';
import type { CourseSlidePanelWidgetProps } from '../../../widgets/ir';

export const courseSlidePanelWidget = defineWidget<CourseSlidePanelWidgetProps>({
  type: 'CourseSlidePanel',
  module: 'course.dsl',
  render: (props, _children, ctx) => (
    <CourseSlidePanel
      className={props.className}
      slide={props.slide}
      snapshot={props.snapshot}
      index={props.index}
      total={props.total}
      visualSide={props.visualSide}
      onPrevious={props.onPreviousAction
        ? () => ctx.dispatchAction(props.onPreviousAction!, { componentType: 'CourseSlidePanel', value: 'previous' })
        : undefined}
      onNext={props.onNextAction
        ? () => ctx.dispatchAction(props.onNextAction!, { componentType: 'CourseSlidePanel', value: 'next' })
        : undefined}
    />
  ),
});
