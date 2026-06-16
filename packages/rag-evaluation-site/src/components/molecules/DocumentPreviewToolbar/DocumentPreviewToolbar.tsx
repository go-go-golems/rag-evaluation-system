import type { HTMLAttributes, ReactNode } from "react";
import { Button } from "../../atoms";
import { Caption, CodeText } from "../../foundation";
import styles from "./DocumentPreviewToolbar.module.css";

export interface DocumentPreviewToolbarProps extends HTMLAttributes<HTMLDivElement> {
	file: ReactNode;
	format?: ReactNode;
	size?: ReactNode;
	onDownload?: () => void;
	downloadLabel?: ReactNode;
	onPrint?: () => void;
	printLabel?: ReactNode;
	rightSlot?: ReactNode;
}

export function DocumentPreviewToolbar({
	file,
	format,
	size,
	onDownload,
	downloadLabel = "⤓ Download",
	onPrint,
	printLabel = "Print",
	rightSlot,
	className,
	...rest
}: DocumentPreviewToolbarProps) {
	return (
		<div
			className={[styles.root, className ?? ""].filter(Boolean).join(" ")}
			data-rag-molecule="DocumentPreviewToolbar"
			{...rest}
		>
			<div className={styles.meta}>
				<CodeText>{file}</CodeText>
				{format && (
					<span className={styles.badge}>
						<Caption tone="inherit">{format}</Caption>
					</span>
				)}
				{size && (
					<span className={styles.badge}>
						<Caption tone="inherit">{size}</Caption>
					</span>
				)}
			</div>
			<div className={styles.actions}>
				{rightSlot}
				{onPrint && (
					<Button size="compact" onClick={onPrint}>
						{printLabel}
					</Button>
				)}
				{onDownload && (
					<Button size="compact" onClick={onDownload}>
						{downloadLabel}
					</Button>
				)}
			</div>
		</div>
	);
}
