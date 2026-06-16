import { useMemo, useState, type HTMLAttributes } from "react";
import type {
	ContextJsonValue,
	ContextStyleSet,
	ContextWindowPart,
	ContextWindowSnapshot,
} from "../../../context";
import { Button } from "../../atoms";
import { Inline, Panel, Stack } from "../../layout";
import {
	ContextLegend,
	ContextStripDiagram,
	ContextGroupedStripDiagram,
	type ContextGroupedStripGroupBy,
} from "../../molecules";
import styles from "./ContextTurnPagerPanel.module.css";

export interface ContextTurnPagerPanelProps extends HTMLAttributes<HTMLDivElement> {
	snapshots: ContextWindowSnapshot[];
	styleSet: ContextStyleSet;
	initialSnapshotId?: string;
	selectedPartId?: string;
	diagram?: "grouped-strip" | "strip";
	groupBy?: ContextGroupedStripGroupBy;
	mode?: "turn-only" | "snapshot";
	includeGlobalParts?: boolean;
	showLegend?: boolean;
	title?: string;
}

function metadataScalar(value: ContextJsonValue | undefined) {
	if (typeof value === "string" || typeof value === "number" || typeof value === "boolean")
		return String(value);
	return "";
}

function partTurnKey(part: ContextWindowPart) {
	return (
		metadataScalar(part.metadata?.turnIndex) ||
		metadataScalar(part.metadata?.turn) ||
		part.sourceId?.match(/turn-(\d+)/)?.[1] ||
		""
	);
}

function snapshotTurnKey(snapshot: ContextWindowSnapshot) {
	return (
		metadataScalar(snapshot.metadata?.turnIndex) || metadataScalar(snapshot.metadata?.turn) || ""
	);
}

function isGlobalPart(part: ContextWindowPart) {
	const turn = partTurnKey(part).toLowerCase();
	const blockType = metadataScalar(part.metadata?.blockType).toLowerCase();
	return (
		turn === "global" ||
		blockType === "system-estimate" ||
		blockType === "system" ||
		part.styleKey === "system"
	);
}

function turnOnlySnapshot(
	snapshot: ContextWindowSnapshot,
	includeGlobalParts: boolean,
): ContextWindowSnapshot {
	const turn = snapshotTurnKey(snapshot);
	if (!turn) return snapshot;
	const parts = snapshot.parts.filter(
		(part) => partTurnKey(part) === turn || (includeGlobalParts && isGlobalPart(part)),
	);
	const tokenTotal = parts.reduce((sum, part) => sum + part.tokens, 0);
	const selectedPartId = parts.some((part) => part.id === snapshot.selectedPartId)
		? snapshot.selectedPartId
		: parts[0]?.id;
	return {
		...snapshot,
		title: snapshot.title || `Turn ${turn} context window`,
		limit: tokenTotal || snapshot.limit,
		subtitle: tokenTotal
			? `${tokenTotal.toLocaleString()} tokens in this turn`
			: "No blocks for this turn",
		selectedPartId,
		parts,
	};
}

function turnLabel(snapshot: ContextWindowSnapshot, index: number) {
	const turn = snapshot.metadata?.turnIndex;
	if (typeof turn === "number" || typeof turn === "string") return `T${turn}`;
	return `T${index + 1}`;
}

export function ContextTurnPagerPanel({
	snapshots,
	styleSet,
	initialSnapshotId,
	selectedPartId,
	diagram = "grouped-strip",
	groupBy = "turn",
	mode = "turn-only",
	includeGlobalParts = true,
	showLegend = true,
	title = "Context turn pager",
	className,
	...rest
}: ContextTurnPagerPanelProps) {
	const initialIndex = Math.max(
		0,
		snapshots.findIndex((snapshot) => snapshot.id === initialSnapshotId),
	);
	const [index, setIndex] = useState(initialIndex === -1 ? 0 : initialIndex);
	const boundedIndex = Math.min(Math.max(index, 0), Math.max(0, snapshots.length - 1));
	const rawSnapshot = snapshots[boundedIndex];
	const snapshot =
		rawSnapshot && mode === "turn-only"
			? turnOnlySnapshot(rawSnapshot, includeGlobalParts)
			: rawSnapshot;
	const selected = selectedPartId ?? snapshot?.selectedPartId;
	const visibleLegend = useMemo(() => {
		if (!snapshot) return styleSet.legend;
		const visibleStyleKeys = new Set(snapshot.parts.map((part) => part.styleKey));
		return styleSet.legend.filter((item) => visibleStyleKeys.has(item.styleKey ?? item.id));
	}, [snapshot, styleSet.legend]);
	const actions = useMemo(
		() => (
			<Inline gap="xs">
				<Button
					size="compact"
					disabled={boundedIndex <= 0}
					onClick={() => setIndex((value) => Math.max(0, value - 1))}
				>
					Prev
				</Button>
				<span className={styles.counter}>
					{snapshots.length ? `${boundedIndex + 1} / ${snapshots.length}` : "0 / 0"}
				</span>
				<Button
					size="compact"
					disabled={boundedIndex >= snapshots.length - 1}
					onClick={() => setIndex((value) => Math.min(snapshots.length - 1, value + 1))}
				>
					Next
				</Button>
			</Inline>
		),
		[boundedIndex, snapshots.length],
	);

	if (!snapshot) {
		return (
			<Panel
				title={title}
				className={className}
				data-rag-organism="ContextTurnPagerPanel"
				{...rest}
			>
				No context snapshots available.
			</Panel>
		);
	}

	return (
		<Panel
			title={title}
			actions={actions}
			className={className}
			data-rag-organism="ContextTurnPagerPanel"
			{...rest}
		>
			<Stack gap="sm">
				<Inline gap="xs" className={styles.turnRail} aria-label="Available turns">
					{snapshots.map((item, itemIndex) => (
						<Button
							key={item.id}
							size="compact"
							selected={itemIndex === boundedIndex}
							aria-pressed={itemIndex === boundedIndex}
							onClick={() => setIndex(itemIndex)}
						>
							{turnLabel(item, itemIndex)}
						</Button>
					))}
				</Inline>
				<div className={styles.heading}>
					<strong>{snapshot.title}</strong>
					{snapshot.subtitle && <span>{snapshot.subtitle}</span>}
				</div>
				{diagram === "strip" ? (
					<ContextStripDiagram snapshot={snapshot} styleSet={styleSet} selectedPartId={selected} />
				) : (
					<ContextGroupedStripDiagram
						snapshot={snapshot}
						styleSet={styleSet}
						selectedPartId={selected}
						groupBy={groupBy}
					/>
				)}
				{showLegend && (
					<ContextLegend items={visibleLegend} styles={styleSet.styles} selectedId={selected} />
				)}
			</Stack>
		</Panel>
	);
}
