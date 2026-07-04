import type { CmsShellWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { CmsShell } from "./CmsShell";

export const cmsShellWidget = defineWidget<CmsShellWidgetProps>({
	type: "CmsShell",
	module: "cms.dsl",
	render: (props, children, ctx) => {
		const onNavigateAction = props.onNavigateAction;
		return (
			<CmsShell
				className={props.className}
				sections={props.sections?.map((section) => ({
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
					onNavigateAction
						? (itemId) =>
								ctx.dispatchAction(onNavigateAction, {
									itemId,
									value: itemId,
									componentType: "CmsShell",
								})
						: undefined
				}
				title={props.title != null ? ctx.renderValue(props.title) : undefined}
				subtitle={props.subtitle != null ? ctx.renderValue(props.subtitle) : undefined}
				headerSlot={props.headerSlot ? ctx.renderNode(props.headerSlot) : undefined}
				sidebarFooter={props.sidebarFooter ? ctx.renderNode(props.sidebarFooter) : undefined}
				contentPadding={props.contentPadding}
			>
				{children}
			</CmsShell>
		);
	},
});
