const widget = require("widget.dsl");
const cards = [
	{ id: "d1", title: "Acme", owner: "Mira", stage: "lead", value: "$12k" },
	{ id: "d2", title: "Globex", owner: "Noah", stage: "won", value: "$25k" },
];
const pipeline = widget.crm.pipeline("Deals", (p) => p.stage("lead", "Lead").stage("won", "Won"));
const page = widget.page("CRM board", (p) =>
	p.section("Pipeline", (s) =>
		s.view(widget.crm.pipelineBoard(pipeline, cards, (board) => board.ariaLabel("Deals"))),
	),
);
