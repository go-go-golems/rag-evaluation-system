const widget = require("widget.dsl");

const rows = [
	{ id: "q1", label: "Query 1", scores: { precision: 0.8, recall: 0.7 } },
	{ id: "q2", label: "Query 2", scores: { precision: 0.6, recall: 0.9 } },
];

const matrix = widget.data
	.matrix(rows, (m) =>
		m
			.column("precision", "Precision")
			.column("recall", "Recall")
			.valueAt(widget.bind.map("scores"))
			.cell({ kind: "value" }),
	)
	.toNode();

const page = widget.page("Matrix heatmap", (p) =>
	p.section("Scores", (s) =>
		s.caption("Resolved matrix cells use the scores map for each query row.").view(matrix),
	),
);
