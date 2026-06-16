import type { Meta, StoryObj } from "@storybook/react-vite";
import { MetadataGrid } from "./MetadataGrid";

const meta = {
	title: "Component Library/Molecules/MetadataGrid",
	component: MetadataGrid,
} satisfies Meta<typeof MetadataGrid>;
export default meta;
type Story = StoryObj<typeof meta>;

export const ResultIdentity: Story = {
	args: {
		items: [
			{ key: "Title", value: "Fast Growing Trees for Privacy Screens" },
			{ key: "Chunk ID", value: "chunk_privacy_004", copyValue: "chunk_privacy_004" },
			{ key: "Score", value: "0.038700" },
		],
	},
};
