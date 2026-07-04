import type { HTMLAttributes } from "react";
import {
	type ContextStyleSet,
	type ContextWindowSnapshot,
	contextVisualStyleToCssVars,
	contextWindowHeadroomStyleKeys,
	contextWindowIsHeadroomPart,
	contextWindowTokenTotal,
	resolveContextVisualStyle,
} from "../../../context";
import { handleContextPartKeyDown } from "../contextKeyboardNavigation";
import styles from "./ContextTreemap.module.css";

export interface ContextTreemapProps extends HTMLAttributes<HTMLDivElement> {
	snapshot: ContextWindowSnapshot;
	styleSet: ContextStyleSet;
	selectedPartId?: string;
	onPartSelect?: (partId: string) => void;
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
export function ContextTreemap({
	snapshot,
	styleSet,
	selectedPartId,
	onPartSelect,
	className,
	...rest
}: ContextTreemapProps) {
	const effectiveSelectedPartId = selectedPartId ?? snapshot.selectedPartId;
	const headroomStyleKeys = contextWindowHeadroomStyleKeys(styleSet);
	const parts = snapshot.parts.filter(
		(part) => !contextWindowIsHeadroomPart(part, { headroomStyleKeys }) && part.tokens > 0,
	);
	const orderedPartIds = parts.map((part) => part.id);
	const total = contextWindowTokenTotal(snapshot, { headroomStyleKeys }) || 1;
	return (
		<div
			className={[styles.root, className ?? ""].filter(Boolean).join(" ")}
			data-rag-molecule="ContextTreemap"
			{...rest}
		>
			<div className={styles.map} role="img" aria-label={`${snapshot.title} token treemap`}>
				{parts.map((part) => {
					const visualStyle = resolveContextVisualStyle(part.styleKey, styleSet);
					const area = Math.max(7, (part.tokens / total) * 100);
					const selected = effectiveSelectedPartId === part.id;
					const interactive = Boolean(onPartSelect);
					return (
						<div
							key={part.id}
							className={[
								styles.tile,
								patternClass(visualStyle.pattern),
								selected ? styles.selected : "",
							]
								.filter(Boolean)
								.join(" ")}
							style={{
								flexBasis: `${area}%`,
								flexGrow: part.tokens,
								...contextVisualStyleToCssVars(visualStyle),
							}}
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
												"both",
												effectiveSelectedPartId,
											)
									: undefined
							}
						>
							<span className={styles.label}>{part.label}</span>
							<span className={styles.tokens}>{formatTokens(part.tokens)}</span>
						</div>
					);
				})}
			</div>
			<div className={styles.caption}>
				{formatTokens(total)} in use · {snapshot.title}
			</div>
		</div>
	);
}
