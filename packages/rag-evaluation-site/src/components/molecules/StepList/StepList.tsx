import type { HTMLAttributes, ReactNode } from "react";
import styles from "./StepList.module.css";

export interface StepListItem {
	id: string;
	index: ReactNode;
	title: ReactNode;
	description?: ReactNode;
	meta?: ReactNode;
}
export interface StepListProps extends HTMLAttributes<HTMLOListElement> {
	items: StepListItem[];
	/** Whether items can be selected. Defaults to true when onItemSelect is provided. */
	selectable?: boolean;
	activeItemId?: string;
	onItemSelect?: (id: string) => void;
}

export function StepList({
	items,
	selectable,
	activeItemId,
	onItemSelect,
	className,
	...rest
}: StepListProps) {
	const isSelectable = selectable ?? onItemSelect != null;
	return (
		<ol
			className={[styles.root, className ?? ""].filter(Boolean).join(" ")}
			data-rag-molecule="StepList"
			data-selectable={isSelectable ? "true" : "false"}
			{...rest}
		>
			{items.map((item) => {
				const active = item.id === activeItemId;
				const content = (
					<>
						<span className={styles.index}>{item.index}</span>
						<span className={styles.body}>
							<span className={styles.title}>{item.title}</span>
							{item.description && <span className={styles.description}>{item.description}</span>}
						</span>
						{item.meta && <span className={styles.meta}>{item.meta}</span>}
					</>
				);
				return (
					<li key={item.id} className={styles.item} data-active={active ? "true" : undefined}>
						{isSelectable ? (
							<button
								type="button"
								className={styles.button}
								onClick={() => onItemSelect?.(item.id)}
							>
								{content}
							</button>
						) : (
							<div className={styles.button}>{content}</div>
						)}
					</li>
				);
			})}
		</ol>
	);
}
