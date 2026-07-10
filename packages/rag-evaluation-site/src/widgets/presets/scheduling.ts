import {
	AVAILABILITY_GLYPHS,
	AVAILABILITY_STATES,
	availabilityStyleSet,
	type CalendarEvent,
	eventStyleSet,
	type MeetingPoll,
	type SlotTally,
	type TimeSlot,
} from "../../scheduling";
import {
	component,
	type JsonObject,
	type MatrixGridWidgetProps,
	type MonthGridMarkerSpec,
	type SegmentedBarSegmentSpec,
	type TimeGridBlockSpec,
	text,
	type WidgetNode,
} from "../ir";

const MONTHS = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];

function formatSlot(slot: TimeSlot): string {
	const [, month, day] = slot.startISO.slice(0, 10).split("-");
	const time = slot.startISO.slice(11, 16);
	return `${MONTHS[Number(month) - 1] ?? "?"} ${Number(day)} · ${time}`;
}

export interface AvailabilityMatrixOptions {
	tallies?: SlotTally[];
	/** Response id whose row is editable (the "You" row). */
	editableResponseId?: string;
}

/**
 * Opinionated `schedule.dsl` preset: emits a MatrixGrid IR node configured with
 * Doodle poll semantics (availability cycle cells, availability palette, tally
 * footer, poll.toggleCell server action). The whole point of the layering — one
 * generic engine, one-arg domain skin.
 */
export function availabilityMatrix(
	poll: MeetingPoll,
	options: AvailabilityMatrixOptions = {},
): WidgetNode {
	const total = poll.responses.length;
	const tallyByOption = new Map((options.tallies ?? []).map((t) => [t.optionId, t]));

	const props: MatrixGridWidgetProps = {
		ariaLabel: poll.title,
		rows: poll.responses as unknown as JsonObject[],
		columns: poll.options.map((option) => {
			const tally = tallyByOption.get(option.id);
			return {
				id: option.id,
				header: text(`${formatSlot(option.slot)}${tally?.isBest ? " ★" : ""}`),
				meta: { yes: tally?.yes ?? 0, total } as JsonObject,
			};
		}),
		valueAt: { mapField: "cells" },
		cell: { kind: "cycle", states: AVAILABILITY_STATES, glyphs: AVAILABILITY_GLYPHS },
		styleSet: availabilityStyleSet,
		rowHeader: { kind: "field", field: "name" },
		editableRowKey: options.editableResponseId,
		getRowKey: { field: "id" },
		footer: { header: text("yes"), cell: { kind: "template", template: "${yes}/${total}" } },
		onCellAction: {
			kind: "server",
			name: "poll.toggleCell",
			payload: {
				pollId: poll.id,
				responseId: { kind: "path", path: "rowKey" },
				optionId: { kind: "path", path: "colId" },
				state: { kind: "path", path: "value" },
			} as JsonObject,
		},
	};

	return component("MatrixGrid", props);
}

/**
 * `schedule.dsl` preset: organizer results — one SegmentedBar per option, ranked
 * order preserved, best slot starred. Emits a Stack of labelled bars.
 */
export function pollResults(poll: MeetingPoll, tallies: SlotTally[]): WidgetNode {
	const byId = new Map(tallies.map((t) => [t.optionId, t]));
	return component(
		"Stack",
		{ gap: "md" },
		poll.options.map((option) => {
			const tally = byId.get(option.id);
			const segments: SegmentedBarSegmentSpec[] = tally
				? [
						{ value: tally.yes, styleKey: "yes", label: text("yes") },
						{ value: tally.ifneedbe, styleKey: "ifneedbe", label: text("maybe") },
						{ value: tally.no, styleKey: "no", label: text("no") },
					]
				: [];
			return component("Stack", { gap: "xs" }, [
				component("Caption", {}, [text(`${formatSlot(option.slot)}${tally?.isBest ? " ★" : ""}`)]),
				component("SegmentedBar", {
					segments,
					styleSet: availabilityStyleSet,
					showCounts: true,
				}),
			]);
		}),
	);
}

/** `calendar.dsl` preset: a month heatmap of event density colored by category. */
export function monthCalendar(events: CalendarEvent[], monthISO: string): WidgetNode {
	const markers: Record<string, MonthGridMarkerSpec> = {};
	for (const event of events) {
		const date = event.startISO.slice(0, 10);
		const existing = markers[date];
		markers[date] = {
			count: (existing?.count ?? 0) + 1,
			styleKey: existing?.styleKey ?? event.colorKey,
		};
	}
	return component("MonthGrid", { monthISO, markers, styleSet: eventStyleSet });
}

/** `calendar.dsl` preset: a week time-grid built from calendar events. */
export function weekCalendar(events: CalendarEvent[], daysISO: string[]): WidgetNode {
	const blocks: TimeGridBlockSpec[] = events.map((event) => ({
		id: event.id,
		dayISO: event.startISO.slice(0, 10),
		startISO: event.startISO,
		endISO: event.endISO,
		styleKey: event.colorKey,
		label: text(event.title),
	}));
	return component("TimeGrid", {
		days: daysISO,
		blocks,
		styleSet: eventStyleSet,
		hourStart: 8,
		hourEnd: 18,
	});
}
