import type { HTMLAttributes, ReactNode } from "react";
import { Button } from "../../atoms";
import { Caption, Text } from "../../foundation";
import styles from "./DocumentListPanel.module.css";

export interface DocumentListItem {
	id: string;
	title: ReactNode;
	format?: ReactNode;
	size?: ReactNode;
	description?: ReactNode;
	icon?: ReactNode;
	disabled?: boolean;
}

export interface DocumentListPanelProps
	extends Omit<HTMLAttributes<HTMLDivElement>, "onSelect" | "title"> {
	title: ReactNode;
	description?: ReactNode;
	items: DocumentListItem[];
	selectedItemId?: string;
	onItemSelect?: (itemId: string) => void;
	onDownloadAll?: () => void;
	downloadAllLabel?: ReactNode;
	emptyMessage?: ReactNode;
	showItemDescriptions?: boolean;
}

function defaultIcon(format?: ReactNode) {
	const value = String(format ?? "").toLowerCase();
	if (value.includes("pdf")) return "▤";
	if (value.includes("json")) return "{ }";
	if (value.includes("markdown") || value.includes("md")) return "¶";
	return "□";
}

export function DocumentListPanel({
	title,
	description,
	items,
	selectedItemId,
	onItemSelect,
	onDownloadAll,
	downloadAllLabel = "⤓ Download all (.zip)",
	emptyMessage = "No documents available.",
	showItemDescriptions = false,
	className,
	...rest
}: DocumentListPanelProps) {
	const hasHeader = Boolean(title || description);

	return (
		<div
			className={[styles.root, className ?? ""].filter(Boolean).join(" ")}
			data-rag-molecule="DocumentListPanel"
			{...rest}
		>
			{hasHeader && (
				<header className={styles.header}>
					{title && (
						<Text as="div" size="label" weight="bold" className={styles.title}>
							{title}
						</Text>
					)}
					{description && <Caption className={styles.description}>{description}</Caption>}
				</header>
			)}
			<div className={styles.list} role="listbox" aria-label="Documents">
				{items.length === 0 && (
					<div className={styles.empty}>
						<Caption>{emptyMessage}</Caption>
					</div>
				)}
				{items.map((item) => {
					const active = item.id === selectedItemId;
					return (
						<button
							key={item.id}
							type="button"
							className={styles.item}
							data-active={active ? "true" : undefined}
							role="option"
							aria-selected={active}
							disabled={item.disabled}
							onClick={() => onItemSelect?.(item.id)}
						>
							<span className={styles.icon} aria-hidden="true">
								{item.icon ?? defaultIcon(item.format)}
							</span>
							<span className={styles.itemBody}>
								<span className={styles.itemTitle}>{item.title}</span>
								<span className={styles.itemMeta}>
									{[item.format, item.size].filter(Boolean).join(" · ")}
								</span>
								{showItemDescriptions && item.description && (
									<span className={styles.itemDescription}>{item.description}</span>
								)}
							</span>
						</button>
					);
				})}
			</div>
			{onDownloadAll && (
				<footer className={styles.footer}>
					<Button variant="primary" onClick={onDownloadAll} className={styles.downloadButton}>
						{downloadAllLabel}
					</Button>
				</footer>
			)}
		</div>
	);
}
