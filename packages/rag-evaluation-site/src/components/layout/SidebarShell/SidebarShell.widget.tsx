import type { SidebarShellWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { SidebarShell } from "./SidebarShell";

export const sidebarShellWidget = defineWidget<SidebarShellWidgetProps>({
	type: "SidebarShell",
	module: "ui.dsl",
	render: (props, children, ctx) => (
		<SidebarShell
			className={props.className}
			sidebar={props.sidebar ? ctx.renderNode(props.sidebar) : undefined}
			sidebarWidth={props.sidebarWidth}
			contentPadding={props.contentPadding}
			header={props.header ? ctx.renderNode(props.header) : undefined}
			footer={props.footer ? ctx.renderNode(props.footer) : undefined}
		>
			{children}
		</SidebarShell>
	),
});
