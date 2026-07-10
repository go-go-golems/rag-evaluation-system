import type { CSSProperties } from "react";
import type { SplitPaneWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { SplitPane } from "./SplitPane";

export const splitPaneWidget = defineWidget<SplitPaneWidgetProps>({
	type: "SplitPane",
	module: "ui.dsl",
	render: (props, _children, ctx) => (
		<SplitPane
			className={props.className}
			style={props.style as CSSProperties | undefined}
			left={ctx.renderNode(props.left)}
			right={ctx.renderNode(props.right)}
			ratio={props.ratio}
			divider={props.divider}
			gutter={props.gutter}
		/>
	),
});
