import { AnchoredCommentRail } from "./AnchoredCommentRail";
import { defineWidget } from "../../../widgets/registry";
import type { AnchoredCommentRailWidgetProps } from "../../../widgets/ir";

export const anchoredCommentRailWidget = defineWidget<AnchoredCommentRailWidgetProps>({
	type: "AnchoredCommentRail",
	module: "context_window.dsl",
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
