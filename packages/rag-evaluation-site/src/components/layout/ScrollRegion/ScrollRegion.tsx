import type { HTMLAttributes, ReactNode } from "react";
import styles from "./ScrollRegion.module.css";

export interface ScrollRegionProps extends HTMLAttributes<HTMLDivElement> {
	axis?: "y" | "x" | "both";
	children?: ReactNode;
}

export function ScrollRegion({ axis = "y", className, children, ...rest }: ScrollRegionProps) {
	return (
		<div
			className={[styles.root, styles[axis], className ?? ""].filter(Boolean).join(" ")}
			data-rag-layout="ScrollRegion"
			{...rest}
		>
			{children}
		</div>
	);
}
