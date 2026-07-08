import type { HTMLAttributes } from "react";
import type { ContextStyleSet } from "../../../context";
import { type CalendarEvent, eventStyleSet } from "../../../scheduling";
import { DateTile } from "../../atoms";
import { Panel, Stack } from "../../layout";
import { type TimeGridBlock, TimeGrid } from "../../molecules";

export interface CalendarWeekPanelProps extends Omit<HTMLAttributes<HTMLDivElement>, "onSelect"> {
	days: string[];
	events: CalendarEvent[];
	styleSet?: ContextStyleSet;
	hourStart?: number;
	hourEnd?: number;
	nowISO?: string;
	selectedEventId?: string;
	title?: string;
	onEventSelect?: (eventId: string) => void;
	onSlotCreate?: (slot: { dayISO: string; hour: number }) => void;
}

function toBlocks(events: CalendarEvent[]): TimeGridBlock[] {
	return events.map((event) => ({
		id: event.id,
		dayISO: event.startISO.slice(0, 10),
		startISO: event.startISO,
		endISO: event.endISO,
		styleKey: event.colorKey,
		label: event.title,
		allDay: event.allDay,
	}));
}

export function CalendarWeekPanel({
	days,
	events,
	styleSet = eventStyleSet,
	hourStart = 8,
	hourEnd = 18,
	nowISO,
	selectedEventId,
	title = "Week",
	onEventSelect,
	onSlotCreate,
	className,
	...rest
}: CalendarWeekPanelProps) {
	return (
		<div className={className} data-rag-organism="CalendarWeekPanel" {...rest}>
			<Panel title={title} density="condensed">
				<Stack gap="sm">
					<TimeGrid
						days={days.map((dayISO) => ({
							dayISO,
							header: <DateTile dateISO={dayISO} size="sm" />,
						}))}
						blocks={toBlocks(events)}
						styleSet={styleSet}
						hourStart={hourStart}
						hourEnd={hourEnd}
						nowISO={nowISO}
						selectedBlockId={selectedEventId}
						onBlockSelect={onEventSelect}
						onSlotCreate={onSlotCreate}
						style={{ maxHeight: 480 }}
					/>
				</Stack>
			</Panel>
		</div>
	);
}
