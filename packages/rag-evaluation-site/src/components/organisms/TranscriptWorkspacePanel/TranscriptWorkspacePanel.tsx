import type { HTMLAttributes } from "react";
import type { ContextStyleSet, TranscriptAnnotation, TranscriptMessage } from "../../../context";
import { AnnotationRailPanel } from "../AnnotationRailPanel";
import { TranscriptReaderPanel } from "../TranscriptReaderPanel";
import styles from "./TranscriptWorkspacePanel.module.css";

export interface TranscriptWorkspacePanelProps
	extends Omit<HTMLAttributes<HTMLDivElement>, "title"> {
	title?: string;
	subtitle?: string;
	messages: TranscriptMessage[];
	annotations?: TranscriptAnnotation[];
	selectedAnnotationId?: string;
	onAnnotationSelect?: (annotationId: string) => void;
	showNotes?: boolean;
	styleSet?: ContextStyleSet;
}

export function TranscriptWorkspacePanel({
	title,
	subtitle,
	messages,
	annotations = [],
	selectedAnnotationId,
	onAnnotationSelect,
	showNotes = annotations.length > 0,
	styleSet,
	className,
	...rest
}: TranscriptWorkspacePanelProps) {
	const hasNotes = showNotes && annotations.length > 0;
	return (
		<div
			className={[styles.root, hasNotes ? styles.withNotes : styles.noNotes, className ?? ""]
				.filter(Boolean)
				.join(" ")}
			data-rag-organism="TranscriptWorkspacePanel"
			{...rest}
		>
			<TranscriptReaderPanel
				title={title}
				subtitle={subtitle}
				messages={messages}
				annotations={hasNotes ? annotations : []}
				selectedAnnotationId={selectedAnnotationId}
				onAnnotationSelect={onAnnotationSelect}
				showAnnotationChips={hasNotes}
				styleSet={styleSet}
			/>
			{hasNotes && (
				<AnnotationRailPanel
					annotations={annotations}
					selectedAnnotationId={selectedAnnotationId}
					onAnnotationSelect={onAnnotationSelect}
					styleSet={styleSet}
				/>
			)}
		</div>
	);
}
