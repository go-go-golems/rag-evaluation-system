import type { HTMLAttributes, ReactNode } from "react";
import { type ContextStyleSet, contextVisualStyleToCssVars } from "../../../context";
import styles from "./SegmentedBar.module.css";

export interface SegmentedBarSegment {
	value: number;
	/** Lookup key into `styleSet.styles`. */
	styleKey: string;
	label?: ReactNode;
}

export interface SegmentedBarMarker {
	/** Position along the bar in the same units as segment values (0..total). */
	at: number;
	styleKey?: string;
	label?: ReactNode;
}

export interface SegmentedBarProps extends HTMLAttributes<HTMLDivElement> {
	segments: SegmentedBarSegment[];
	styleSet: ContextStyleSet;
	/** Denominator for proportions. Defaults to the sum of segment values. */
	total?: number;
	/** Render a counts row underneath the bar. */
	showCounts?: boolean;
	markers?: SegmentedBarMarker[];
	size?: "sm" | "md" | "lg";
	onSegmentSelect?: (styleKey: string, index: number) => void;
}

function pct(part: number, total: number): number {
	if (total <= 0) return 0;
	return Math.min(100, Math.max(0, (part / total) * 100));
}

export function SegmentedBar({
	segments,
	styleSet,
	total,
	showCounts = false,
	markers,
	size = "md",
	onSegmentSelect,
	className,
	...rest
}: SegmentedBarProps) {
	const sum = segments.reduce((acc, s) => acc + Math.max(0, s.value), 0);
	const denom = total ?? sum;
	const interactive = Boolean(onSegmentSelect);

	return (
		<div
			className={[styles.root, styles[size], className ?? ""].filter(Boolean).join(" ")}
			data-rag-molecule="SegmentedBar"
			{...rest}
		>
			<div className={styles.track}>
				{segments.map((segment, index) => {
					const visualStyle = styleSet.styles[segment.styleKey] ?? styleSet.fallbackStyle;
					const width = pct(Math.max(0, segment.value), denom);
					if (width <= 0) return null;
					const key = `${segment.styleKey}-${index}`;
					const commonStyle = visualStyle ? contextVisualStyleToCssVars(visualStyle) : undefined;
					const title = typeof segment.label === "string" ? segment.label : segment.styleKey;
					if (interactive) {
						return (
							<button
								key={key}
								type="button"
								className={styles.segment}
								style={{ ...commonStyle, width: `${width}%` }}
								data-style-key={segment.styleKey}
								title={title}
								onClick={() => onSegmentSelect?.(segment.styleKey, index)}
							/>
						);
					}
					return (
						<span
							key={key}
							className={styles.segment}
							style={{ ...commonStyle, width: `${width}%` }}
							data-style-key={segment.styleKey}
							title={title}
						/>
					);
				})}
				{markers?.map((marker, index) => {
					const visualStyle = marker.styleKey ? styleSet.styles[marker.styleKey] : undefined;
					return (
						<span
							key={`marker-${marker.at}-${index}`}
							className={styles.marker}
							style={{
								left: `${pct(marker.at, denom)}%`,
								...(visualStyle ? contextVisualStyleToCssVars(visualStyle) : undefined),
							}}
						>
							{marker.label ? <span className={styles.markerLabel}>{marker.label}</span> : null}
						</span>
					);
				})}
			</div>
			{showCounts ? (
				<div className={styles.counts}>
					{segments.map((segment, index) => {
						const visualStyle = styleSet.styles[segment.styleKey] ?? styleSet.fallbackStyle;
						return (
							<span key={`count-${segment.styleKey}-${index}`} className={styles.count}>
								<span
									className={styles.swatch}
									style={visualStyle ? contextVisualStyleToCssVars(visualStyle) : undefined}
									aria-hidden="true"
								/>
								<span className={styles.countValue}>{segment.value}</span>
								{segment.label ? <span className={styles.countLabel}>{segment.label}</span> : null}
							</span>
						);
					})}
				</div>
			) : null}
		</div>
	);
}
