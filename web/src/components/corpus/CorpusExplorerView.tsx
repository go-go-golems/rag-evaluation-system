import React, { useState, useMemo, useCallback, useEffect } from "react";
import { IconButton } from "@go-go-golems/rag-evaluation-site";
import { Caption } from "@go-go-golems/rag-evaluation-site";
import { Panel, ScrollRegion } from "@go-go-golems/rag-evaluation-site";
import {
	useListCorpusSourcesQuery,
	useListCorpusDocumentsQuery,
	useGetCorpusDocumentQuery,
	useListChunkingStrategiesQuery,
	useGetDocumentProcessingIdentitiesQuery,
	useGetDocumentProcessingCoverageQuery,
	CorpusIdentityArgs,
	DocumentProcessingCoverageArgs,
	DocumentProcessingCoverageItem,
} from "../../services/api";
import { IdentityBar } from "./IdentityBar/IdentityBar";
import { DocProcessingIdentityBar } from "./ArtifactIdentityBar/ArtifactIdentityBar";
import { SourcePanel } from "./SourcePanel/SourcePanel";
import { DocumentBrowser } from "./DocumentBrowser/DocumentBrowser";
import { DocumentInspector } from "./DocumentInspector/DocumentInspector";

export interface ChunkNavigationTarget {
	documentId: string;
	chunkId: string;
	sourceId: string;
}

interface CorpusExplorerViewProps {
	initialTarget?: ChunkNavigationTarget | null;
	onTargetConsumed?: () => void;
}

const DEFAULT_IDENTITY: CorpusIdentityArgs = {
	strategy_id: "fixed-1200-150",
	provider_type: "openai",
	model: "text-embedding-3-small",
	dimensions: 1536,
};

const PAGE_SIZE = 100;

const DEFAULT_PREPROCESSING_IDENTITY: DocumentProcessingCoverageArgs = {
	artifact_type: "clean_text",
	prompt_version: "v1",
	provider: "fake",
	model: "fake-document-processor",
};

export const CorpusExplorerView: React.FC<CorpusExplorerViewProps> = ({
	initialTarget,
	onTargetConsumed,
}) => {
	const [identity, setIdentity] = useState<CorpusIdentityArgs>(DEFAULT_IDENTITY);
	const [artifactIdentity, setArtifactIdentity] = useState<DocumentProcessingCoverageArgs>(
		DEFAULT_PREPROCESSING_IDENTITY,
	);
	const [sourceId, setSourceId] = useState("");
	const [documentId, setDocumentId] = useState("");
	const [highlightChunkId, setHighlightChunkId] = useState<string | null>(null);
	const [pageLimit, setPageLimit] = useState(PAGE_SIZE);

	// Handle navigation from Search workbench
	useEffect(() => {
		if (initialTarget) {
			setSourceId(initialTarget.sourceId);
			setDocumentId(initialTarget.documentId);
			setHighlightChunkId(initialTarget.chunkId);
			onTargetConsumed?.();
		}
	}, [initialTarget, onTargetConsumed]);

	const { data: strategies = [] } = useListChunkingStrategiesQuery();
	const { data: sources = [], isLoading: sourcesLoading } = useListCorpusSourcesQuery(identity);
	const { data: documents = [], isLoading: docsLoading } = useListCorpusDocumentsQuery(
		{ ...identity, source_id: sourceId, limit: pageLimit },
		{ skip: !sourceId },
	);
	const { data: detail, isLoading: detailLoading } = useGetCorpusDocumentQuery(
		{ ...identity, document_id: documentId, include_text: true },
		{ skip: !documentId },
	);

	// Preprocessing identity selector
	const { data: preprocessingIdentitiesData } = useGetDocumentProcessingIdentitiesQuery();
	const preprocessingIdentities = preprocessingIdentitiesData?.items ?? [];

	// When identities load, auto-select the first one if the default has no coverage
	useEffect(() => {
		if (preprocessingIdentities.length > 0 && artifactIdentity.provider === "fake") {
			const first = preprocessingIdentities[0]!;
			setArtifactIdentity({
				artifact_type: first.artifact_type,
				prompt_version: first.prompt_version,
				provider: first.provider,
				model: first.model,
			});
		}
	}, [preprocessingIdentities, artifactIdentity.provider]);

	// Preprocessing coverage — uses selected artifact identity
	const { data: preprocessingCoverageData } =
		useGetDocumentProcessingCoverageQuery(artifactIdentity);
	const preprocessingCoverage = useMemo(() => {
		const map: Record<string, DocumentProcessingCoverageItem> = {};
		for (const item of preprocessingCoverageData?.items ?? []) {
			map[item.source_id] = item;
		}
		return map;
	}, [preprocessingCoverageData]);

	const selectedSource = useMemo(
		() => sources.find((s) => s.source_id === sourceId),
		[sources, sourceId],
	);

	const selectedDoc = useMemo(
		() => documents.find((d) => d.id === documentId),
		[documents, documentId],
	);

	const handleSelectSource = useCallback((id: string) => {
		setSourceId(id);
		setDocumentId("");
		setHighlightChunkId(null);
		setPageLimit(PAGE_SIZE);
	}, []);

	const handleSelectDocument = useCallback((id: string) => {
		setDocumentId(id);
		setHighlightChunkId(null);
	}, []);

	const handleLoadMore = useCallback(() => {
		setPageLimit((prev) => prev + PAGE_SIZE);
	}, []);

	const totalDocs = useMemo(() => sources.reduce((s, x) => s + x.document_count, 0), [sources]);
	const totalChunks = useMemo(() => sources.reduce((s, x) => s + x.chunk_count, 0), [sources]);
	const totalEmbedded = useMemo(() => sources.reduce((s, x) => s + x.embedded_count, 0), [sources]);

	const hasMore = selectedSource ? documents.length < selectedSource.document_count : false;

	return (
		<div style={{ display: "flex", flexDirection: "column", gap: 8 }}>
			<IdentityBar
				identity={identity}
				onChange={setIdentity}
				strategies={strategies}
				totalDocs={totalDocs}
				totalChunks={totalChunks}
				totalEmbedded={totalEmbedded}
			/>

			{preprocessingIdentities.length > 0 && (
				<DocProcessingIdentityBar
					identities={preprocessingIdentities}
					selected={artifactIdentity}
					onChange={setArtifactIdentity}
				/>
			)}

			<div style={{ display: "flex", gap: 8, minHeight: 500 }}>
				<SourcePanel
					sources={sources}
					isLoading={sourcesLoading}
					selectedId={sourceId}
					onSelect={handleSelectSource}
					preprocessingCoverage={preprocessingCoverage}
				/>

				<DocumentBrowser
					documents={documents}
					isLoading={docsLoading}
					sourceName={selectedSource?.source_name ?? null}
					totalDocs={selectedSource?.document_count ?? 0}
					selectedDocId={documentId}
					onSelectDocument={handleSelectDocument}
					onLoadMore={handleLoadMore}
					hasMore={hasMore}
				/>

				<Panel
					title={selectedDoc ? selectedDoc.title : "Document Inspector"}
					actions={
						selectedDoc && (
							<IconButton
								label="Copy document ID"
								onClick={() => navigator.clipboard.writeText(selectedDoc.id)}
							>
								#{selectedDoc.id}
							</IconButton>
						)
					}
					density="condensed"
					fill
					style={{ flex: 1, minWidth: 0 }}
				>
					<ScrollRegion axis="y" style={{ maxHeight: 600 }}>
						{!documentId ? (
							<Caption>Select a document to inspect.</Caption>
						) : detailLoading ? (
							<Caption>Loading...</Caption>
						) : detail ? (
							<DocumentInspector
								detail={detail}
								chunks={detail.chunks}
								identity={identity}
								highlightChunkId={highlightChunkId}
							/>
						) : (
							<Caption>Document not found.</Caption>
						)}
					</ScrollRegion>
				</Panel>
			</div>
		</div>
	);
};
