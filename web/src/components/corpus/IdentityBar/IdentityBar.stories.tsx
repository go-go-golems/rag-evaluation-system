import { useState } from "react";
import type { Meta, StoryObj } from "@storybook/react-vite";
import { IdentityBar } from "./IdentityBar";
import type { CorpusIdentityArgs } from "../../../services/api";

const strategies = [
	{ id: "fixed-1200-150" },
	{ id: "semantic-sections-v2" },
	{ id: "markdown-headings" },
];

const initialIdentity: CorpusIdentityArgs = {
	strategy_id: "fixed-1200-150",
	provider_type: "openai",
	model: "text-embedding-3-small",
	dimensions: 1536,
};

const meta = {
	title: "Component Library/Corpus/IdentityBar",
	component: IdentityBar,
	args: {
		identity: initialIdentity,
		strategies,
		totalDocs: 204,
		totalChunks: 1324,
		totalEmbedded: 1308,
		onChange: () => {},
	},
} satisfies Meta<typeof IdentityBar>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Interactive: Story = {
	render: (args) => {
		const [identity, setIdentity] = useState(args.identity);
		return <IdentityBar {...args} identity={identity} onChange={setIdentity} />;
	},
};

export const SparseCoverage: Story = {
	args: {
		totalDocs: 204,
		totalChunks: 1324,
		totalEmbedded: 92,
	},
	render: (args) => {
		const [identity, setIdentity] = useState(args.identity);
		return <IdentityBar {...args} identity={identity} onChange={setIdentity} />;
	},
};

export const NoStrategies: Story = {
	args: {
		strategies: [],
		identity: { ...initialIdentity, strategy_id: "" },
		totalDocs: 0,
		totalChunks: 0,
		totalEmbedded: 0,
	},
	render: (args) => {
		const [identity, setIdentity] = useState(args.identity);
		return <IdentityBar {...args} identity={identity} onChange={setIdentity} />;
	},
};
