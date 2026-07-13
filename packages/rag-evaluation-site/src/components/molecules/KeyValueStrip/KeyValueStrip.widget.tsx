import type { KeyValueStripWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { KeyValueStrip } from "./KeyValueStrip";

export const keyValueStripWidget = defineWidget<KeyValueStripWidgetProps>({
	type: "KeyValueStrip",
	module: "widget.dsl",
	render: (props, _children, ctx) => (
		<KeyValueStrip
			className={props.className}
			items={props.items.map((item) => ({
				key: ctx.renderValue(item.key),
				value: ctx.renderValue(item.value),
			}))}
		/>
	),
});
