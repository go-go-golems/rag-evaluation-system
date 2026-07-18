import type { HTMLAttributes, ReactNode } from "react";
import styles from "./Panel.module.css";

export interface PanelProps extends Omit<HTMLAttributes<HTMLDivElement>, "title"> {
	title?: ReactNode;
	actions?: ReactNode;
	density?: "normal" | "condensed";
	fill?: boolean;
	titleTone?: "default" | "accent";
	children?: ReactNode;
}

export function Panel({
	title,
	actions,
	density = "normal",
	fill = false,
	titleTone = "default",
	className,
	children,
	...rest
}: PanelProps) {
	return (
		<section
			className={[
				styles.root,
				density === "condensed" ? styles.rootCondensed : "",
				fill ? styles.fill : "",
				className ?? "",
			]
				.filter(Boolean)
				.join(" ")}
			data-rag-layout="Panel"
			{...rest}
		>
			{(title || actions) && (
				<div
					className={[styles.header, titleTone === "accent" ? styles.headerAccent : ""]
						.filter(Boolean)
						.join(" ")}
				>
					<span>{title}</span>
					{actions}
				</div>
			)}
			<div className={density === "condensed" ? styles.condensed : styles.body}>{children}</div>
		</section>
	);
}
