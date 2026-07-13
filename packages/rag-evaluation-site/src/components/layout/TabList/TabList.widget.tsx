import type { TabListWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { TabList } from "./TabList";

export const tabListWidget = defineWidget<TabListWidgetProps>({
	type: "TabList",
	module: "widget.dsl",
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
