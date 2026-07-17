import type { MonthGridWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { MonthGrid } from "./MonthGrid";

export const monthGridWidget = defineWidget<MonthGridWidgetProps>({
	type: "MonthGrid",
	module: "widget.dsl",
	render: (props, _children, ctx) => (
		<MonthGrid
			className={props.className}
			monthISO={props.monthISO}
			markers={
				props.markers
					? Object.fromEntries(
							Object.entries(props.markers).map(([date, m]) => [
								date,
								{
									count: m.count,
									styleKey: m.styleKey,
									label: m.label != null ? ctx.renderValue(m.label) : undefined,
								},
							]),
						)
					: undefined
			}
			styleSet={props.styleSet}
			selectedDateISO={props.selectedDateISO}
			todayISO={props.todayISO}
			minDateISO={props.minDateISO}
			maxDateISO={props.maxDateISO}
			weekStartsOn={props.weekStartsOn}
			showHeader={props.showHeader}
			onDaySelect={
				props.onDaySelectAction
					? (dateISO) =>
							ctx.dispatchAction(props.onDaySelectAction!, {
								dateISO,
								value: dateISO,
								componentType: "MonthGrid",
							})
					: undefined
			}
			onMonthChange={
				props.onMonthChangeAction
					? (monthISO) =>
							ctx.dispatchAction(props.onMonthChangeAction!, {
								monthISO,
								value: monthISO,
								componentType: "MonthGrid",
							})
					: undefined
			}
		/>
	),
});
