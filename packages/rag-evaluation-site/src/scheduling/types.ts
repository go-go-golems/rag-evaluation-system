/**
 * Scheduling domain DTOs. Pure data — no React, no Widget IR, no server calls.
 * This is the scheduling counterpart to `src/context/types.ts`: engines,
 * adapters, presets, panels, and stories all import from here.
 */

export type AvailabilityState = "yes" | "ifneedbe" | "no" | "unknown";

export interface TimeSlot {
	id: string;
	startISO: string;
	endISO: string;
	tz: string;
	allDay?: boolean;
	label?: string;
}

export interface PollOption {
	id: string;
	slot: TimeSlot;
	note?: string;
}

export interface ParticipantResponse {
	id: string;
	name: string;
	avatarUrl?: string;
	comment?: string;
	cells: Record<string, AvailabilityState>;
	submittedAtISO?: string;
}

export interface PollSettings {
	allowIfNeedBe: boolean;
	hideVotesUntilResponded?: boolean;
	limitPerSlot?: number;
	deadlineISO?: string;
	finalOptionId?: string;
}

export interface MeetingPoll {
	id: string;
	title: string;
	description?: string;
	location?: string;
	organizer: { name: string; avatarUrl?: string };
	options: PollOption[];
	responses: ParticipantResponse[];
	settings: PollSettings;
	status: "open" | "finalized" | "closed";
}

/** Server-computed. Engines only render it; they never derive it. */
export interface SlotTally {
	optionId: string;
	yes: number;
	ifneedbe: number;
	no: number;
	score: number;
	isBest?: boolean;
	atCapacity?: boolean;
}

export interface CalendarEvent {
	id: string;
	title: string;
	startISO: string;
	endISO: string;
	allDay?: boolean;
	/** Palette lookup key (never a raw color). */
	colorKey: string;
	location?: string;
	attendees?: { name: string; avatarUrl?: string }[];
}

// ── 1:1 booking (Calendly-style) ────────────────────────────────────────────

export interface BookingType {
	id: string;
	title: string;
	durationMin: number;
	location?: string;
	host: { name: string; avatarUrl?: string };
}

export interface BookableDay {
	dateISO: string;
	slotCount: number;
	disabled?: boolean;
}

export interface BookableSlot {
	id: string;
	startISO: string;
	endISO: string;
	tz: string;
	disabled?: boolean;
}
