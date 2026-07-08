import type { HTMLAttributes } from "react";
import type { ContextStyleSet } from "../../../context";
import { type CalendarEvent, eventStyleSet } from "../../../scheduling";
import { Caption, Text } from "../../foundation";
import { Inline, Panel, Stack } from "../../layout";
import { type MonthGridDayMarker, MonthGrid } from "../../molecules";

export interface CalendarMonthPanelProps extends Omit<HTMLAttributes<HTMLDivElement>, "onSelect"> {
	monthISO: string;
	events: CalendarEvent[];
	styleSet?: ContextStyleSet;
	selectedDateISO?: string;
	todayISO?: string;
	title?: string;
	onDaySelect?: (dateISO: string) => void;
	onMonthChange?: (monthISO: string) => void;
	onEventSelect?: (eventId: string) => void;
}

function markersFor(events: CalendarEvent[]): Record<string, MonthGridDayMarker> {
	const markers: Record<string, MonthGridDayMarker> = {};
	for (const event of events) {
		const date = event.startISO.slice(0, 10);
		const existing = markers[date];
		markers[date] = {
			count: (existing?.count ?? 0) + 1,
			styleKey: existing?.styleKey ?? event.colorKey,
		};
	}
	return markers;
}

export function CalendarMonthPanel({
	monthISO,
	events,
	styleSet = eventStyleSet,
	selectedDateISO,
	todayISO,
	title = "Calendar",
	onDaySelect,
	onMonthChange,
	onEventSelect,
	className,
	...rest
}: CalendarMonthPanelProps) {
	const dayEvents = selectedDateISO
		? events.filter((e) => e.startISO.slice(0, 10) === selectedDateISO)
		: [];

	return (
		<div className={className} data-rag-organism="CalendarMonthPanel" {...rest}>
			<Panel title={title} density="condensed">
				<Stack gap="sm">
					<MonthGrid
						monthISO={monthISO}
						markers={markersFor(events)}
						styleSet={styleSet}
						selectedDateISO={selectedDateISO}
						todayISO={todayISO}
						onDaySelect={onDaySelect}
						onMonthChange={onMonthChange}
					/>
					{selectedDateISO ? (
						<Stack gap="xs">
							<Caption tone="muted">{selectedDateISO}</Caption>
							{dayEvents.length ? (
								dayEvents.map((event) => (
									<Inline key={event.id} gap="xs">
										<span
											aria-hidden="true"
											style={{
												width: 8,
												height: 8,
												background: styleSet.styles[event.colorKey]?.fill ?? "var(--mac-accent)",
												display: "inline-block",
											}}
										/>
										<button
											type="button"
											onClick={() => onEventSelect?.(event.id)}
											style={{
												border: 0,
												background: "transparent",
												padding: 0,
												cursor: onEventSelect ? "pointer" : "default",
												font: "inherit",
											}}
										>
											<Text size="compact">{`${event.startISO.slice(11, 16)} ${event.title}`}</Text>
										</button>
									</Inline>
								))
							) : (
								<Caption tone="muted">No events</Caption>
							)}
						</Stack>
					) : null}
				</Stack>
			</Panel>
		</div>
	);
}
