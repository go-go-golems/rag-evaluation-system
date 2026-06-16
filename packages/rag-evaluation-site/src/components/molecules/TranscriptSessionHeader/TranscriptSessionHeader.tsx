import type { HTMLAttributes, ReactNode } from "react";
import { Caption, Text } from "../../foundation";
import styles from "./TranscriptSessionHeader.module.css";

export interface TranscriptSessionHeaderProps extends Omit<HTMLAttributes<HTMLElement>, "title"> {
	title: ReactNode;
	subtitle?: ReactNode;
	messageCount?: number;
	annotationCount?: number;
	tokenTotal?: number;
	rightSlot?: ReactNode;
}

export function TranscriptSessionHeader({
	title,
	subtitle,
	messageCount,
	annotationCount,
	tokenTotal,
	rightSlot,
	className,
	...rest
}: TranscriptSessionHeaderProps) {
	const stats = [
		messageCount == null ? undefined : `${messageCount} message${messageCount === 1 ? "" : "s"}`,
		annotationCount == null
			? undefined
			: `${annotationCount} annotation${annotationCount === 1 ? "" : "s"}`,
	]
		.filter(Boolean)
		.join(" · ");

	return (
		<header
			className={[styles.root, className ?? ""].filter(Boolean).join(" ")}
			data-rag-molecule="TranscriptSessionHeader"
			{...rest}
		>
			<div className={styles.titleRow}>
				<Text as="div" size="metric" weight="bold">
					Session: {title}
				</Text>
				{tokenTotal != null && <Caption>{tokenTotal.toLocaleString()} tok cumulative</Caption>}
			</div>
			{(subtitle || stats || rightSlot) && (
				<div className={styles.metaRow}>
					<Caption>{[subtitle, stats].filter(Boolean).join(" · ")}</Caption>
					{rightSlot && <div className={styles.rightSlot}>{rightSlot}</div>}
				</div>
			)}
		</header>
	);
}
