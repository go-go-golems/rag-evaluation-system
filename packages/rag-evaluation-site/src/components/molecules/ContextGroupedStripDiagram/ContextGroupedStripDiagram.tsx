import { type HTMLAttributes, type ReactNode, useState } from "react";
import {
	type ContextStyleSet,
	type ContextWindowPart,
	type ContextWindowSnapshot,
	contextVisualStyleToCssVars,
	resolveContextVisualStyle,
} from "../../../context";
import { Text } from "../../foundation";
import { handleContextPartKeyDown } from "../contextKeyboardNavigation";
import styles from "./ContextGroupedStripDiagram.module.css";

export type ContextGroupedStripGroupBy = "turn" | "styleKey" | "sourceId";

export interface ContextGroupedStripDiagramProps extends HTMLAttributes<HTMLDivElement> {
	snapshot: ContextWindowSnapshot;
	styleSet: ContextStyleSet;
	selectedPartId?: string;
	groupBy?: ContextGroupedStripGroupBy;
	showGroupLabels?: boolean;
	showPartLabels?: boolean;
	showSelection?: boolean;
	onPartSelect?: (partId: string) => void;
	renderPartTooltip?: (part: ContextWindowPart) => ReactNode;
}

interface Group {
	id: string;
	label: string;
	tokens: number;
	parts: ContextWindowPart[];
}

function formatTokens(tokens: number) {
	return `${tokens.toLocaleString()} tok`;
}
function patternClass(pattern: string | undefined) {
	return styles[`pattern_${pattern ?? "none"}`] ?? styles.pattern_none;
}
function styleName(styleSet: ContextStyleSet, styleKey: string) {
	return styleSet.legend.find((item) => (item.styleKey ?? item.id) === styleKey)?.label ?? styleKey;
}
function metadataString(part: ContextWindowPart, key: string) {
	const value = part.metadata?.[key];
	return typeof value === "string" || typeof value === "number" || typeof value === "boolean"
		? String(value)
		: "";
}
function groupKey(part: ContextWindowPart, groupBy: ContextGroupedStripGroupBy) {
	if (groupBy === "styleKey") return part.styleKey;
	if (groupBy === "sourceId") return part.sourceId || part.id;
	return (
		metadataString(part, "turnIndex") ||
		metadataString(part, "turn") ||
		part.sourceId?.match(/turn-(\d+)/)?.[1] ||
		"global"
	);
}
function groupLabel(
	key: string,
	groupBy: ContextGroupedStripGroupBy,
	firstPart: ContextWindowPart,
	styleSet: ContextStyleSet,
) {
	if (groupBy === "styleKey") return styleName(styleSet, key);
	if (groupBy === "sourceId") return firstPart.sourceId || firstPart.label;
	return key === "global" ? "global" : `turn ${key}`;
}
function groupedParts(
	parts: ContextWindowPart[],
	groupBy: ContextGroupedStripGroupBy,
	styleSet: ContextStyleSet,
): Group[] {
	const groups: Group[] = [];
	const byKey = new Map<string, Group>();
	for (const part of parts) {
		const key = groupKey(part, groupBy);
		let group = byKey.get(key);
		if (!group) {
			group = { id: key, label: groupLabel(key, groupBy, part, styleSet), tokens: 0, parts: [] };
			byKey.set(key, group);
			groups.push(group);
		}
		group.parts.push(part);
		group.tokens += part.tokens;
	}
	return groups;
}
function metadataEntries(part: ContextWindowPart) {
	return Object.entries(part.metadata ?? {})
		.filter(([, value]) => value !== undefined && value !== null && value !== "")
		.slice(0, 6);
}

function defaultPartTooltip(part: ContextWindowPart, styleSet: ContextStyleSet) {
	const preview = part.contentPreview || part.content || "";
	const metadata = metadataEntries(part);
	return (
		<div className={styles.tooltipContent}>
			<div className={styles.tooltipTitleRow}>
				<Text as="strong" size="compact">
					{part.label}
				</Text>
				<Text as="span" size="metadata" tone="muted">
					{formatTokens(part.tokens)} · {styleName(styleSet, part.styleKey)}
				</Text>
			</div>
			{part.note && (
				<Text as="span" size="metadata">
					{part.note}
				</Text>
			)}
			{preview && <pre className={styles.tooltipPreview}>{preview}</pre>}
			{metadata.length > 0 && (
				<dl className={styles.tooltipMetadata}>
					{metadata.map(([key, value]) => (
						<div key={key} className={styles.tooltipMetadataItem}>
							<dt>{key}</dt>
							<dd>{String(value)}</dd>
						</div>
					))}
				</dl>
			)}
		</div>
	);
}
export function ContextGroupedStripDiagram({
	snapshot,
	styleSet,
	selectedPartId,
	groupBy = "turn",
	showGroupLabels = true,
	showPartLabels = false,
	showSelection,
	onPartSelect,
	renderPartTooltip,
	className,
	...rest
}: ContextGroupedStripDiagramProps) {
	const [tooltipPartId, setTooltipPartId] = useState<string | undefined>();
	const effectiveSelectedPartId = selectedPartId ?? snapshot.selectedPartId;
	const totalWidth =
		snapshot.limit || snapshot.parts.reduce((sum, part) => sum + part.tokens, 0) || 1;
	const groups = groupedParts(snapshot.parts, groupBy, styleSet);
	const orderedPartIds = groups.flatMap((group) => group.parts.map((part) => part.id));

	const tooltipPart = snapshot.parts.find((part) => part.id === tooltipPartId);

	return (
		<div
			className={[styles.root, className ?? ""].filter(Boolean).join(" ")}
			data-rag-molecule="ContextGroupedStripDiagram"
			{...rest}
		>
			<div
				className={styles.groups}
				role="img"
				aria-label={`${snapshot.title} grouped context strip`}
			>
				{groups.map((group) => {
					const groupWidth = Math.max(2, (group.tokens / totalWidth) * 100);
					const groupTotal = group.tokens || 1;
					return (
						<div key={group.id} className={styles.group} style={{ width: `${groupWidth}%` }}>
							{showGroupLabels && groupWidth >= 6 && (
								<div className={styles.groupLabel}>{group.label}</div>
							)}
							<div className={styles.groupParts}>
								{group.parts.map((part) => {
									const visualStyle = resolveContextVisualStyle(part.styleKey, styleSet);
									const width = Math.max(2, (part.tokens / groupTotal) * 100);
									const interactive = Boolean(onPartSelect);
									const selected =
										(showSelection ?? interactive) && effectiveSelectedPartId === part.id;
									return (
										<div
											key={part.id}
											className={[
												styles.segment,
												patternClass(visualStyle.pattern),
												selected ? styles.selected : "",
											]
												.filter(Boolean)
												.join(" ")}
											style={{ width: `${width}%`, ...contextVisualStyleToCssVars(visualStyle) }}
											onPointerEnter={() => setTooltipPartId(part.id)}
											onPointerLeave={() => setTooltipPartId(undefined)}
											onFocus={() => setTooltipPartId(part.id)}
											onBlur={() => setTooltipPartId(undefined)}
											role={interactive ? "button" : undefined}
											tabIndex={interactive ? 0 : undefined}
											aria-pressed={interactive ? selected : undefined}
											onClick={interactive ? () => onPartSelect?.(part.id) : undefined}
											onKeyDown={
												interactive
													? (event) =>
															handleContextPartKeyDown(
																event,
																part.id,
																orderedPartIds,
																onPartSelect,
																"horizontal",
																effectiveSelectedPartId,
															)
													: undefined
											}
										>
											{showPartLabels && width >= 12 && (
												<span className={styles.partLabel}>{part.label}</span>
											)}
										</div>
									);
								})}
							</div>
						</div>
					);
				})}
			</div>
			<div className={styles.caption}>
				{snapshot.title} · {snapshot.limit.toLocaleString()} token window · grouped by {groupBy}
			</div>
			<div className={styles.tooltipPanel} role="tooltip" aria-live="polite">
				{tooltipPart ? (
					renderPartTooltip ? (
						renderPartTooltip(tooltipPart)
					) : (
						defaultPartTooltip(tooltipPart, styleSet)
					)
				) : (
					<Text as="span" size="metadata" tone="muted">
						Hover a block for details.
					</Text>
				)}
			</div>
		</div>
	);
}
