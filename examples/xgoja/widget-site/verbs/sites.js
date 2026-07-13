__package__({ name: "sites", short: "WidgetRenderer xgoja sites" });

__verb__("demo", {
	name: "demo",
	output: "text",
	short: "Serve the native widget.dsl v3 interactive collection demo",
	tags: ["http", "widget", "db", "actions"],
});
function demo() {
	const express = require("express");
	const assets = require("fs:assets");
	const db = require("db");
	const widget = require("widget.dsl");

	db.exec(
		"CREATE TABLE queries (id INTEGER PRIMARY KEY, name TEXT, status TEXT, priority INTEGER, owner TEXT, notes TEXT)",
	);
	[
		["Fast Growing Trees", "succeeded", 3, "botany", "Seeded from the xgoja demo"],
		["Arborvitae Spacing", "running", 2, "landscape", "Needs another pass"],
		["Broken Import Example", "failed", 1, "ingest", "Retry after source cleanup"],
	].forEach((row) => {
		db.exec(
			"INSERT INTO queries (name, status, priority, owner, notes) VALUES (?, ?, ?, ?, ?)",
			...row,
		);
	});

	const app = express.app();
	app.spaFromAssetsModule("/", assets, "/app/public", {
		excludePrefixes: ["/api", "/healthz", "/favicon.ico"],
	});
	app
		.get("/favicon.ico")
		.public()
		.handle((_req, res) => res.status(204).end());
	app
		.get("/healthz")
		.public()
		.handle((_req, res) =>
			res.json({ ok: true, site: "rag-widget-xgoja-site", module: "widget.dsl" }),
		);

	function rows() {
		return db.query(
			"SELECT id, name, status, priority, owner, notes FROM queries ORDER BY priority DESC, id ASC",
		);
	}
	function page() {
		const queryRows = rows();
		const fields = widget.data.fields((f) =>
			f
				.key("id", { label: "ID" })
				.primary("name", { label: "Query" })
				.status("status", { label: "Status" })
				.count("priority", { label: "Priority" })
				.short("owner", { label: "Owner" })
				.prose("notes", { label: "Notes" }),
		);
		const collection = widget.data.collection("queries", queryRows, (c) =>
			c
				.schema(fields.build())
				.empty("No queued queries.")
				.search((search) =>
					search
						.query("q", { placeholder: "Search queued work" })
						.resultCount(queryRows.length)
						.submit(
							widget.act.navigate("/pages/demo", {
								query: { q: widget.bind.context("query") },
								omitEmpty: true,
							}),
						),
				)
				.paginate((pager) =>
					pager
						.current(1)
						.size(20)
						.total(queryRows.length)
						.sizes(20, 50, 100)
						.onChange(
							widget.act.navigate("/pages/demo", {
								query: {
									page: widget.bind.context("page"),
									pageSize: widget.bind.context("pageSize"),
								},
								preserveQuery: ["q"],
								omitEmpty: true,
							}),
						),
				)
				.table((table) =>
					table
						.keyboard((keys) => keys.mode("rows").selection("manual").vimAliases(true))
						.rowSelect(
							widget.act.server("cycle-status", { payload: { id: widget.bind.context("row.id") } }),
						)
						.command("cycle", (command) =>
							command
								.key("s")
								.label("Cycle status")
								.action(
									widget.act.server("cycle-status", {
										payload: { id: widget.bind.context("row.id") },
									}),
								),
						)
						.command("archive", (command) =>
							command
								.key("r")
								.label("Archive")
								.danger()
								.action(
									widget.act.server("archive-query", {
										confirm: "Archive ${row.name}?",
										payload: { id: widget.bind.context("row.id") },
									}),
								),
						)
						.command("note", (command) =>
							command.key("t").label("Edit note").action(widget.act.openOverlay("edit-note")),
						)
						.styleWhen("status", "succeeded", "success")
						.styleWhen("status", "failed", "danger"),
				),
		);
		const noteDialog = widget.ui.formDialog("edit-note", (dialog) =>
			dialog
				.title("Edit query note")
				.initialFocus("notes")
				.body(
					widget.ui.formRow(
						"Notes",
						widget.ui.textareaInput({ name: "notes", required: true, rows: 5 }),
					),
				)
				.submit(
					widget.act.server("save-note", {
						payload: {
							id: widget.bind.context("row.id"),
							notes: widget.bind.context("form.notes"),
						},
					}),
				),
		);
		return widget
			.page("Widget DSL v3 work queue", (p) =>
				p
					.section("Queue", (section) =>
						section
							.caption("Arrow keys move; Enter selects; S cycles; R archives; T edits notes.")
							.view(collection.toNode()),
					)
					.view(noteDialog),
			)
			.toPage();
	}

	const pageIds = [
		"index",
		"demo",
		"actions",
		"semantic",
		"upload",
		"transcripts",
		"slides",
		"handouts",
		"course-examples",
		"context",
		"course",
		"handout",
	];
	pageIds.forEach((id) => {
		app
			.get(`/api/widget/pages/${id}`)
			.public()
			.handle((_req, res) => res.json(page()));
	});

	function payload(req) {
		return (req.body && req.body.payload) || {};
	}
	function nextStatus(status) {
		return status === "pending"
			? "running"
			: status === "running"
				? "succeeded"
				: status === "succeeded"
					? "failed"
					: "pending";
	}
	app
		.post("/api/widget/actions/cycle-status")
		.public()
		.handle((req, res) => {
			const id = Number(payload(req).id);
			const found = db.query("SELECT status FROM queries WHERE id = ?", id);
			if (!found.length) return res.status(404).json({ ok: false, error: "Query not found" });
			const status = nextStatus(found[0].status);
			db.exec("UPDATE queries SET status = ? WHERE id = ?", status, id);
			return res.json({ ok: true, refresh: true, toast: `Query #${id} → ${status}` });
		});
	app
		.post("/api/widget/actions/archive-query")
		.public()
		.handle((req, res) => {
			const id = Number(payload(req).id);
			db.exec("DELETE FROM queries WHERE id = ?", id);
			res.json({ ok: true, refresh: true, toast: `Archived query #${id}` });
		});
	app
		.post("/api/widget/actions/save-note")
		.public()
		.handle((req, res) => {
			const data = payload(req);
			if (!String(data.notes || "").trim())
				return res.status(422).json({
					ok: false,
					error: "A note is required.",
					fieldErrors: { notes: "Enter a note." },
				});
			db.exec("UPDATE queries SET notes = ? WHERE id = ?", String(data.notes), Number(data.id));
			res.json({ ok: true, refresh: true, toast: "Note saved" });
		});
}
