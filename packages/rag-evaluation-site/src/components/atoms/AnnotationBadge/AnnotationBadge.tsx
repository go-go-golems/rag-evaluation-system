import type { HTMLAttributes } from "react";
import type { ContextVisualStyle } from "../../../context";
import { ContextStyleSwatch } from "../ContextStyleSwatch";
import styles from "./AnnotationBadge.module.css";

export interface AnnotationBadgeProps extends HTMLAttributes<HTMLSpanElement> {
	visualStyle: ContextVisualStyle;
	label: string;
	selected?: boolean;
}

export function AnnotationBadge({
	visualStyle,
	label,
	selected = false,
	className,
	...rest
}: AnnotationBadgeProps) {
	return (
		<span
			className={[styles.root, selected ? styles.selected : "", className ?? ""]
				.filter(Boolean)
				.join(" ")}
			data-rag-atom="AnnotationBadge"
			{...rest}
		>
			<ContextStyleSwatch visualStyle={visualStyle} size="sm" selected={selected} />
			<span>{label}</span>
		</span>
	);
}
