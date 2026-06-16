import type { HTMLAttributes, ReactNode } from "react";
import styles from "./Stack.module.css";

export type StackGap = "xs" | "sm" | "md" | "lg" | "xl";
export type StackAlign = "stretch" | "start" | "center" | "end";

export interface StackProps extends HTMLAttributes<HTMLDivElement> {
	gap?: StackGap;
	align?: StackAlign;
	children?: ReactNode;
}

const gapClass: Record<StackGap, string> = {
	xs: styles.gapXs,
	sm: styles.gapSm,
	md: styles.gapMd,
	lg: styles.gapLg,
	xl: styles.gapXl,
};

const alignClass: Record<StackAlign, string> = {
	stretch: styles.alignStretch,
	start: styles.alignStart,
	center: styles.alignCenter,
	end: styles.alignEnd,
};

export function Stack({ gap = "md", align = "stretch", className, children, ...rest }: StackProps) {
	return (
		<div
			className={[styles.root, gapClass[gap], alignClass[align], className ?? ""]
				.filter(Boolean)
				.join(" ")}
			data-rag-layout="Stack"
			{...rest}
		>
			{children}
		</div>
	);
}
