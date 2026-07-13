import type { CSSProperties } from "react";
import type { TimeGridWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { TimeGrid } from "./TimeGrid";

export const timeGridWidget = defineWidget<TimeGridWidgetProps>({
	type: "TimeGrid",
	module: "widget.dsl",
	render: (props, _children, ctx) => (
		<TimeGrid
			className={props.className}
			style={props.style as CSSProperties | undefined}
			days={props.days.map((d) =>
				typeof d === "string"
					? d
					: { dayISO: d.dayISO, header: d.header != null ? ctx.renderValue(d.header) : undefined },
			)}
			blocks={props.blocks.map((b) => ({ ...b, label: ctx.renderValue(b.label) }))}
			styleSet={props.styleSet}
			hourStart={props.hourStart}
			hourEnd={props.hourEnd}
			hourHeight={props.hourHeight}
			nowISO={props.nowISO}
			selectedBlockId={props.selectedBlockId}
			onBlockSelect={
				props.onBlockSelectAction
					? (blockId) =>
							ctx.dispatchAction(props.onBlockSelectAction!, {
								blockId,
								value: blockId,
								componentType: "TimeGrid",
							})
					: undefined
			}
			onSlotCreate={
				props.onSlotCreateAction
					? (slot) =>
							ctx.dispatchAction(props.onSlotCreateAction!, {
								dayISO: slot.dayISO,
								hour: slot.hour,
								value: slot.dayISO,
								componentType: "TimeGrid",
							})
					: undefined
			}
		/>
	),
});
