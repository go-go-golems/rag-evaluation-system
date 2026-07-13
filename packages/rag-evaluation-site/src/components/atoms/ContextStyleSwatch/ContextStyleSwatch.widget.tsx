import type { ContextStyleSwatchWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { ContextStyleSwatch } from "./ContextStyleSwatch";

export const contextStyleSwatchWidget = defineWidget<ContextStyleSwatchWidgetProps>({
	type: "ContextStyleSwatch",
	module: "widget.dsl",
	render: (props) => (
		<ContextStyleSwatch
			className={props.className}
			visualStyle={props.visualStyle}
			size={props.size}
			selected={props.selected}
		/>
	),
});
