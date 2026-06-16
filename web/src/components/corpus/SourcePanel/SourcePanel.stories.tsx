import { useState } from "react";
import type { Meta, StoryObj } from "@storybook/react-vite";
import { SourcePanel } from "./SourcePanel";
import type { CorpusSourceSummary, DocumentProcessingCoverageItem } from "../../../services/api";

const meta = {
	title: "Component Library/Corpus/SourcePanel",
	component: SourcePanel,
} satisfies Meta<typeof SourcePanel>;

export default meta;
type Story = StoryObj;

const sources: CorpusSourceSummary[] = [
	{
		source_id: "ttc-guides",
		source_name: "TTC Guides",
		source_type: "web",
		document_count: 128,
		word_count: 242_000,
		chunk_count: 920,
		embedded_count: 920,
		missing_embedding_count: 0,
	},
	{
		source_id: "ttc-articles",
		source_name: "TTC Articles",
		source_type: "web",
		document_count: 84,
		word_count: 131_400,
		chunk_count: 512,
		embedded_count: 240,
		missing_embedding_count: 272,
	},
	{
		source_id: "ttc-products",
		source_name: "TTC Products",
		source_type: "dump",
		document_count: 1200,
		word_count: 412_000,
		chunk_count: 1200,
		embedded_count: 0,
		missing_embedding_count: 1200,
	},
];

const preprocessing: Record<string, DocumentProcessingCoverageItem> = {
	"ttc-guides": {
		source_id: "ttc-guides",
		document_count: 128,
		artifact_count: 128,
		fresh_count: 126,
		failed_count: 0,
		missing_count: 0,
	},
	"ttc-articles": {
		source_id: "ttc-articles",
		document_count: 84,
		artifact_count: 42,
		fresh_count: 40,
		failed_count: 2,
		missing_count: 42,
	},
	"ttc-products": {
		source_id: "ttc-products",
		document_count: 1200,
		artifact_count: 0,
		fresh_count: 0,
		failed_count: 0,
		missing_count: 1200,
	},
};

export const Interactive: Story = {
	render: () => {
		const [selectedId, setSelectedId] = useState("ttc-guides");
		return (
			<SourcePanel
				sources={sources}
				selectedId={selectedId}
				onSelect={setSelectedId}
				isLoading={false}
				preprocessingCoverage={preprocessing}
			/>
		);
	},
};

export const Loading: Story = {
	render: () => (
		<SourcePanel
			sources={[]}
			selectedId=""
			onSelect={() => undefined}
			isLoading
			preprocessingCoverage={{}}
		/>
	),
};
