import type { CSSProperties, HTMLAttributes } from "react";
import {
	type ContextStyleSet,
	type ContextVisualStyle,
	contextVisualStyleToCssVars,
	resolveContextVisualStyle,
	type TranscriptAnnotation,
	type TranscriptMessage,
	type TranscriptRole,
	useTranscriptStyleSet,
} from "../../../context";
import { Caption } from "../../foundation";
import styles from "./TranscriptMessageCard.module.css";

export interface TranscriptMessageCardProps extends HTMLAttributes<HTMLElement> {
	message: TranscriptMessage;
	annotations?: TranscriptAnnotation[];
	selectedAnnotationId?: string;
	onAnnotationSelect?: (annotationId: string) => void;
	showAnnotationChips?: boolean;
	styleSet?: ContextStyleSet;
}

const ROLE_LABEL: Record<TranscriptRole, string> = {
	system: "SYSTEM",
	developer: "DEVELOPER",
	user: "USER",
	assistant: "ASSISTANT",
	tool: "TOOL",
	other: "MESSAGE",
};

const ROLE_GLYPH: Record<TranscriptRole, string> = {
	system: "◆",
	developer: "◇",
	user: "▸",
	assistant: "✦",
	tool: "⚙",
	other: "•",
};

function annotationIndex(annotations: TranscriptAnnotation[], annotationId: string) {
	return annotations.findIndex((annotation) => annotation.id === annotationId) + 1;
}

function roleStyleKey(role: TranscriptRole) {
	return role === "other" ? "other" : role;
}

function patternClass(style: ContextVisualStyle) {
	return styles[`pattern_${style.pattern ?? "none"}`] ?? styles.pattern_none;
}

function toneVars(style: ContextVisualStyle): CSSProperties {
	return contextVisualStyleToCssVars(style);
}

export function TranscriptMessageCard({
	message,
	annotations = [],
	selectedAnnotationId,
	onAnnotationSelect,
	showAnnotationChips = true,
	styleSet,
	className,
	style,
	...rest
}: TranscriptMessageCardProps) {
	const resolvedStyleSet = useTranscriptStyleSet(styleSet);
	const messageAnnotations = annotations.filter(
		(annotation) =>
			annotation.targetMessageId === message.id || message.annotationIds?.includes(annotation.id),
	);
	const hasSelectedAnnotation = messageAnnotations.some(
		(annotation) => annotation.id === selectedAnnotationId,
	);
	const title = message.name
		? `${ROLE_LABEL[message.role]} · ${message.name}`
		: ROLE_LABEL[message.role];
	const visualStyle = resolveContextVisualStyle(roleStyleKey(message.role), resolvedStyleSet);
	const codeStyle = resolveContextVisualStyle("code", resolvedStyleSet);
	const bodyClassName = [styles.body, message.role === "tool" ? styles.toolBody : ""]
		.filter(Boolean)
		.join(" ");

	return (
		<article
			className={[
				styles.root,
				styles[message.role] ?? styles.other,
				patternClass(visualStyle),
				hasSelectedAnnotation ? styles.selected : "",
				className ?? "",
			]
				.filter(Boolean)
				.join(" ")}
			data-rag-molecule="TranscriptMessageCard"
			data-rag-message-id={message.id}
			data-rag-role={message.role}
			data-rag-has-notes={messageAnnotations.length > 0 ? "true" : undefined}
			style={
				{
					...toneVars(visualStyle),
					"--tx-code-fill": codeStyle.fill,
					"--tx-code-line": codeStyle.line ?? codeStyle.stroke ?? "var(--mac-border)",
					...style,
				} as CSSProperties
			}
			{...rest}
		>
			<header className={styles.titleBar}>
				<div className={styles.titleCluster}>
					<span className={styles.roleGlyph} aria-hidden="true">
						{ROLE_GLYPH[message.role]}
					</span>
					<span className={styles.roleLabel}>{title}</span>
				</div>
				<div className={styles.metaCluster}>
					{showAnnotationChips &&
						messageAnnotations.map((annotation) => {
							const index = annotationIndex(annotations, annotation.id);
							const selected = selectedAnnotationId === annotation.id;
							return (
								<button
									key={annotation.id}
									type="button"
									className={styles.noteChip}
									data-selected={selected ? "true" : undefined}
									style={toneVars(resolveContextVisualStyle(annotation.styleKey, resolvedStyleSet))}
									onClick={() => onAnnotationSelect?.(annotation.id)}
								>
									note {index > 0 ? index : ""}
								</button>
							);
						})}
					<Caption tone="inherit" className={styles.tokenChip}>
						{message.tokens?.toLocaleString() ?? "—"} tok
					</Caption>
				</div>
			</header>
			<div className={bodyClassName}>{message.text}</div>
		</article>
	);
}
