import React, { useState, useEffect, useCallback } from "react";
import { PipelineView } from "./components/pipeline/PipelineView";
import { EmbeddingsView } from "./components/embeddings/EmbeddingsView";
import { SearchView } from "./components/search/SearchView";
import { EvaluationView } from "./components/evaluation/EvaluationView";
import { CorpusExplorerView } from "./components/corpus/CorpusExplorerView";
import { WorkflowsView } from "./components/workflows/WorkflowsView";
import { DslPreviewPage } from "./components/pages";
import { AppNav, AppShell } from "@go-go-golems/rag-evaluation-site";

const views = [
	{ id: "search", label: "Search" },
	{ id: "corpus", label: "Corpus" },
	{ id: "workflows", label: "Workflows" },
	{ id: "pipeline", label: "Pipeline" },
	{ id: "embeddings", label: "Embeddings" },
	{ id: "evaluation", label: "Evaluation" },
	{ id: "dsl", label: "DSL" },
];

export interface ChunkNavigationTarget {
	documentId: string;
	chunkId: string;
	sourceId: string;
}

export const App: React.FC = () => {
	const [activeView, setActiveView] = useState("search");
	const [chunkTarget, setChunkTarget] = useState<ChunkNavigationTarget | null>(null);

	// Listen for navigation events from SearchView
	const handleNavigateToChunk = useCallback((e: Event) => {
		const detail = (e as CustomEvent<ChunkNavigationTarget>).detail;
		setChunkTarget(detail);
		setActiveView("corpus");
	}, []);

	// Listen for navigation from Corpus Explorer to Workflows
	const handleNavigateToWorkflows = useCallback((_e: Event) => {
		setActiveView("workflows");
	}, []);

	useEffect(() => {
		window.addEventListener("rag:navigate-to-chunk", handleNavigateToChunk);
		window.addEventListener("rag:navigate-to-workflows", handleNavigateToWorkflows);
		return () => {
			window.removeEventListener("rag:navigate-to-chunk", handleNavigateToChunk);
			window.removeEventListener("rag:navigate-to-workflows", handleNavigateToWorkflows);
		};
	}, [handleNavigateToChunk, handleNavigateToWorkflows]);

	// Clear target once consumed
	const handleChunkTargetConsumed = useCallback(() => {
		setChunkTarget(null);
	}, []);

	const renderView = () => {
		switch (activeView) {
			case "corpus":
				return (
					<CorpusExplorerView
						initialTarget={chunkTarget}
						onTargetConsumed={handleChunkTargetConsumed}
					/>
				);
			case "workflows":
				return <WorkflowsView />;
			case "pipeline":
				return <PipelineView />;
			case "embeddings":
				return <EmbeddingsView />;
			case "search":
				return <SearchView />;
			case "evaluation":
				return <EvaluationView />;
			case "dsl":
				return <DslPreviewPage pageId="demo" />;
			default:
				return <SearchView />;
		}
	};

	return (
		<AppShell
			header={
				<AppNav
					brand="◉ RAG Eval"
					items={views}
					activeItemId={activeView}
					onItemSelect={setActiveView}
				/>
			}
		>
			{renderView()}
		</AppShell>
	);
};
