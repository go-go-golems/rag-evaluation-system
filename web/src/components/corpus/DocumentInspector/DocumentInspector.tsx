import React, { useState } from "react";
import { Button, ErrorCallout, IconButton } from "@go-go-golems/rag-evaluation-site";
import { Caption, StatusText } from "@go-go-golems/rag-evaluation-site";
import { Panel, ScrollRegion, TabList } from "@go-go-golems/rag-evaluation-site";
import {
	DataTable,
	MetadataGrid,
	type DataTableColumn,
	type MetadataGridItem,
} from "@go-go-golems/rag-evaluation-site";
import {
	CorpusChunk,
	CorpusDocumentDetail,
	CorpusIdentityArgs,
	useGetDocumentProcessingArtifactsQuery,
	DocumentProcessingArtifact,
} from "../../../services/api";
import { ChunkTimelineBar } from "../ChunkTimelineBar/ChunkTimelineBar";
import styles from "./DocumentInspector.module.css";

function timeAgoShort(iso: string): string {
	const ms = Date.now() - new Date(iso).getTime();
	if (ms < 60000) return `${Math.floor(ms / 1000)}s`;
	if (ms < 3600000) return `${Math.floor(ms / 60000)}m`;
	if (ms < 86400000) return `${Math.floor(ms / 3600000)}h`;
	return `${Math.floor(ms / 86400000)}d`;
}

const ArtifactDetail: React.FC<{ artifacts: DocumentProcessingArtifact[] }> = ({ artifacts }) => {
	const [selectedIdx, setSelectedIdx] = useState<number | null>(null);
	const selected = selectedIdx !== null ? artifacts[selectedIdx] : null;

	return selected ? (
		<Panel
			className={styles.artifactPanel}
			title={`Artifact: ${selected.artifact_type} (${selected.prompt_version})`}
			actions={
				<IconButton label="Close artifact detail" onClick={() => setSelectedIdx(null)}>
					✕
				</IconButton>
			}
			density="condensed"
		>
			<MetadataGrid
				items={[
					{ key: "Provider", value: `${selected.provider}/${selected.model}` },
					{
						key: "Status",
						value: (
							<StatusText
								status={
									selected.status === "fresh"
										? "done"
										: selected.status === "failed"
											? "error"
											: "pending"
								}
							>
								{selected.status}
							</StatusText>
						),
					},
					{ key: "Hash", value: selected.input_hash },
				]}
			/>
			{selected.output_text && (
				<section className={styles.section}>
					<div className={styles.sectionTitle}>Output Text</div>
					<div className={styles.textContent}>{selected.output_text}</div>
				</section>
			)}
			{selected.error_message && (
				<ErrorCallout className={styles.error}>{selected.error_message}</ErrorCallout>
			)}
		</Panel>
	) : null;
};

interface DocumentInspectorProps {
	detail: CorpusDocumentDetail;
	chunks: CorpusChunk[];
	identity: CorpusIdentityArgs;
	highlightChunkId?: string | null;
}

type InspectorTab = "overview" | "text" | "chunks" | "coverage" | "artifacts";

const tabs = [
	{ id: "overview", label: "overview" },
	{ id: "text", label: "text" },
	{ id: "chunks", label: "chunks" },
	{ id: "coverage", label: "coverage" },
	{ id: "artifacts", label: "artifacts" },
] satisfies { id: InspectorTab; label: string }[];

export const DocumentInspector: React.FC<DocumentInspectorProps> = ({
	detail,
	chunks,
	identity,
	highlightChunkId,
}) => {
	const [activeTab, setActiveTab] = useState<InspectorTab>(
		highlightChunkId ? "chunks" : "overview",
	);
	const [selectedChunkIdx, setSelectedChunkIdx] = useState<number | null>(() => {
		if (!highlightChunkId) return null;
		const idx = (chunks ?? []).findIndex((c) => c.id === highlightChunkId);
		return idx >= 0 ? idx : null;
	});

	const safeChunks = chunks ?? [];
	const doc = detail.document;
	const metaKeys = Object.keys(doc.metadata || {});
	const embeddedCount = safeChunks.filter((c) => c.embedding?.present).length;
	const missingCount = safeChunks.length - embeddedCount;
	const identityLabel = `${identity.provider_type || "?"}/${identity.model || "?"} @ ${identity.dimensions || "?"}`;

	const { data: artifactsData, isLoading: artifactsLoading } =
		useGetDocumentProcessingArtifactsQuery(doc.id, { skip: activeTab !== "artifacts" });
	const artifacts = artifactsData?.items ?? [];

	const overviewItems: MetadataGridItem[] = [
		{ key: "ID", value: doc.id, copyValue: doc.id },
		{ key: "Source", value: doc.source_id },
		{
			key: "URL",
			value: doc.url ? (
				<a href={doc.url} target="_blank" rel="noreferrer" className={styles.link}>
					{doc.url}
				</a>
			) : (
				"—"
			),
		},
		{ key: "Words", value: doc.word_count.toLocaleString() },
		{ key: "Chunks", value: safeChunks.length },
		{
			key: "Embedded",
			value: (
				<Caption
					tone={
						embeddedCount === safeChunks.length && safeChunks.length > 0
							? "success"
							: embeddedCount > 0
								? "warning"
								: "muted"
					}
				>
					{embeddedCount}/{safeChunks.length}
				</Caption>
			),
		},
		{
			key: "Status",
			value: (
				<StatusText status={doc.status === "chunked" ? "done" : doc.status}>
					{doc.status}
				</StatusText>
			),
		},
	];

	const chunkColumns: DataTableColumn<CorpusChunk>[] = [
		{ id: "index", header: "#", cell: (chunk) => chunk.chunk_index },
		{ id: "range", header: "Range", cell: (chunk) => `${chunk.start_offset}–${chunk.end_offset}` },
		{ id: "tokens", header: "Tokens", align: "end", cell: (chunk) => chunk.token_count },
		{
			id: "embed",
			header: "Embed",
			align: "center",
			cell: (chunk) => (
				<Caption tone={chunk.embedding?.present ? "success" : "muted"}>
					{chunk.embedding?.present ? "●" : "○"}
				</Caption>
			),
		},
		{
			id: "enrich",
			header: "Enrich",
			align: "center",
			cell: (chunk) => (
				<Caption
					tone={chunk.enrichment?.present ? "success" : "muted"}
					title={chunk.enrichment?.short_summary ?? ""}
				>
					{chunk.enrichment?.present ? "●" : "○"}
				</Caption>
			),
		},
		{
			id: "id",
			header: "ID",
			cell: (chunk) => (
				<IconButton
					label="Copy chunk ID"
					onClick={(e) => {
						e.stopPropagation();
						navigator.clipboard.writeText(chunk.id);
					}}
				>
					{chunk.id.slice(0, 12)}…
				</IconButton>
			),
		},
	];

	const missingColumns: DataTableColumn<CorpusChunk>[] = [
		{ id: "index", header: "#", cell: (chunk) => chunk.chunk_index },
		{ id: "range", header: "Range", cell: (chunk) => `${chunk.start_offset}–${chunk.end_offset}` },
		{ id: "tokens", header: "Tokens", align: "end", cell: (chunk) => chunk.token_count },
		{ id: "id", header: "ID", cell: (chunk) => `${chunk.id.slice(0, 16)}…` },
	];

	const artifactColumns: DataTableColumn<DocumentProcessingArtifact>[] = [
		{ id: "type", header: "Type", cell: (a) => a.artifact_type },
		{ id: "version", header: "Version", cell: (a) => a.prompt_version },
		{
			id: "provider",
			header: "Provider",
			cell: (a) => (
				<Caption>
					{a.provider}/{a.model}
				</Caption>
			),
		},
		{
			id: "status",
			header: "Status",
			cell: (a) => (
				<StatusText
					status={a.status === "fresh" ? "done" : a.status === "failed" ? "error" : "pending"}
				>
					{a.status}
				</StatusText>
			),
		},
		{ id: "age", header: "Age", cell: (a) => <Caption>{timeAgoShort(a.updated_at)}</Caption> },
	];

	return (
		<div className={styles.root}>
			<TabList
				items={tabs}
				activeId={activeTab}
				onChange={setActiveTab}
				ariaLabel="Document inspector tabs"
			/>
			<ScrollRegion axis="y" className={styles.body}>
				{activeTab === "overview" && (
					<>
						<section className={styles.section}>
							<MetadataGrid items={overviewItems} />
						</section>
						{metaKeys.length > 0 && (
							<section className={styles.section}>
								<div className={styles.sectionTitle}>Metadata</div>
								<MetadataGrid
									items={metaKeys.map((key) => ({ key, value: String(doc.metadata[key] ?? "") }))}
								/>
							</section>
						)}
					</>
				)}

				{activeTab === "text" && (
					<section className={styles.section}>
						<div className={styles.sectionTitle}>
							Extracted Text ({doc.word_count.toLocaleString()} words)
						</div>
						<div className={styles.textContent}>
							{doc.content_text || <span className={styles.empty}>No content text available.</span>}
						</div>
					</section>
				)}

				{activeTab === "chunks" && (
					<>
						<section className={styles.section}>
							<div className={styles.sectionTitle}>
								Chunks ({safeChunks.length}) — {embeddedCount} embedded, {missingCount} missing ·{" "}
								{safeChunks.filter((c) => c.enrichment?.present).length} enriched
							</div>
							<ChunkTimelineBar
								chunks={safeChunks}
								selectedIdx={selectedChunkIdx}
								onSelect={setSelectedChunkIdx}
							/>
							<DataTable
								rows={safeChunks}
								columns={chunkColumns}
								getRowKey={(chunk) => chunk.id}
								selectedKey={selectedChunkIdx !== null ? safeChunks[selectedChunkIdx]?.id : null}
								onRowSelect={(chunk) => {
									const idx = safeChunks.findIndex((c) => c.id === chunk.id);
									setSelectedChunkIdx(idx === selectedChunkIdx ? null : idx);
								}}
							/>
						</section>
						{selectedChunkIdx !== null && safeChunks[selectedChunkIdx] && (
							<section className={`${styles.section} ${styles.chunkDetail}`}>
								<div className={styles.sectionTitle}>
									Chunk #{safeChunks[selectedChunkIdx].chunk_index} —{" "}
									{safeChunks[selectedChunkIdx].token_count} tokens
								</div>
								<div className={styles.textContent}>{safeChunks[selectedChunkIdx].text}</div>
							</section>
						)}
					</>
				)}

				{activeTab === "coverage" && (
					<>
						<section className={styles.section}>
							<div className={styles.sectionTitle}>Embedding Coverage — {identityLabel}</div>
							<MetadataGrid
								items={[
									{ key: "Total", value: safeChunks.length },
									{ key: "Embedded", value: <Caption tone="success">{embeddedCount}</Caption> },
									{ key: "Missing", value: <Caption tone="warning">{missingCount}</Caption> },
									{
										key: "Coverage",
										value: `${safeChunks.length > 0 ? Math.round((embeddedCount / safeChunks.length) * 100) : 0}%`,
									},
								]}
							/>
						</section>
						<section className={styles.section}>
							<div className={styles.sectionTitle}>Per-chunk status</div>
							<div className={styles.coverageStrip}>
								{safeChunks.map((chunk) => (
									<span
										key={chunk.id}
										className={`${styles.coverageDot} ${chunk.embedding?.present ? styles.present : styles.missing}`}
										title={`#${chunk.chunk_index} ${chunk.embedding?.present ? "embedded" : "missing"}`}
									/>
								))}
							</div>
							<Caption>● embedded · ○ missing</Caption>
						</section>
						{missingCount > 0 && (
							<section className={styles.section}>
								<div className={styles.sectionTitle}>Missing chunks</div>
								<DataTable
									rows={safeChunks.filter((c) => !c.embedding?.present)}
									columns={missingColumns}
									getRowKey={(chunk) => chunk.id}
								/>
							</section>
						)}
					</>
				)}

				{activeTab === "artifacts" && (
					<section className={styles.section}>
						<div className={styles.sectionTitle}>
							Document Processing Artifacts ({artifacts.length})
						</div>
						{artifactsLoading ? (
							<Caption>Loading...</Caption>
						) : artifacts.length === 0 ? (
							<div className={styles.empty}>
								No preprocessing artifacts for this document.
								<Button
									className={styles.inlineAction}
									onClick={() => window.dispatchEvent(new CustomEvent("rag:navigate-to-workflows"))}
								>
									Submit Workflow →
								</Button>
							</div>
						) : (
							<DataTable
								rows={artifacts}
								columns={artifactColumns}
								getRowKey={(a) => `${a.artifact_type}-${a.prompt_version}-${a.provider}-${a.model}`}
							/>
						)}
						{artifacts.length > 0 && <ArtifactDetail artifacts={artifacts} />}
					</section>
				)}
			</ScrollRegion>
		</div>
	);
};
