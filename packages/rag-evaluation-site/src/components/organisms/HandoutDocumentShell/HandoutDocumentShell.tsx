import type { HTMLAttributes, ReactNode } from "react";
import type { ContextHandoutDocument, ContextStyleSet } from "../../../context";
import { Caption } from "../../foundation";
import { DocumentListPanel, DocumentPreviewToolbar, MarkdownArticle } from "../../molecules";
import { RichArticle } from "../RichArticle";
import styles from "./HandoutDocumentShell.module.css";

export interface HandoutDocumentShellProps
	extends Omit<HTMLAttributes<HTMLDivElement>, "onSelect" | "title"> {
	intro: ReactNode;
	documents: ContextHandoutDocument[];
	selectedDocumentId?: string;
	onDocumentSelect?: (documentId: string) => void;
	onDownload?: (documentId: string) => void;
	onDownloadAll?: () => void;
	onPrint?: (documentId: string) => void;
	styleSet?: ContextStyleSet;
	title?: ReactNode;
	emptyMessage?: ReactNode;
}

function selectedDocument(documents: ContextHandoutDocument[], selectedDocumentId?: string) {
	return documents.find((document) => document.id === selectedDocumentId) ?? documents[0];
}

export function HandoutDocumentShell({
	intro,
	documents,
	selectedDocumentId,
	onDocumentSelect,
	onDownload,
	onDownloadAll,
	onPrint,
	styleSet,
	title = "Handout",
	emptyMessage = "Select a document to preview it.",
	className,
	...rest
}: HandoutDocumentShellProps) {
	const activeDocument = selectedDocument(documents, selectedDocumentId);
	return (
		<div
			className={[styles.root, className ?? ""].filter(Boolean).join(" ")}
			data-rag-organism="HandoutDocumentShell"
			{...rest}
		>
			<aside className={styles.sidebar}>
				<DocumentListPanel
					title={title}
					description={intro}
					items={documents.map((document) => ({
						id: document.id,
						title: document.title,
						format: document.format,
						size: document.size,
						description: document.description,
					}))}
					selectedItemId={activeDocument?.id}
					onItemSelect={onDocumentSelect}
					onDownloadAll={onDownloadAll}
				/>
			</aside>
			<section className={styles.preview}>
				{activeDocument ? (
					<>
						<DocumentPreviewToolbar
							file={activeDocument.file}
							format={activeDocument.format}
							size={activeDocument.size}
							onDownload={onDownload ? () => onDownload(activeDocument.id) : undefined}
							downloadLabel="Download Markdown"
							onPrint={onPrint ? () => onPrint(activeDocument.id) : undefined}
						/>
						<div className={styles.articleWrap}>
							{activeDocument.blocks && activeDocument.blocks.length > 0 ? (
								<RichArticle blocks={activeDocument.blocks} styleSet={styleSet} />
							) : activeDocument.format.toLowerCase().includes("markdown") ? (
								<MarkdownArticle source={activeDocument.body} />
							) : (
								<div className={styles.placeholder}>
									<Caption>{activeDocument.description}</Caption>
								</div>
							)}
						</div>
					</>
				) : (
					<div className={styles.placeholder}>
						<Caption>{emptyMessage}</Caption>
					</div>
				)}
			</section>
		</div>
	);
}
