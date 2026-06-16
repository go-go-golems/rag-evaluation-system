import type { ReactNode } from "react";
import styles from "./TabList.module.css";

export interface TabListItem<T extends string> {
	id: T;
	label: ReactNode;
}
export interface TabListProps<T extends string> {
	items: TabListItem<T>[];
	activeId: T;
	onChange: (id: T) => void;
	ariaLabel?: string;
}

export function TabList<T extends string>({
	items,
	activeId,
	onChange,
	ariaLabel = "Tabs",
}: TabListProps<T>) {
	return (
		<div className={styles.root} role="tablist" aria-label={ariaLabel} data-rag-layout="TabList">
			{items.map((item) => (
				<button
					key={item.id}
					type="button"
					role="tab"
					aria-selected={item.id === activeId}
					className={[styles.item, item.id === activeId ? styles.active : ""]
						.filter(Boolean)
						.join(" ")}
					onClick={() => onChange(item.id)}
				>
					{item.label}
				</button>
			))}
		</div>
	);
}
