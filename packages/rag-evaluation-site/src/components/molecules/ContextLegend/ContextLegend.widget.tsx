import type { ContextLegendWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { ContextLegend } from "./ContextLegend";

export const contextLegendWidget = defineWidget<ContextLegendWidgetProps>({
	type: "ContextLegend",
	module: "widget.dsl",
	render: (props) => (
		<ContextLegend
			className={props.className}
			items={props.items}
			styles={props.styles}
			size={props.size}
			selectedId={props.selectedId}
		/>
	),
});
