import type { HTMLAttributes, ReactNode } from "react";
import styles from "./FormRow.module.css";

export type FormRowOrientation = "inline" | "stacked";

export interface FormRowProps extends HTMLAttributes<HTMLDivElement> {
	label: ReactNode;
	control: ReactNode;
	description?: ReactNode;
	hint?: ReactNode;
	error?: ReactNode;
	success?: ReactNode;
	counter?: ReactNode;
	required?: boolean;
	orientation?: FormRowOrientation;
}

export function FormRow({
	label,
	control,
	description,
	hint,
	error,
	success,
	counter,
	required = false,
	orientation = "inline",
	className,
	...rest
}: FormRowProps) {
	const state = error ? "error" : success ? "success" : undefined;
	const hasHintRow = Boolean(hint || counter);

	return (
		<div
			className={[
				styles.root,
				orientation === "stacked" ? styles.stacked : styles.inline,
				className ?? "",
			]
				.filter(Boolean)
				.join(" ")}
			data-rag-layout="FormRow"
			data-orientation={orientation}
			data-state={state}
			{...rest}
		>
			<div className={styles.label}>
				{label}
				{required && (
					<span className={styles.required} aria-hidden="true">
						{" "}
						*
					</span>
				)}
			</div>
			<div className={styles.control}>{control}</div>
			{description && <div className={styles.description}>{description}</div>}
			{hasHintRow && (
				<div className={styles.hintRow}>
					{hint && <span className={styles.hint}>{hint}</span>}
					{counter && <span className={styles.counter}>{counter}</span>}
				</div>
			)}
			{success && <div className={styles.success}>{success}</div>}
			{error && (
				<div className={styles.error} role="alert">
					{error}
				</div>
			)}
		</div>
	);
}
