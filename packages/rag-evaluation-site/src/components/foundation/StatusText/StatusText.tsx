import type { HTMLAttributes, ReactNode } from "react";
import styles from "./StatusText.module.css";

export type RagStatus =
	| "pending"
	| "ready"
	| "running"
	| "succeeded"
	| "done"
	| "partial"
	| "warning"
	| "failed"
	| "error"
	| "canceled";

export interface StatusTextProps extends HTMLAttributes<HTMLSpanElement> {
	status: RagStatus | string;
	icon?: boolean;
	children?: ReactNode;
}

const statusIcon: Record<string, string> = {
	pending: "◌",
	ready: "◌",
	running: "●",
	succeeded: "✔",
	done: "●",
	partial: "◐",
	warning: "⚠",
	failed: "✘",
	error: "✘",
	canceled: "⊘",
};

export function StatusText({
	status,
	icon = false,
	className,
	children,
	...rest
}: StatusTextProps) {
	const normalized = status.toLowerCase();
	const statusClass = styles[normalized] ?? "";
	const label = children ?? status;
	return (
		<span
			className={[styles.root, statusClass, className ?? ""].filter(Boolean).join(" ")}
			{...rest}
		>
			{icon ? `${statusIcon[normalized] ?? "?"} ` : ""}
			{label}
		</span>
	);
}
