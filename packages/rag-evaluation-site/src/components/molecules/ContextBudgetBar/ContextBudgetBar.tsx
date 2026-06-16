import type { HTMLAttributes } from "react";
import {
	contextVisualStyleToCssVars,
	contextWindowFillRatio,
	contextWindowHeadroomStyleKeys,
	contextWindowIsHeadroomPart,
	contextWindowTokenTotal,
	resolveContextVisualStyle,
	type ContextStyleSet,
	type ContextWindowPart,
	type ContextWindowSnapshot,
} from "../../../context";
import { Caption } from "../../foundation";
import { ContextStyleSwatch } from "../../atoms";
import { handleContextPartKeyDown } from "../contextKeyboardNavigation";
import styles from "./ContextBudgetBar.module.css";

export interface ContextBudgetBarProps extends HTMLAttributes<HTMLDivElement> {
	snapshot: ContextWindowSnapshot;
	styleSet: ContextStyleSet;
	showLegend?: boolean;
	selectedPartId?: string;
	onPartSelect?: (partId: string) => void;
}

function formatTokens(tokens: number) {
	return `${Math.round(tokens).toLocaleString()} tok`;
}
function usedParts(parts: ContextWindowPart[], headroomStyleKeys: Iterable<string>) {
	return parts.filter(
		(part) => !contextWindowIsHeadroomPart(part, { headroomStyleKeys }) && part.tokens > 0,
	);
}
function styleName(styleSet: ContextStyleSet, styleKey: string) {
	return styleSet.legend.find((item) => (item.styleKey ?? item.id) === styleKey)?.label ?? styleKey;
}
function patternClass(pattern: string | undefined) {
	return styles[`pattern_${pattern ?? "none"}`] ?? styles.pattern_none;
}
export function ContextBudgetBar({
	snapshot,
	styleSet,
	showLegend = true,
	selectedPartId,
	onPartSelect,
	className,
	...rest
}: ContextBudgetBarProps) {
	const headroomStyleKeys = contextWindowHeadroomStyleKeys(styleSet);
	const total = contextWindowTokenTotal(snapshot, { headroomStyleKeys });
	const ratio = contextWindowFillRatio(snapshot, { headroomStyleKeys });
	const overBudget = total > snapshot.limit;
	const nearBudget = !overBudget && ratio >= 0.8;
	const parts = usedParts(snapshot.parts, headroomStyleKeys);
	const orderedPartIds = parts.map((part) => part.id);
	const effectiveSelectedPartId = selectedPartId ?? snapshot.selectedPartId;

	return (
		<div
			className={[
				styles.root,
				overBudget ? styles.overBudget : nearBudget ? styles.nearBudget : "",
				className ?? "",
			]
				.filter(Boolean)
				.join(" ")}
			data-rag-molecule="ContextBudgetBar"
			{...rest}
		>
			<div className={styles.header}>
				<Caption transform="uppercase">{snapshot.title}</Caption>
				<Caption tone={overBudget ? "danger" : nearBudget ? "warning" : "muted"}>
					{formatTokens(total)} / {formatTokens(snapshot.limit)} · {Math.round(ratio * 100)}%
				</Caption>
			</div>
			<div
				className={styles.track}
				role="img"
				aria-label={`${snapshot.title}: ${formatTokens(total)} of ${formatTokens(snapshot.limit)} used`}
			>
				{parts.map((part) => {
					const visualStyle = resolveContextVisualStyle(part.styleKey, styleSet);
					const width =
						snapshot.limit > 0 ? Math.max(0.5, (part.tokens / snapshot.limit) * 100) : 0;
					const selected = effectiveSelectedPartId === part.id;
					const interactive = Boolean(onPartSelect);
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
							title={`${part.label}: ${formatTokens(part.tokens)} (${styleName(styleSet, part.styleKey)})`}
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
						/>
					);
				})}
				{overBudget && <div className={styles.limitMarker} title="model context limit" />}
			</div>
			{showLegend && (
				<div className={styles.legend}>
					{parts.map((part) => {
						const visualStyle = resolveContextVisualStyle(part.styleKey, styleSet);
						return (
							<span
								key={part.id}
								className={styles.legendItem}
								data-selected={effectiveSelectedPartId === part.id ? "true" : undefined}
							>
								<ContextStyleSwatch visualStyle={visualStyle} size="sm" />
								<span>{part.label}</span>
								<span className={styles.tokens}>{formatTokens(part.tokens)}</span>
							</span>
						);
					})}
				</div>
			)}
		</div>
	);
}
