import type { HTMLAttributes, ReactNode } from "react";
import { type ContextStyleSet, contextVisualStyleToCssVars } from "../../../context";
import { buildMonthGridCells, parseMonth, shiftMonth } from "./MonthGrid.logic";
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

	const weekdayHeader = Array.from({ length: 7 }, (_, i) => WEEKDAY_LABELS[(weekStartsOn + i) % 7]);

	const cells: MonthGridDayPayload[] = buildMonthGridCells({
		monthISO,
		weekStartsOn,
		todayISO,
		selectedDateISO,
		minDateISO,
		maxDateISO,
	}).map((cell) => ({
		...cell,
		marker: markers?.[cell.dateISO],
		onSelect: () => onDaySelect?.(cell.dateISO),
	}));

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
