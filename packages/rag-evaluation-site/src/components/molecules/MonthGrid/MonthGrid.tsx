import type { HTMLAttributes, ReactNode } from "react";
import { type ContextStyleSet, contextVisualStyleToCssVars } from "../../../context";
import styles from "./MonthGrid.module.css";

export interface MonthGridDayMarker {
	count?: number;
	/** Lookup key into `styleSet.styles` for a heat/color background. */
	styleKey?: string;
	label?: ReactNode;
}

/**
 * The stable payload every day cell receives. Mirrors MatrixGrid's cell
 * contract: MonthGrid owns the calendar geometry; the day renderer owns how a
 * day looks. Any renderer honoring this shape is a valid day cell.
 */
export interface MonthGridDayPayload {
	dateISO: string;
	dayOfMonth: number;
	inMonth: boolean;
	isToday: boolean;
	selected: boolean;
	disabled: boolean;
	marker?: MonthGridDayMarker;
	onSelect: () => void;
}

export interface MonthGridProps extends Omit<HTMLAttributes<HTMLDivElement>, "onSelect"> {
	/** Any date in the target month (`2026-07`, `2026-07-01`, or full ISO). */
	monthISO: string;
	markers?: Record<string, MonthGridDayMarker>;
	styleSet?: ContextStyleSet;
	selectedDateISO?: string;
	/** ISO date treated as "today"; omit to disable the today highlight. */
	todayISO?: string;
	minDateISO?: string;
	maxDateISO?: string;
	/** 0 = Sunday, 1 = Monday (default). */
	weekStartsOn?: 0 | 1;
	showHeader?: boolean;
	/** Mode A — custom day renderer. Omit for the default day cell. */
	renderDay?: (payload: MonthGridDayPayload) => ReactNode;
	onDaySelect?: (dateISO: string) => void;
	/** Called with the new `YYYY-MM` when the header prev/next is used. */
	onMonthChange?: (monthISO: string) => void;
}

const MONTHS_FULL = [
	"January",
	"February",
	"March",
	"April",
	"May",
	"June",
	"July",
	"August",
	"September",
	"October",
	"November",
	"December",
];
const WEEKDAY_LABELS = ["Su", "Mo", "Tu", "We", "Th", "Fr", "Sa"];

function parseMonth(monthISO: string): { year: number; month: number } {
	const [y, m] = monthISO.slice(0, 7).split("-");
	return { year: Number(y), month: Number(m) - 1 };
}

function pad2(n: number): string {
	return String(n).padStart(2, "0");
}

function isoDate(year: number, month: number, day: number): string {
	// month is 0-based; Date normalizes over/underflow so adjacent-month days work.
	const d = new Date(Date.UTC(year, month, day));
	return `${d.getUTCFullYear()}-${pad2(d.getUTCMonth() + 1)}-${pad2(d.getUTCDate())}`;
}

function shiftMonth(year: number, month: number, delta: number): string {
	const d = new Date(Date.UTC(year, month + delta, 1));
	return `${d.getUTCFullYear()}-${pad2(d.getUTCMonth() + 1)}`;
}

export function MonthGrid({
	monthISO,
	markers,
	styleSet,
	selectedDateISO,
	todayISO,
	minDateISO,
	maxDateISO,
	weekStartsOn = 1,
	showHeader = true,
	renderDay,
	onDaySelect,
	onMonthChange,
	className,
	...rest
}: MonthGridProps) {
	const { year, month } = parseMonth(monthISO);
	const valid = Number.isFinite(year) && Number.isFinite(month);

	const firstWeekday = new Date(Date.UTC(year, month, 1)).getUTCDay();
	const leading = (firstWeekday - weekStartsOn + 7) % 7;
	const daysInMonth = new Date(Date.UTC(year, month + 1, 0)).getUTCDate();
	const weekCount = Math.ceil((leading + daysInMonth) / 7);
	const cellCount = weekCount * 7;

	const weekdayHeader = Array.from({ length: 7 }, (_, i) => WEEKDAY_LABELS[(weekStartsOn + i) % 7]);

	const cells: MonthGridDayPayload[] = [];
	for (let i = 0; i < cellCount; i++) {
		const dayNumber = i - leading + 1;
		const dateISO = isoDate(year, month, dayNumber);
		const inMonth = dayNumber >= 1 && dayNumber <= daysInMonth;
		const outOfRange =
			(minDateISO != null && dateISO < minDateISO) || (maxDateISO != null && dateISO > maxDateISO);
		const d = new Date(Date.UTC(year, month, dayNumber));
		cells.push({
			dateISO,
			dayOfMonth: d.getUTCDate(),
			inMonth,
			isToday: todayISO != null && dateISO === todayISO,
			selected: selectedDateISO != null && dateISO === selectedDateISO,
			disabled: !inMonth || outOfRange,
			marker: markers?.[dateISO],
			onSelect: () => onDaySelect?.(dateISO),
		});
	}

	return (
		<div
			className={[styles.root, className ?? ""].filter(Boolean).join(" ")}
			data-rag-molecule="MonthGrid"
			{...rest}
		>
			{showHeader ? (
				<div className={styles.header}>
					<button
						type="button"
						className={styles.nav}
						aria-label="Previous month"
						disabled={!valid || !onMonthChange}
						onClick={() => onMonthChange?.(shiftMonth(year, month, -1))}
					>
						‹
					</button>
					<span className={styles.title}>{valid ? `${MONTHS_FULL[month]} ${year}` : "—"}</span>
					<button
						type="button"
						className={styles.nav}
						aria-label="Next month"
						disabled={!valid || !onMonthChange}
						onClick={() => onMonthChange?.(shiftMonth(year, month, 1))}
					>
						›
					</button>
				</div>
			) : null}
			<div className={styles.weekdays}>
				{weekdayHeader.map((label, i) => (
					<span key={`${label}-${i}`} className={styles.weekday}>
						{label}
					</span>
				))}
			</div>
			<div className={styles.grid} role="grid">
				{cells.map((cell) =>
					renderDay ? (
						<div key={cell.dateISO} role="gridcell">
							{renderDay(cell)}
						</div>
					) : (
						<DefaultDay key={cell.dateISO} cell={cell} styleSet={styleSet} />
					),
				)}
			</div>
		</div>
	);
}

function DefaultDay({ cell, styleSet }: { cell: MonthGridDayPayload; styleSet?: ContextStyleSet }) {
	const visualStyle =
		cell.marker?.styleKey && styleSet ? styleSet.styles[cell.marker.styleKey] : undefined;
	return (
		<button
			type="button"
			className={styles.day}
			role="gridcell"
			data-in-month={cell.inMonth || undefined}
			data-today={cell.isToday || undefined}
			data-selected={cell.selected || undefined}
			aria-label={cell.dateISO}
			aria-pressed={cell.selected}
			disabled={cell.disabled}
			onClick={cell.onSelect}
			style={visualStyle ? contextVisualStyleToCssVars(visualStyle) : undefined}
		>
			<span className={styles.dayNumber}>{cell.dayOfMonth}</span>
			{cell.marker ? (
				<span className={styles.marker} aria-hidden="true">
					{cell.marker.label ?? (cell.marker.count != null ? cell.marker.count : "•")}
				</span>
			) : null}
		</button>
	);
}
