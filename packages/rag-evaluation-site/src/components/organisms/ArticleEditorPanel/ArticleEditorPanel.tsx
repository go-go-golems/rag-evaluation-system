import type { ReactNode } from "react";
import type { CmsAsset, CmsContentStatus } from "../../../cms/types";
import { Button, MediaThumb, SelectInput, TextInput } from "../../atoms";
import { FormRow, ScrollRegion, SplitPane } from "../../layout";
import { MarkdownArticle, MarkdownEditor, TagListInput } from "../../molecules";
import { FormPanel, type FormPanelStatus } from "../FormPanel";
import styles from "./ArticleEditorPanel.module.css";

export interface ArticleEditorDraft {
	title: string;
	slug: string;
	status: CmsContentStatus;
	tags: string[];
	excerpt?: string;
	coverAssetId?: string;
	body: string;
}

export interface ArticleEditorPanelProps {
	draft: ArticleEditorDraft;
	onDraftChange: (draft: ArticleEditorDraft) => void;
	onSave?: () => void;
	onPublish?: () => void;
	status?: FormPanelStatus;
	statusMessage?: ReactNode;
	coverAsset?: CmsAsset;
	onPickCoverAsset?: () => void;
	onInsertAsset?: () => void;
	preview?: "live" | "hidden";
	tagSuggestions?: string[];
	title?: ReactNode;
	className?: string;
}

const STATUS_OPTIONS: CmsContentStatus[] = ["draft", "published", "scheduled", "archived"];

export function ArticleEditorPanel({
	draft,
	onDraftChange,
	onSave,
	onPublish,
	status = "idle",
	statusMessage,
	coverAsset,
	onPickCoverAsset,
	onInsertAsset,
	preview = "live",
	tagSuggestions,
	title = "Article",
	className,
}: ArticleEditorPanelProps) {
	const patch = (changes: Partial<ArticleEditorDraft>) => onDraftChange({ ...draft, ...changes });

	const editor = (
		<MarkdownEditor
			value={draft.body}
			onValueChange={(body) => patch({ body })}
			onInsertAsset={onInsertAsset}
			minRows={16}
			disabled={status === "saving"}
		/>
	);

	return (
		<FormPanel
			className={[styles.root, className ?? ""].filter(Boolean).join(" ")}
			title={title}
			status={status}
			statusMessage={statusMessage}
			submitLabel="Save"
			actions={
				onPublish && (
					<Button size="compact" disabled={status === "saving"} onClick={onPublish}>
						● Publish
					</Button>
				)
			}
			onSubmit={(event) => {
				event.preventDefault();
				onSave?.();
			}}
			data-rag-organism="ArticleEditorPanel"
		>
			<FormRow
				label="Title"
				required
				control={
					<TextInput
						value={draft.title}
						onChange={(event) => patch({ title: event.target.value })}
						placeholder="Article title"
					/>
				}
			/>
			<FormRow
				label="Slug"
				control={
					<TextInput
						className={styles.slugInput}
						value={draft.slug}
						onChange={(event) => patch({ slug: event.target.value })}
						placeholder="article-slug"
					/>
				}
				hint="lowercase, hyphenated; used in URLs"
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
			{onPickCoverAsset && (
				<FormRow
					label="Cover"
					control={
						<span className={styles.coverControl}>
							<MediaThumb
								className={styles.coverThumb}
								src={coverAsset?.src}
								alt={coverAsset?.alt ?? coverAsset?.title ?? ""}
							/>
							<Button size="compact" onClick={onPickCoverAsset}>
								Choose…
							</Button>
						</span>
					}
				/>
			)}
			{preview === "live" ? (
				<SplitPane
					className={styles.editorPane}
					divider
					left={editor}
					right={
						<ScrollRegion className={styles.previewScroll}>
							<MarkdownArticle className={styles.preview} source={draft.body} />
						</ScrollRegion>
					}
				/>
			) : (
				editor
			)}
		</FormPanel>
	);
}
