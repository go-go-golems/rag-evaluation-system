import type { HTMLAttributes, ReactNode } from "react";
import styles from "./Caption.module.css";

export type CaptionTone = "muted" | "accent" | "success" | "warning" | "danger" | "inherit";
export type CaptionTransform = "none" | "uppercase";

export interface CaptionProps extends HTMLAttributes<HTMLSpanElement> {
	tone?: CaptionTone;
	transform?: CaptionTransform;
	truncate?: boolean;
	children?: ReactNode;
}

const toneClass: Record<CaptionTone, string> = {
	muted: styles.toneMuted,
	accent: styles.toneAccent,
	success: styles.toneSuccess,
	warning: styles.toneWarning,
	danger: styles.toneDanger,
	inherit: styles.toneInherit,
};

export function Caption({
	tone = "muted",
	transform = "none",
	truncate = false,
	className,
	children,
	...rest
}: CaptionProps) {
	return (
		<span
			className={[
				styles.root,
				toneClass[tone],
				transform === "uppercase" ? styles.uppercase : "",
				truncate ? styles.truncate : "",
				className ?? "",
			]
				.filter(Boolean)
				.join(" ")}
			data-rag-foundation="Caption"
			{...rest}
		>
			{children}
		</span>
	);
}
