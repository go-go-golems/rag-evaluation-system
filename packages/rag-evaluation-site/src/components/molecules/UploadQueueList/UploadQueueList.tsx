import type { HTMLAttributes } from "react";
import { formatCmsAssetSize } from "../../../cms/fixtures";
import type { UploadItemStatus, UploadQueueItem } from "../../../cms/types";
import type { MeterBarTone } from "../../atoms";
import { IconButton, MeterBar } from "../../atoms";
import { Caption, CodeText, StatusText } from "../../foundation";
import styles from "./UploadQueueList.module.css";

export interface UploadQueueListProps extends HTMLAttributes<HTMLUListElement> {
	items: UploadQueueItem[];
	onCancel?: (itemId: string) => void;
	onRetry?: (itemId: string) => void;
	onDismiss?: (itemId: string) => void;
}

const STATUS_LABEL: Record<UploadItemStatus, string> = {
	queued: "pending",
	uploading: "running",
	done: "done",
	error: "error",
	canceled: "canceled",
};

function toneFor(status: UploadItemStatus): MeterBarTone {
	if (status === "error") return "danger";
	if (status === "done") return "success";
	return "accent";
}

export function UploadQueueList({
	items,
	onCancel,
	onRetry,
	onDismiss,
	className,
	...rest
}: UploadQueueListProps) {
	return (
		<ul
			className={[styles.root, className ?? ""].filter(Boolean).join(" ")}
			data-rag-molecule="UploadQueueList"
			{...rest}
		>
			{items.map((item) => {
				const cancelable = item.status === "queued" || item.status === "uploading";
				const dismissable = item.status === "done" || item.status === "canceled";
				return (
					<li key={item.id} className={styles.row} data-status={item.status}>
						<div className={styles.fileCell}>
							<CodeText className={styles.filename}>{item.filename}</CodeText>
							<Caption>{formatCmsAssetSize(item.size)}</Caption>
						</div>
						<MeterBar className={styles.meter} value={item.progress} tone={toneFor(item.status)} />
						<StatusText className={styles.status} status={STATUS_LABEL[item.status]} icon />
						<span className={styles.actions}>
							{cancelable && onCancel && (
								<IconButton
									size="large"
									variant="boxed"
									label={`Cancel upload of ${item.filename}`}
									onClick={() => onCancel(item.id)}
								>
									×
								</IconButton>
							)}
							{item.status === "error" && onRetry && (
								<IconButton
									size="large"
									variant="boxed"
									label={`Retry upload of ${item.filename}`}
									onClick={() => onRetry(item.id)}
								>
									↻
								</IconButton>
							)}
							{(dismissable || item.status === "error") && onDismiss && (
								<IconButton
									size="large"
									variant="boxed"
									label={`Dismiss ${item.filename}`}
									onClick={() => onDismiss(item.id)}
								>
									–
								</IconButton>
							)}
						</span>
						{item.error && (
							<Caption tone="danger" className={styles.error}>
								{item.error}
							</Caption>
						)}
					</li>
				);
			})}
		</ul>
	);
}
