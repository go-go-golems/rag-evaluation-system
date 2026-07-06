import type { Meta, StoryObj } from "@storybook/react-vite";
import { useState } from "react";
import { cmsArticleFixtures, cmsAssetFixtures, cmsTagSuggestions } from "../../../cms";
import type { AssetDetailDraft } from "./AssetDetailPanel";
import { AssetDetailPanel } from "./AssetDetailPanel";

function assetById(id: string) {
	const found = cmsAssetFixtures.find((asset) => asset.id === id);
	if (!found) throw new Error(`missing asset fixture ${id}`);
	return found;
}

const imageAsset = assetById("asset-budget-sketch");
const fileAsset = assetById("asset-slides-pdf");
const brokenAsset = assetById("asset-broken");

function draftFor(asset: typeof imageAsset): AssetDetailDraft {
	return {
		title: asset.title,
		alt: asset.alt ?? "",
		tags: asset.tags,
		status: asset.status,
	};
}

const meta = {
	title: "Component Library/Organisms/AssetDetailPanel",
	component: AssetDetailPanel,
	args: {
		asset: imageAsset,
		draft: draftFor(imageAsset),
		onDraftChange: () => {},
		onSave: () => {},
		onDownload: () => {},
		onDelete: () => {},
		tagSuggestions: cmsTagSuggestions,
		style: { maxWidth: 960 },
	},
} satisfies Meta<typeof AssetDetailPanel>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Image: Story = {};

export const File: Story = {
	args: { asset: fileAsset, draft: draftFor(fileAsset) },
};

export const BrokenImage: Story = {
	args: { asset: brokenAsset, draft: draftFor(brokenAsset) },
};

export const WithUsage: Story = {
	args: {
		usedBy: cmsArticleFixtures.slice(0, 2),
		onUsageSelect: () => {},
	},
};

export const Interactive: Story = {
	render: () => {
		const [draft, setDraft] = useState(draftFor(imageAsset));
		return (
			<AssetDetailPanel
				style={{ maxWidth: 960 }}
				asset={imageAsset}
				draft={draft}
				onDraftChange={setDraft}
				onSave={() => {}}
				usedBy={cmsArticleFixtures.slice(0, 2)}
				onUsageSelect={() => {}}
				onDownload={() => {}}
				onDelete={() => {}}
				tagSuggestions={cmsTagSuggestions}
			/>
		);
	},
};
