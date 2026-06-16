import type { HTMLAttributes, ReactNode } from "react";
import styles from "./KeyPointList.module.css";

export interface KeyPointListItem {
	id?: string;
	index?: ReactNode;
	title?: ReactNode;
	text: ReactNode;
	meta?: ReactNode;
}

export interface KeyPointListProps extends HTMLAttributes<HTMLOListElement> {
	items: Array<KeyPointListItem | ReactNode>;
	markerTone?: "accent" | "muted";
}

function normalizeItem(
	item: KeyPointListItem | ReactNode,
	fallbackIndex: number,
): KeyPointListItem {
	if (typeof item === "object" && item !== null && !Array.isArray(item) && "text" in item) {
		return item as KeyPointListItem;
	}
	return { text: item, index: String(fallbackIndex + 1).padStart(2, "0") };
}

export function KeyPointList({
	items,
	markerTone = "accent",
	className,
	...rest
}: KeyPointListProps) {
	return (
		<ol
			className={[styles.root, className ?? ""].filter(Boolean).join(" ")}
			data-rag-molecule="KeyPointList"
			data-marker-tone={markerTone}
			{...rest}
		>
			{items.map((rawItem, i) => {
				const item = normalizeItem(rawItem, i);
				const index = item.index ?? String(i + 1).padStart(2, "0");
				return (
					<li key={item.id ?? i} className={styles.item}>
						<span className={styles.index}>{index}</span>
						<span className={styles.body}>
							{item.title && <span className={styles.title}>{item.title}</span>}
							<span className={styles.text}>{item.text}</span>
						</span>
						{item.meta && <span className={styles.meta}>{item.meta}</span>}
					</li>
				);
			})}
		</ol>
	);
}
