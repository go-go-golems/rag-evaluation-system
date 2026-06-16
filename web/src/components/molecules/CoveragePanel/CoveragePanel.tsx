import { Caption, StatusText } from "@go-go-golems/rag-evaluation-site";
import { Panel } from "@go-go-golems/rag-evaluation-site";
import { DataTable, MetadataGrid, type DataTableColumn } from "@go-go-golems/rag-evaluation-site";
import type { EmbeddingCoverageResult } from "../../../services/api";
import styles from "./CoveragePanel.module.css";

export interface CoveragePanelProps {
	coverage: EmbeddingCoverageResult;
	coveragePct: number;
}

type SourceCoverage = NonNullable<EmbeddingCoverageResult["sources"]>[number];

export function CoveragePanel({ coverage, coveragePct }: CoveragePanelProps) {
	const isSparse = coveragePct < 50;
	const columns: DataTableColumn<SourceCoverage>[] = [
		{ id: "source", header: "Source", maxWidth: 120, cell: (source) => source.source_name },
		{
			id: "coverage",
			header: "Cov",
			align: "end",
			cell: (source) => (
				<Caption tone={source.coverage_pct >= 50 ? "success" : "warning"}>
					{Math.round(source.coverage_pct)}%
				</Caption>
			),
		},
		{ id: "embedded", header: "Emb", align: "end", cell: (source) => source.embedded_chunks },
		{ id: "total", header: "Tot", align: "end", cell: (source) => source.total_chunks },
	];

	return (
		<Panel
			title="Coverage"
			actions={<StatusText status={isSparse ? "partial" : "done"}>{coveragePct}%</StatusText>}
			density="condensed"
			data-rag-component="CoveragePanel"
		>
			<MetadataGrid
				density="compact"
				items={[
					{ key: "Embedded", value: <Caption tone="success">{coverage.embedded_chunks}</Caption> },
					{ key: "Total", value: coverage.total_chunks },
					{
						key: "Missing",
						value: (
							<Caption tone="danger">{coverage.total_chunks - coverage.embedded_chunks}</Caption>
						),
					},
				]}
			/>
			{coverage.sources && coverage.sources.length > 0 && (
				<DataTable
					rows={coverage.sources}
					columns={columns}
					getRowKey={(source) => source.source_id}
				/>
			)}
			{isSparse && (
				<div className={styles.warning}>
					⚠ Vector search only compares embedded chunks. Coverage is sparse — results validate
					behavior but not full-corpus quality.
				</div>
			)}
		</Panel>
	);
}
