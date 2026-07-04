import type { Meta, StoryObj } from "@storybook/react-vite";
import { useState } from "react";
import type { CmsAsset } from "../../../cms";
import { cmsAssetFixtures } from "../../../cms";
import { ContentStatusBadge } from "../../atoms";
import { TileGrid } from "../../layout";
import { AssetTile } from "./AssetTile";

function assetById(id: string): CmsAsset {
	const found = cmsAssetFixtures.find((asset) => asset.id === id);
	if (!found) throw new Error(`missing asset fixture ${id}`);
	return found;
}

const budgetSketch = assetById("asset-budget-sketch");
const hero = assetById("asset-hero");

const meta = {
	title: "Component Library/Molecules/AssetTile",
	component: AssetTile,
	args: { asset: budgetSketch, style: { width: 180 } },
} satisfies Meta<typeof AssetTile>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const Selected: Story = {
	args: { selected: true },
};

export const BrokenImage: Story = {
	args: { asset: assetById("asset-broken") },
};

export const FileKind: Story = {
	args: { asset: assetById("asset-slides-pdf") },
};

export const LongTitle: Story = {
	args: { asset: assetById("asset-treemap") },
};

export const WithStatusFooter: Story = {
	args: {
		asset: hero,
		footerSlot: <ContentStatusBadge status={hero.status} icon={false} />,
	},
};

export const Interactive: Story = {
	render: () => {
		const [selectedId, setSelectedId] = useState<string>();
		return (
			<TileGrid minTileWidth={160} style={{ maxWidth: 560 }}>
				{cmsAssetFixtures.slice(0, 6).map((asset) => (
					<AssetTile
						key={asset.id}
						asset={asset}
						selected={asset.id === selectedId}
						onSelect={setSelectedId}
					/>
				))}
			</TileGrid>
		);
	},
};
