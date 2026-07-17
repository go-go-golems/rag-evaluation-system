import type { IconButtonWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { IconButton } from "./IconButton";

export const iconButtonWidget = defineWidget<IconButtonWidgetProps>({
	type: "IconButton",
	module: "widget.dsl",
	render: (props, children, ctx) => (
		<IconButton
			className={props.className}
			size={props.size}
			variant={props.variant}
			label={props.label}
			disabled={props.disabled}
			onClick={ctx.bindAction(props.action, { componentType: "IconButton" })}
		>
			{children}
		</IconButton>
	),
});
