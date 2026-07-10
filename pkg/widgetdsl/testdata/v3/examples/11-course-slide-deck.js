const widget = require("widget.dsl");

const deck = {
	index: 1,
	slides: [
		{ id: "s1", title: "Budget", notes: ["Define the context envelope", "Show budget pressure"] },
		{ id: "s2", title: "Reclaim", notes: ["Expire stale turns", "Keep cited evidence pinned"] },
	],
	snapshot: { id: "ctx", title: "Slide context", limit: 4096, parts: [] },
};

const page = widget.page("Slide deck", (p) =>
	p.section("Slides", (s) =>
		s.view(
			widget.course.slideDeck(deck, (d) =>
				d
					.mode("speaker")
					.onPrevious(widget.course.intent.previousSlide())
					.onNext(widget.course.intent.nextSlide())
					.onPresent(widget.course.intent.presentSlide()),
			),
		),
	),
);
