import React, { useState } from "react";
import { Button, TextInput } from "@go-go-golems/rag-evaluation-site";
import { Caption, StatusText } from "@go-go-golems/rag-evaluation-site";
import { Panel, ScrollRegion } from "@go-go-golems/rag-evaluation-site";
import { DataTable, type DataTableColumn } from "@go-go-golems/rag-evaluation-site";
import { CorpusDocumentRow } from "../../../services/api";
import styles from "./DocumentBrowser.module.css";

interface DocumentBrowserProps {
	documents: CorpusDocumentRow[];
	isLoading: boolean;
	sourceName: string | null;
	totalDocs: number;
	selectedDocId: string;
	onSelectDocument: (id: string) => void;
	onLoadMore: () => void;
	hasMore: boolean;
}

export const DocumentBrowser: React.FC<DocumentBrowserProps> = ({
	documents,
	isLoading,
	sourceName,
	totalDocs,
	selectedDocId,
	onSelectDocument,
	onLoadMore,
	hasMore,
}) => {
	const [search, setSearch] = useState("");

	const filtered = search
		? documents.filter(
				(d) =>
					d.title.toLowerCase().includes(search.toLowerCase()) ||
					d.id.toLowerCase().includes(search.toLowerCase()),
			)
		: documents;

	const columns: DataTableColumn<CorpusDocumentRow>[] = [
		{ id: "title", header: "Title", maxWidth: 300, cell: (doc) => doc.title },
		{ id: "words", header: "Words", align: "end", cell: (doc) => doc.word_count.toLocaleString() },
		{ id: "chunks", header: "Chunks", align: "end", cell: (doc) => doc.chunk_count },
		{
			id: "embed",
			header: "Embed",
			align: "end",
			cell: (doc) =>
				doc.chunk_count > 0 ? (
					<Caption
						tone={
							doc.embedded_count === doc.chunk_count
								? "success"
								: doc.embedded_count > 0
									? "warning"
									: "muted"
						}
					>
						{doc.embedded_count}/{doc.chunk_count}
					</Caption>
				) : (
					"—"
				),
		},
		{
			id: "status",
			header: "Status",
			cell: (doc) => (
				<StatusText status={doc.status === "chunked" ? "done" : doc.status}>
					{doc.status}
				</StatusText>
			),
		},
	];

	return (
		<Panel
			className={styles.root}
			title={sourceName ? `${sourceName} — Documents` : "Documents"}
			actions={sourceName && <Caption>{totalDocs} docs</Caption>}
			density="condensed"
		>
			{sourceName && (
				<div className={styles.searchBar}>
					<TextInput
						className={styles.searchInput}
						placeholder="Search title or ID..."
						value={search}
						onChange={(e) => setSearch(e.target.value)}
					/>
				</div>
			)}
			<ScrollRegion axis="y" style={{ maxHeight: 570 }}>
				{!sourceName ? (
					<Caption>Select a source to browse documents.</Caption>
				) : isLoading ? (
					<Caption>Loading...</Caption>
				) : filtered.length === 0 ? (
					<Caption>{search ? "No matching documents." : "No documents found."}</Caption>
				) : (
					<>
						<DataTable
							rows={filtered}
							columns={columns}
							getRowKey={(doc) => doc.id}
							selectedKey={selectedDocId}
							onRowSelect={(doc) => onSelectDocument(doc.id)}
						/>
						{hasMore && (
							<div className={styles.loadMore}>
								<Button onClick={onLoadMore}>
									Load more ({totalDocs - documents.length} remaining)
								</Button>
							</div>
						)}
					</>
				)}
			</ScrollRegion>
		</Panel>
	);
};
