import { ContextLegend } from "./ContextLegend";
import { defineWidget } from "../../../widgets/registry";
import type { ContextLegendWidgetProps } from "../../../widgets/ir";

export const contextLegendWidget = defineWidget<ContextLegendWidgetProps>({
	type: "ContextLegend",
	module: "context_window.dsl",
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
