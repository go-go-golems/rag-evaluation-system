import { Stack } from "./Stack";
import { defineWidget } from "../../../widgets/registry";
import type { StackWidgetProps } from "../../../widgets/ir";

export const stackWidget = defineWidget<StackWidgetProps>({
	type: "Stack",
	module: "ui.dsl",
	render: (props, children) => (
		<Stack className={props.className} gap={props.gap} align={props.align}>
			{children}
		</Stack>
	),
});
