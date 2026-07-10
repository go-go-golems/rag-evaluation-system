import type { SegmentedBarWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { SegmentedBar } from "./SegmentedBar";

export const segmentedBarWidget = defineWidget<SegmentedBarWidgetProps>({
	type: "SegmentedBar",
	module: "ui.dsl",
	render: (props, _children, ctx) => (
		<SegmentedBar
			className={props.className}
			segments={props.segments.map((s) => ({
				value: s.value,
				styleKey: s.styleKey,
				label: s.label != null ? ctx.renderValue(s.label) : undefined,
			}))}
			styleSet={props.styleSet}
			total={props.total}
			showCounts={props.showCounts}
			markers={props.markers?.map((m) => ({
				at: m.at,
				styleKey: m.styleKey,
				label: m.label != null ? ctx.renderValue(m.label) : undefined,
			}))}
			size={props.size}
			onSegmentSelect={
				props.onSegmentAction
					? (styleKey, index) =>
							ctx.dispatchAction(props.onSegmentAction!, {
								styleKey,
								index,
								value: styleKey,
								componentType: "SegmentedBar",
							})
					: undefined
			}
		/>
	),
});
