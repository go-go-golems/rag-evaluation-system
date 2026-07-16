import { useEffect, useMemo, useState } from "react";
import { Caption, Text } from "@go-go-golems/rag-evaluation-site";
import { Panel, Stack } from "@go-go-golems/rag-evaluation-site";
import {
	type ExperimentSpecificationInput,
	useCreateExperimentRunMutation,
	useCreateExperimentSpecificationMutation,
	useGetExperimentComparisonQuery,
	useGetExperimentSpecificationQuery,
	useGetLabCatalogQuery,
	useListExperimentRunTracesQuery,
	useListExperimentRunsQuery,
	useListExperimentSpecificationsQuery,
} from "../../../services/api";
import styles from "./EvaluationPage.module.css";

const candidateDatasetID = "candidate:ttc-baseline-v1";

function shortID(id: string) {
	return id.length > 22 ? `${id.slice(0, 18)}…` : id;
}

type RerankingTrace = {
	identity?: { kind?: string; model?: string };
	candidates?: Array<{ candidateId?: string; preRerankRank?: number; retrievalScore?: number }>;
	results?: Array<{ candidateId?: string; rank?: number; score?: number }>;
};

function rerankingTrace(trace: Record<string, unknown>): RerankingTrace | undefined {
	const value = trace.reranking;
	return value && typeof value === "object" ? (value as RerankingTrace) : undefined;
}

function RerankingTraceView({ trace }: { trace: Record<string, unknown> }) {
	const reranking = rerankingTrace(trace);
	if (!reranking) return null;
	const resultByID = new Map((reranking.results ?? []).map((result) => [result.candidateId, result]));
	return (
		<section className={styles.reranking} aria-label="Reranking trace">
			<strong>Cross-encoder reranking</strong>
			<Caption>{reranking.identity?.kind ?? "unknown"} · {reranking.identity?.model ?? "unknown model"}</Caption>
			<div className={styles.rerankTable} role="table" aria-label="Candidate rank changes">
				<div role="row" className={styles.rerankHeader}><span>Candidate</span><span>Before</span><span>Retrieval</span><span>After</span><span>Rerank score</span></div>
				{(reranking.candidates ?? []).map((candidate) => {
					const result = resultByID.get(candidate.candidateId);
					return <div role="row" className={styles.rerankRow} key={candidate.candidateId}>
						<span title={candidate.candidateId}>{shortID(candidate.candidateId ?? "unknown")}</span>
						<span>{candidate.preRerankRank ?? "—"}</span><span>{candidate.retrievalScore?.toFixed(4) ?? "—"}</span>
						<span>{result?.rank ?? "truncated"}</span><span>{result?.score?.toFixed(4) ?? "—"}</span>
					</div>;
				})}
			</div>
		</section>
	);
}

export function EvaluationPage() {
	const catalog = useGetLabCatalogQuery();
	const specifications = useListExperimentSpecificationsQuery();
	const runs = useListExperimentRunsQuery();
	const [createSpecification, createSpecificationState] =
		useCreateExperimentSpecificationMutation();
	const [createRun, createRunState] = useCreateExperimentRunMutation();
	const [snapshotID, setSnapshotID] = useState("");
	const [chunkSetID, setChunkSetID] = useState("");
	const [bm25ArtifactID, setBM25ArtifactID] = useState("");
	const [embeddingSetID, setEmbeddingSetID] = useState("");
	const [configText, setConfigText] = useState(
		'{"limit":10,"rrf_k":60,"channels":["bm25","vector"]}',
	);
	const [selectedRunID, setSelectedRunID] = useState("");
	const [compareRunID, setCompareRunID] = useState("");
	const [formError, setFormError] = useState("");

	const chunkSets = useMemo(
		() => (catalog.data?.chunk_sets ?? []).filter((item) => item.corpus_snapshot_id === snapshotID),
		[catalog.data?.chunk_sets, snapshotID],
	);
	const bm25Artifacts = useMemo(
		() => (catalog.data?.bm25_artifacts ?? []).filter((item) => item.chunk_set_id === chunkSetID),
		[catalog.data?.bm25_artifacts, chunkSetID],
	);
	const embeddingSets = useMemo(
		() => (catalog.data?.embedding_sets ?? []).filter((item) => item.chunk_set_id === chunkSetID),
		[catalog.data?.embedding_sets, chunkSetID],
	);

	useEffect(() => {
		if (!snapshotID && catalog.data?.snapshots[0]) setSnapshotID(catalog.data.snapshots[0].id);
	}, [catalog.data?.snapshots, snapshotID]);
	useEffect(() => {
		if (!chunkSets.some((item) => item.id === chunkSetID)) setChunkSetID(chunkSets[0]?.id ?? "");
	}, [chunkSets, chunkSetID]);
	useEffect(() => {
		if (!bm25Artifacts.some((item) => item.id === bm25ArtifactID))
			setBM25ArtifactID(bm25Artifacts[0]?.id ?? "");
	}, [bm25Artifacts, bm25ArtifactID]);
	useEffect(() => {
		if (!embeddingSets.some((item) => item.id === embeddingSetID))
			setEmbeddingSetID(embeddingSets[0]?.id ?? "");
	}, [embeddingSets, embeddingSetID]);

	const traces = useListExperimentRunTracesQuery(selectedRunID, { skip: !selectedRunID });
	const selectedRun = runs.data?.find((run) => run.id === selectedRunID);
	const selectedSpecification = useGetExperimentSpecificationQuery(
		selectedRun?.experiment_spec_id ?? "",
		{
			skip: !selectedRun,
		},
	);
	const comparison = useGetExperimentComparisonQuery(
		{ left: selectedRunID, right: compareRunID },
		{ skip: !selectedRunID || !compareRunID },
	);

	async function submitSpecification(event: React.FormEvent<HTMLFormElement>) {
		event.preventDefault();
		setFormError("");
		let config: Record<string, unknown>;
		try {
			config = JSON.parse(configText) as Record<string, unknown>;
		} catch {
			setFormError("Configuration must be a JSON object.");
			return;
		}
		const input: ExperimentSpecificationInput = {
			corpus_snapshot_id: snapshotID,
			chunk_set_id: chunkSetID,
			bm25_artifact_id: bm25ArtifactID || undefined,
			embedding_set_id: embeddingSetID || undefined,
			evaluation_dataset_id: candidateDatasetID,
			config,
		};
		try {
			await createSpecification(input).unwrap();
		} catch (error) {
			setFormError(
				error instanceof Error ? error.message : "Could not create experiment specification.",
			);
		}
	}

	async function startRun(specificationID: string) {
		try {
			const run = await createRun(specificationID).unwrap();
			setSelectedRunID(run.id);
		} catch (error) {
			setFormError(error instanceof Error ? error.message : "Could not create run.");
		}
	}

	return (
		<Stack gap="md" data-rag-page="EvaluationPage">
			<Panel title="RAG Laboratory">
				<Stack gap="sm">
					<Text>
						Create immutable retrieval specifications, inspect append-only runs, and compare
						query-level evidence. The TTC baseline cards are source-validated candidates; they are
						not labeled as human-frozen truth.
					</Text>
					<Caption>
						A specification identifies immutable artifacts. A run records events, traces, and a
						terminal summary without overwriting earlier evidence.
					</Caption>
				</Stack>
			</Panel>

			<Panel title="Create specification">
				<form className={styles.form} onSubmit={submitSpecification}>
					<label>
						Corpus snapshot
						<select value={snapshotID} onChange={(event) => setSnapshotID(event.target.value)}>
							<option value="">Choose snapshot</option>
							{catalog.data?.snapshots.map((item) => (
								<option value={item.id} key={item.id}>
									{shortID(item.id)} · {item.document_count} documents
								</option>
							))}
						</select>
					</label>
					<label>
						Chunk set
						<select value={chunkSetID} onChange={(event) => setChunkSetID(event.target.value)}>
							<option value="">Choose chunk set</option>
							{chunkSets.map((item) => (
								<option value={item.id} key={item.id}>
									{shortID(item.id)} · {item.chunk_count} chunks
								</option>
							))}
						</select>
					</label>
					<label>
						BM25 artifact (optional)
						<select
							value={bm25ArtifactID}
							onChange={(event) => setBM25ArtifactID(event.target.value)}
						>
							<option value="">No lexical channel</option>
							{bm25Artifacts.map((item) => (
								<option value={item.id} key={item.id}>
									{shortID(item.id)} · {item.chunk_count} chunks
								</option>
							))}
						</select>
					</label>
					<label>
						Embedding set (optional)
						<select
							value={embeddingSetID}
							onChange={(event) => setEmbeddingSetID(event.target.value)}
						>
							<option value="">No vector channel</option>
							{embeddingSets.map((item) => (
								<option value={item.id} key={item.id}>
									{shortID(item.id)} · {item.embedding_count} vectors
								</option>
							))}
						</select>
					</label>
					<label className={styles.wide}>
						Retrieval configuration JSON
						<textarea
							rows={3}
							value={configText}
							onChange={(event) => setConfigText(event.target.value)}
						/>
					</label>
					<div className={styles.actions}>
						<button
							type="submit"
							disabled={
								!snapshotID ||
								!chunkSetID ||
								(!bm25ArtifactID && !embeddingSetID) ||
								createSpecificationState.isLoading
							}
						>
							{createSpecificationState.isLoading ? "Creating…" : "Create immutable specification"}
						</button>
						<Caption>Dataset: {candidateDatasetID}</Caption>
					</div>
					{formError && <Text tone="danger">{formError}</Text>}
				</form>
			</Panel>

			<Panel title="Specifications and runs">
				<div className={styles.columns}>
					<section>
						<h3>Specifications</h3>
						{specifications.data?.map((specification) => (
							<div className={styles.card} key={specification.id}>
								<strong title={specification.id}>{shortID(specification.id)}</strong>
								<Caption>
									{shortID(specification.chunk_set_id)} ·{" "}
									{specification.config.channels
										? String(specification.config.channels)
										: "configured retrieval"}
								</Caption>
								<button
									type="button"
									onClick={() => startRun(specification.id)}
									disabled={createRunState.isLoading}
								>
									Start append-only run
								</button>
							</div>
						)) ?? <Caption>No immutable specifications yet.</Caption>}
					</section>
					<section>
						<h3>Runs</h3>
						{runs.data?.map((run) => (
							<div
								className={`${styles.card} ${selectedRunID === run.id ? styles.selected : ""}`}
								key={run.id}
							>
								<button
									type="button"
									className={styles.runButton}
									onClick={() => setSelectedRunID(run.id)}
								>
									<strong title={run.id}>{shortID(run.id)}</strong>
									<br />
									<Caption>
										{run.status} · {run.events.length} event(s)
									</Caption>
								</button>
							</div>
						)) ?? <Caption>No runs yet.</Caption>}
					</section>
				</div>
			</Panel>

			{selectedRunID && (
				<Panel title={`Trace inspector — ${shortID(selectedRunID)}`}>
					<Stack gap="sm">
						<Caption>
							One immutable trace per query card. A newly created run will be empty until an
							executor/importer records its traces.
						</Caption>
						{selectedRun && (
							<div className={styles.specificationLink}>
								<Caption>
									Specification {shortID(selectedRun.experiment_spec_id)} ·{" "}
									{selectedSpecification.data?.schema_version ?? "loading exported manifest…"}
								</Caption>
								<a
									href={`/api/v1/lab/specifications/${encodeURIComponent(selectedRun.experiment_spec_id)}`}
									target="_blank"
									rel="noreferrer"
								>
									Open exported immutable specification JSON
								</a>
							</div>
						)}
						<div className={styles.traceGrid}>
							{traces.data?.map((trace) => (
								<article className={styles.card} key={trace.query_card_id}>
									<strong>{trace.query_card_id}</strong>
									<Caption>{JSON.stringify(trace.timing)}</Caption>
									<RerankingTraceView trace={trace.trace} />
									<pre>{JSON.stringify(trace.trace, null, 2).slice(0, 1200)}</pre>
								</article>
							)) ?? <Caption>Loading traces…</Caption>}
						</div>
						<label className={styles.compare}>
							Compare selected run with
							<select
								value={compareRunID}
								onChange={(event) => setCompareRunID(event.target.value)}
							>
								<option value="">Choose another run</option>
								{runs.data
									?.filter((run) => run.id !== selectedRunID)
									.map((run) => (
										<option value={run.id} key={run.id}>
											{shortID(run.id)} · {run.status}
										</option>
									))}
							</select>
						</label>
						{comparison.data && (
							<Text>
								Comparison loaded: {comparison.data.left_traces.length} left traces and{" "}
								{comparison.data.right_traces.length} right traces.
							</Text>
						)}
					</Stack>
				</Panel>
			)}
		</Stack>
	);
}
