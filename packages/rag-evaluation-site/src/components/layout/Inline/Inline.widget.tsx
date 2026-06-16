import { Inline } from "./Inline";
import { defineWidget } from "../../../widgets/registry";
import type { InlineWidgetProps } from "../../../widgets/ir";

export const inlineWidget = defineWidget<InlineWidgetProps>({
	type: "Inline",
	module: "ui.dsl",
	render: (props, children) => (
		<Inline className={props.className} gap={props.gap} justify={props.justify} wrap={props.wrap}>
			{children}
		</Inline>
	),
});
