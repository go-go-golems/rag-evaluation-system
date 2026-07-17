import type { PanelWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { Panel } from "./Panel";

export const panelWidget = defineWidget<PanelWidgetProps>({
	type: "Panel",
	module: "widget.dsl",
	render: (props, children, ctx) => (
		<Panel
			className={props.className}
			title={ctx.renderValue(props.title)}
			actions={ctx.renderValue(props.actions)}
			density={props.density}
			titleTone={props.titleTone}
			fill={props.fill}
		>
			{children}
		</Panel>
	),
});
