import type { HTMLAttributes } from "react";
import type { ContextStyleSet, TranscriptAnnotation } from "../../../context";
import { Caption, Text } from "../../foundation";
import { AnnotationNoteCard } from "../../molecules";
import styles from "./AnnotationRailPanel.module.css";

export interface AnnotationRailPanelProps extends Omit<HTMLAttributes<HTMLElement>, "title"> {
	title?: string;
	description?: string;
	annotations: TranscriptAnnotation[];
	selectedAnnotationId?: string;
	onAnnotationSelect?: (annotationId: string) => void;
	styleSet?: ContextStyleSet;
}

export function AnnotationRailPanel({
	title = "Context-engineering notes",
	description = "What each message costs the window. Click a note or a message to focus it.",
	annotations,
	selectedAnnotationId,
	onAnnotationSelect,
	styleSet,
	className,
	...rest
}: AnnotationRailPanelProps) {
	return (
		<aside
			className={[styles.root, className ?? ""].filter(Boolean).join(" ")}
			data-rag-organism="AnnotationRailPanel"
			{...rest}
		>
			<header className={styles.header}>
				<Text as="div" size="label">
					{title}
				</Text>
				<Caption className={styles.description}>{description}</Caption>
			</header>
			<div className={styles.list}>
				{annotations.map((annotation, index) => (
					<button
						key={annotation.id}
						type="button"
						className={styles.button}
						onClick={() => onAnnotationSelect?.(annotation.id)}
					>
						<AnnotationNoteCard
							annotation={annotation}
							styleSet={styleSet}
							selected={annotation.id === selectedAnnotationId}
							index={index + 1}
						/>
					</button>
				))}
			</div>
		</aside>
	);
}
