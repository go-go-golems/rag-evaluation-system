import type { HTMLAttributes, ReactNode } from "react";
import styles from "./Inline.module.css";

export type InlineGap = "xs" | "sm" | "md" | "lg";
export type InlineJustify = "start" | "between" | "end";

export interface InlineProps extends HTMLAttributes<HTMLDivElement> {
	gap?: InlineGap;
	justify?: InlineJustify;
	wrap?: boolean;
	children?: ReactNode;
}

const gapClass: Record<InlineGap, string> = {
	xs: styles.gapXs,
	sm: styles.gapSm,
	md: styles.gapMd,
	lg: styles.gapLg,
};
const justifyClass: Record<InlineJustify, string> = {
	start: styles.justifyStart,
	between: styles.justifyBetween,
	end: styles.justifyEnd,
};

export function Inline({
	gap = "sm",
	justify = "start",
	wrap = true,
	className,
	children,
	...rest
}: InlineProps) {
	return (
		<div
			className={[
				styles.root,
				gapClass[gap],
				justifyClass[justify],
				wrap ? styles.wrap : styles.noWrap,
				className ?? "",
			]
				.filter(Boolean)
				.join(" ")}
			data-rag-layout="Inline"
			{...rest}
		>
			{children}
		</div>
	);
}
