const widget = require("widget.dsl");
const cards = [
	{ id: "d1", title: "Acme", owner: "Mira", stage: "lead", value: "$12k" },
	{ id: "d2", title: "Globex", owner: "Noah", stage: "won", value: "$25k" },
];
const props = {
	columns: [
		{ id: "lead", header: "Lead" },
		{ id: "won", header: "Won" },
	],
	cards,
	columnField: "stage",
	getCardId: "id",
	card: {
		title: { kind: "field", field: "title" },
		subtitle: { kind: "field", field: "owner" },
		meta: { kind: "field", field: "value" },
	},
	ariaLabel: "Deals",
};
const page = widget.page("CRM board", (p) =>
	p.section("Pipeline", (s) => s.view(widget.raw.component("BoardEngine", props))),
);
