import type { SidebarNavWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { SidebarNav } from "./SidebarNav";

export const sidebarNavWidget = defineWidget<SidebarNavWidgetProps>({
	type: "SidebarNav",
	module: "ui.dsl",
	render: (props, _children, ctx) => (
		<SidebarNav
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
			ariaLabel={props.ariaLabel}
			onItemSelect={
				props.onItemSelectAction
					? (itemId) =>
							ctx.dispatchAction(props.onItemSelectAction!, {
								itemId,
								value: itemId,
								componentType: "SidebarNav",
							})
					: undefined
			}
		/>
	),
});
