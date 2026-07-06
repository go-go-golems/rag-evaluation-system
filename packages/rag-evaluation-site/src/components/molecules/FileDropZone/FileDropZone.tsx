import {
	type DragEvent,
	type HTMLAttributes,
	type KeyboardEvent,
	type MouseEvent,
	type ReactNode,
	useRef,
	useState,
} from "react";
import { UploadGlyph } from "../../atoms";
import { Caption, Text } from "../../foundation";
import styles from "./FileDropZone.module.css";

export interface FileDropZoneProps
	extends Omit<HTMLAttributes<HTMLDivElement>, "title" | "onDrop"> {
	title?: ReactNode;
	description?: ReactNode;
	icon?: ReactNode;
	accept?: string;
	multiple?: boolean;
	disabled?: boolean;
	active?: boolean;
	inputAriaLabel?: string;
	onFilesSelected?: (files: File[]) => void;
}

function filesFromList(files: FileList | null) {
	return files ? Array.from(files) : [];
}

export function FileDropZone({
	title = "Drop a file here",
	description = "or choose a file",
	icon,
	accept,
	multiple = false,
	disabled = false,
	active = false,
	inputAriaLabel = "Choose file to upload",
	onFilesSelected,
	className,
	onDragEnter,
	onDragOver,
	onDragLeave,
	onKeyDown,
	...rest
}: FileDropZoneProps) {
	const inputRef = useRef<HTMLInputElement>(null);
	const [dragActive, setDragActive] = useState(false);
	const isActive = active || dragActive;

	function openFilePicker() {
		if (!disabled) inputRef.current?.click();
	}

	function handleFiles(files: File[]) {
		if (!disabled && files.length > 0) onFilesSelected?.(files);
	}

	function handleDragEnter(event: DragEvent<HTMLDivElement>) {
		onDragEnter?.(event);
		if (event.defaultPrevented || disabled) return;
		event.preventDefault();
		setDragActive(true);
	}

	function handleDragOver(event: DragEvent<HTMLDivElement>) {
		onDragOver?.(event);
		if (event.defaultPrevented || disabled) return;
		event.preventDefault();
		setDragActive(true);
	}

	function handleDragLeave(event: DragEvent<HTMLDivElement>) {
		onDragLeave?.(event);
		if (disabled) return;
		if (!event.currentTarget.contains(event.relatedTarget as Node | null)) setDragActive(false);
	}

	function handleDrop(event: DragEvent<HTMLDivElement>) {
		event.preventDefault();
		setDragActive(false);
		handleFiles(filesFromList(event.dataTransfer.files));
	}

	function handleClick(event: MouseEvent<HTMLDivElement>) {
		if (event.target === inputRef.current) return;
		openFilePicker();
	}

	function handleKeyDown(event: KeyboardEvent<HTMLDivElement>) {
		onKeyDown?.(event);
		if (event.defaultPrevented || disabled) return;
		if (event.key === "Enter" || event.key === " ") {
			event.preventDefault();
			openFilePicker();
		}
	}

	return (
		<div
			className={[styles.root, className ?? ""].filter(Boolean).join(" ")}
			role="button"
			tabIndex={disabled ? undefined : 0}
			aria-disabled={disabled || undefined}
			data-active={isActive ? "true" : undefined}
			data-disabled={disabled ? "true" : undefined}
			data-rag-molecule="FileDropZone"
			onClick={handleClick}
			onDragEnter={handleDragEnter}
			onDragOver={handleDragOver}
			onDragLeave={handleDragLeave}
			onDrop={handleDrop}
			onKeyDown={handleKeyDown}
			{...rest}
		>
			<input
				ref={inputRef}
				className={styles.input}
				type="file"
				accept={accept}
				multiple={multiple}
				disabled={disabled}
				aria-label={inputAriaLabel}
				tabIndex={-1}
				onChange={(event) => {
					handleFiles(filesFromList(event.currentTarget.files));
					event.currentTarget.value = "";
				}}
			/>
			<div className={styles.icon} aria-hidden="true">
				{icon ?? <UploadGlyph title="" />}
			</div>
			<Text as="div" size="body" weight="bold" className={styles.title}>
				{title}
			</Text>
			{description && <Caption className={styles.description}>{description}</Caption>}
		</div>
	);
}
