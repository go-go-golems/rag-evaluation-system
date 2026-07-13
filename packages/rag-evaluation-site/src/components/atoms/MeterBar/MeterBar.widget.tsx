import type { MeterBarWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { MeterBar } from "./MeterBar";

export const meterBarWidget = defineWidget<MeterBarWidgetProps>({
	type: "MeterBar",
	module: "widget.dsl",
	render: (props, _children, ctx) => (
		<MeterBar
			className={props.className}
			value={props.value}
			tone={props.tone}
			label={props.label != null ? ctx.renderValue(props.label) : undefined}
		/>
	),
});
