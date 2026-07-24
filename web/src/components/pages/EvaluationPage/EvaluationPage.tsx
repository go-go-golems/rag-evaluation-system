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
						RAG authoring, compilation, artifact resolution, and task semantics belong to rag-eval.
						Researchctl owns experiment lifecycle; Scraper Workflow V3 owns durable execution.
					</Text>
					<Caption>
						Compile canonical cells into a Researchctl plan, then execute them through the RAG
						Workflow runner.
					</Caption>
					<pre>{`artifact_root="$PWD/laboratory/artifacts"
rag-eval study compile experiments/rag-sol2/study.js \\
  --inputs experiments/rag-sol2/inputs.json \\
  --ttc-database /path/to/rag-eval.db \\
  --artifact-root "$artifact_root" \\
  --output-dir "$artifact_root/inputs/rag-study" \\
  --experiment-id EXP-RAG

researchctl experiment run-plan "$artifact_root/inputs/rag-study/researchctl-plan.js" \\
  --project project.js \\
  --runner-command rag-workflow-runner \\
  --runner-name scraper-workflow-runner \\
  --runner-version v1`}</pre>
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
