import type { ReactNode } from "react";
import styles from "./MetadataGrid.module.css";

export interface MetadataGridItem {
	key: ReactNode;
	value: ReactNode;
	copyValue?: string;
}

export interface MetadataGridProps {
	items: MetadataGridItem[];
	density?: "normal" | "compact";
	className?: string;
	onCopy?: (value: string) => void;
}

export function MetadataGrid({ items, density = "normal", className, onCopy }: MetadataGridProps) {
	const copy = (value: string) => {
		if (onCopy) {
			onCopy(value);
			return;
		}
		navigator.clipboard?.writeText(value).catch(() => {});
	};

	return (
		<dl
			className={[styles.root, density === "compact" ? styles.compact : "", className ?? ""]
				.filter(Boolean)
				.join(" ")}
			data-rag-component="MetadataGrid"
		>
			{items.map((item, index) => (
				<div key={index} className={styles.row}>
					<dt className={styles.key}>{item.key}</dt>
					<dd className={styles.value}>
						{item.value}
						{item.copyValue && (
							<button
								type="button"
								className={styles.copyButton}
								onClick={() => copy(item.copyValue!)}
							>
								⧉
							</button>
						)}
					</dd>
				</div>
			))}
		</dl>
	);
}
