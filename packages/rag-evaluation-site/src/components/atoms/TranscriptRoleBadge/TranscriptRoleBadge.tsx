import type { HTMLAttributes } from "react";
import type { TranscriptRole } from "../../../context";
import styles from "./TranscriptRoleBadge.module.css";
export interface TranscriptRoleBadgeProps extends HTMLAttributes<HTMLSpanElement> {
	role: TranscriptRole;
	name?: string;
}
export function TranscriptRoleBadge({ role, name, className, ...rest }: TranscriptRoleBadgeProps) {
	return (
		<span
			className={[styles.root, styles[role] ?? styles.other, className ?? ""]
				.filter(Boolean)
				.join(" ")}
			data-rag-atom="TranscriptRoleBadge"
			data-role={role}
			{...rest}
		>
			{name ?? role}
		</span>
	);
}
