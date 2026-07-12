import type { HTMLAttributes, ReactNode } from "react";
import { Button } from "../../atoms";
import { Inline, Stack } from "../../layout";
import { FileDropZone } from "../../molecules";

export interface ContextUploadDropAreaProps
	extends Omit<HTMLAttributes<HTMLDivElement>, "title" | "onDrop"> {
	title?: ReactNode;
	description?: ReactNode;
	accept?: string;
	disabled?: boolean;
	active?: boolean;
	onFilesSelected?: (files: File[]) => void;
	items?: Array<{ id: string; label?: ReactNode }>;
	onDelete?: (itemId: string) => void;
}

export function ContextUploadDropArea({
	title = "Drop a .json file here",
	description = "or paste below · max 200k tokens",
	accept = "application/json,.json",
	disabled = false,
	active = false,
	onFilesSelected,
	items = [],
	onDelete,
	...rest
}: ContextUploadDropAreaProps) {
	return (
		<Stack gap="sm" data-rag-organism="ContextUploadDropArea" {...rest}>
			<FileDropZone
				title={title}
				description={description}
				accept={accept}
				disabled={disabled}
				active={active}
				onFilesSelected={onFilesSelected}
				inputAriaLabel="Choose context-window JSON file"
			/>
			{items.map((item) => (
				<Inline key={item.id} justify="between" style={{ alignItems: "center" }}>
					<span>{item.label ?? item.id}</span>
					{onDelete && (
						<Button size="compact" onClick={() => onDelete(item.id)}>
							Delete
						</Button>
					)}
				</Inline>
			))}
		</Stack>
	);
}
