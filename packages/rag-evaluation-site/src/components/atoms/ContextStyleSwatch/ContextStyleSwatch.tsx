import type { HTMLAttributes } from "react";
import {
	contextVisualStyleToCssVars,
	type ContextVisualStyle,
	type ContextStyleSize,
} from "../../../context";
import styles from "./ContextStyleSwatch.module.css";

export interface ContextStyleSwatchProps extends HTMLAttributes<HTMLSpanElement> {
	visualStyle: ContextVisualStyle;
	size?: ContextStyleSize;
	selected?: boolean;
}

export function ContextStyleSwatch({
	visualStyle,
	size = "md",
	selected = false,
	className,
	style,
	...rest
}: ContextStyleSwatchProps) {
	const pattern = visualStyle.pattern ?? "none";
	return (
		<span
			className={[
				styles.root,
				styles[size],
				styles[`pattern_${pattern}`] ?? "",
				visualStyle.dashed ? styles.dashed : "",
				visualStyle.dotted ? styles.dotted : "",
				selected ? styles.selected : "",
				className ?? "",
			]
				.filter(Boolean)
				.join(" ")}
			data-rag-atom="ContextStyleSwatch"
			data-pattern={pattern}
			style={{ ...contextVisualStyleToCssVars(visualStyle), ...style }}
			{...rest}
		/>
	);
}
