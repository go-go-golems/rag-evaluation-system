const { createCalendarHelpers } = require("./calendar");
const { createWidgetHelpers } = require("./widget-helpers");

function createPages({ widget, store }) {
	const helpers = createWidgetHelpers(widget);
	const {
		act,
		statusText,
		emptyState,
		formRow,
		textInput,
		textareaInput,
		selectInput,
		collectionTable,
		applyPageMeta,
		asPage,
	} = helpers;
	const calendar = createCalendarHelpers({ widget, act, emptyState, nowISO: store.nowISO });

	function indexPage() {
		const polls = store.allPolls();
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

	function pollPage(pollId, query = {}) {
		const poll = store.getPoll(pollId);
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

		const options = store.pollOptions(pollId);
		const participants = store.pollParticipants(pollId);
		const votes = store.pollVotes(pollId);
		const voteMap = buildVoteMap(votes);
		const tally = buildTally(options, participants, voteMap);
		const availabilityPoll = buildAvailabilityPoll(poll, options, participants, voteMap);
		const resultRows = buildResultRows(tally);
		const summaryTallies = buildSummaryTallies(tally);
		const availabilityGrid = widget.schedule.availabilityPoll(availabilityPoll, (b) =>
			b.readOnly(),
		);
		const resultsTable = resultTable(resultRows);
		const summaryGrid = widget.schedule.pollSummary(availabilityPoll, summaryTallies);
		const calendarEvents = calendar.calendarEventsForPoll(
			poll,
			options,
			participants,
			voteMap,
			tally,
		);
		const availabilityForm = voteForm(poll, options);

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
					.section("Calendar view", (s) =>
						s
							.caption(
								"The same offered slots are rendered as calendar widgets. The month view uses compact markers; the detail column lists who answered yes, maybe, or no for the selected day.",
							)
							.view(calendar.calendarView(poll, calendarEvents, query)),
					)
					.section("Add your availability", (s) => s.view(availabilityForm));
			}),
		);
	}

	function buildVoteMap(votes) {
		const voteMap = {};
		votes.forEach((v) => {
			const pid = v.participant_id;
			if (!voteMap[pid]) voteMap[pid] = {};
			voteMap[pid][v.option_id] = v.value;
		});
		return voteMap;
	}

	function availabilityValue(value) {
		if (value === "yes") return "available";
		if (value === "no") return "unavailable";
		if (value === "maybe") return "maybe";
		return "unknown";
	}

	function buildAvailabilityPoll(poll, options, participants, voteMap) {
		return {
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
	}

	function buildTally(options, participants, voteMap) {
		return options.map((opt) => {
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
	}

	function buildResultRows(tally) {
		let bestScore = -1;
		tally.forEach((t) => {
			if (t.score > bestScore) bestScore = t.score;
		});
		return tally.map((t) => ({
			id: t.option.id,
			label: t.option.label,
			yes: t.yes,
			maybe: t.maybe,
			no: t.no,
			score: t.score,
			verdict: t.score === bestScore && bestScore > 0 ? "succeeded" : "pending",
		}));
	}

	function resultTable(resultRows) {
		return collectionTable("results", resultRows, (f) =>
			f
				.key("id", { label: "ID" })
				.primary("label", { label: "Time slot" })
				.count("score", { label: "Score" })
				.status("verdict", { label: "Best" }),
		);
	}

	function buildSummaryTallies(tally) {
		return [
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
	}

	function voteForm(poll, options) {
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
		return widget.ui.form(
			{
				title: "Add your availability",
				method: "post",
				formAction: `/api/form/cast-vote?poll=${poll.id}`,
				submitLabel: "Submit availability",
			},
			...formChildren,
		);
	}

	return { indexPage, createPage, pollPage };
}

module.exports = { createPages };
