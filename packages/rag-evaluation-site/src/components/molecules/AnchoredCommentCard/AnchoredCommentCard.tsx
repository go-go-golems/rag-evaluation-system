import type { HTMLAttributes } from "react";
import type { AnchoredComment } from "../../../context";
import { Caption } from "../../foundation";
import styles from "./AnchoredCommentCard.module.css";

export interface AnchoredCommentCardProps extends HTMLAttributes<HTMLDivElement> {
	comment: AnchoredComment;
	index?: number;
	selected?: boolean;
	compact?: boolean;
	onDismiss?: () => void;
}

export function AnchoredCommentCard({
	comment,
	index,
	selected = false,
	compact = false,
	onDismiss,
	className,
	...rest
}: AnchoredCommentCardProps) {
	return (
		<article
			className={[
				styles.root,
				selected ? styles.selected : "",
				compact ? styles.compact : "",
				comment.status === "resolved" ? styles.resolved : "",
				className ?? "",
			]
				.filter(Boolean)
				.join(" ")}
			data-rag-molecule="AnchoredCommentCard"
			data-status={comment.status ?? "open"}
			{...rest}
		>
			<header className={styles.header}>
				{index != null && <span className={styles.index}>{index}</span>}
				<span className={styles.author}>{comment.author.toUpperCase()}</span>
				<span className={styles.spacer} />
				{comment.time && <span className={styles.time}>{comment.time}</span>}
				{onDismiss && (
					<button
						type="button"
						className={styles.dismiss}
						onClick={onDismiss}
						aria-label="Dismiss comment"
					>
						✕
					</button>
				)}
			</header>
			<p className={styles.text}>{comment.text}</p>
			{comment.status === "resolved" && <Caption tone="success">resolved</Caption>}
		</article>
	);
}
