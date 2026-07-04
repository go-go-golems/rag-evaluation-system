import type { HTMLAttributes, ReactNode } from "react";
import styles from "./MeterBar.module.css";

export type MeterBarTone = "accent" | "success" | "danger";

export interface MeterBarProps extends HTMLAttributes<HTMLDivElement> {
	value: number;
	tone?: MeterBarTone;
	label?: ReactNode;
}

export function MeterBar({ value, tone = "accent", label, className, ...rest }: MeterBarProps) {
	const clamped = Math.min(1, Math.max(0, value));
	return (
		<div
			className={[styles.root, className ?? ""].filter(Boolean).join(" ")}
			data-rag-atom="MeterBar"
			data-tone={tone}
			{...rest}
		>
			<div
				className={styles.track}
				role="progressbar"
				aria-valuemin={0}
				aria-valuemax={100}
				aria-valuenow={Math.round(clamped * 100)}
			>
				<div
					className={[styles.fill, styles[tone]].join(" ")}
					style={{ width: `${clamped * 100}%` }}
				/>
			</div>
			{label != null && <span className={styles.label}>{label}</span>}
		</div>
	);
}
