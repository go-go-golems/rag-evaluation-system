import type { FigureBlockWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { FigureBlock } from "./FigureBlock";

export const figureBlockWidget = defineWidget<FigureBlockWidgetProps>({
	type: "FigureBlock",
	module: "ui.dsl",
	render: (props, children, ctx) => (
		<FigureBlock
			className={props.className}
			as={props.as}
			caption={ctx.renderValue(props.caption)}
			legend={props.legend ? ctx.renderNode(props.legend) : undefined}
			frame={props.frame}
		>
			{children}
		</FigureBlock>
	),
});
