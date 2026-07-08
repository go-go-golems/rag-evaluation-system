import type { ButtonHTMLAttributes, ReactNode } from "react";
import { type ContextStyleSet, contextVisualStyleToCssVars } from "../../../context";
import styles from "./CycleCell.module.css";

export type CycleCellSize = "sm" | "md";

export interface CycleCellProps
	extends Omit<ButtonHTMLAttributes<HTMLButtonElement>, "onClick" | "value"> {
	/** Current state id. Must be one of `states`. */
	value: string;
	/** Ordered ring of state ids; clicking advances to the next one. */
	states: string[];
	/** state id -> glyph. Falls back to the state's first character. */
	glyphs?: Record<string, ReactNode>;
	/**
	 * Optional palette. The current state id is looked up as a styleKey in
	 * `styleSet.styles`; the resolved fill/label color themes the cell.
	 */
	styleSet?: ContextStyleSet;
	size?: CycleCellSize;
	selected?: boolean;
	readOnly?: boolean;
	/** Fired with the next state id when the cell is activated. */
	onCycle?: (next: string) => void;
}

function nextState(states: string[], current: string): string {
	if (states.length === 0) return current;
	const idx = states.indexOf(current);
	return states[(idx + 1) % states.length] ?? current;
}

export function CycleCell({
	value,
	states,
	glyphs,
	styleSet,
	size = "md",
	selected = false,
	readOnly = false,
	onCycle,
	className,
	style,
	disabled,
	...rest
}: CycleCellProps) {
	const visualStyle = styleSet?.styles[value] ?? styleSet?.fallbackStyle;
	const glyph = glyphs?.[value] ?? value.charAt(0).toUpperCase();
	const interactive = !readOnly && !disabled;

	return (
		<button
			type="button"
			className={[styles.root, styles[size], selected ? styles.selected : "", className ?? ""]
				.filter(Boolean)
				.join(" ")}
			data-rag-atom="CycleCell"
			data-state={value}
			data-readonly={readOnly || undefined}
			aria-label={value}
			aria-pressed={selected}
			disabled={!interactive}
			onClick={interactive ? () => onCycle?.(nextState(states, value)) : undefined}
			style={visualStyle ? { ...contextVisualStyleToCssVars(visualStyle), ...style } : style}
			{...rest}
		>
			<span className={styles.glyph}>{glyph}</span>
		</button>
	);
}
