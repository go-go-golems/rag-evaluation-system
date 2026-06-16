import { useState } from "react";
import { Button, Caption, Panel, Text } from "@go-go-golems/rag-evaluation-site";
import { DashboardGrid, ScrollRegion, Stack } from "@go-go-golems/rag-evaluation-site";
import {
	ContextDiagramPanel,
	type ContextDiagramView,
	contextWindowSnapshots,
} from "@go-go-golems/rag-evaluation-site";
import styles from "./ContextVisualizerPage.module.css";

const SNAPSHOTS = contextWindowSnapshots;
const VIEWS: { id: ContextDiagramView; label: string }[] = [
	{ id: "strip", label: "Strip" },
	{ id: "stack", label: "Stack" },
	{ id: "budget", label: "Budget" },
	{ id: "treemap", label: "Treemap" },
];

export function ContextVisualizerPage() {
	const [selectedSnapshotId, setSelectedSnapshotId] = useState(SNAPSHOTS[1]!.id);
	const [view, setView] = useState<ContextDiagramView>("stack");

	const snapshot = SNAPSHOTS.find((s) => s.id === selectedSnapshotId) ?? SNAPSHOTS[1]!;

	return (
		<DashboardGrid className={styles.root} data-rag-page="ContextVisualizerPage">
			{/* Left: snapshot selector */}
			<Stack gap="sm" className={styles.controls}>
				<Panel title="Snapshots" density="condensed">
					<ScrollRegion>
						<Stack gap="xs">
							{SNAPSHOTS.map((s) => (
								<Button
									key={s.id}
									size="compact"
									selected={s.id === selectedSnapshotId}
									onClick={() => setSelectedSnapshotId(s.id)}
									className={styles.selectorButton}
								>
									<Stack gap="xs">
										<Text size="label" tone="primary">
											{s.title}
										</Text>
										{s.subtitle && <Caption>{s.subtitle}</Caption>}
									</Stack>
								</Button>
							))}
						</Stack>
					</ScrollRegion>
				</Panel>
				<Panel title="Views" density="condensed">
					<Stack gap="xs">
						{VIEWS.map((v) => (
							<Button
								key={v.id}
								size="compact"
								selected={v.id === view}
								onClick={() => setView(v.id)}
							>
								{v.label}
							</Button>
						))}
					</Stack>
				</Panel>
			</Stack>

			{/* Center: diagram */}
			<ScrollRegion className={styles.diagramArea}>
				<ContextDiagramPanel
					key={`${snapshot.id}-${view}`}
					snapshot={snapshot}
					initialView={view}
					showPartDetails
				/>
			</ScrollRegion>
		</DashboardGrid>
	);
}

export default ContextVisualizerPage;
