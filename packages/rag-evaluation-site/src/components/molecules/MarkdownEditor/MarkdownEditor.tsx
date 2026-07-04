import type { ReactNode } from "react";
import { useRef } from "react";
import { Button, TextareaInput } from "../../atoms";
import { Caption } from "../../foundation";
import styles from "./MarkdownEditor.module.css";

export interface MarkdownEditorProps {
	value: string;
	onValueChange: (value: string) => void;
	onInsertAsset?: () => void;
	minRows?: number;
	maxLength?: number;
	disabled?: boolean;
	toolbarSlot?: ReactNode;
	textareaAriaLabel?: string;
	/** Form-post participation: the textarea carries this name in native submits. */
	name?: string;
	className?: string;
}

interface WrapSpec {
	label: ReactNode;
	title: string;
	before: string;
	after: string;
	block?: boolean;
}

const WRAPS: WrapSpec[] = [
	{ label: <strong>B</strong>, title: "Bold", before: "**", after: "**" },
	{ label: "`code`", title: "Inline code", before: "`", after: "`" },
	{ label: "[link]", title: "Link", before: "[", after: "](https://)" },
	{ label: "H2", title: "Heading", before: "## ", after: "", block: true },
	{ label: "• list", title: "List item", before: "- ", after: "", block: true },
];

export function MarkdownEditor({
	value,
	onValueChange,
	onInsertAsset,
	minRows = 16,
	maxLength,
	disabled = false,
	toolbarSlot,
	textareaAriaLabel = "Markdown source",
	name,
	className,
}: MarkdownEditorProps) {
	const textareaRef = useRef<HTMLTextAreaElement>(null);

	const applyWrap = (wrap: WrapSpec) => {
		const el = textareaRef.current;
		if (!el) return;
		const start = el.selectionStart ?? value.length;
		const end = el.selectionEnd ?? value.length;
		let insertAt = start;
		let before = wrap.before;
		if (wrap.block) {
			insertAt = value.lastIndexOf("\n", Math.max(0, start - 1)) + 1;
			before = wrap.before;
			const next = `${value.slice(0, insertAt)}${before}${value.slice(insertAt)}`;
			onValueChange(next);
			requestAnimationFrame(() => {
				el.focus();
				el.setSelectionRange(start + before.length, end + before.length);
			});
			return;
		}
		const next = `${value.slice(0, start)}${before}${value.slice(start, end)}${wrap.after}${value.slice(end)}`;
		onValueChange(next);
		requestAnimationFrame(() => {
			el.focus();
			el.setSelectionRange(start + before.length, end + before.length);
		});
	};

	return (
		<div
			className={[styles.root, className ?? ""].filter(Boolean).join(" ")}
			data-rag-molecule="MarkdownEditor"
			data-disabled={disabled || undefined}
		>
			<div className={styles.toolbar}>
				{WRAPS.map((wrap) => (
					<Button
						key={wrap.title}
						size="compact"
						title={wrap.title}
						disabled={disabled}
						onClick={() => applyWrap(wrap)}
					>
						{wrap.label}
					</Button>
				))}
				{onInsertAsset && (
					<Button
						size="compact"
						title="Insert image from media library"
						disabled={disabled}
						onClick={onInsertAsset}
					>
						▨ image
					</Button>
				)}
				{toolbarSlot}
			</div>
			<TextareaInput
				ref={textareaRef}
				name={name}
				className={styles.textarea}
				value={value}
				rows={minRows}
				maxLength={maxLength}
				disabled={disabled}
				aria-label={textareaAriaLabel}
				onChange={(event) => onValueChange(event.target.value)}
			/>
			{maxLength != null && (
				<Caption className={styles.counter}>
					{value.length} / {maxLength}
				</Caption>
			)}
		</div>
	);
}
