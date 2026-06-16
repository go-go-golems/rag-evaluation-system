import { ContextStyleSwatch } from "./ContextStyleSwatch";
import { defineWidget } from "../../../widgets/registry";
import type { ContextStyleSwatchWidgetProps } from "../../../widgets/ir";

export const contextStyleSwatchWidget = defineWidget<ContextStyleSwatchWidgetProps>({
	type: "ContextStyleSwatch",
	module: "context_window.dsl",
	render: (props) => (
		<ContextStyleSwatch
			className={props.className}
			visualStyle={props.visualStyle}
			size={props.size}
			selected={props.selected}
		/>
	),
});
