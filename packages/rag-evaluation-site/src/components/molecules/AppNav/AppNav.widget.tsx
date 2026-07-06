import type { AppNavWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { AppNav } from "./AppNav";

export const appNavWidget = defineWidget<AppNavWidgetProps>({
	type: "AppNav",
	module: "ui.dsl",
	render: (props, _children, ctx) => (
		<AppNav
			brand={ctx.renderValue(props.brand)}
			items={props.items.map((item) => ({ id: item.id, label: ctx.renderValue(item.label) }))}
			activeItemId={props.activeItemId}
			onItemSelect={(itemId) => {
				const item = props.items.find((candidate) => candidate.id === itemId);
				if (item?.action)
					ctx.dispatchAction(item.action, { value: itemId, componentType: "AppNav" });
			}}
		/>
	),
});
