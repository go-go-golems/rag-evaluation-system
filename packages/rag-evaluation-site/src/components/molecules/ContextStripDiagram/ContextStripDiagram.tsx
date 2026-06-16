import { useState, type HTMLAttributes, type ReactNode } from "react";
import {
	contextVisualStyleToCssVars,
	resolveContextVisualStyle,
	type ContextStyleSet,
	type ContextWindowPart,
	type ContextWindowSnapshot,
} from "../../../context";
import { Text } from "../../foundation";
import { handleContextPartKeyDown } from "../contextKeyboardNavigation";
import styles from "./ContextStripDiagram.module.css";

export interface ContextStripDiagramProps extends HTMLAttributes<HTMLDivElement> {
	snapshot: ContextWindowSnapshot;
	styleSet: ContextStyleSet;
	selectedPartId?: string;
	showLabels?: boolean;
	showSelection?: boolean;
	onPartSelect?: (partId: string) => void;
	renderPartTooltip?: (part: ContextWindowPart) => ReactNode;
}

function formatTokens(tokens: number) {
	return `${tokens.toLocaleString()} tok`;
}
function styleName(styleSet: ContextStyleSet, styleKey: string) {
	return styleSet.legend.find((item) => (item.styleKey ?? item.id) === styleKey)?.label ?? styleKey;
}
function patternClass(pattern: string | undefined) {
	return styles[`pattern_${pattern ?? "none"}`] ?? styles.pattern_none;
}
function defaultPartTooltip(part: ContextWindowPart, styleSet: ContextStyleSet) {
	return (
		<div className={styles.tooltipContent}>
			<Text as="strong" size="compact">
				{part.label}
			</Text>
			<Text as="span" size="metadata" tone="muted">
				{formatTokens(part.tokens)} · {styleName(styleSet, part.styleKey)}
			</Text>
			{part.note && (
				<Text as="span" size="metadata">
					{part.note}
				</Text>
			)}
		</div>
	);
}

export function ContextStripDiagram({
	snapshot,
	styleSet,
	selectedPartId,
	showLabels = true,
	showSelection,
	onPartSelect,
	renderPartTooltip,
	className,
	...rest
}: ContextStripDiagramProps) {
	const [tooltipPartId, setTooltipPartId] = useState<string | undefined>();
	const effectiveSelectedPartId = selectedPartId ?? snapshot.selectedPartId;
	const totalWidth = snapshot.limit || snapshot.parts.reduce((sum, part) => sum + part.tokens, 0);
	const orderedPartIds = snapshot.parts.map((part) => part.id);

	const tooltipPart = snapshot.parts.find((part) => part.id === tooltipPartId);

	return (
		<div
			className={[styles.root, className ?? ""].filter(Boolean).join(" ")}
			data-rag-molecule="ContextStripDiagram"
			{...rest}
		>
			<div className={styles.strip} role="img" aria-label={`${snapshot.title} context strip`}>
				{snapshot.parts.map((part) => {
					const visualStyle = resolveContextVisualStyle(part.styleKey, styleSet);
					const width = Math.max(2, (part.tokens / totalWidth) * 100);
					const interactive = Boolean(onPartSelect);
					const selected = (showSelection ?? interactive) && effectiveSelectedPartId === part.id;
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
							{showLabels && width >= 7 && <span className={styles.label}>{part.label}</span>}
						</div>
					);
				})}
			</div>
			<div className={styles.caption}>
				{snapshot.title} · {snapshot.limit.toLocaleString()} token window
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
