const widget = require("widget.dsl");

const transcript = {
	title: "Debug session",
	subtitle: "Retrieval failure triage",
	messages: [
		{
			id: "m1",
			role: "user",
			text: "Why did retrieval fail?",
			tokens: 42,
		},
		{
			id: "m2",
			role: "assistant",
			text: "The corpus filter was too narrow, so the retriever excluded the documents that contained the answer.",
			tokens: 96,
			annotationIds: ["a1"],
		},
	],
	annotations: [
		{
			id: "a1",
			targetMessageId: "m2",
			styleKey: "assistant",
			label: "Root cause",
			text: "The filter was scoped to one corpus instead of the full evaluation fixture set.",
			confidence: 0.91,
		},
	],
};

const page = widget.page("Transcript workspace", (p) =>
	p.section("Workspace", (s) =>
		s.view(
			widget.context.workspace(transcript, (w) =>
				w
					.selectedAnnotation("a1")
					.onAnnotationSelect(
						widget.context.intent.selectAnnotation(widget.bind.context("annotation.id")),
					),
			),
		),
	),
);
