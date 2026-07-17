import type { ButtonWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { Button } from "./Button";

export const buttonWidget = defineWidget<ButtonWidgetProps>({
	type: "Button",
	module: "widget.dsl",
	render: (props, children, ctx) => (
		<Button
			className={props.className}
			variant={props.variant}
			size={props.size}
			selected={props.selected}
			disabled={props.disabled}
			type={props.type ?? "button"}
			onClick={ctx.bindAction(props.action, { componentType: "Button" })}
		>
			{children}
		</Button>
	),
});
