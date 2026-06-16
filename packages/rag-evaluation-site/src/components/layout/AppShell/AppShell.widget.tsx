import { AppShell } from "./AppShell";
import { defineWidget } from "../../../widgets/registry";
import type { AppShellWidgetProps } from "../../../widgets/ir";

export const appShellWidget = defineWidget<AppShellWidgetProps>({
	type: "AppShell",
	module: "ui.dsl",
	render: (props, children, ctx) => (
		<AppShell
			className={props.className}
			header={props.header ? ctx.renderNode(props.header) : undefined}
			sidebar={props.sidebar ? ctx.renderNode(props.sidebar) : undefined}
		>
			{children}
		</AppShell>
	),
});
