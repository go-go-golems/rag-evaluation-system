import type { FormPanelWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { FormPanel } from "./FormPanel";

export const formPanelWidget = defineWidget<FormPanelWidgetProps>({
	type: "FormPanel",
	module: "ui.dsl",
	render: (props, children, ctx) => (
		<FormPanel
			className={props.className}
			title={ctx.renderValue(props.title)}
			subtitle={ctx.renderValue(props.subtitle)}
			actions={ctx.renderValue(props.actions)}
			action={props.formAction}
			method={props.method}
			status={props.status}
			statusMessage={ctx.renderValue(props.statusMessage)}
			submitLabel={ctx.renderValue(props.submitLabel)}
			resetLabel={ctx.renderValue(props.resetLabel)}
			footer={ctx.renderValue(props.footer)}
			disabled={props.disabled}
		>
			{children}
		</FormPanel>
	),
});
