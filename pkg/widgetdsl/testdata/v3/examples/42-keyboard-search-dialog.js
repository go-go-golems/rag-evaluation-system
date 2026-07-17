const widget = require("widget.dsl");

const rows = [
	{ id: "job-1", title: "Go API engineer", status: "shortlisted", starred: true },
	{ id: "job-2", title: "ESP32 firmware", status: "new", starred: false },
];
const fields = widget.data.fields((f) =>
	f
		.key("id", { label: "ID" })
		.primary("title", { label: "Job" })
		.status("status")
		.short("starred", { label: "Starred" }),
);
const navigate = widget.act.navigate("/pages/jobs", {
	query: {
		q: widget.bind.context("query"),
		page: widget.bind.context("page"),
		pageSize: widget.bind.context("pageSize"),
	},
	preserveQuery: ["status", "tag"],
	omitEmpty: true,
});
const jobs = widget.data.collection("jobs", rows, (collection) =>
	collection
		.schema(fields.build())
		.search((search) =>
			search
				.value("go")
				.query("q", { placeholder: "Search jobs" })
				.resultCount(47)
				.submit(navigate)
				.clear(navigate),
		)
		.paginate((pager) => pager.current(2).size(20).total(47).sizes(20, 50, 100).onChange(navigate))
		.table((table) =>
			table
				.keyboard((keys) => keys.mode("rows").selection("manual").vimAliases(true))
				.rowSelect(
					widget.act.navigate("/pages/jobs", {
						query: { job: widget.bind.context("row.id") },
						preserveQuery: ["q", "status", "tag"],
					}),
				)
				.command("star", (command) =>
					command
						.key("s")
						.label("Toggle star")
						.action(
							widget.act.server("job.star", { payload: { id: widget.bind.context("row.id") } }),
						),
				)
				.command("reject", (command) =>
					command
						.key("r")
						.label("Reject")
						.danger()
						.action(
							widget.act.server("job.reject", { payload: { id: widget.bind.context("row.id") } }),
						),
				)
				.command("tag", (command) =>
					command.key("t").label("Add tag").action(widget.act.openOverlay("add-tag")),
				)
				.styleWhen("status", "shortlisted", "success")
				.styleWhen("status", "rejected", "muted"),
		),
);
const addTag = widget.ui.formDialog("add-tag", (dialog) =>
	dialog
		.title(widget.bind.template("Add tag"))
		.initialFocus("tag")
		.body(widget.ui.formRow("Tag", widget.ui.textInput({ name: "tag", required: true })))
		.submit(
			widget.act.server("job.add-tag", {
				payload: { jobId: widget.bind.context("row.id"), tag: widget.bind.context("form.tag") },
			}),
		),
);
const page = widget.page("Keyboard triage", (page) =>
	page.section("Jobs", (section) => section.view(jobs)).view(addTag),
);
