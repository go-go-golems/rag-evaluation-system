__package__({ name: "capabilities", short: "xgoja module capability probes" });

function moduleProbe() {
	const checks = [];
	for (const name of [
		"fs",
		"yaml",
		"db",
		"markdown",
		"sanitize",
		"extract",
		"express",
		"geppetto",
	]) {
		try {
			const mod = require(name);
			checks.push({ name, ok: true, keys: Object.keys(mod).sort().slice(0, 20) });
		} catch (err) {
			checks.push({ name, ok: false, error: String(err && err.message ? err.message : err) });
		}
	}
	return checks;
}

__verb__("moduleProbe", {
	short: "Show which xgoja modules are require-able in this experiment binary",
});

function serveProbe(message) {
	const express = require("express");
	const app = express.app();
	const body = message || "rag-eval-js express probe";
	app.get("/", (_req, res) => {
		res.json({
			ok: true,
			message: body,
			hint: "Run with: ./dist/rag-eval-js run --keep-alive jsverbs/capabilities.js",
		});
	});
	return {
		ok: true,
		route: "/",
		message: body,
		note: "The express module starts the xgoja HTTP server when the runtime is initialized with HTTP enabled. For long-running use prefer a dedicated run script with --keep-alive.",
	};
}

__verb__("serveProbe", {
	short: "Register a tiny Express-style HTTP route as an API smoke test",
	fields: {
		message: {
			type: "string",
			default: "rag-eval-js express probe",
			help: "Message returned by the route",
		},
	},
});
