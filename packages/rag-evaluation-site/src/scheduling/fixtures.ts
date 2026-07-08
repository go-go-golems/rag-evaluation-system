import type {
	BookableDay,
	BookableSlot,
	BookingType,
	CalendarEvent,
	MeetingPoll,
	SlotTally,
} from "./types";

function slot(id: string, dateISO: string, from: string, to: string) {
	return {
		id,
		startISO: `${dateISO}T${from}:00`,
		endISO: `${dateISO}T${to}:00`,
		tz: "Europe/Berlin",
	};
}

export const sampleTeamSyncPoll: MeetingPoll = {
	id: "poll-teamsync",
	title: "Team sync — pick a time",
	location: "Zoom",
	organizer: { name: "Manuel" },
	options: [
		{ id: "s1", slot: slot("s1", "2026-07-09", "14:00", "15:00") },
		{ id: "s2", slot: slot("s2", "2026-07-10", "10:00", "11:00") },
		{ id: "s3", slot: slot("s3", "2026-07-10", "16:00", "17:00") },
		{ id: "s4", slot: slot("s4", "2026-07-11", "09:00", "10:00") },
	],
	responses: [
		{ id: "alice", name: "Alice", cells: { s1: "yes", s2: "ifneedbe", s3: "no", s4: "yes" } },
		{ id: "bob", name: "Bob", cells: { s1: "yes", s2: "yes", s3: "no", s4: "unknown" } },
		{ id: "chen", name: "Chen", cells: { s1: "yes", s2: "yes", s3: "yes", s4: "yes" } },
		{ id: "you", name: "You", cells: { s1: "yes", s2: "ifneedbe", s3: "no", s4: "unknown" } },
	],
	settings: { allowIfNeedBe: true, deadlineISO: "2026-07-08T23:59:00" },
	status: "open",
};

/** Server would compute these; hard-coded here for stories. */
export const sampleTeamSyncTallies: SlotTally[] = [
	{ optionId: "s1", yes: 4, ifneedbe: 0, no: 0, score: 4, isBest: true },
	{ optionId: "s2", yes: 2, ifneedbe: 2, no: 0, score: 3 },
	{ optionId: "s3", yes: 1, ifneedbe: 0, no: 3, score: 1 },
	{ optionId: "s4", yes: 2, ifneedbe: 0, no: 0, score: 2 },
];

export const sampleWeekEvents: CalendarEvent[] = [
	{
		id: "e1",
		title: "Standup",
		startISO: "2026-07-06T09:00",
		endISO: "2026-07-06T09:30",
		colorKey: "meeting",
	},
	{
		id: "e2",
		title: "Focus block",
		startISO: "2026-07-06T11:00",
		endISO: "2026-07-06T12:30",
		colorKey: "focus",
	},
	{
		id: "e3",
		title: "Team sync",
		startISO: "2026-07-09T14:00",
		endISO: "2026-07-09T15:00",
		colorKey: "meeting",
	},
	{
		id: "e4",
		title: "Gym",
		startISO: "2026-07-10T16:00",
		endISO: "2026-07-10T17:30",
		colorKey: "personal",
	},
];

// ── 1:1 booking fixtures ────────────────────────────────────────────────────

export const sampleBookingType: BookingType = {
	id: "bt-30",
	title: "30 min meeting",
	durationMin: 30,
	location: "Google Meet",
	host: { name: "Manuel Odendahl" },
};

export const sampleBookableDays: BookableDay[] = [
	{ dateISO: "2026-07-09", slotCount: 6 },
	{ dateISO: "2026-07-10", slotCount: 4 },
	{ dateISO: "2026-07-13", slotCount: 5 },
	{ dateISO: "2026-07-14", slotCount: 2 },
];

export const sampleBookableSlots: BookableSlot[] = [
	{ id: "s0900", startISO: "2026-07-09T09:00", endISO: "2026-07-09T09:30", tz: "Europe/Berlin" },
	{ id: "s0930", startISO: "2026-07-09T09:30", endISO: "2026-07-09T10:00", tz: "Europe/Berlin" },
	{ id: "s1000", startISO: "2026-07-09T10:00", endISO: "2026-07-09T10:30", tz: "Europe/Berlin" },
	{
		id: "s1030",
		startISO: "2026-07-09T10:30",
		endISO: "2026-07-09T11:00",
		tz: "Europe/Berlin",
		disabled: true,
	},
	{ id: "s1400", startISO: "2026-07-09T14:00", endISO: "2026-07-09T14:30", tz: "Europe/Berlin" },
	{ id: "s1430", startISO: "2026-07-09T14:30", endISO: "2026-07-09T15:00", tz: "Europe/Berlin" },
];
