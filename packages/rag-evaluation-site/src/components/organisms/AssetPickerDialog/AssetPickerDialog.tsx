import type { ReactNode } from "react";
import { useState } from "react";
import type { CmsAsset } from "../../../cms/types";
import { Button } from "../../atoms";
import { DialogShell, type DialogShellMode } from "../../layout";
import type { MediaLibraryKindFilter } from "../MediaLibraryPanel";
import { MediaLibraryPanel } from "../MediaLibraryPanel";
import styles from "./AssetPickerDialog.module.css";

export interface AssetPickerDialogProps {
	open: boolean;
	assets: CmsAsset[];
	onConfirm: (asset: CmsAsset) => void;
	onCancel: () => void;
	title?: ReactNode;
	confirmLabel?: ReactNode;
	query?: string;
	onQueryChange?: (query: string) => void;
	kindFilter?: MediaLibraryKindFilter;
	onKindFilterChange?: (kind: MediaLibraryKindFilter) => void;
	page?: number;
	pageCount?: number;
	onPageChange?: (page: number) => void;
	mode?: DialogShellMode;
	className?: string;
}

export function AssetPickerDialog({
	open,
	assets,
	onConfirm,
	onCancel,
	title = "Choose an asset",
	confirmLabel = "Use asset",
	query,
	onQueryChange,
	kindFilter,
	onKindFilterChange,
	page,
	pageCount,
	onPageChange,
	mode = "modal",
	className,
}: AssetPickerDialogProps) {
	const [selectedId, setSelectedId] = useState<string>();
	const selected = assets.find((asset) => asset.id === selectedId);

	const confirm = (asset: CmsAsset | undefined) => {
		if (!asset) return;
		onConfirm(asset);
		setSelectedId(undefined);
	};

	return (
		<DialogShell
			open={open}
			onClose={onCancel}
			title={title}
			size="lg"
			mode={mode}
			className={className}
			data-rag-organism="AssetPickerDialog"
			footer={
				<>
					<Button size="compact" onClick={onCancel}>
						Cancel
					</Button>
					<Button
						size="compact"
						variant="primary"
						disabled={!selected}
						onClick={() => confirm(selected)}
					>
						{confirmLabel}
					</Button>
				</>
			}
		>
			<MediaLibraryPanel
				className={styles.library}
				title={undefined}
				assets={assets}
				selectionMode="single"
				selectedAssetIds={selectedId ? [selectedId] : []}
				onAssetSelect={setSelectedId}
				onAssetOpen={(assetId) => confirm(assets.find((asset) => asset.id === assetId))}
				showStatusBadges={false}
				query={query}
				onQueryChange={onQueryChange}
				kindFilter={kindFilter}
				onKindFilterChange={onKindFilterChange}
				page={page}
				pageCount={pageCount}
				onPageChange={onPageChange}
				emptyMessage="No assets available"
			/>
		</DialogShell>
	);
}
