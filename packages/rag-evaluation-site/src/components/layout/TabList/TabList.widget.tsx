import { TabList } from "./TabList";
import { defineWidget } from "../../../widgets/registry";
import type { TabListWidgetProps } from "../../../widgets/ir";

export const tabListWidget = defineWidget<TabListWidgetProps>({
	type: "TabList",
	module: "ui.dsl",
	render: (props, _children, ctx) => (
		<TabList
			items={props.items.map((item) => ({ id: item.id, label: ctx.renderValue(item.label) }))}
			activeId={props.activeId}
			ariaLabel={props.ariaLabel}
			onChange={(id) => {
				if (props.onChange)
					ctx.dispatchAction(props.onChange, { value: id, componentType: "TabList" });
			}}
		/>
	),
});
