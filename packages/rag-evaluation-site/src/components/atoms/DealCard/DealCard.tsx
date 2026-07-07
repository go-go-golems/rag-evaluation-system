import type { CSSProperties, ReactNode } from "react";
import styles from "./DealCard.module.css";

export type DealCardStatus = "open" | "won" | "lost";

/**
 * The swappable unit of a pipeline board — small, high-frequency, read at a
 * glance. Purely presentational: the BoardEngine owns geometry, drag, and
 * selection and wraps this. Slots are nodes so the same card renders from a
 * TS preset or from serialized IR CellSpecs.
 */
export interface DealCardProps {
	title: ReactNode;
	/** e.g. a formatted amount. */
	subtitle?: ReactNode;
	/** e.g. owner avatar + name, due indicator, tag dots. */
	meta?: ReactNode;
	status?: DealCardStatus;
	/** Left accent bar color, resolved from a ContextStyleSet by the caller. */
	accentStyle?: CSSProperties;
	selected?: boolean;
	dragging?: boolean;
}

export function DealCard({
	title,
	subtitle,
	meta,
	status = "open",
	accentStyle,
	selected = false,
	dragging = false,
}: DealCardProps) {
	return (
		<div
			className={styles.root}
			data-rag-atom="DealCard"
			data-status={status}
			data-selected={selected || undefined}
			data-dragging={dragging || undefined}
			style={accentStyle}
		>
			<div className={styles.accent} aria-hidden="true" />
			<div className={styles.body}>
				<div className={styles.title}>{title}</div>
				{subtitle != null ? <div className={styles.subtitle}>{subtitle}</div> : null}
				{meta != null ? <div className={styles.meta}>{meta}</div> : null}
			</div>
			{status !== "open" ? (
				<span className={styles.badge} data-status={status}>
					{status === "won" ? "✓ won" : "✕ lost"}
				</span>
			) : null}
		</div>
	);
}
