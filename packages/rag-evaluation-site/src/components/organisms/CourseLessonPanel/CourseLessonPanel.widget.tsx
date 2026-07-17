import type { CourseLessonPanelWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { CourseLessonPanel } from "./CourseLessonPanel";

export const courseLessonPanelWidget = defineWidget<CourseLessonPanelWidgetProps>({
	type: "CourseLessonPanel",
	module: "widget.dsl",
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
