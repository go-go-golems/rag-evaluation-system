const widget = require("widget.dsl");
const poll = {
	title: "Summary",
	options: [
		{ id: "mon", label: "Mon" },
		{ id: "tue", label: "Tue" },
	],
};
const tallies = [
	{ id: "available", label: "Available", counts: { mon: 3, tue: 5 } },
	{ id: "maybe", label: "Maybe", counts: { mon: 2, tue: 1 } },
];
const page = widget.page("Poll summary", (p) =>
	p.section("Tallies", (s) => s.view(widget.schedule.pollSummary(poll, tallies))),
);
