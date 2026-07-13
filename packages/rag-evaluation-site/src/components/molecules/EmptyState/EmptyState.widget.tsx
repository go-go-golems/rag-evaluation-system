import type { EmptyStateWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { EmptyState } from "./EmptyState";

export const emptyStateWidget = defineWidget<EmptyStateWidgetProps>({
	type: "EmptyState",
	module: "widget.dsl",
	render: (props, _children, ctx) => (
		<EmptyState
			className={props.className}
			glyph={props.glyph != null ? ctx.renderValue(props.glyph) : undefined}
			title={ctx.renderValue(props.title)}
			hint={props.hint != null ? ctx.renderValue(props.hint) : undefined}
			action={props.actionSlot ? ctx.renderNode(props.actionSlot) : undefined}
			framed={props.framed}
		/>
	),
});
