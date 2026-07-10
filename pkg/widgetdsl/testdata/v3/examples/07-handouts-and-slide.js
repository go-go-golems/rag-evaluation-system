const widget = require("widget.dsl");

const docs = [
	{
		id: "deck",
		title: "Slide deck",
		format: "Markdown",
		body: "# Slide deck\n\nUse this deck to explain context-window pressure and reclaim policies.",
	},
	{
		id: "lab",
		title: "Lab guide",
		format: "Markdown",
		body: "# Lab guide\n\n1. Inspect the prompt parts.\n2. Pick a reclaim policy.\n3. Validate with the token diagram.",
	},
];

const slides = [
	{
		id: "s1",
		title: "Context windows",
		body: "Think in segments, budgets, and reclaim policies.",
		notes: ["Introduce token limits", "Connect the diagram to the handout"],
	},
	{
		id: "s2",
		title: "Lab workflow",
		body: "Move from inspection to an explicit reclaim decision.",
		notes: ["Open the lab guide", "Walk through the three validation steps"],
	},
];

function clampSlide(value) {
	const parsed = Number(value || 0);
	if (!Number.isFinite(parsed)) return 0;
	return Math.max(0, Math.min(slides.length - 1, Math.trunc(parsed)));
}

function renderPage(query = {}) {
	const selectedDoc = query.doc === "lab" ? "lab" : "deck";
	const slideIndex = clampSlide(query.slide);
	const previousSlide = Math.max(0, slideIndex - 1);
	const nextSlide = Math.min(slides.length - 1, slideIndex + 1);

	const handouts = widget.course.handouts(
		{
			intro: "Printable materials for the workshop.",
			docs,
		},
		(h) =>
			h
				.selected(selectedDoc)
				.onSelect(widget.act.navigate(`?doc=${"${documentId}"}&slide=${slideIndex}`))
				.onDownload(widget.course.intent.downloadHandout(widget.bind.context("documentId")))
				.onPrint(widget.course.intent.printHandout(widget.bind.context("documentId"))),
	);

	const slide = widget.course.slideDeck(
		{
			index: slideIndex,
			slides,
			snapshot: { id: "ctx", title: "Slide context", limit: 2000, parts: [] },
		},
		(d) =>
			d
				.mode("speaker")
				.onPrevious(widget.act.navigate(`?doc=${selectedDoc}&slide=${previousSlide}`))
				.onNext(widget.act.navigate(`?doc=${selectedDoc}&slide=${nextSlide}`))
				.onPresent(widget.course.intent.presentSlide(slides[slideIndex].id)),
	);

	return widget.page("Handouts and slide", (p) =>
		p.section("Handouts", (s) => s.view(handouts)).section("Slide", (s) => s.view(slide)),
	);
}

const page = renderPage({});
