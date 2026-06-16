import { createElement, type HTMLAttributes, type ReactNode } from "react";
import styles from "./Text.module.css";

export type TextAs = "p" | "span" | "div" | "strong";
export type TextSize = "body" | "compact" | "metadata" | "label" | "metric";
export type TextTone =
	| "primary"
	| "muted"
	| "inverse"
	| "accent"
	| "success"
	| "warning"
	| "danger"
	| "inherit";
export type TextWeight = "regular" | "bold";
export type TextAlign = "start" | "center" | "end";

export interface TextProps extends HTMLAttributes<HTMLElement> {
	as?: TextAs;
	size?: TextSize;
	tone?: TextTone;
	weight?: TextWeight;
	align?: TextAlign;
	truncate?: boolean;
	children?: ReactNode;
}

const sizeClass: Record<TextSize, string> = {
	body: styles.sizeBody,
	compact: styles.sizeCompact,
	metadata: styles.sizeMetadata,
	label: styles.sizeLabel,
	metric: styles.sizeMetric,
};

const toneClass: Record<TextTone, string> = {
	primary: styles.tonePrimary,
	muted: styles.toneMuted,
	inverse: styles.toneInverse,
	accent: styles.toneAccent,
	success: styles.toneSuccess,
	warning: styles.toneWarning,
	danger: styles.toneDanger,
	inherit: styles.toneInherit,
};

const weightClass: Record<TextWeight, string> = {
	regular: styles.weightRegular,
	bold: styles.weightBold,
};

const alignClass: Record<TextAlign, string> = {
	start: styles.alignStart,
	center: styles.alignCenter,
	end: styles.alignEnd,
};

export function Text({
	as = "p",
	size = "body",
	tone = "primary",
	weight = "regular",
	align = "start",
	truncate = false,
	className,
	children,
	...rest
}: TextProps) {
	const rootClassName = [
		styles.root,
		sizeClass[size],
		toneClass[tone],
		weightClass[weight],
		alignClass[align],
		truncate ? styles.truncate : "",
		className ?? "",
	]
		.filter(Boolean)
		.join(" ");

	return createElement(as, { className: rootClassName, ...rest }, children);
}
