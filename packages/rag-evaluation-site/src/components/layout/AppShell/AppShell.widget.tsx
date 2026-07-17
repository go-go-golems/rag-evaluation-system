import type { AppShellWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { AppShell } from "./AppShell";

export const appShellWidget = defineWidget<AppShellWidgetProps>({
	type: "AppShell",
	module: "widget.dsl",
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
