import type { HTMLAttributes, ReactNode } from "react";
import { type MeterBarTone, MeterBar } from "../../atoms";
import { Caption } from "../../foundation";
import styles from "./StatTile.module.css";

export type StatTrend = "up" | "down" | "flat";

export interface StatTileProps extends HTMLAttributes<HTMLDivElement> {
	label: ReactNode;
	value: ReactNode;
	/** Signed change; its sign drives the default trend arrow. */
	delta?: number;
	/** Overrides the derived `${delta}%` delta text. */
	deltaLabel?: ReactNode;
	/** Explicit trend; otherwise derived from the sign of `delta`. */
	trend?: StatTrend;
	/** Inline proportion 0..1 rendered as a MeterBar under the value. */
	progress?: number;
	tone?: MeterBarTone;
}

const TREND_GLYPH: Record<StatTrend, string> = { up: "▲", down: "▼", flat: "■" };

function deriveTrend(delta?: number): StatTrend | undefined {
	if (delta == null) return undefined;
	if (delta > 0) return "up";
	if (delta < 0) return "down";
	return "flat";
}

/**
 * A labeled number with an optional delta and inline progress bar, for
 * dashboards ("Open pipeline $2.1M ▲12%"). Reuses MeterBar for the track. Lay
 * several out with the TileGrid / DashboardGrid layout primitives.
 */
export function StatTile({
	label,
	value,
	delta,
	deltaLabel,
	trend,
	progress,
	tone = "accent",
	className,
	...rest
}: StatTileProps) {
	const resolvedTrend = trend ?? deriveTrend(delta);
	const deltaText = deltaLabel ?? (delta != null ? `${delta > 0 ? "+" : ""}${delta}%` : null);
	return (
		<div
			className={[styles.root, className ?? ""].filter(Boolean).join(" ")}
			data-rag-molecule="StatTile"
			{...rest}
		>
			<Caption>{label}</Caption>
			<div className={styles.value}>{value}</div>
			{deltaText != null && resolvedTrend ? (
				<div className={styles.delta} data-trend={resolvedTrend}>
					<span className={styles.arrow} aria-hidden="true">
						{TREND_GLYPH[resolvedTrend]}
					</span>
					<span>{deltaText}</span>
				</div>
			) : null}
			{progress != null ? <MeterBar className={styles.bar} value={progress} tone={tone} /> : null}
		</div>
	);
}
