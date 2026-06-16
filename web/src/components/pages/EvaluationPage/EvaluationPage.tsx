import { Caption, Text } from "@go-go-golems/rag-evaluation-site";
import { Panel, Stack } from "@go-go-golems/rag-evaluation-site";
import styles from "./EvaluationPage.module.css";

export interface EvaluationPageProps {
	title?: string;
	capabilities?: string[];
}

const defaultCapabilities = [
	"Create and manage evaluation query sets",
	"Run evaluations across different search configurations",
	"Compare Recall@K, MRR, nDCG@K metrics",
	"See per-query pass/fail matrices",
];

export function EvaluationPage({
	title = "Evaluation Dashboard",
	capabilities = defaultCapabilities,
}: EvaluationPageProps) {
	return (
		<Panel title={title} data-rag-page="EvaluationPage">
			<Stack gap="sm">
				<Text tone="muted">Evaluation dashboard coming soon. This view will let you:</Text>
				<ul className={styles.list}>
					{capabilities.map((capability) => (
						<li key={capability}>{capability}</li>
					))}
				</ul>
				<Caption>
					Placeholder page boundary; wire evaluation DTOs here when the evaluation API surface is
					promoted.
				</Caption>
			</Stack>
		</Panel>
	);
}
