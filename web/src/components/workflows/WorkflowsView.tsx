import { Button, Caption, Panel, Stack, Text } from "@go-go-golems/rag-evaluation-site";
import type React from "react";
import { useEffect, useState } from "react";
import {
	type SubmitIntakeRequest,
	useCancelIntakeRunMutation,
	useGetIntakeObservationsQuery,
	useGetIntakeRunQuery,
	useListIntakeRunsQuery,
	useSubmitIntakeRunMutation,
} from "../../services/api";
import styles from "./WorkflowsView.module.css";

const SubmitIntake: React.FC<{ onSubmitted: (id: string) => void }> = ({ onSubmitted }) => {
	const [documents, setDocuments] = useState("");
	const [sources, setSources] = useState("");
	const [skipEmbeddings, setSkipEmbeddings] = useState(true);
	const [submit, state] = useSubmitIntakeRunMutation();
	const run = async () => {
		const request: SubmitIntakeRequest = {
			document_ids: documents
				.split(",")
				.map((v) => v.trim())
				.filter(Boolean),
			source_ids: sources
				.split(",")
				.map((v) => v.trim())
				.filter(Boolean),
			skip_preprocessing: true,
			skip_chunk_enrichment: true,
			skip_embeddings: skipEmbeddings,
		};
		const result = await submit(request).unwrap();
		onSubmitted(result.submission.runId);
	};
	return (
		<Panel title="Submit Workflow V3 intake">
			<Stack gap="sm">
				<label>
					Document IDs{" "}
					<input
						className={styles.fullInput}
						value={documents}
						onChange={(e) => setDocuments(e.target.value)}
						placeholder="comma-separated; optional with source IDs"
					/>
				</label>
				<label>
					Source IDs{" "}
					<input
						className={styles.fullInput}
						value={sources}
						onChange={(e) => setSources(e.target.value)}
						placeholder="comma-separated"
					/>
				</label>
				<label>
					<input
						type="checkbox"
						checked={skipEmbeddings}
						onChange={(e) => setSkipEmbeddings(e.target.checked)}
					/>{" "}
					Skip embeddings
				</label>
				<Button onClick={run} disabled={state.isLoading || (!documents.trim() && !sources.trim())}>
					Submit intake
				</Button>
				{state.error && (
					<Text>Submission failed. Inspect server logs for the stable error code.</Text>
				)}
			</Stack>
		</Panel>
	);
};

export const WorkflowsView: React.FC = () => {
	const [status, setStatus] = useState("");
	const [selected, setSelected] = useState<string | null>(null);
	const runs = useListIntakeRunsQuery({ status, limit: 50 }, { pollingInterval: 3000 });
	const run = useGetIntakeRunQuery(selected ?? "", { skip: !selected, pollingInterval: 2000 });
	const terminal = ["succeeded", "failed", "canceled"].includes(run.data?.snapshot.status ?? "");
	const observations = useGetIntakeObservationsQuery(selected ?? "", {
		skip: !selected || !terminal,
		pollingInterval: 3000,
	});
	const [cancel] = useCancelIntakeRunMutation();
	useEffect(() => {
		const first = runs.data?.[0];
		if (!selected && first) setSelected(first.runId);
	}, [runs.data, selected]);
	return (
		<div className={styles.layout}>
			<Stack gap="md">
				<SubmitIntake onSubmitted={setSelected} />
				<Panel title="Intake runs">
					<label>
						Status{" "}
						<select value={status} onChange={(e) => setStatus(e.target.value)}>
							<option value="">all</option>
							<option>running</option>
							<option>succeeded</option>
							<option>failed</option>
							<option>canceled</option>
						</select>
					</label>
					<div className={styles.runList}>
						{runs.isLoading && <Text>Loading intake runs…</Text>}
						{runs.isError && <Text>Unable to load intake runs.</Text>}
						{runs.data?.length === 0 && <Text>No intake runs match this filter.</Text>}
						{runs.data?.map((item) => (
							<button
								type="button"
								key={item.runId}
								className={item.runId === selected ? styles.selectedRun : styles.run}
								onClick={() => setSelected(item.runId)}
							>
								<strong>{item.runId}</strong>
								<span>{item.status}</span>
								<Caption>{item.planDigest.slice(0, 20)}…</Caption>
							</button>
						))}
					</div>
				</Panel>
			</Stack>
			<Panel title={selected ? `Workflow V3 run ${selected}` : "Workflow V3 run"}>
				{run.data ? (
					<Stack gap="sm">
						<Text>Status: {run.data.snapshot.status}</Text>
						<Text>Plan: {run.data.snapshot.planDigest}</Text>
						{run.data.snapshot.status === "running" && (
							<Button onClick={() => selected && cancel(selected)}>Cancel run</Button>
						)}
						<table className={styles.attempts}>
							<thead>
								<tr>
									<th>Node</th>
									<th>Attempt</th>
									<th>Status</th>
									<th>Resource</th>
									<th>Failure</th>
								</tr>
							</thead>
							<tbody>
								{(run.data.snapshot.attempts ?? []).map((attempt) => (
									<tr key={`${attempt.nodeKey}-${attempt.number}`}>
										<td>{attempt.nodeKey}</td>
										<td>{attempt.number}</td>
										<td>{attempt.status}</td>
										<td>{attempt.resourceClass}</td>
										<td>{attempt.failure?.code ?? "—"}</td>
									</tr>
								))}
							</tbody>
						</table>
						{observations.data && (
							<details>
								<summary>Canonical observations</summary>
								<pre>{JSON.stringify(observations.data, null, 2)}</pre>
							</details>
						)}
					</Stack>
				) : (
					<Text>
						{runs.data?.length === 0
							? "Submit an intake run to inspect its Workflow V3 execution."
							: "Select an intake run."}
					</Text>
				)}
			</Panel>
		</div>
	);
};
