// Doodle-style scheduling site.
//
// Create a poll with proposed time slots; participants record their availability
// (yes / maybe / no) per slot; a results grid tallies the best slot.
//
// - Data lives in SQLite (the `db` module, file-backed via xgoja.v2.yaml).
// - Pages are Widget IR built with `widget.dsl` v3 from the rag-widget-site
//   provider and rendered by the React RagEvaluationSiteApp SPA.
// - HTTP is served by the planned-route Express API:
//     app.get(pattern).public().handle((ctx, res) => ...)
//   Query:  ctx.request.query.<name>   Body: ctx.body.<name>   Params: ctx.params.<name>
//   Typed form inputs submit via a native <form> POST (application/x-www-form-urlencoded),
//   which the host parses into ctx.body; handlers respond with res.redirect(303, ...).

const { createPages } = require("./lib/pages");
const { createStore } = require("./lib/store");

__package__({ name: "doodle", short: "Doodle-style scheduling site" });

__verb__("site", {
	name: "site",
	output: "text",
	short: "Serve a Doodle-style scheduling site backed by SQLite and widget.dsl v3",
	tags: ["http", "widget", "db", "doodle"],
});

function site() {
	const express = require("express");
	const assets = require("fs:assets");
	const db = require("db");
	const widget = require("widget.dsl");

	const store = createStore(db);
	const pages = createPages({ widget, store });

	const app = express.app();
	app.spaFromAssetsModule("/", assets, "/app/public", {
		excludePrefixes: ["/api", "/healthz", "/favicon.ico"],
	});

	app
		.get("/favicon.ico")
		.public()
		.handle((_ctx, res) => res.status(204).end());
	app
		.get("/healthz")
		.public()
		.handle((_ctx, res) => res.json({ ok: true, site: "doodle-site" }));

	app
		.get("/api/widget/pages/index")
		.public()
		.handle((_ctx, res) => res.json(pages.indexPage()));
	app
		.get("/api/widget/pages/create")
		.public()
		.handle((_ctx, res) => res.json(pages.createPage()));
	app
		.get("/api/widget/pages/poll")
		.public()
		.handle((ctx, res) => {
			const query = ctx.request?.query || {};
			const pollId = Number(query.poll || 0);
			res.json(pages.pollPage(pollId, query));
		});

	app
		.post("/api/form/create-poll")
		.public()
		.handle((ctx, res) => {
			const body = ctx.body || {};
			const title = String(body.title || "").trim();
			const slots = String(body.slots || "")
				.split(/\r?\n/)
				.map((s) => s.trim())
				.filter((s) => s.length > 0);
			if (!title || slots.length === 0) {
				return res.redirect(303, "/pages/create");
			}
			const pollId = store.createPoll(
				title,
				String(body.description || "").trim(),
				String(body.location || "").trim(),
				slots,
			);
			res.redirect(303, `/pages/poll?poll=${pollId}`);
		});

	app
		.post("/api/form/cast-vote")
		.public()
		.handle((ctx, res) => {
			const body = ctx.body || {};
			const pollId = Number(ctx.request?.query?.poll || 0);
			const poll = store.getPoll(pollId);
			const name = String(body.name || "").trim();
			if (!poll || !name) {
				return res.redirect(303, `/pages/poll?poll=${pollId}`);
			}
			store.castVote(pollId, name, body);
			res.redirect(303, `/pages/poll?poll=${pollId}`);
		});
}
