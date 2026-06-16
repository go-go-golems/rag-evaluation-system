import type { HTMLAttributes, ReactNode } from "react";
import styles from "./KeyValueStrip.module.css";

export interface KeyValueStripItem {
	key: ReactNode;
	value: ReactNode;
}
export interface KeyValueStripProps extends HTMLAttributes<HTMLDListElement> {
	items: KeyValueStripItem[];
}

export function KeyValueStrip({ items, className, ...rest }: KeyValueStripProps) {
	return (
		<dl
			className={[styles.root, className ?? ""].filter(Boolean).join(" ")}
			data-rag-molecule="KeyValueStrip"
			{...rest}
		>
			{items.map((item, index) => (
				<div className={styles.cell} key={index}>
					<dt>{item.key}</dt>
					<dd>{item.value}</dd>
				</div>
			))}
		</dl>
	);
}
