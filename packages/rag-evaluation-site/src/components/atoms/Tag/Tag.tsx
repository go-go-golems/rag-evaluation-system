import type { HTMLAttributes } from "react";
import { IconButton } from "../IconButton";
import styles from "./Tag.module.css";

export interface TagProps extends HTMLAttributes<HTMLSpanElement> {
	label: string;
	selected?: boolean;
	disabled?: boolean;
	onRemove?: () => void;
}

export function Tag({
	label,
	selected = false,
	disabled = false,
	onRemove,
	className,
	...rest
}: TagProps) {
	return (
		<span
			className={[
				styles.root,
				selected ? styles.selected : "",
				disabled ? styles.disabled : "",
				className ?? "",
			]
				.filter(Boolean)
				.join(" ")}
			data-rag-atom="Tag"
			data-active={selected || undefined}
			data-disabled={disabled || undefined}
			{...rest}
		>
			<span className={styles.label}>{label}</span>
			{onRemove && (
				<IconButton
					className={styles.remove}
					label={`Remove tag ${label}`}
					disabled={disabled}
					onClick={onRemove}
				>
					×
				</IconButton>
			)}
		</span>
	);
}
