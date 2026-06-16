import type { HTMLAttributes, ReactNode } from "react";
import styles from "./SplitPane.module.css";

export interface SplitPaneProps extends HTMLAttributes<HTMLDivElement> {
	left: ReactNode;
	right: ReactNode;
	ratio?: "balanced" | "leftNarrow" | "rightNarrow" | "course" | "sidebar";
	divider?: boolean;
	gutter?: "none" | "md" | "lg";
}

export function SplitPane({
	left,
	right,
	ratio = "balanced",
	divider = false,
	gutter = "none",
	className,
	...rest
}: SplitPaneProps) {
	return (
		<div
			className={[
				styles.root,
				styles[ratio],
				divider ? styles.divider : "",
				gutter !== "none" ? styles[`gutter${gutter[0]!.toUpperCase()}${gutter.slice(1)}`] : "",
				className ?? "",
			]
				.filter(Boolean)
				.join(" ")}
			data-rag-layout="SplitPane"
			{...rest}
		>
			<div className={styles.pane}>{left}</div>
			<div className={styles.pane}>{right}</div>
		</div>
	);
}
