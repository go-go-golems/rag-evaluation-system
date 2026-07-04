import type { TextWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { Text } from "./Text";

export const textWidget = defineWidget<TextWidgetProps>({
	type: "Text",
	module: "ui.dsl",
	render: (props, children) => (
		<Text
			className={props.className}
			as={props.as}
			size={props.size}
			tone={props.tone}
			weight={props.weight}
			align={props.align}
			truncate={props.truncate}
		>
			{children}
		</Text>
	),
});
