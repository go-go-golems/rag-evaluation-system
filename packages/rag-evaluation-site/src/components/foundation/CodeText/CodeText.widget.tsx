import type { CodeTextWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { CodeText } from "./CodeText";

export const codeTextWidget = defineWidget<CodeTextWidgetProps>({
	type: "CodeText",
	module: "ui.dsl",
	render: (props, children) => (
		<CodeText
			className={props.className}
			as={props.as}
			tone={props.tone}
			display={props.display}
			copyable={props.copyable}
		>
			{children}
		</CodeText>
	),
});
