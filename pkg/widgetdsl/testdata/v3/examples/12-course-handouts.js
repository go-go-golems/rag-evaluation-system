const widget = require("widget.dsl");

const page = widget.page("Course handouts", (p) =>
	p.section("Handouts", (s) =>
		s.view(
			widget.course.handouts(
				{
					intro: "Downloadable workshop material",
					docs: [
						{ id: "deck", title: "Deck", format: "Markdown", body: "# Deck" },
						{ id: "lab", title: "Lab guide", format: "Markdown", body: "# Lab" },
					],
				},
				(h) =>
					h
						.selected("lab")
						.onDownload(widget.course.intent.downloadHandout(widget.bind.context("doc.id"))),
			),
		),
	),
);
