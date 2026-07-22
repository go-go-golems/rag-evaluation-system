import { Caption, Panel, Stack, Text } from "@go-go-golems/rag-evaluation-site";
import { useGetRAGArtifactCatalogQuery } from "../../../services/api";

export function EvaluationPage() {
	const catalog = useGetRAGArtifactCatalogQuery();
	const snapshots = catalog.data?.snapshots.length ?? 0;
	const chunkSets = catalog.data?.chunk_sets.length ?? 0;
	const embeddingSets = catalog.data?.embedding_sets.length ?? 0;
	const bm25Artifacts = catalog.data?.bm25_artifacts.length ?? 0;

	return (
		<Stack gap="lg">
			<Panel title="Canonical RAG studies">
				<Stack gap="md">
					<Text>
						RAG authoring, compilation, artifact resolution, and worker semantics belong to
						rag-eval. Researchctl provides only the domain-neutral run lifecycle.
					</Text>
					<Caption>
						Validate and explain a study before submitting its canonical cells through the RAG-owned
						adapter.
					</Caption>
					<pre>{`rag-eval study validate experiments/rag-sol2/study.js \\
  --inputs experiments/rag-sol2/inputs.json \\
  --ttc-database /path/to/rag-eval.db

rag-eval study run experiments/rag-sol2/study.js \\
  --project project.yaml \\
  --experiment-id EXP-RAG \\
  --inputs experiments/rag-sol2/inputs.json \\
  --ttc-database /path/to/rag-eval.db \\
  --researchctl-command researchctl \\
  --worker-command rag-worker`}</pre>
				</Stack>
			</Panel>

			<Panel title="Read-only domain artifact catalog">
				{catalog.isLoading ? (
					<Text>Loading immutable catalog identities…</Text>
				) : catalog.error ? (
					<Text>Catalog unavailable. Verify the rag-eval database and server logs.</Text>
				) : (
					<Stack gap="sm">
						<Text>{snapshots} corpus snapshots</Text>
						<Text>{chunkSets} chunk sets</Text>
						<Text>{embeddingSets} embedding sets</Text>
						<Text>{bm25Artifacts} BM25 artifacts</Text>
					</Stack>
				)}
				<Caption>
					This page can inspect domain artifact availability, but it cannot allocate or mutate
					laboratory runs.
				</Caption>
			</Panel>
		</Stack>
	);
}
