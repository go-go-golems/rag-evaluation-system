import type { StepListWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { StepList } from "./StepList";

export const stepListWidget = defineWidget<StepListWidgetProps>({
	type: "StepList",
	module: "ui.dsl",
	render: (props, _children, ctx) => (
		<StepList
			className={props.className}
			items={props.items.map((item) => ({
				...item,
				index: ctx.renderValue(item.index),
				title: ctx.renderValue(item.title),
				description: ctx.renderValue(item.description),
				meta: ctx.renderValue(item.meta),
			}))}
			activeItemId={props.activeItemId}
			onItemSelect={
				props.onItemSelectAction
					? (itemId) =>
							ctx.dispatchAction(props.onItemSelectAction!, {
								itemId,
								value: itemId,
								componentType: "StepList",
							})
					: undefined
			}
		/>
	),
});
