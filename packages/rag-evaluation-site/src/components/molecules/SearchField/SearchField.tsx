import type { InputHTMLAttributes } from "react";
import { IconButton, TextInput } from "../../atoms";
import styles from "./SearchField.module.css";

export interface SearchFieldProps
	extends Omit<InputHTMLAttributes<HTMLInputElement>, "onChange" | "onSubmit" | "value"> {
	value: string;
	onValueChange: (value: string) => void;
	onSubmit?: (value: string) => void;
	onClear?: () => void;
}

export function SearchField({
	value,
	onValueChange,
	onSubmit,
	onClear,
	placeholder = "search…",
	disabled,
	className,
	...rest
}: SearchFieldProps) {
	const clear = () => {
		onValueChange("");
		onClear?.();
	};

	return (
		<div
			className={[styles.root, className ?? ""].filter(Boolean).join(" ")}
			data-rag-molecule="SearchField"
			data-disabled={disabled || undefined}
		>
			<span className={styles.glyph} aria-hidden="true">
				⌕
			</span>
			<TextInput
				className={styles.input}
				type="search"
				value={value}
				placeholder={placeholder}
				disabled={disabled}
				onChange={(event) => onValueChange(event.target.value)}
				onKeyDown={(event) => {
					if (event.key === "Enter") onSubmit?.(value);
					if (event.key === "Escape" && value) clear();
				}}
				{...rest}
			/>
			{value && !disabled && (
				<IconButton className={styles.clear} label="Clear search" onClick={clear}>
					×
				</IconButton>
			)}
		</div>
	);
}
