const widget = require("widget.dsl");

const handouts = widget.course.handouts(
	{
		intro: "Printable materials for the workshop.",
		docs: [
			{ id: "deck", title: "Slide deck", format: "PDF" },
			{ id: "lab", title: "Lab guide", format: "Markdown" },
		],
	},
	(h) =>
		h
			.selected("deck")
			.onSelect(widget.course.intent.selectHandout(widget.bind.context("doc.id")))
			.onDownload(widget.course.intent.downloadHandout(widget.bind.context("doc.id")))
			.onPrint(widget.course.intent.printHandout(widget.bind.context("doc.id"))),
);

const slide = widget.course.slideDeck(
	{
		index: 0,
		slides: [
			{
				id: "s1",
				title: "Context windows",
				body: "Think in segments, budgets, and reclaim policies.",
			},
		],
		snapshot: { id: "ctx", title: "Slide context", limit: 2000, parts: [] },
	},
	(d) =>
		d
			.mode("speaker")
			.onPrevious(widget.course.intent.previousSlide())
			.onNext(widget.course.intent.nextSlide())
			.onPresent(widget.course.intent.presentSlide("s1")),
);

const page = widget.page("Handouts and slide", (p) =>
	p.section("Handouts", (s) => s.view(handouts)).section("Slide", (s) => s.view(slide)),
);
