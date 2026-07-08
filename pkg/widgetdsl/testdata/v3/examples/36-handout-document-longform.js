const widget = require("widget.dsl");
const docs = [
	{
		id: "guide",
		title: "Guide",
		format: "Markdown",
		body: `# Guide

This is a longer handout body with **markdown** and checklists.

- Read
- Practice
- Review`,
	},
];
const page = widget.page("Longform handout", (p) =>
	p.section("Handout", (s) => s.view(widget.course.handouts({ docs }, (h) => h.selected("guide")))),
);
