import type { Meta, StoryObj } from "@storybook/react-vite";
import { useState } from "react";
import { cmsArticleFixtures, cmsAssetFixtures } from "../../../cms";
import { Caption } from "../../foundation";
import { ArticleListPanel } from "../ArticleListPanel";
import { MediaLibraryPanel } from "../MediaLibraryPanel";
import { CmsShell } from "./CmsShell";

const meta = {
	title: "Component Library/Organisms/CmsShell",
	component: CmsShell,
	args: {
		activeItemId: "articles",
		onNavigate: () => {},
	},
} satisfies Meta<typeof CmsShell>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
	args: {
		children: <ArticleListPanel articles={cmsArticleFixtures} onRowAction={() => {}} />,
	},
};

export const MediaActive: Story = {
	args: {
		activeItemId: "media",
		children: <MediaLibraryPanel assets={cmsAssetFixtures} />,
	},
};

export const WithFooter: Story = {
	args: {
		sidebarFooter: <Caption>content studio · v0.1</Caption>,
		children: <ArticleListPanel articles={cmsArticleFixtures.slice(0, 3)} />,
	},
};

export const Interactive: Story = {
	render: () => {
		const [active, setActive] = useState("articles");
		return (
			<CmsShell activeItemId={active} onNavigate={setActive}>
				{active === "media" ? (
					<MediaLibraryPanel assets={cmsAssetFixtures} />
				) : (
					<ArticleListPanel articles={cmsArticleFixtures} onRowAction={() => {}} />
				)}
			</CmsShell>
		);
	},
};
