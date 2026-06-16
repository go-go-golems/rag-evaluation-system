import { useState, type ComponentProps } from "react";
import type { Meta, StoryObj } from "@storybook/react-vite";
import { Caption } from "@go-go-golems/rag-evaluation-site";
import { Stack } from "@go-go-golems/rag-evaluation-site";
import { ChunkTimelineBar } from "./ChunkTimelineBar";
import type { CorpusChunk } from "../../../services/api";

const baseChunks: CorpusChunk[] = [
	{
		id: "chunk-0",
		strategy_id: "fixed-1200-150",
		chunk_index: 0,
		start_offset: 0,
		end_offset: 900,
		token_count: 210,
		text: "Introductory guidance on privacy tree selection.",
		embedding: {
			present: true,
			provider_type: "openai",
			model: "text-embedding-3-small",
			dimensions: 1536,
		},
	},
	{
		id: "chunk-1",
		strategy_id: "fixed-1200-150",
		chunk_index: 1,
		start_offset: 760,
		end_offset: 1680,
		token_count: 225,
		text: "Evergreen tree options and spacing considerations.",
		embedding: {
			present: false,
			provider_type: "openai",
			model: "text-embedding-3-small",
			dimensions: 1536,
		},
	},
	{
		id: "chunk-2",
		strategy_id: "fixed-1200-150",
		chunk_index: 2,
		start_offset: 1510,
		end_offset: 2500,
		token_count: 242,
		text: "Maintenance notes and cultivar tradeoffs.",
		embedding: {
			present: true,
			provider_type: "openai",
			model: "text-embedding-3-small",
			dimensions: 1536,
		},
	},
	{
		id: "chunk-3",
		strategy_id: "fixed-1200-150",
		chunk_index: 3,
		start_offset: 2350,
		end_offset: 3100,
		token_count: 190,
		text: "Closing notes and references.",
		embedding: {
			present: false,
			provider_type: "openai",
			model: "text-embedding-3-small",
			dimensions: 1536,
		},
	},
];

const meta = {
	title: "Component Library/Corpus/ChunkTimelineBar",
	component: ChunkTimelineBar,
	args: {
		chunks: baseChunks,
		selectedIdx: null,
		onSelect: () => {},
	},
} satisfies Meta<typeof ChunkTimelineBar>;

export default meta;
type Story = StoryObj<typeof meta>;

function InteractiveTimeline(args: Omit<ComponentProps<typeof ChunkTimelineBar>, "onSelect">) {
	const [selectedIdx, setSelectedIdx] = useState<number | null>(args.selectedIdx ?? null);
	return (
		<Stack gap="xs">
			<ChunkTimelineBar {...args} selectedIdx={selectedIdx} onSelect={setSelectedIdx} />
			<Caption>Selected: {selectedIdx === null ? "none" : `chunk ${selectedIdx}`}</Caption>
		</Stack>
	);
}

export const MixedCoverage: Story = {
	render: (args) => <InteractiveTimeline {...args} />,
};

export const AllEmbedded: Story = {
	args: {
		chunks: baseChunks.map((chunk) => ({
			...chunk,
			embedding: {
				present: true,
				provider_type: "openai",
				model: "text-embedding-3-small",
				dimensions: 1536,
			},
		})),
	},
	render: (args) => <InteractiveTimeline {...args} />,
};

export const SelectedChunk: Story = {
	args: {
		selectedIdx: 1,
	},
	render: (args) => <InteractiveTimeline {...args} />,
};

export const Empty: Story = {
	args: {
		chunks: [],
	},
	render: (args) => (
		<Stack gap="xs">
			<ChunkTimelineBar {...args} onSelect={() => {}} />
			<Caption>Empty chunk lists intentionally render no timeline.</Caption>
		</Stack>
	),
};
