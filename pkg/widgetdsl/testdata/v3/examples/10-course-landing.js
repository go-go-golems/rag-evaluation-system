const widget = require("widget.dsl");

const course = {
	title: "Context Engineering",
	subtitle: "Four lessons",
	tagline: "Learn how to budget, visualize, and reclaim context in production RAG systems.",
	when: "July cohort",
	where: "Remote workshop",
	format: "Live lab",
	outcomes: [
		"Model context as budgeted segments",
		"Pick a visual representation for long prompts",
		"Design reclaim policies for overloaded windows",
	],
	agenda: [
		{ id: "intro", number: "01", title: "Intro", duration: "20 min" },
		{ id: "lab", number: "02", title: "Lab", duration: "60 min" },
	],
};

const page = widget.page("Course landing", (p) =>
	p.section("Landing", (s) =>
		s.view(
			widget.course.landing(course, (l) =>
				l
					.activeAgenda("lab")
					.onAgendaSelect(widget.course.intent.editAgenda(widget.bind.context("agenda.id"))),
			),
		),
	),
);
