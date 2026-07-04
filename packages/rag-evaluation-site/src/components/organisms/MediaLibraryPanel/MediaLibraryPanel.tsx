import type { HTMLAttributes, ReactNode } from "react";
import type { CmsAsset, CmsAssetKind, UploadQueueItem } from "../../../cms/types";
import { ContentStatusBadge } from "../../atoms";
import { Panel, ScrollRegion, Stack, TabList, TileGrid } from "../../layout";
import {
	AssetTile,
	EmptyState,
	FileDropZone,
	Pagination,
	SearchField,
	UploadQueueList,
} from "../../molecules";
import styles from "./MediaLibraryPanel.module.css";

export type MediaLibraryKindFilter = "all" | CmsAssetKind;
export type MediaLibrarySelectionMode = "none" | "single" | "multi";

export interface MediaLibraryPanelProps
	extends Omit<HTMLAttributes<HTMLDivElement>, "onSelect" | "title"> {
	assets: CmsAsset[];
	selectedAssetIds?: string[];
	selectionMode?: MediaLibrarySelectionMode;
	onAssetSelect?: (assetId: string) => void;
	onAssetOpen?: (assetId: string) => void;
	query?: string;
	onQueryChange?: (query: string) => void;
	onQuerySubmit?: (query: string) => void;
	kindFilter?: MediaLibraryKindFilter;
	onKindFilterChange?: (kind: MediaLibraryKindFilter) => void;
	page?: number;
	pageCount?: number;
	onPageChange?: (page: number) => void;
	onFilesSelected?: (files: File[]) => void;
	uploads?: UploadQueueItem[];
	onUploadCancel?: (itemId: string) => void;
	onUploadRetry?: (itemId: string) => void;
	onUploadDismiss?: (itemId: string) => void;
	showStatusBadges?: boolean;
	emptyMessage?: ReactNode;
	emptyAction?: ReactNode;
	title?: ReactNode;
	minTileWidth?: number;
}

const KIND_TABS: Array<{ id: MediaLibraryKindFilter; label: ReactNode }> = [
	{ id: "all", label: "All" },
	{ id: "image", label: "Images" },
	{ id: "file", label: "Files" },
];

export function MediaLibraryPanel({
	assets,
	selectedAssetIds = [],
	selectionMode = "single",
	onAssetSelect,
	onAssetOpen,
	query,
	onQueryChange,
	onQuerySubmit,
	kindFilter = "all",
	onKindFilterChange,
	page,
	pageCount,
	onPageChange,
	onFilesSelected,
	uploads,
	onUploadCancel,
	onUploadRetry,
	onUploadDismiss,
	showStatusBadges = true,
	emptyMessage = "No assets yet",
	emptyAction,
	title = "Media",
	minTileWidth = 160,
	className,
	...rest
}: MediaLibraryPanelProps) {
	const showToolbar = onQueryChange || (onPageChange && page != null && pageCount != null);
	const showUploads = uploads != null && uploads.length > 0;

	return (
		<Panel
			className={[styles.root, className ?? ""].filter(Boolean).join(" ")}
			title={title}
			fill
			data-rag-organism="MediaLibraryPanel"
			data-selection-mode={selectionMode}
			{...rest}
		>
			<Stack gap="sm" className={styles.body}>
				{onKindFilterChange && (
					<TabList
						items={KIND_TABS}
						activeId={kindFilter}
						onChange={onKindFilterChange}
						ariaLabel="Asset kind filter"
					/>
				)}
				{showToolbar && (
					<div className={styles.toolbar}>
						{onQueryChange && (
							<SearchField
								className={styles.search}
								value={query ?? ""}
								onValueChange={onQueryChange}
								onSubmit={onQuerySubmit}
								placeholder="search assets…"
							/>
						)}
						{onPageChange && page != null && pageCount != null && (
							<Pagination page={page} pageCount={pageCount} onPageChange={onPageChange} />
						)}
					</div>
				)}
				{onFilesSelected && (
					<FileDropZone
						className={styles.dropZone}
						title="Drop files here"
						description="or click to choose · images and documents"
						multiple
						onFilesSelected={onFilesSelected}
					/>
				)}
				{showUploads && (
					<UploadQueueList
						items={uploads}
						onCancel={onUploadCancel}
						onRetry={onUploadRetry}
						onDismiss={onUploadDismiss}
					/>
				)}
				{assets.length === 0 ? (
					<EmptyState glyph="▨" title={emptyMessage} action={emptyAction} />
				) : (
					<ScrollRegion className={styles.scroll}>
						<TileGrid minTileWidth={minTileWidth}>
							{assets.map((asset) => (
								<AssetTile
									key={asset.id}
									asset={asset}
									selected={selectedAssetIds.includes(asset.id)}
									onSelect={selectionMode === "none" ? undefined : onAssetSelect}
									onOpen={onAssetOpen}
									footerSlot={
										showStatusBadges ? (
											<ContentStatusBadge status={asset.status} icon={false} />
										) : undefined
									}
								/>
							))}
						</TileGrid>
					</ScrollRegion>
				)}
			</Stack>
		</Panel>
	);
}
