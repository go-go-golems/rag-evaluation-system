import type { TagWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { Tag } from "./Tag";

export const tagWidget = defineWidget<TagWidgetProps>({
	type: "Tag",
	module: "widget.dsl",
	render: (props, _children, ctx) => {
		const onRemoveAction = props.onRemoveAction;
		return (
			<Tag
				className={props.className}
				label={props.label}
				selected={props.selected}
				disabled={props.disabled}
				onRemove={
					onRemoveAction
						? () =>
								ctx.dispatchAction(onRemoveAction, {
									tag: props.label,
									value: props.label,
									componentType: "Tag",
								})
						: undefined
				}
			/>
		);
	},
});
