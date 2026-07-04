import type { StackWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { Stack } from "./Stack";

export const stackWidget = defineWidget<StackWidgetProps>({
	type: "Stack",
	module: "ui.dsl",
	render: (props, children) => (
		<Stack className={props.className} gap={props.gap} align={props.align}>
			{children}
		</Stack>
	),
});
