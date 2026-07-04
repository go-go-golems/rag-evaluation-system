import type { KeyboardEvent } from "react";
import { useId, useState } from "react";
import { Tag, TextInput } from "../../atoms";
import styles from "./TagListInput.module.css";

export interface TagListInputProps {
	tags: string[];
	onAdd?: (tag: string) => void;
	onRemove?: (tag: string) => void;
	suggestions?: string[];
	placeholder?: string;
	/** Form-post participation: emits a hidden input carrying tags.join(","). */
	name?: string;
	disabled?: boolean;
	className?: string;
}

export function TagListInput({
	tags,
	onAdd,
	onRemove,
	suggestions,
	placeholder = "add tag…",
	name,
	disabled = false,
	className,
}: TagListInputProps) {
	const [draft, setDraft] = useState("");
	const listId = useId();

	const commit = () => {
		const value = draft.trim().replace(/,+$/, "");
		if (!value) return;
		if (!tags.includes(value)) onAdd?.(value);
		setDraft("");
	};

	const onKeyDown = (event: KeyboardEvent<HTMLInputElement>) => {
		if (event.key === "Enter" || event.key === ",") {
			event.preventDefault();
			commit();
		} else if (event.key === "Backspace" && !draft && tags.length > 0) {
			onRemove?.(tags[tags.length - 1] as string);
		}
	};

	return (
		<div
			className={[styles.root, className ?? ""].filter(Boolean).join(" ")}
			data-rag-molecule="TagListInput"
			data-disabled={disabled || undefined}
		>
			{name && <input type="hidden" name={name} value={tags.join(",")} />}
			{tags.map((tag) => (
				<Tag
					key={tag}
					label={tag}
					disabled={disabled}
					onRemove={onRemove ? () => onRemove(tag) : undefined}
				/>
			))}
			{onAdd && (
				<>
					<TextInput
						className={styles.input}
						value={draft}
						placeholder={placeholder}
						disabled={disabled}
						list={suggestions ? listId : undefined}
						aria-label="Add tag"
						onChange={(event) => setDraft(event.target.value)}
						onKeyDown={onKeyDown}
						onBlur={commit}
					/>
					{suggestions && (
						<datalist id={listId}>
							{suggestions.map((suggestion) => (
								<option key={suggestion} value={suggestion} />
							))}
						</datalist>
					)}
				</>
			)}
		</div>
	);
}
