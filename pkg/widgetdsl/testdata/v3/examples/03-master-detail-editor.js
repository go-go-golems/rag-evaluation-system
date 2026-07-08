const widget = require("widget.dsl");

const rows = [
	{ id: "agenda-1", title: "Welcome", duration: 15, status: "locked" },
	{ id: "agenda-2", title: "Hands-on lab", duration: 45, status: "editable" },
];

const schema = widget.data
	.fields("agenda", (f) =>
		f.key("id").primary("title").count("duration", { label: "Minutes" }).status("status"),
	)
	.build();

const editor = widget.data
	.collection("agenda", rows, (c) =>
		c
			.schema(schema)
			.select(widget.data.selection("single", { keyField: "id", selected: "agenda-2" }))
			.masterDetail()
			.edit((e) =>
				e
					.submitPost("/api/course/agenda")
					.reorder(widget.act.server("course.agenda.reorder"))
					.remove(
						widget.act.server("course.agenda.remove", { confirm: "Remove this agenda item?" }),
					),
			),
	)
	.toNode();

const page = widget.page("Master-detail editor", (p) =>
	p.section("Agenda editor", (s) =>
		s.caption("A master-detail collection with edit actions.").view(editor),
	),
);
