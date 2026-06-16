import { KeyValueStrip } from "./KeyValueStrip";
import { defineWidget } from "../../../widgets/registry";
import type { KeyValueStripWidgetProps } from "../../../widgets/ir";

export const keyValueStripWidget = defineWidget<KeyValueStripWidgetProps>({
	type: "KeyValueStrip",
	module: "ui.dsl",
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
