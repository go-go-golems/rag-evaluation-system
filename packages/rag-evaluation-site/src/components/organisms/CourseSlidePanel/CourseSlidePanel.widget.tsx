import type { CourseSlidePanelWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { CourseSlidePanel } from "./CourseSlidePanel";

export const courseSlidePanelWidget = defineWidget<CourseSlidePanelWidgetProps>({
	type: "CourseSlidePanel",
	module: "widget.dsl",
	render: (props, _children, ctx) => (
		<CourseSlidePanel
			className={props.className}
			slide={props.slide}
			snapshot={props.snapshot}
			styleSet={props.styleSet}
			index={props.index}
			total={props.total}
			visualSide={props.visualSide}
			mode={props.mode}
			onPrevious={
				props.onPreviousAction
					? () =>
							ctx.dispatchAction(props.onPreviousAction!, {
								componentType: "CourseSlidePanel",
								value: "previous",
							})
					: undefined
			}
			onNext={
				props.onNextAction
					? () =>
							ctx.dispatchAction(props.onNextAction!, {
								componentType: "CourseSlidePanel",
								value: "next",
							})
					: undefined
			}
			onPresent={
				props.onPresentAction
					? () =>
							ctx.dispatchAction(props.onPresentAction!, {
								componentType: "CourseSlidePanel",
								value: "present",
							})
					: undefined
			}
			onFullscreen={
				props.onFullscreenAction
					? () =>
							ctx.dispatchAction(props.onFullscreenAction!, {
								componentType: "CourseSlidePanel",
								value: "fullscreen",
							})
					: undefined
			}
		/>
	),
});
