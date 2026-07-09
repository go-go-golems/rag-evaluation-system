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

	// ---------------------------------------------------------------- schema + seed
	db.exec(
		"CREATE TABLE IF NOT EXISTS polls (id INTEGER PRIMARY KEY AUTOINCREMENT, slug TEXT UNIQUE, title TEXT NOT NULL, description TEXT, location TEXT, created_at TEXT NOT NULL)",
	);
	db.exec(
		"CREATE TABLE IF NOT EXISTS options (id INTEGER PRIMARY KEY AUTOINCREMENT, poll_id INTEGER NOT NULL, label TEXT NOT NULL, sort INTEGER NOT NULL)",
	);
	db.exec(
		"CREATE TABLE IF NOT EXISTS participants (id INTEGER PRIMARY KEY AUTOINCREMENT, poll_id INTEGER NOT NULL, name TEXT NOT NULL, created_at TEXT NOT NULL)",
	);
	db.exec(
		"CREATE TABLE IF NOT EXISTS votes (participant_id INTEGER NOT NULL, option_id INTEGER NOT NULL, value TEXT NOT NULL, PRIMARY KEY (participant_id, option_id))",
	);

	function nowISO() {
		return new Date().toISOString();
	}

	function slugify(title) {
		const base = String(title || "poll")
			.toLowerCase()
			.replace(/[^a-z0-9]+/g, "-")
			.replace(/(^-+|-+$)/g, "")
			.slice(0, 40);
		return base || "poll";
	}

	function createPoll(title, description, location, slotLabels) {
		db.exec(
			"INSERT INTO polls (slug, title, description, location, created_at) VALUES (?, ?, ?, ?, ?)",
			null,
			title,
			description || "",
			location || "",
			nowISO(),
		);
		const pollId = Number(db.query("SELECT last_insert_rowid() AS id")[0].id);
		db.exec("UPDATE polls SET slug = ? WHERE id = ?", `${slugify(title)}-${pollId}`, pollId);
		slotLabels.forEach((label, index) => {
			db.exec("INSERT INTO options (poll_id, label, sort) VALUES (?, ?, ?)", pollId, label, index);
		});
		return pollId;
	}

	if (Number(db.query("SELECT COUNT(*) AS n FROM polls")[0].n) === 0) {
		const seedId = createPoll(
			"Team offsite dinner",
			"Pick the evening that works best for the crew dinner.",
			"Trattoria Luca, downtown",
			["Thu Jul 9 · 19:00", "Fri Jul 10 · 19:00", "Sat Jul 11 · 18:30", "Sat Jul 11 · 20:00"],
		);
		const opts = db.query("SELECT id FROM options WHERE poll_id = ? ORDER BY sort", seedId);
		function seedParticipant(name, values) {
			db.exec(
				"INSERT INTO participants (poll_id, name, created_at) VALUES (?, ?, ?)",
				seedId,
				name,
				nowISO(),
			);
			const pid = Number(db.query("SELECT last_insert_rowid() AS id")[0].id);
			values.forEach((value, i) => {
				db.exec(
					"INSERT INTO votes (participant_id, option_id, value) VALUES (?, ?, ?)",
					pid,
					opts[i].id,
					value,
				);
			});
		}
		seedParticipant("Ada", ["yes", "no", "yes", "maybe"]);
		seedParticipant("Grace", ["maybe", "no", "yes", "yes"]);
		seedParticipant("Linus", ["no", "yes", "yes", "no"]);
	}

	// ---------------------------------------------------------------- data helpers
	function allPolls() {
		return db.query(
			"SELECT p.id, p.slug, p.title, p.location, p.created_at, " +
				"(SELECT COUNT(*) FROM options o WHERE o.poll_id = p.id) AS slots, " +
				"(SELECT COUNT(*) FROM participants pt WHERE pt.poll_id = p.id) AS people " +
				"FROM polls p ORDER BY p.id DESC",
		);
	}

	function getPoll(pollId) {
		const rows = db.query(
			"SELECT id, slug, title, description, location, created_at FROM polls WHERE id = ?",
			pollId,
		);
		return rows.length ? rows[0] : null;
	}

	function pollOptions(pollId) {
		return db.query(
			"SELECT id, label, sort FROM options WHERE poll_id = ? ORDER BY sort, id",
			pollId,
		);
	}

	function pollParticipants(pollId) {
		return db.query("SELECT id, name FROM participants WHERE poll_id = ? ORDER BY id", pollId);
	}

	function pollVotes(pollId) {
		return db.query(
			"SELECT v.participant_id, v.option_id, v.value FROM votes v " +
				"JOIN participants pt ON pt.id = v.participant_id WHERE pt.poll_id = ?",
			pollId,
		);
	}

	// ---------------------------------------------------------------- widget.dsl helpers
	const act = widget.act;
	const navItems = [
		{ id: "index", label: "All polls", action: act.navigate("/pages/index") },
		{ id: "create", label: "New poll", action: act.navigate("/pages/create") },
	];
	function statusText(status, text) {
		return widget.ui.status(status, text);
	}

	function emptyState(title, description) {
		return widget.ui.emptyState(title, description);
	}

	function formRow(label, control, options = {}) {
		return widget.ui.formRow(label, control, options);
	}

	function textInput(props) {
		return widget.ui.textInput(props);
	}

	function textareaInput(props) {
		return widget.ui.textareaInput(props);
	}

	function selectInput(props) {
		return widget.ui.selectInput(props);
	}

	function collectionTable(name, rows, configureFields, options = {}) {
		const schema = widget.data.fields(name, configureFields).build();
		return widget.data
			.collection(name, rows, (c) => {
				c.schema(schema)
					.empty(options.empty || "No rows")
					.table();
			})
			.toNode();
	}

	function applyPageMeta(p, id, activeNavItemId) {
		return p
			.id(id)
			.meta("activeNavItemId", activeNavItemId)
			.meta("navItems", navItems)
			.meta("maxWidth", "wide");
	}

	function asPage(pageBuilder) {
		return pageBuilder.toPage();
	}

	function indexPage() {
		const polls = allPolls();
		const totalPeople = polls.reduce((acc, p) => acc + Number(p.people || 0), 0);
		const rows = polls.map((p) => ({
			id: p.id,
			title: p.title,
			location: p.location || "—",
			slots: Number(p.slots || 0),
			people: Number(p.people || 0),
			created: String(p.created_at || "").slice(0, 10),
		}));

		const table = collectionTable(
			"polls",
			rows,
			(f) =>
				f
					.key("id", { label: "ID" })
					.primary("title", { label: "Poll" })
					.short("location", { label: "Where" })
					.count("slots", { label: "Slots" })
					.count("people", { label: "Responses" })
					.date("created", { label: "Created" }),
			{ empty: "No polls yet" },
		);

		return asPage(
			widget.page("Doodle · scheduling polls", (p) => {
				applyPageMeta(p, "index", "index")
					.section("Scheduling polls", (s) =>
						s
							.view(statusText("succeeded", `${polls.length} active poll(s)`))
							.caption(
								"Create a poll with a few time slots, share it, and let people mark their availability. Data is stored in SQLite.",
							)
							.view(
								widget.ui.button("+ New poll", act.navigate("/pages/create"), {
									variant: "primary",
								}),
							),
					)
					.section("Metrics", (s) =>
						s
							.metric("Polls", String(polls.length), { status: "ready" })
							.metric("Total responses", String(totalPeople), { status: "succeeded" })
							.metric(
								"Time slots",
								String(polls.reduce((acc, poll) => acc + Number(poll.slots || 0), 0)),
								{ status: "running" },
							),
					)
					.section("Polls", (s) =>
						s.view(
							polls.length === 0
								? emptyState("No polls yet", "Create your first scheduling poll.")
								: table,
						),
					)
					.section("Open a poll", (s) =>
						s.view(
							polls.length === 0
								? widget.ui.caption("No poll links yet.", { tone: "muted" })
								: widget.ui.inline(
										{ gap: "sm", wrap: true },
										...rows.map((r) =>
											widget.ui.button(`${r.title} →`, act.navigate(`/pages/poll?poll=${r.id}`), {
												variant: "secondary",
											}),
										),
									),
						),
					);
			}),
		);
	}

	function createPage() {
		const form = widget.ui.form(
			{
				title: "Event details",
				method: "post",
				formAction: "/api/form/create-poll",
				submitLabel: "Create poll",
			},
			formRow(
				"Title",
				textInput({ name: "title", placeholder: "Team offsite dinner", required: true }),
				{ required: true },
			),
			formRow(
				"Description",
				textareaInput({
					name: "description",
					placeholder: "Optional context for invitees",
					rows: 2,
				}),
			),
			formRow("Location", textInput({ name: "location", placeholder: "Trattoria Luca, downtown" })),
			formRow(
				"Time slots (one per line)",
				textareaInput({
					name: "slots",
					placeholder: "Thu Jul 9 · 19:00\nFri Jul 10 · 19:00\nSat Jul 11 · 18:30",
					rows: 5,
					required: true,
				}),
				{ required: true },
			),
		);

		return asPage(
			widget.page("New scheduling poll", (p) => {
				applyPageMeta(p, "create", "create")
					.section("Create a poll", (s) =>
						s.caption(
							"Give the event a title and list one time slot per line. Everything is stored in SQLite.",
						),
					)
					.section("Event details", (s) => s.view(form))
					.section("Navigation", (s) =>
						s.view(widget.ui.button("← Back to all polls", act.navigate("/pages/index"))),
					);
			}),
		);
	}

	function pollPage(pollId) {
		const poll = getPoll(pollId);
		if (!poll) {
			return asPage(
				widget.page("Poll not found", (p) => {
					applyPageMeta(p, "poll", "index").section("Poll not found", (s) =>
						s.view(statusText("failed", `No poll with id ${pollId}`)).view(
							widget.ui.button("← All polls", act.navigate("/pages/index"), {
								variant: "primary",
							}),
						),
					);
				}),
			);
		}

		const options = pollOptions(pollId);
		const participants = pollParticipants(pollId);
		const votes = pollVotes(pollId);

		const voteMap = {};
		votes.forEach((v) => {
			const pid = v.participant_id;
			if (!voteMap[pid]) voteMap[pid] = {};
			voteMap[pid][v.option_id] = v.value;
		});

		function availabilityValue(value) {
			if (value === "yes") return "available";
			if (value === "no") return "unavailable";
			if (value === "maybe") return "maybe";
			return "unknown";
		}

		const availabilityPoll = {
			title: poll.title,
			options: options.map((opt) => ({ id: String(opt.id), label: opt.label })),
			responses: participants.map((pt) => {
				const availability = {};
				options.forEach((opt) => {
					availability[String(opt.id)] = availabilityValue(
						(voteMap[pt.id] && voteMap[pt.id][opt.id]) || "",
					);
				});
				return { id: String(pt.id), name: pt.name, availability };
			}),
		};
		const availabilityGrid = widget.schedule.availabilityPoll(availabilityPoll, (b) =>
			b.readOnly(),
		);

		const tally = options.map((opt) => {
			let yes = 0;
			let maybe = 0;
			let no = 0;
			participants.forEach((pt) => {
				const value = voteMap[pt.id] && voteMap[pt.id][opt.id];
				if (value === "yes") yes += 1;
				else if (value === "maybe") maybe += 1;
				else if (value === "no") no += 1;
			});
			return { option: opt, yes, maybe, no, score: yes * 2 + maybe };
		});
		let bestScore = -1;
		tally.forEach((t) => {
			if (t.score > bestScore) bestScore = t.score;
		});
		const resultRows = tally.map((t) => ({
			id: t.option.id,
			label: t.option.label,
			yes: t.yes,
			maybe: t.maybe,
			no: t.no,
			score: t.score,
			verdict: t.score === bestScore && bestScore > 0 ? "succeeded" : "pending",
		}));
		const resultsTable = collectionTable("results", resultRows, (f) =>
			f
				.key("id", { label: "ID" })
				.primary("label", { label: "Time slot" })
				.count("score", { label: "Score" })
				.status("verdict", { label: "Best" }),
		);
		const summaryTallies = [
			{
				id: "available",
				label: "Available",
				counts: Object.fromEntries(tally.map((t) => [String(t.option.id), t.yes])),
			},
			{
				id: "maybe",
				label: "Maybe",
				counts: Object.fromEntries(tally.map((t) => [String(t.option.id), t.maybe])),
			},
			{
				id: "unavailable",
				label: "Unavailable",
				counts: Object.fromEntries(tally.map((t) => [String(t.option.id), t.no])),
			},
		];
		const summaryGrid = widget.schedule.pollSummary(availabilityPoll, summaryTallies);

		const formChildren = [
			formRow(
				"Your name",
				textInput({ name: "name", placeholder: "e.g. Katherine", required: true }),
				{
					required: true,
				},
			),
		];
		options.forEach((opt) => {
			formChildren.push(
				formRow(
					opt.label,
					selectInput({
						name: `opt_${opt.id}`,
						defaultValue: "maybe",
						options: [
							{ value: "yes", label: "✓ Yes" },
							{ value: "maybe", label: "~ Maybe" },
							{ value: "no", label: "✗ No" },
						],
					}),
				),
			);
		});
		const availabilityForm = widget.ui.form(
			{
				title: "Add your availability",
				method: "post",
				formAction: `/api/form/cast-vote?poll=${poll.id}`,
				submitLabel: "Submit availability",
			},
			...formChildren,
		);

		return asPage(
			widget.page(poll.title, (p) => {
				applyPageMeta(p, "poll", "index")
					.section(poll.title, (s) =>
						s
							.metadata({
								Location: poll.location || "—",
								Responses: String(participants.length),
								Slots: String(options.length),
								"Share link": `/pages/poll?poll=${poll.id}`,
							})
							.caption(poll.description || "")
							.view(widget.ui.button("← All polls", act.navigate("/pages/index"))),
					)
					.section("Availability grid", (s) =>
						s.view(
							participants.length === 0
								? emptyState("No responses yet", "Be the first to add your availability below.")
								: availabilityGrid,
						),
					)
					.section("Results by slot", (s) => s.view(summaryGrid).view(resultsTable))
					.section("Add your availability", (s) => s.view(availabilityForm));
			}),
		);
	}

	// ---------------------------------------------------------------- HTTP wiring
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
		.handle((_ctx, res) => res.json(indexPage()));
	app
		.get("/api/widget/pages/create")
		.public()
		.handle((_ctx, res) => res.json(createPage()));
	app
		.get("/api/widget/pages/poll")
		.public()
		.handle((ctx, res) => {
			const pollId = Number(ctx.request?.query?.poll || 0);
			res.json(pollPage(pollId));
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
			const pollId = createPoll(
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
			const poll = getPoll(pollId);
			const name = String(body.name || "").trim();
			if (!poll || !name) {
				return res.redirect(303, `/pages/poll?poll=${pollId}`);
			}
			db.exec(
				"INSERT INTO participants (poll_id, name, created_at) VALUES (?, ?, ?)",
				pollId,
				name,
				nowISO(),
			);
			const participantId = Number(db.query("SELECT last_insert_rowid() AS id")[0].id);
			pollOptions(pollId).forEach((opt) => {
				const rawValue = body[`opt_${opt.id}`];
				const value =
					rawValue === "yes" || rawValue === "no" || rawValue === "maybe" ? rawValue : "no";
				db.exec(
					"INSERT INTO votes (participant_id, option_id, value) VALUES (?, ?, ?)",
					participantId,
					opt.id,
					value,
				);
			});
			res.redirect(303, `/pages/poll?poll=${pollId}`);
		});
}
