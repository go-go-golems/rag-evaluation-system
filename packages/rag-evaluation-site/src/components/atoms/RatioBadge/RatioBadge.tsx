import type { HTMLAttributes } from "react";
import styles from "./RatioBadge.module.css";

export type RatioBadgeTone = "neutral" | "positive" | "warning" | "muted";

export interface RatioBadgeProps extends HTMLAttributes<HTMLSpanElement> {
	count: number;
	total: number;
	/** Force a tone; otherwise it is derived from `count / total`. */
	tone?: RatioBadgeTone;
	/** Optional label rendered before the ratio, e.g. "yes". */
	label?: string;
	/** Show a small proportional dot track after the numbers. */
	showTrack?: boolean;
}

function deriveTone(ratio: number): RatioBadgeTone {
	if (ratio >= 0.75) return "positive";
	if (ratio >= 0.4) return "neutral";
	if (ratio > 0) return "warning";
	return "muted";
}

export function RatioBadge({
	count,
	total,
	tone,
	label,
	showTrack = false,
	className,
	...rest
}: RatioBadgeProps) {
	const safeTotal = Math.max(0, total);
	const safeCount = Math.min(Math.max(0, count), safeTotal);
	const ratio = safeTotal === 0 ? 0 : safeCount / safeTotal;
	const resolvedTone = tone ?? deriveTone(ratio);

	return (
		<span
			className={[styles.root, className ?? ""].filter(Boolean).join(" ")}
			data-rag-atom="RatioBadge"
			data-tone={resolvedTone}
			aria-label={`${label ? `${label}: ` : ""}${safeCount} of ${safeTotal}`}
			{...rest}
		>
			{label ? <span className={styles.label}>{label}</span> : null}
			<span className={styles.value}>
				<span className={styles.count}>{safeCount}</span>
				<span className={styles.slash}>/</span>
				<span className={styles.total}>{safeTotal}</span>
			</span>
			{showTrack ? (
				<span className={styles.track} aria-hidden="true">
					<span className={styles.fill} style={{ width: `${ratio * 100}%` }} />
				</span>
			) : null}
		</span>
	);
}
