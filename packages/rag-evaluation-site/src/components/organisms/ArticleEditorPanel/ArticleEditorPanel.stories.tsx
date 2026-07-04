import type { Meta, StoryObj } from "@storybook/react-vite";
import { useState } from "react";
import type { CmsAsset } from "../../../cms";
import { cmsArticleDetailFixture, cmsAssetFixtures, cmsTagSuggestions } from "../../../cms";
import { AssetPickerDialog } from "../AssetPickerDialog";
import type { ArticleEditorDraft } from "./ArticleEditorPanel";
import { ArticleEditorPanel } from "./ArticleEditorPanel";

const sampleDraft: ArticleEditorDraft = {
	title: cmsArticleDetailFixture.title,
	slug: cmsArticleDetailFixture.slug,
	status: cmsArticleDetailFixture.status,
	tags: cmsArticleDetailFixture.tags,
	coverAssetId: cmsArticleDetailFixture.coverAssetId,
	body:
		cmsArticleDetailFixture.blocks[0]?.kind === "markdown"
			? cmsArticleDetailFixture.blocks[0].source
			: "",
};

const coverAsset = cmsAssetFixtures.find(
	(asset) => asset.id === cmsArticleDetailFixture.coverAssetId,
);

const meta = {
	title: "Component Library/Organisms/ArticleEditorPanel",
	component: ArticleEditorPanel,
	args: {
		draft: sampleDraft,
		onDraftChange: () => {},
		onSave: () => {},
		onPublish: () => {},
		coverAsset,
		onPickCoverAsset: () => {},
		onInsertAsset: () => {},
		tagSuggestions: cmsTagSuggestions,
	},
} satisfies Meta<typeof ArticleEditorPanel>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Draft: Story = {};

export const Saving: Story = {
	args: { status: "saving" },
};

export const Success: Story = {
	args: { status: "success", statusMessage: "Saved." },
};

export const ErrorState: Story = {
	args: { status: "error", statusMessage: "409 slug already exists" },
};

export const PreviewHidden: Story = {
	args: { preview: "hidden" },
};

export const Interactive: Story = {
	render: () => {
		const [draft, setDraft] = useState(sampleDraft);
		const [pickerOpen, setPickerOpen] = useState(false);
		const [cover, setCover] = useState<CmsAsset | undefined>(coverAsset);
		return (
			<>
				<ArticleEditorPanel
					draft={draft}
					onDraftChange={setDraft}
					onSave={() => {}}
					onPublish={() => setDraft((current) => ({ ...current, status: "published" }))}
					coverAsset={cover}
					onPickCoverAsset={() => setPickerOpen(true)}
					onInsertAsset={() => setPickerOpen(true)}
					tagSuggestions={cmsTagSuggestions}
				/>
				<AssetPickerDialog
					open={pickerOpen}
					assets={cmsAssetFixtures.filter((asset) => asset.kind === "image")}
					onConfirm={(asset) => {
						setCover(asset);
						setDraft((current) => ({
							...current,
							coverAssetId: asset.id,
							body: `${current.body}\n\n![${asset.alt ?? asset.title}](${asset.src})`,
						}));
						setPickerOpen(false);
					}}
					onCancel={() => setPickerOpen(false)}
				/>
			</>
		);
	},
};
