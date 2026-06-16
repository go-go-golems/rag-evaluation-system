import { CourseLessonPanel } from "./CourseLessonPanel";
import { defineWidget } from "../../../widgets/registry";
import type { CourseLessonPanelWidgetProps } from "../../../widgets/ir";

export const courseLessonPanelWidget = defineWidget<CourseLessonPanelWidgetProps>({
	type: "CourseLessonPanel",
	module: "course.dsl",
	render: (props, _children, ctx) => (
		<CourseLessonPanel
			className={props.className}
			course={props.course}
			activeAgendaItemId={props.activeAgendaItemId}
			onAgendaItemSelect={
				props.onAgendaItemSelectAction
					? (agendaItemId) =>
							ctx.dispatchAction(props.onAgendaItemSelectAction!, {
								agendaItemId,
								value: agendaItemId,
								componentType: "CourseLessonPanel",
							})
					: undefined
			}
			onPrimaryCta={
				props.onPrimaryCtaAction
					? () =>
							ctx.dispatchAction(props.onPrimaryCtaAction!, {
								componentType: "CourseLessonPanel",
								cta: "primary",
							})
					: undefined
			}
			onSecondaryCta={
				props.onSecondaryCtaAction
					? () =>
							ctx.dispatchAction(props.onSecondaryCtaAction!, {
								componentType: "CourseLessonPanel",
								cta: "secondary",
							})
					: undefined
			}
		/>
	),
});
