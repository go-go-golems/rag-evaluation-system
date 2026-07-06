import type { HTMLAttributes, ReactNode } from "react";
import { Caption } from "../../foundation";
import styles from "./SectionBlock.module.css";

export type SectionBlockLevel = 1 | 2 | 3;

export interface SectionBlockProps extends HTMLAttributes<HTMLElement> {
	as?: "section" | "div" | "article";
	label?: ReactNode;
	caption?: ReactNode;
	actions?: ReactNode;
	level?: SectionBlockLevel;
	rule?: boolean;
	density?: "normal" | "spacious" | "flush";
	divider?: "none" | "top" | "bottom" | "both";
	children?: ReactNode;
}

export function SectionBlock({
	as: Element = "section",
	label,
	caption,
	actions,
	level = 1,
	rule = false,
	density = "normal",
	divider = "none",
	className,
	children,
	...rest
}: SectionBlockProps) {
	return (
		<Element
			className={[styles.root, styles[density], styles[`divider-${divider}`], className ?? ""]
				.filter(Boolean)
				.join(" ")}
			data-rag-layout="SectionBlock"
			data-level={level}
			{...rest}
		>
			{(label || actions) && (
				<div className={[styles.labelRow, rule ? styles.rule : ""].filter(Boolean).join(" ")}>
					{label && (
						<div className={[styles.label, styles[`level${level}`]].join(" ")}>{label}</div>
					)}
					{actions && <div className={styles.actions}>{actions}</div>}
				</div>
			)}
			{caption && <Caption className={styles.caption}>{caption}</Caption>}
			<div className={styles.body}>{children}</div>
		</Element>
	);
}
