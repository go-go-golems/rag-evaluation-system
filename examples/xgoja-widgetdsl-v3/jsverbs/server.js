const express = require("express");
const assets = require("fs:assets");
const widget = require("widget.dsl");

const app = express.app();

const exampleIds = assets
	.readdirSync("/examples")
	.filter((name) => name.endsWith(".js"))
	.sort()
	.map((name) => name.replace(/\.js$/, ""));

function titleFromId(id) {
	const parts = id.replace(/^\d+-/, "").split("-");
	return parts
		.map((part) => (part === "cms" ? "CMS" : part.slice(0, 1).toUpperCase() + part.slice(1)))
		.join(" ");
}

function navItems() {
	return [{ id: "index", label: "Index" }].concat(
		exampleIds.map((id) => ({ id, label: titleFromId(id) })),
	);
}

function attachShell(page, active) {
	if (!page.shell) {
		page.shell = widget.app.shell((shell) =>
			shell
				.brand(page.title || "widget.dsl v3 examples")
				.navigation((navigation) =>
					navigation
						.placement("top")
						.active(active)
						.ariaLabel("Examples")
						.section("examples", "Examples", (items) => {
							navItems().forEach((item) =>
								items.item(item.id, item.label, widget.act.navigate(`/pages/${item.id}`)),
							);
						}),
				)
				.content((content) => content.maxWidth("wide").padding("none")),
		);
	}
	return page;
}

function indexPage() {
	return widget
		.page("widget.dsl v3 examples", (p) =>
			p.section("Examples", (s) => {
				s.caption(
					"These pages are rendered by xgoja from the committed widget.dsl v3 example scripts.",
				);
				for (const id of exampleIds) {
					s.view(
						widget.ui.button(titleFromId(id), widget.act.navigate(`/pages/${id}`), {
							variant: "ghost",
						}),
					);
				}
			}),
		)
		.toPage();
}

function renderExample(id, query) {
	if (!exampleIds.includes(id)) return null;
	const source = assets.readFileSync(`/examples/${id}.js`, "utf8");
	const fn = new Function(
		"require",
		"query",
		`${source}\nconst rendered = typeof renderPage === "function" ? renderPage(query || {}) : page;\nreturn rendered && typeof rendered.toPage === "function" ? rendered.toPage() : rendered;`,
	);
	return fn(require, query || {});
}

app
	.get("/healthz")
	.public()
	.handle((_ctx, res) => res.json({ ok: true, service: "widgetdsl-v3-examples" }));

app
	.get("/api/widget/pages/:id")
	.public()
	.handle((ctx, res) => {
		const id = ctx.params.id || "index";
		const query = (ctx.request && ctx.request.query) || {};
		const page = id === "index" ? indexPage() : renderExample(id, query);
		if (!page) {
			res.status(404).json({ error: { code: "page_not_found", message: `Unknown page ${id}` } });
			return;
		}
		res.json(attachShell(page, id));
	});

app
	.post("/api/widget/actions/:name")
	.public()
	.handle((ctx, res) => {
		res.json({ refresh: false, toast: `Preview action accepted: ${ctx.params.name}` });
	});

app.spaFromAssetsModule("/", assets, "/app", {
	excludePrefixes: ["/api", "/healthz", "/favicon.ico"],
});
