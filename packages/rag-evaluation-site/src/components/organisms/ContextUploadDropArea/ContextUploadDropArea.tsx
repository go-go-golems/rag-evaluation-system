import type { HTMLAttributes, ReactNode } from "react";
import { FileDropZone } from "../../molecules";

export interface ContextUploadDropAreaProps
	extends Omit<HTMLAttributes<HTMLDivElement>, "title" | "onDrop"> {
	title?: ReactNode;
	description?: ReactNode;
	accept?: string;
	disabled?: boolean;
	active?: boolean;
	onFilesSelected?: (files: File[]) => void;
}

export function ContextUploadDropArea({
	title = "Drop a .json file here",
	description = "or paste below · max 200k tokens",
	accept = "application/json,.json",
	disabled = false,
	active = false,
	onFilesSelected,
	...rest
}: ContextUploadDropAreaProps) {
	return (
		<FileDropZone
			title={title}
			description={description}
			accept={accept}
			disabled={disabled}
			active={active}
			onFilesSelected={onFilesSelected}
			inputAriaLabel="Choose context-window JSON file"
			data-rag-organism="ContextUploadDropArea"
			{...rest}
		/>
	);
}
