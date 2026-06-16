import type { HTMLAttributes, ReactNode } from "react";
import styles from "./PersonSummary.module.css";

export interface PersonSummaryProps extends HTMLAttributes<HTMLDivElement> {
	name: ReactNode;
	subtitle?: ReactNode;
	bio?: ReactNode;
	avatar?: ReactNode;
}

export function PersonSummary({
	name,
	subtitle,
	bio,
	avatar,
	className,
	...rest
}: PersonSummaryProps) {
	const initial = typeof name === "string" && name.length > 0 ? name[0] : "?";
	return (
		<div
			className={[styles.root, className ?? ""].filter(Boolean).join(" ")}
			data-rag-molecule="PersonSummary"
			{...rest}
		>
			<div className={styles.avatar}>{avatar ?? initial}</div>
			<div className={styles.body}>
				<div className={styles.name}>{name}</div>
				{subtitle && <div className={styles.role}>{subtitle}</div>}
				{bio && <div className={styles.bio}>{bio}</div>}
			</div>
		</div>
	);
}
