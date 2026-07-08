const widget = require("widget.dsl");

const docs = [
	{
		id: "overview",
		title: "Overview",
		format: "Markdown",
		body: "# Overview\n\nRead this before starting the lab.",
	},
	{
		id: "lab",
		title: "Lab guide",
		format: "Markdown",
		body: "# Lab\n\nFollow the worksheet and capture your findings.",
	},
];

function renderPage(query = {}) {
	const selected = query.item === "lab" ? "lab" : "overview";
	return widget.page("Course handouts", (p) =>
		p.section("Handouts", (s) =>
			s.view(
				widget.course.handouts(
					{
						intro: "Downloadable workshop material",
						docs,
					},
					(h) =>
						h
							.selected(selected)
							.onSelect(widget.act.navigate("?item=${documentId}"))
							.onDownload(widget.course.intent.downloadHandout(widget.bind.context("document.id"))),
				),
			),
		),
	);
}

const page = renderPage({});
