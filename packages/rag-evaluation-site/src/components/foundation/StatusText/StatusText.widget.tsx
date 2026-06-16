import { StatusText } from "./StatusText";
import { defineWidget } from "../../../widgets/registry";
import type { StatusTextWidgetProps } from "../../../widgets/ir";

export const statusTextWidget = defineWidget<StatusTextWidgetProps>({
	type: "StatusText",
	module: "ui.dsl",
	render: (props, children) => (
		<StatusText className={props.className} status={props.status} icon={props.icon}>
			{children}
		</StatusText>
	),
});
