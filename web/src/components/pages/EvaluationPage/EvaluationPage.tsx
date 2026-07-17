import { Caption, Panel, Stack, Text } from "@go-go-golems/rag-evaluation-site";
import { useGetLabCatalogQuery } from "../../../services/api";

export function EvaluationPage() {
	const catalog = useGetLabCatalogQuery();
	const snapshots = catalog.data?.snapshots.length ?? 0;
	const chunkSets = catalog.data?.chunk_sets.length ?? 0;
	const embeddingSets = catalog.data?.embedding_sets.length ?? 0;
	const bm25Artifacts = catalog.data?.bm25_artifacts.length ?? 0;

	return (
		<Stack gap="lg">
			<Panel title="Researchctl RAG laboratory">
				<Stack gap="md">
					<Text>
						Native RAG execution, run lifecycle, traces, metrics, export, and import now belong
						exclusively to the project-local researchctl laboratory.
					</Text>
					<Caption>
						The former writable rag-eval experiment specification and run endpoints were removed
						after external-import and native-rerun parity validation.
					</Caption>
					<pre>{`researchctl experiment run-rag experiment.js \\
  --project project.yaml \\
  --experiment-id EXP-RAG \\
  --inputs inputs.json \\
  --ttc-database /path/to/rag-eval.db \\
  --runner /path/to/rag-lab-worker

researchctl lab runs list --project project.yaml --output json
researchctl lab runs show RUN_ID --project project.yaml --output json`}</pre>
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
