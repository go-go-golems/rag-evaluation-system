import type { HTMLAttributes, ReactNode } from "react";
import { useState } from "react";
import { cmsAssetMeta, formatCmsAssetSize } from "../../../cms/fixtures";
import type { CmsArticleSummary, CmsAsset, CmsContentStatus } from "../../../cms/types";
import { Button, MediaThumb, SelectInput, TextInput } from "../../atoms";
import { Caption } from "../../foundation";
import { FormRow, SplitPane, Stack } from "../../layout";
import { DocumentPreviewToolbar, MetadataGrid, TagListInput } from "../../molecules";
import { ConfirmDialog } from "../ConfirmDialog";
import { FormPanel, type FormPanelStatus } from "../FormPanel";
import styles from "./AssetDetailPanel.module.css";

export interface AssetDetailDraft {
	title: string;
	alt: string;
	tags: string[];
	status: CmsContentStatus;
}

export interface AssetDetailPanelProps extends HTMLAttributes<HTMLDivElement> {
	asset: CmsAsset;
	draft: AssetDetailDraft;
	onDraftChange: (draft: AssetDetailDraft) => void;
	onSave?: () => void;
	status?: FormPanelStatus;
	statusMessage?: ReactNode;
	usedBy?: CmsArticleSummary[];
	onUsageSelect?: (articleId: string) => void;
	onDownload?: () => void;
	onDelete?: () => void;
	tagSuggestions?: string[];
}

const STATUS_OPTIONS: CmsContentStatus[] = ["draft", "published", "scheduled", "archived"];

export function AssetDetailPanel({
	asset,
	draft,
	onDraftChange,
	onSave,
	status = "idle",
	statusMessage,
	usedBy,
	onUsageSelect,
	onDownload,
	onDelete,
	tagSuggestions,
	className,
	...rest
}: AssetDetailPanelProps) {
	const [confirmingDelete, setConfirmingDelete] = useState(false);
	const patch = (changes: Partial<AssetDetailDraft>) => onDraftChange({ ...draft, ...changes });

	const metadataItems = [
		{ key: "id", value: asset.id, copyValue: asset.id },
		{ key: "mime", value: asset.mime },
		{ key: "size", value: formatCmsAssetSize(asset.size) },
		...(asset.width && asset.height
			? [{ key: "dimensions", value: `${asset.width} × ${asset.height}` }]
			: []),
		{ key: "src", value: asset.src, copyValue: asset.src },
		{ key: "created", value: asset.createdAt.slice(0, 10) },
		{ key: "updated", value: asset.updatedAt.slice(0, 10) },
	];

	return (
		<div
			className={[styles.root, className ?? ""].filter(Boolean).join(" ")}
			data-rag-organism="AssetDetailPanel"
			data-rag-asset-id={asset.id}
			{...rest}
		>
			<DocumentPreviewToolbar
				file={asset.filename}
				format={cmsAssetMeta(asset)}
				onDownload={onDownload}
				rightSlot={
					onDelete && (
						<Button size="compact" onClick={() => setConfirmingDelete(true)}>
							× Delete
						</Button>
					)
				}
			/>
			<SplitPane
				className={styles.body}
				ratio="rightNarrow"
				divider
				gutter="md"
				left={
					<Stack gap="md">
						<MediaThumb
							className={styles.preview}
							src={asset.kind === "image" ? asset.src : undefined}
							alt={asset.alt ?? asset.title}
							aspect="natural"
							fit="contain"
							fallbackGlyph="▤"
							fallbackLabel={asset.kind === "file" ? asset.filename : undefined}
						/>
						<MetadataGrid items={metadataItems} density="compact" />
					</Stack>
				}
				right={
					<Stack gap="md">
						<FormPanel
							title="Details"
							status={status}
							statusMessage={statusMessage}
							submitLabel="Save"
							onSubmit={(event) => {
								event.preventDefault();
								onSave?.();
							}}
						>
							<FormRow
								label="Title"
								control={
									<TextInput
										value={draft.title}
										onChange={(event) => patch({ title: event.target.value })}
									/>
								}
							/>
							<FormRow
								label="Alt text"
								control={
									<TextInput
										value={draft.alt}
										onChange={(event) => patch({ alt: event.target.value })}
										placeholder="describe the image"
									/>
								}
								hint="used by screen readers and broken-image fallbacks"
							/>
							<FormRow
								label="Status"
								control={
									<SelectInput
										value={draft.status}
										onChange={(event) => patch({ status: event.target.value as CmsContentStatus })}
									>
										{STATUS_OPTIONS.map((option) => (
											<option key={option} value={option}>
												{option}
											</option>
										))}
									</SelectInput>
								}
							/>
							<FormRow
								label="Tags"
								control={
									<TagListInput
										tags={draft.tags}
										suggestions={tagSuggestions}
										onAdd={(tag) => patch({ tags: [...draft.tags, tag] })}
										onRemove={(tag) => patch({ tags: draft.tags.filter((entry) => entry !== tag) })}
									/>
								}
							/>
						</FormPanel>
						{usedBy && (
							<Stack gap="xs" data-rag-asset-usage>
								<Caption transform="uppercase">Used in {usedBy.length} article(s)</Caption>
								{usedBy.length === 0 ? (
									<Caption>Not referenced by any article.</Caption>
								) : (
									<ul className={styles.usageList}>
										{usedBy.map((article) => (
											<li key={article.id}>
												<button
													type="button"
													className={styles.usageLink}
													onClick={() => onUsageSelect?.(article.id)}
												>
													{article.title}
												</button>
											</li>
										))}
									</ul>
								)}
							</Stack>
						)}
					</Stack>
				}
			/>
			<ConfirmDialog
				open={confirmingDelete}
				title="Delete asset"
				message={`Delete “${asset.filename}” permanently?`}
				detail={
					usedBy && usedBy.length > 0
						? `Used in ${usedBy.length} article(s) — those references will break.`
						: undefined
				}
				confirmLabel="Delete"
				destructive
				onConfirm={() => {
					setConfirmingDelete(false);
					onDelete?.();
				}}
				onCancel={() => setConfirmingDelete(false)}
			/>
		</div>
	);
}
