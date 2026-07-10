import type { HTMLAttributes } from "react";
import { type BookableDay, type BookableSlot, type BookingType } from "../../../scheduling";
import { Button } from "../../atoms";
import { Caption, Text } from "../../foundation";
import { Panel, SplitPane, Stack, TileGrid } from "../../layout";
import { type MonthGridDayMarker, KeyValueStrip, MonthGrid, PersonSummary } from "../../molecules";

export interface BookingPagePanelProps extends Omit<HTMLAttributes<HTMLDivElement>, "onSelect"> {
	bookingType: BookingType;
	monthISO: string;
	days: BookableDay[];
	selectedDateISO?: string;
	slots: BookableSlot[];
	selectedSlotId?: string;
	tz?: string;
	onDaySelect?: (dateISO: string) => void;
	onSlotSelect?: (slotId: string) => void;
	onMonthChange?: (monthISO: string) => void;
	onConfirm?: () => void;
}

function markersFor(days: BookableDay[]): Record<string, MonthGridDayMarker> {
	const markers: Record<string, MonthGridDayMarker> = {};
	for (const day of days) {
		if (day.slotCount > 0 && !day.disabled) markers[day.dateISO] = { count: day.slotCount };
	}
	return markers;
}

export function BookingPagePanel({
	bookingType,
	monthISO,
	days,
	selectedDateISO,
	slots,
	selectedSlotId,
	tz,
	onDaySelect,
	onSlotSelect,
	onMonthChange,
	onConfirm,
	className,
	...rest
}: BookingPagePanelProps) {
	const selectedSlot = slots.find((s) => s.id === selectedSlotId);

	return (
		<div className={className} data-rag-organism="BookingPagePanel" {...rest}>
			<SplitPane
				ratio="sidebar"
				divider
				left={
					<Panel density="condensed">
						<Stack gap="sm">
							<PersonSummary name={bookingType.host.name} subtitle={bookingType.title} />
							<KeyValueStrip
								items={[
									{ key: "Duration", value: `${bookingType.durationMin} min` },
									...(bookingType.location
										? [{ key: "Location", value: bookingType.location }]
										: []),
								]}
							/>
							<MonthGrid
								monthISO={monthISO}
								markers={markersFor(days)}
								selectedDateISO={selectedDateISO}
								onDaySelect={onDaySelect}
								onMonthChange={onMonthChange}
							/>
						</Stack>
					</Panel>
				}
				right={
					<Panel
						title={selectedDateISO ? `Select a time — ${selectedDateISO}` : "Select a day"}
						density="condensed"
					>
						<Stack gap="sm">
							<Caption tone="muted">{`🌐 ${tz ? tz : "timezone"}`}</Caption>
							{selectedDateISO && slots.length ? (
								<TileGrid minTileWidth={96} gap="sm">
									{slots.map((slot) => (
										<Button
											key={slot.id}
											selected={slot.id === selectedSlotId}
											disabled={slot.disabled}
											onClick={() => onSlotSelect?.(slot.id)}
										>
											{slot.startISO.slice(11, 16)}
										</Button>
									))}
								</TileGrid>
							) : (
								<Text size="compact" tone="muted">
									{selectedDateISO ? "No times available." : "Pick a day to see times."}
								</Text>
							)}
							<Button variant="primary" disabled={!selectedSlot} onClick={() => onConfirm?.()}>
								{selectedSlot ? `Confirm ${selectedSlot.startISO.slice(11, 16)} →` : "Confirm"}
							</Button>
						</Stack>
					</Panel>
				}
			/>
		</div>
	);
}
