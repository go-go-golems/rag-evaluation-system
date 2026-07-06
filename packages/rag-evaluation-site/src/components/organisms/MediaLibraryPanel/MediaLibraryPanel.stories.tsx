import type { Meta, StoryObj } from "@storybook/react-vite";
import { useMemo, useState } from "react";
import { cmsAssetFixtures, cmsUploadQueueFixtures } from "../../../cms";
import { Button } from "../../atoms";
import type { MediaLibraryKindFilter } from "./MediaLibraryPanel";
import { MediaLibraryPanel } from "./MediaLibraryPanel";

const meta = {
	title: "Component Library/Organisms/MediaLibraryPanel",
	component: MediaLibraryPanel,
	args: {
		assets: cmsAssetFixtures,
		style: { maxWidth: 880 },
	},
} satisfies Meta<typeof MediaLibraryPanel>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Populated: Story = {
	args: {
		selectedAssetIds: [cmsAssetFixtures[1]?.id ?? ""],
		query: "",
		onQueryChange: () => {},
		kindFilter: "all",
		onKindFilterChange: () => {},
		page: 1,
		pageCount: 3,
		onPageChange: () => {},
	},
};

export const Empty: Story = {
	args: {
		assets: [],
		emptyMessage: "No assets yet",
		emptyAction: <Button variant="primary">⤓ Upload files</Button>,
	},
};

export const Uploading: Story = {
	args: {
		assets: cmsAssetFixtures.slice(0, 4),
		onFilesSelected: () => {},
		uploads: cmsUploadQueueFixtures,
		onUploadCancel: () => {},
		onUploadRetry: () => {},
		onUploadDismiss: () => {},
	},
};

export const MultiSelect: Story = {
	args: {
		selectionMode: "multi",
		selectedAssetIds: [cmsAssetFixtures[0]?.id ?? "", cmsAssetFixtures[4]?.id ?? ""],
	},
};

export const PickerMode: Story = {
	args: {
		title: "Choose an asset",
		assets: cmsAssetFixtures.filter((asset) => asset.kind === "image"),
		showStatusBadges: false,
		selectedAssetIds: [cmsAssetFixtures[0]?.id ?? ""],
	},
};

export const DenseFiftyAssets: Story = {
	args: {
		assets: Array.from({ length: 50 }, (_, index) => ({
			...(cmsAssetFixtures[index % cmsAssetFixtures.length] as (typeof cmsAssetFixtures)[number]),
			id: `dense-${index}`,
			title: `Asset ${index + 1}`,
		})),
		minTileWidth: 120,
	},
};

export const Interactive: Story = {
	render: () => {
		const [selected, setSelected] = useState<string[]>([]);
		const [query, setQuery] = useState("");
		const [kind, setKind] = useState<MediaLibraryKindFilter>("all");
		const [page, setPage] = useState(1);
		const filtered = useMemo(
			() =>
				cmsAssetFixtures.filter(
					(asset) =>
						(kind === "all" || asset.kind === kind) &&
						(query === "" || asset.title.toLowerCase().includes(query.toLowerCase())),
				),
			[query, kind],
		);
		return (
			<MediaLibraryPanel
				style={{ maxWidth: 880 }}
				assets={filtered}
				selectionMode="multi"
				selectedAssetIds={selected}
				onAssetSelect={(id) =>
					setSelected((current) =>
						current.includes(id) ? current.filter((entry) => entry !== id) : [...current, id],
					)
				}
				query={query}
				onQueryChange={(value) => {
					setQuery(value);
					setPage(1);
				}}
				kindFilter={kind}
				onKindFilterChange={setKind}
				page={page}
				pageCount={2}
				onPageChange={setPage}
				emptyMessage="No assets match this filter"
			/>
		);
	},
};
