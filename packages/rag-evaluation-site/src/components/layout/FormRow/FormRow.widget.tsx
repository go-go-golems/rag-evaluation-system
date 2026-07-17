import type { FormRowWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { FormRow } from "./FormRow";

export const formRowWidget = defineWidget<FormRowWidgetProps>({
	type: "FormRow",
	module: "widget.dsl",
	render: (props, _children, ctx) => (
		<FormRow
			className={props.className}
			label={ctx.renderValue(props.label)}
			control={ctx.renderNode(props.control)}
			description={ctx.renderValue(props.description)}
			hint={ctx.renderValue(props.hint)}
			error={ctx.renderValue(props.error)}
			success={ctx.renderValue(props.success)}
			counter={ctx.renderValue(props.counter)}
			required={props.required}
			orientation={props.orientation}
		/>
	),
});
