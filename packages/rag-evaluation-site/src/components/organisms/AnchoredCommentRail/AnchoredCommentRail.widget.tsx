import type { AnchoredCommentRailWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { AnchoredCommentRail } from "./AnchoredCommentRail";

export const anchoredCommentRailWidget = defineWidget<AnchoredCommentRailWidgetProps>({
	type: "AnchoredCommentRail",
	module: "widget.dsl",
	render: (props, _children, ctx) => (
		<AnchoredCommentRail
			className={props.className}
			title={props.title}
			comments={props.comments}
			selectedCommentId={props.selectedCommentId}
			onCommentSelect={
				props.onCommentSelectAction
					? (commentId) =>
							ctx.dispatchAction(props.onCommentSelectAction!, {
								commentId,
								value: commentId,
								componentType: "AnchoredCommentRail",
							})
					: undefined
			}
		/>
	),
});
