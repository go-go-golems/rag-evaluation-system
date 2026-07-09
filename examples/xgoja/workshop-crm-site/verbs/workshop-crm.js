// Workshop CRM reference host: SQLite data, widget.dsl v3 pages, and the
// embedded RagEvaluationSite SPA. Routes intentionally use native forms so the
// full write path remains inspectable without a bespoke action dispatcher.
const { createPages } = require("./lib/pages");
const { createStore } = require("./lib/store");

__package__({ name: "workshopcrm", short: "Workshop CRM reference host" });

__verb__("site", {
	name: "site",
	output: "text",
	short: "Serve the SQLite-backed workshop CRM widget.dsl v3 reference host",
	tags: ["http", "widget", "crm", "sqlite"],
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
		.handle((_ctx, res) => res.json({ ok: true, site: "workshop-crm-site" }));

	app
		.get("/api/widget/pages/index")
		.public()
		.handle((_ctx, res) => res.json(pages.indexPage()));
	app
		.get("/api/widget/pages/pipeline")
		.public()
		.handle((_ctx, res) => res.json(pages.pipelinePage()));
	app
		.get("/api/widget/pages/lead")
		.public()
		.handle((_ctx, res) => res.json(pages.leadPage()));
	app
		.get("/api/widget/pages/opportunity")
		.public()
		.handle((ctx, res) => res.json(pages.opportunityPage(Number(ctx.request?.query?.deal || 0))));
	app
		.get("/api/widget/pages/availability")
		.public()
		.handle((ctx, res) => res.json(pages.availabilityPage(Number(ctx.request?.query?.deal || 0))));
	app
		.get("/api/widget/pages/runs")
		.public()
		.handle((_ctx, res) => res.json(pages.runsPage()));

	app
		.post("/api/form/create-lead")
		.public()
		.handle((ctx, res) => {
			const body = ctx.body || {};
			const organization = String(body.organization || "").trim();
			const contact = String(body.contact || "").trim();
			const email = String(body.email || "").trim();
			if (!organization || !contact || !email) return res.redirect(303, "/pages/lead");
			const dealId = store.createLead({
				organization,
				contact,
				email,
				amount: body.amount,
				format: body.format,
			});
			res.redirect(303, `/pages/opportunity?deal=${dealId}`);
		});

	app
		.post("/api/widget/actions/crm.deal.move")
		.public()
		.handle((ctx, res) => {
			const payload = ctx.body?.payload || {};
			const moved = store.moveDeal(Number(payload.dealId || 0), String(payload.toStage || ""));
			res.json({
				ok: moved,
				refresh: moved,
				toast: moved ? "Opportunity moved" : "Unable to move opportunity",
			});
		});

	app
		.post("/api/form/schedule-run")
		.public()
		.handle((ctx, res) => {
			const dealId = Number(ctx.request?.query?.deal || 0);
			const optionId = Number((ctx.body || {}).option || 0);
			if (!store.scheduleRun(dealId, optionId))
				return res.redirect(303, `/pages/availability?deal=${dealId}`);
			res.redirect(303, "/pages/runs");
		});
}
