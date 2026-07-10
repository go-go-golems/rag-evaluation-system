const widget = require("widget.dsl");

const rows = [
	{ id: "course-101", title: "Foundations", owner: "Mira", status: "draft" },
	{ id: "course-202", title: "Context diagrams", owner: "Noah", status: "published" },
];

const schema = widget.data
	.fields("courses", (f) =>
		f.key("id").primary("title").short("owner", { label: "Owner" }).status("status"),
	)
	.build();

const table = widget.data
	.collection("courses", rows, (c) =>
		c
			.schema(schema)
			.select(widget.data.selection.urlParam("course", "course-202"))
			.table((t) =>
				t.rowSelect(
					widget.act.navigate("/courses/${row.id}", {
						payload: { id: widget.bind.context("row.id") },
					}),
				),
			),
	)
	.toNode();

const page = widget.page("Selectable table", (p) =>
	p.section("Courses", (s) =>
		s.caption("Row selection is expressed as serializable ActionSpec data.").view(table),
	),
);
