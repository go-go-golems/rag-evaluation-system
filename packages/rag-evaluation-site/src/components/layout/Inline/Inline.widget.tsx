import type { InlineWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { Inline } from "./Inline";

export const inlineWidget = defineWidget<InlineWidgetProps>({
	type: "Inline",
	module: "widget.dsl",
	render: (props, children) => (
		<Inline className={props.className} gap={props.gap} justify={props.justify} wrap={props.wrap}>
			{children}
		</Inline>
	),
});
