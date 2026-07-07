import type { StatTileWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { StatTile } from "./StatTile";

export const statTileWidget = defineWidget<StatTileWidgetProps>({
	type: "StatTile",
	module: "data.dsl",
	render: (props, _children, ctx) => (
		<StatTile
			label={ctx.renderValue(props.label)}
			value={ctx.renderValue(props.value)}
			delta={props.delta}
			deltaLabel={props.deltaLabel != null ? ctx.renderValue(props.deltaLabel) : undefined}
			trend={props.trend}
			progress={props.progress}
			tone={props.tone}
			onClick={
				props.onAction
					? () => ctx.dispatchAction(props.onAction!, { componentType: "StatTile" })
					: undefined
			}
		/>
	),
});
