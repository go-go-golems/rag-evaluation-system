import type { HTMLAttributes, ReactNode } from "react";
import { Caption } from "../../foundation";
import styles from "./EmptyState.module.css";

export interface EmptyStateProps extends Omit<HTMLAttributes<HTMLDivElement>, "title"> {
	glyph?: ReactNode;
	title: ReactNode;
	hint?: ReactNode;
	action?: ReactNode;
	framed?: boolean;
}

export function EmptyState({
	glyph = "□",
	title,
	hint,
	action,
	framed = false,
	className,
	...rest
}: EmptyStateProps) {
	return (
		<div
			className={[styles.root, framed ? styles.framed : "", className ?? ""]
				.filter(Boolean)
				.join(" ")}
			data-rag-molecule="EmptyState"
			{...rest}
		>
			<span className={styles.glyph} aria-hidden="true">
				{glyph}
			</span>
			<span className={styles.title}>{title}</span>
			{hint && <Caption className={styles.hint}>{hint}</Caption>}
			{action && <span className={styles.action}>{action}</span>}
		</div>
	);
}
