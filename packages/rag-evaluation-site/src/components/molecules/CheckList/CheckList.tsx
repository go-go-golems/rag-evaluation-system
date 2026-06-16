import type { HTMLAttributes, ReactNode } from "react";
import styles from "./CheckList.module.css";

export interface CheckListItem {
	id?: string;
	text: ReactNode;
}
export interface CheckListProps extends HTMLAttributes<HTMLUListElement> {
	items: Array<CheckListItem | ReactNode>;
	marker?: ReactNode;
}

export function CheckList({ items, marker = "✔", className, ...rest }: CheckListProps) {
	return (
		<ul
			className={[styles.root, className ?? ""].filter(Boolean).join(" ")}
			data-rag-molecule="CheckList"
			{...rest}
		>
			{items.map((item, index) => {
				const text = typeof item === "object" && item !== null && "text" in item ? item.text : item;
				const key =
					typeof item === "object" && item !== null && "id" in item && item.id ? item.id : index;
				return (
					<li key={key} className={styles.item}>
						<span className={styles.marker}>{marker}</span>
						<span>{text}</span>
					</li>
				);
			})}
		</ul>
	);
}
