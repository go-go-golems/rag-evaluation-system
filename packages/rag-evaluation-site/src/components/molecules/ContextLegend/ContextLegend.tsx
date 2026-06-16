import type { HTMLAttributes } from "react";
import type { ContextLegendItemSpec, ContextStyleSize, ContextVisualStyle } from "../../../context";
import { ContextStyleSwatch } from "../../atoms";
import styles from "./ContextLegend.module.css";

export interface ContextLegendProps extends HTMLAttributes<HTMLDivElement> {
	items: ContextLegendItemSpec[];
	styles: Record<string, ContextVisualStyle>;
	size?: ContextStyleSize;
	selectedId?: string;
}

export function ContextLegend({
	items,
	styles: visualStyles,
	size = "md",
	selectedId,
	className,
	...rest
}: ContextLegendProps) {
	const visibleItems = [...items]
		.filter((item) => !item.hidden)
		.sort((a, b) => (a.order ?? 0) - (b.order ?? 0));
	return (
		<div
			className={[styles.root, styles[size], className ?? ""].filter(Boolean).join(" ")}
			data-rag-molecule="ContextLegend"
			{...rest}
		>
			{visibleItems.map((item) => {
				const styleKey = item.styleKey ?? item.id;
				const visualStyle = visualStyles[styleKey];
				if (!visualStyle) return null;
				return (
					<div
						key={item.id}
						className={styles.item}
						data-selected={selectedId === item.id || selectedId === styleKey ? "true" : undefined}
						title={item.description}
					>
						<ContextStyleSwatch
							visualStyle={visualStyle}
							size={size}
							selected={selectedId === item.id || selectedId === styleKey}
						/>
						<span className={styles.label}>{item.label}</span>
					</div>
				);
			})}
		</div>
	);
}
