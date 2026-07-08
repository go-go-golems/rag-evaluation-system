import type { HTMLAttributes } from "react";
import styles from "./DateTile.module.css";

export type DateTileEmphasis = "default" | "muted" | "accent";
export type DateTileSize = "sm" | "md" | "lg";

export interface DateTileProps extends HTMLAttributes<HTMLDivElement> {
	/** ISO date. Date-only (`2026-07-09`) or full ISO; only the calendar day is shown. */
	dateISO: string;
	emphasis?: DateTileEmphasis;
	size?: DateTileSize;
	/** Hide the weekday row. */
	hideWeekday?: boolean;
}

const MONTHS = ["JAN", "FEB", "MAR", "APR", "MAY", "JUN", "JUL", "AUG", "SEP", "OCT", "NOV", "DEC"];
const WEEKDAYS = ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"];

function toUtcDate(dateISO: string): Date {
	// Normalize to a UTC instant so a bare `YYYY-MM-DD` never drifts by timezone.
	const iso = dateISO.length === 10 ? `${dateISO}T00:00:00Z` : dateISO;
	return new Date(iso);
}

export function DateTile({
	dateISO,
	emphasis = "default",
	size = "md",
	hideWeekday = false,
	className,
	...rest
}: DateTileProps) {
	const d = toUtcDate(dateISO);
	const valid = !Number.isNaN(d.getTime());
	const month = valid ? MONTHS[d.getUTCMonth()] : "—";
	const day = valid ? String(d.getUTCDate()) : "—";
	const weekday = valid ? WEEKDAYS[d.getUTCDay()] : "";

	return (
		<div
			className={[styles.root, styles[size], className ?? ""].filter(Boolean).join(" ")}
			data-rag-atom="DateTile"
			data-emphasis={emphasis}
			role="img"
			aria-label={valid ? `${weekday} ${month} ${day}` : "Invalid date"}
			{...rest}
		>
			<span className={styles.month}>{month}</span>
			<span className={styles.day}>{day}</span>
			{!hideWeekday && weekday ? <span className={styles.weekday}>{weekday}</span> : null}
		</div>
	);
}
