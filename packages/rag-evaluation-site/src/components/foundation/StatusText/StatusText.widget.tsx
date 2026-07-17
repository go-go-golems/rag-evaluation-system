import type { StatusTextWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { StatusText } from "./StatusText";

export const statusTextWidget = defineWidget<StatusTextWidgetProps>({
	type: "StatusText",
	module: "widget.dsl",
	render: (props, children) => (
		<StatusText className={props.className} status={props.status} icon={props.icon}>
			{children}
		</StatusText>
	),
});
