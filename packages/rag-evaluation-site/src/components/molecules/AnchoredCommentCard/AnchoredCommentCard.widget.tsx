import type { AnchoredCommentCardWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { AnchoredCommentCard } from "./AnchoredCommentCard";

export const anchoredCommentCardWidget = defineWidget<AnchoredCommentCardWidgetProps>({
	type: "AnchoredCommentCard",
	module: "context_window.dsl",
	render: (props, _children, ctx) => (
		<AnchoredCommentCard
			className={props.className}
			comment={props.comment}
			index={props.index}
			selected={props.selected}
			compact={props.compact}
			onDismiss={
				props.onDismissAction
					? () =>
							ctx.dispatchAction(props.onDismissAction!, {
								commentId: props.comment.id,
								value: props.comment.id,
								componentType: "AnchoredCommentCard",
							})
					: undefined
			}
		/>
	),
});
