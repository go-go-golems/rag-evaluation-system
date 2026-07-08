// Doodle-style scheduling site.
//
// Create a poll with proposed time slots; participants record their availability
// (yes / maybe / no) per slot; a results grid tallies the best slot.
//
// - Data lives in SQLite (the `db` module, file-backed via xgoja.v2.yaml).
// - Pages are Widget IR built with the rag Widget DSL (`ui.dsl` + `data.dsl`) and
//   rendered by the React RagEvaluationSiteApp SPA.
// - HTTP is served by the planned-route Express API:
//     app.get(pattern).public().handle((ctx, res) => ...)
//   Query:  ctx.request.query.<name>   Body: ctx.body.<name>   Params: ctx.params.<name>
//   Typed form inputs submit via a native <form> POST (application/x-www-form-urlencoded),
//   which the host parses into ctx.body; handlers respond with res.redirect(303, ...).

__package__({ name: "doodle", short: "Doodle-style scheduling site" });

__verb__("site", {
	name: "site",
	output: "text",
	short: "Serve a Doodle-style scheduling site backed by SQLite and the Widget DSL",
	tags: ["http", "widget", "db", "doodle"],
});

function site() {
	const express = require("express");
	const assets = require("fs:assets");
	const db = require("db");
	const ui = require("ui.dsl");
	const data = require("data.dsl");

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
		// Unique-ish slug from title + id (id guarantees uniqueness).
		db.exec("UPDATE polls SET slug = ? WHERE id = ?", slugify(title) + "-" + pollId, pollId);
		slotLabels.forEach((label, index) => {
			db.exec("INSERT INTO options (poll_id, label, sort) VALUES (?, ?, ?)", pollId, label, index);
		});
		return pollId;
	}

	// Seed one example poll on first run so the site is not empty.
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

	// ---------------------------------------------------------------- rendering
	const navItems = [
		{ id: "index", label: "All polls", action: ui.action.navigate("/pages/index") },
		{ id: "create", label: "New poll", action: ui.action.navigate("/pages/create") },
	];

	function pageMeta(activeNavItemId) {
		return { activeNavItemId: activeNavItemId, navItems: navItems, maxWidth: "wide" };
	}

	const VOTE_GLYPH = { yes: "✓ yes", maybe: "~ maybe", no: "✗ no", "": "— no reply" };

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

		return ui.page({
			schemaVersion: "0.1.0",
			id: "index",
			title: "Doodle · scheduling polls",
			meta: pageMeta("index"),
			sections: [
				ui.panel(
					{ title: "Scheduling polls" },
					ui.statusText({ status: "succeeded", icon: true }, polls.length + " active poll(s)"),
					ui.caption(
						{ tone: "muted" },
						"Create a poll with a few time slots, share it, and let people mark their availability. Data is stored in SQLite.",
					),
					ui.inline(
						{ gap: "sm", wrap: true },
						ui.button(
							{ variant: "primary", action: ui.action.navigate("/pages/create") },
							"+ New poll",
						),
					),
				),
				ui.recipes.metrics({
					items: [
						{ label: "Polls", value: polls.length, status: "ready" },
						{ label: "Total responses", value: totalPeople, status: "succeeded" },
						{
							label: "Time slots",
							value: polls.reduce((acc, p) => acc + Number(p.slots || 0), 0),
							status: "running",
						},
					],
				}),
				ui.panel(
					{ title: "Polls" },
					polls.length === 0
						? ui.emptyState({
								title: "No polls yet",
								description: "Create your first scheduling poll.",
							})
						: data.dataTable({
								rows: rows,
								getRowKey: "id",
								columns: [
									{ id: "title", header: "Poll", cell: data.cell.field("title") },
									{
										id: "location",
										header: "Where",
										cell: data.cell.caption("location", { tone: "muted" }),
									},
									{ id: "slots", header: "Slots", align: "right", cell: data.cell.number("slots") },
									{
										id: "people",
										header: "Responses",
										align: "right",
										cell: data.cell.number("people"),
									},
									{
										id: "created",
										header: "Created",
										cell: data.cell.caption("created", { tone: "muted" }),
									},
								],
							}),
				),
				polls.length === 0
					? ui.panel({ title: "", density: "condensed" }, ui.caption({ tone: "muted" }, " "))
					: ui.panel(
							{ title: "Open a poll" },
							ui.inline(
								{ gap: "sm", wrap: true },
								...rows.map((r) =>
									ui.button(
										{
											variant: "secondary",
											action: ui.action.navigate("/pages/poll?poll=" + r.id),
										},
										r.title + " →",
									),
								),
							),
						),
			],
		});
	}

	function createPage() {
		return ui.page({
			schemaVersion: "0.1.0",
			id: "create",
			title: "New scheduling poll",
			meta: pageMeta("create"),
			sections: [
				ui.panel(
					{ title: "Create a poll" },
					ui.caption(
						{ tone: "muted" },
						"Give the event a title and list one time slot per line. Everything is stored in SQLite.",
					),
				),
				ui.formPanel(
					{
						title: "Event details",
						method: "post",
						formAction: "/api/form/create-poll",
						submitLabel: "Create poll",
					},
					ui.formRow({
						label: "Title",
						required: true,
						control: ui.textInput({
							name: "title",
							placeholder: "Team offsite dinner",
							required: true,
							readOnly: false,
						}),
					}),
					ui.formRow({
						label: "Description",
						control: ui.textareaInput({
							name: "description",
							placeholder: "Optional context for invitees",
							rows: 2,
							readOnly: false,
						}),
					}),
					ui.formRow({
						label: "Location",
						control: ui.textInput({
							name: "location",
							placeholder: "Trattoria Luca, downtown",
							readOnly: false,
						}),
					}),
					ui.formRow({
						label: "Time slots (one per line)",
						required: true,
						control: ui.textareaInput({
							name: "slots",
							placeholder: "Thu Jul 9 · 19:00\nFri Jul 10 · 19:00\nSat Jul 11 · 18:30",
							rows: 5,
							required: true,
							readOnly: false,
						}),
					}),
				),
				ui.panel(
					{ title: "", density: "condensed" },
					ui.button({ action: ui.action.navigate("/pages/index") }, "← Back to all polls"),
				),
			],
		});
	}

	function pollPage(pollId) {
		const poll = getPoll(pollId);
		if (!poll) {
			return ui.page({
				schemaVersion: "0.1.0",
				id: "poll",
				title: "Poll not found",
				meta: pageMeta("index"),
				sections: [
					ui.panel(
						{ title: "Poll not found" },
						ui.statusText({ status: "failed", icon: true }, "No poll with id " + pollId),
						ui.button(
							{ variant: "primary", action: ui.action.navigate("/pages/index") },
							"← All polls",
						),
					),
				],
			});
		}

		const options = pollOptions(pollId);
		const participants = pollParticipants(pollId);
		const votes = pollVotes(pollId);

		// index votes: voteMap[participantId][optionId] = value
		const voteMap = {};
		votes.forEach((v) => {
			const pid = v.participant_id;
			if (!voteMap[pid]) voteMap[pid] = {};
			voteMap[pid][v.option_id] = v.value;
		});

		// Availability grid: one row per participant, one column per option.
		const gridColumns = [{ id: "name", header: "Who", cell: data.cell.field("name") }];
		options.forEach((opt) => {
			gridColumns.push({
				id: "opt_" + opt.id,
				header: opt.label,
				cell: data.cell.field("opt_" + opt.id),
			});
		});
		const gridRows = participants.map((pt) => {
			const row = { id: pt.id, name: pt.name };
			options.forEach((opt) => {
				const value = (voteMap[pt.id] && voteMap[pt.id][opt.id]) || "";
				row["opt_" + opt.id] = VOTE_GLYPH[value] || VOTE_GLYPH[""];
			});
			return row;
		});

		// Tally per option and pick the best slot (yes = 2 pts, maybe = 1 pt).
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
			return { option: opt, yes: yes, maybe: maybe, no: no, score: yes * 2 + maybe };
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
			verdictText: t.score === bestScore && bestScore > 0 ? "★ best slot" : " ",
		}));

		// "Add your availability" form: name + one select per option.
		const formChildren = [
			ui.formRow({
				label: "Your name",
				required: true,
				control: ui.textInput({
					name: "name",
					placeholder: "e.g. Katherine",
					required: true,
					readOnly: false,
				}),
			}),
		];
		options.forEach((opt) => {
			formChildren.push(
				ui.formRow({
					label: opt.label,
					control: ui.selectInput({
						name: "opt_" + opt.id,
						defaultValue: "maybe",
						options: [
							{ value: "yes", label: "✓ Yes" },
							{ value: "maybe", label: "~ Maybe" },
							{ value: "no", label: "✗ No" },
						],
					}),
				}),
			);
		});

		return ui.page({
			schemaVersion: "0.1.0",
			id: "poll",
			title: poll.title,
			meta: pageMeta("index"),
			sections: [
				ui.panel(
					{ title: poll.title },
					ui.metadataGrid({
						density: "compact",
						items: [
							{ key: "Location", value: poll.location || "—" },
							{ key: "Responses", value: String(participants.length) },
							{ key: "Slots", value: String(options.length) },
							{
								key: "Share link",
								value: "/pages/poll?poll=" + poll.id,
								copyValue: "/pages/poll?poll=" + poll.id,
							},
						],
					}),
					poll.description
						? ui.caption({ tone: "muted" }, poll.description)
						: ui.caption({ tone: "muted" }, ""),
					ui.button({ action: ui.action.navigate("/pages/index") }, "← All polls"),
				),
				ui.panel(
					{ title: "Availability grid" },
					participants.length === 0
						? ui.emptyState({
								title: "No responses yet",
								description: "Be the first to add your availability below.",
							})
						: data.dataTable({ rows: gridRows, getRowKey: "id", columns: gridColumns }),
				),
				ui.panel(
					{ title: "Results by slot" },
					data.dataTable({
						rows: resultRows,
						getRowKey: "id",
						columns: [
							{ id: "label", header: "Time slot", cell: data.cell.field("label") },
							{ id: "yes", header: "Yes", align: "right", cell: data.cell.number("yes") },
							{ id: "maybe", header: "Maybe", align: "right", cell: data.cell.number("maybe") },
							{ id: "no", header: "No", align: "right", cell: data.cell.number("no") },
							{ id: "score", header: "Score", align: "right", cell: data.cell.number("score") },
							{ id: "verdict", header: "", cell: data.cell.status("verdict", { icon: true }) },
						],
					}),
				),
				ui.formPanel(
					{
						title: "Add your availability",
						method: "post",
						formAction: "/api/form/cast-vote?poll=" + poll.id,
						submitLabel: "Submit availability",
					},
					...formChildren,
				),
			],
		});
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
			const pollId = Number((ctx.request && ctx.request.query && ctx.request.query.poll) || 0);
			res.json(pollPage(pollId));
		});

	// Native form POST handlers (application/x-www-form-urlencoded -> ctx.body).
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
			res.redirect(303, "/pages/poll?poll=" + pollId);
		});

	app
		.post("/api/form/cast-vote")
		.public()
		.handle((ctx, res) => {
			const body = ctx.body || {};
			const pollId = Number((ctx.request && ctx.request.query && ctx.request.query.poll) || 0);
			const poll = getPoll(pollId);
			const name = String(body.name || "").trim();
			if (!poll || !name) {
				return res.redirect(303, "/pages/poll?poll=" + pollId);
			}
			db.exec(
				"INSERT INTO participants (poll_id, name, created_at) VALUES (?, ?, ?)",
				pollId,
				name,
				nowISO(),
			);
			const participantId = Number(db.query("SELECT last_insert_rowid() AS id")[0].id);
			pollOptions(pollId).forEach((opt) => {
				const raw = body["opt_" + opt.id];
				const value = raw === "yes" || raw === "no" || raw === "maybe" ? raw : "no";
				db.exec(
					"INSERT INTO votes (participant_id, option_id, value) VALUES (?, ?, ?)",
					participantId,
					opt.id,
					value,
				);
			});
			res.redirect(303, "/pages/poll?poll=" + pollId);
		});
}
