import type { CSSProperties, HTMLAttributes, ReactNode } from "react";
import styles from "./FieldGrid.module.css";

export type FieldGridGap = "sm" | "md";

export interface FieldGridProps extends HTMLAttributes<HTMLDivElement> {
	columns?: number;
	gap?: FieldGridGap;
	children?: ReactNode;
}

export function FieldGrid({
	columns = 2,
	gap = "md",
	className,
	style,
	children,
	...rest
}: FieldGridProps) {
	return (
		<div
			className={[styles.root, gap === "sm" ? styles.gapSm : styles.gapMd, className ?? ""]
				.filter(Boolean)
				.join(" ")}
			style={{ "--rag-field-columns": String(columns), ...style } as CSSProperties}
			data-rag-layout="FieldGrid"
			{...rest}
		>
			{children}
		</div>
	);
}
