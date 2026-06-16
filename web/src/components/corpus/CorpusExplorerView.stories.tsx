import type { Meta, StoryObj } from "@storybook/react-vite";
import { MockApiProvider } from "../../storybook/MockApiProvider";
import { CorpusExplorerView } from "./CorpusExplorerView";

const meta = {
	title: "Views/Corpus/CorpusExplorerView",
	component: CorpusExplorerView,
	decorators: [
		(Story) => (
			<MockApiProvider>
				<Story />
			</MockApiProvider>
		),
	],
	parameters: { layout: "fullscreen" },
} satisfies Meta<typeof CorpusExplorerView>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const DeepLinkedChunk: Story = {
	args: {
		initialTarget: {
			sourceId: "src-tree-guides",
			documentId: "doc-privacy-trees",
			chunkId: "chunk_privacy_002",
		},
	},
};
