import type { HTMLAttributes, ReactNode } from "react";
import { Caption } from "../../foundation";
import styles from "./SectionBlock.module.css";

export interface SectionBlockProps extends HTMLAttributes<HTMLElement> {
	as?: "section" | "div" | "article";
	label?: ReactNode;
	caption?: ReactNode;
	density?: "normal" | "spacious";
	divider?: "none" | "top" | "bottom" | "both";
	children?: ReactNode;
}

export function SectionBlock({
	as: Element = "section",
	label,
	caption,
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
			{...rest}
		>
			{label && <div className={styles.label}>{label}</div>}
			{caption && <Caption className={styles.caption}>{caption}</Caption>}
			<div className={styles.body}>{children}</div>
		</Element>
	);
}
