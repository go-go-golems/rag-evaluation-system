import type { CourseStudioShellWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { CourseStudioShell } from "./CourseStudioShell";

export const courseStudioShellWidget = defineWidget<CourseStudioShellWidgetProps>({
	type: "CourseStudioShell",
	module: "widget.dsl",
	render: (props, children, ctx) => (
		<CourseStudioShell
			className={props.className}
			sections={props.sections.map((section) => ({
				...section,
				label: ctx.renderValue(section.label),
				items: section.items.map((item) => ({
					...item,
					label: ctx.renderValue(item.label),
					icon: ctx.renderValue(item.icon),
					badge: ctx.renderValue(item.badge),
				})),
			}))}
			activeItemId={props.activeItemId}
			onNavigate={
				props.onNavigateAction
					? (itemId) =>
							ctx.dispatchAction(props.onNavigateAction!, {
								itemId,
								item: { id: itemId },
								value: itemId,
								componentType: "CourseStudioShell",
							})
					: undefined
			}
			title={ctx.renderValue(props.title)}
			subtitle={ctx.renderValue(props.subtitle)}
			sidebarFooter={props.sidebarFooter ? ctx.renderNode(props.sidebarFooter) : undefined}
			contentPadding={props.contentPadding}
		>
			{children}
		</CourseStudioShell>
	),
});
