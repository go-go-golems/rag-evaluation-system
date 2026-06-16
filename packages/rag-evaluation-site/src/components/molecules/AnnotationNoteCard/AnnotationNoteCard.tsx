import type { CSSProperties, HTMLAttributes } from "react";
import {
	contextVisualStyleToCssVars,
	resolveContextVisualStyle,
	transcriptDefaultStyleSet,
	type ContextStyleSet,
	type ContextVisualStyle,
	type TranscriptAnnotation,
} from "../../../context";
import { ContextStyleSwatch } from "../../atoms";
import { Caption } from "../../foundation";
import styles from "./AnnotationNoteCard.module.css";

export interface AnnotationNoteCardProps extends HTMLAttributes<HTMLElement> {
	annotation: TranscriptAnnotation;
	styleSet?: ContextStyleSet;
	selected?: boolean;
	index?: number;
}

function patternClass(style: ContextVisualStyle) {
	return styles[`pattern_${style.pattern ?? "none"}`] ?? styles.pattern_none;
}

export function AnnotationNoteCard({
	annotation,
	styleSet = transcriptDefaultStyleSet,
	selected = false,
	index,
	className,
	style,
	...rest
}: AnnotationNoteCardProps) {
	const visualStyle = resolveContextVisualStyle(annotation.styleKey, styleSet);
	return (
		<aside
			className={[
				styles.root,
				patternClass(visualStyle),
				selected ? styles.selected : "",
				className ?? "",
			]
				.filter(Boolean)
				.join(" ")}
			data-rag-molecule="AnnotationNoteCard"
			data-rag-note-id={annotation.id}
			data-rag-target-message-id={annotation.targetMessageId}
			style={{ ...contextVisualStyleToCssVars(visualStyle), ...style } as CSSProperties}
			{...rest}
		>
			<header className={styles.titleBar}>
				<span className={styles.noteLabel}>NOTE {index ?? ""}</span>
				<ContextStyleSwatch visualStyle={visualStyle} size="sm" selected={selected} />
			</header>
			<div className={styles.body}>
				<div className={styles.title}>{annotation.label}</div>
				<p className={styles.text}>{annotation.text}</p>
				{annotation.confidence != null && (
					<Caption tone="inherit" className={styles.confidenceChip}>
						confidence {Math.round(annotation.confidence * 100)}%
					</Caption>
				)}
			</div>
		</aside>
	);
}
