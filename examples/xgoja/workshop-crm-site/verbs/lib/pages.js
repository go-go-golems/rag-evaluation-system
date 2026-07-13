function createPages({ widget, store }) {
	const act = widget.act;
	const pipeline = widget.crm.pipeline({ id: "workshops", name: "AI engineering workshops" }, (p) =>
		p
			.stage("lead", "New lead", { colorKey: "lead", probability: 0.1 })
			.stage("proposal", "Proposal", { colorKey: "focus", probability: 0.5 })
			.stage("won", "Scheduled", { colorKey: "success", probability: 1 }),
	);
	const dealFields = widget.crm.fields("Workshop opportunity", (f) =>
		f
			.text("organization", { label: "Organization", group: "Customer" })
			.text("contact", { label: "Primary contact", group: "Customer" })
			.email("email", { label: "Email", group: "Customer" })
			.currency("amount", { label: "Expected value", group: "Commercial", unit: "USD" })
			.select("format", {
				label: "Format",
				group: "Workshop",
				options: [{ value: "onsite-2d", label: "2-day onsite" }],
			})
			.date("workshopDate", { label: "Scheduled date", group: "Workshop" }),
	);

	const navItems = [
		{ id: "index", label: "Command center" },
		{ id: "pipeline", label: "Pipeline" },
		{ id: "lead", label: "New lead" },
		{ id: "runs", label: "Workshop runs" },
	];

	function appShell(active) {
		return widget.app.shell((shell) =>
			shell
				.brand("Workshop CRM")
				.navigation((navigation) =>
					navigation
						.placement("top")
						.active(active)
						.ariaLabel("Workspace")
						.section("workspace", "Workspace", (items) => {
							navItems.forEach((item) =>
								items.item(item.id, item.label, act.navigate(`/pages/${item.id}`)),
							);
						}),
				)
				.content((content) => content.maxWidth("wide").padding("none")),
		);
	}

	function asPage(title, id, configure) {
		return widget
			.page(title, (p) => {
				p.id(id).shell(appShell(id));
				configure(p);
			})
			.toPage();
	}

	function boardDeals(deals) {
		return deals.map((deal) => ({
			id: String(deal.id),
			title: deal.title,
			stageId: deal.stageId,
			amount: Number(deal.amount),
			ownerId: deal.contact,
			status: deal.stageId === "won" ? "success" : "focus",
		}));
	}

	function summaries(deals) {
		return ["lead", "proposal", "won"].map((stageId) => {
			const inStage = deals.filter((deal) => deal.stageId === stageId);
			return {
				stageId,
				count: inStage.length,
				amountTotal: inStage.reduce((sum, deal) => sum + Number(deal.amount), 0),
			};
		});
	}

	function indexPage() {
		const deals = store.allDeals();
		const scheduled = deals.filter((deal) => deal.stageId === "won");
		return asPage("Workshop CRM", "index", (p) =>
			p
				.section("Workshop command center", (s) =>
					s
						.caption(
							"A SQLite-backed lead-to-workshop-run reference host using widget.crm and existing Widget IR engines.",
						)
						.view(
							widget.ui.button("Capture a lead", act.navigate("/pages/lead"), {
								variant: "primary",
							}),
						),
				)
				.section("Commercial health", (s) =>
					s
						.view(
							widget.ui.inline(
								{ gap: "md", wrap: true },
								widget.crm.stat("Open opportunities", String(deals.length)),
								widget.crm.stat("Scheduled runs", String(scheduled.length)),
								widget.crm.stat(
									"Pipeline value",
									`$${deals.reduce((sum, deal) => sum + Number(deal.amount), 0).toLocaleString()}`,
								),
							),
						)
						.view(widget.crm.funnel(pipeline, summaries(deals))),
				)
				.section("Next actions", (s) =>
					s.view(
						widget.ui.inline(
							{ gap: "sm", wrap: true },
							...deals.map((deal) =>
								widget.ui.button(
									`${deal.organization} →`,
									act.navigate(`/pages/opportunity?deal=${deal.id}`),
								),
							),
						),
					),
				),
		);
	}

	function pipelinePage() {
		const deals = store.allDeals();
		return asPage("Sales pipeline", "pipeline", (p) =>
			p
				.section("Workshop sales pipeline", (s) =>
					s
						.caption(
							"Open an opportunity by selecting its card. Drag a card to persist its stage change through the CRM action route.",
						)
						.view(
							widget.crm.pipelineBoard(pipeline, boardDeals(deals), (b) =>
								b
									.summaries(summaries(deals))
									.onMove(widget.crm.intent.moveDeal("${cardId}", "${to}"))
									.onOpen(widget.crm.intent.openDeal("${cardId}")),
							),
						),
				)
				.section("Open opportunities", (s) =>
					s.view(
						widget.ui.stack(
							{ gap: "sm" },
							...deals.map((deal) =>
								widget.ui.button(
									`${deal.title} · ${deal.stageId}`,
									act.navigate(`/pages/opportunity?deal=${deal.id}`),
								),
							),
						),
					),
				),
		);
	}

	function leadPage() {
		const form = widget.ui.form(
			{
				title: "Workshop lead",
				method: "post",
				formAction: "/api/form/create-lead",
				submitLabel: "Create opportunity",
			},
			widget.ui.formRow(
				"Organization",
				widget.ui.textInput({ name: "organization", required: true, placeholder: "Acme Robotics" }),
				{ required: true },
			),
			widget.ui.formRow(
				"Contact",
				widget.ui.textInput({ name: "contact", required: true, placeholder: "Maya Chen" }),
				{ required: true },
			),
			widget.ui.formRow(
				"Email",
				widget.ui.textInput({
					name: "email",
					type: "email",
					required: true,
					placeholder: "maya@example.com",
				}),
				{ required: true },
			),
			widget.ui.formRow(
				"Expected value",
				widget.ui.textInput({ name: "amount", type: "number", value: "18000" }),
			),
			widget.ui.formRow(
				"Format",
				widget.ui.selectInput({
					name: "format",
					value: "onsite-2d",
					options: [
						{ value: "onsite-2d", label: "2-day onsite" },
						{ value: "remote-1d", label: "1-day remote" },
					],
				}),
			),
		);
		return asPage("Capture workshop lead", "lead", (p) =>
			p.section("New lead", (s) =>
				s
					.caption(
						"Create the organization, buyer contact, and CRM opportunity in one SQLite transaction.",
					)
					.view(form),
			),
		);
	}

	function opportunityPage(dealId) {
		const deal = store.getDeal(dealId);
		if (!deal) {
			return asPage("Opportunity not found", "pipeline", (p) =>
				p.section("Opportunity not found", (s) =>
					s.view(widget.ui.button("Back to pipeline", act.navigate("/pages/pipeline"))),
				),
			);
		}
		const activities = store.activitiesForDeal(dealId).map((activity) => ({
			...activity,
			actor: { id: "workshop-crm", name: "Workshop CRM" },
		}));
		const values = {
			organization: deal.organization,
			contact: deal.contact,
			email: deal.email,
			amount: Number(deal.amount),
			format: deal.format,
			workshopDate: deal.workshopDate || "",
		};
		return asPage(deal.title, "pipeline", (p) =>
			p
				.section(deal.title, (s) =>
					s
						.metadata({
							Stage: deal.stageId,
							Organization: deal.organization,
							Contact: deal.contact,
						})
						.view(
							widget.ui.inline(
								{ gap: "sm" },
								widget.ui.button(
									"Availability",
									act.navigate(`/pages/availability?deal=${deal.id}`),
									{
										variant: "primary",
									},
								),
								widget.ui.button("Pipeline", act.navigate("/pages/pipeline")),
							),
						),
				)
				.section("Opportunity record", (s) => s.view(widget.crm.recordFields(values, dealFields)))
				.section("Activity", (s) => s.view(widget.data.activityFeed(activities))),
		);
	}

	function availabilityPage(dealId) {
		const deal = store.getDeal(dealId);
		if (!deal) return opportunityPage(dealId);
		store.createAvailability(dealId);
		const options = store.availabilityForDeal(dealId);
		const poll = {
			title: `${deal.organization} delivery dates`,
			options: options.map((option) => ({ id: String(option.id), label: option.label })),
			responses: [],
		};
		const form = widget.ui.form(
			{
				title: "Schedule workshop run",
				method: "post",
				formAction: `/api/form/schedule-run?deal=${deal.id}`,
				submitLabel: "Schedule selected date",
			},
			widget.ui.formRow(
				"Delivery date",
				widget.ui.selectInput({
					name: "option",
					options: options.map((option) => ({ value: String(option.id), label: option.label })),
				}),
				{ required: true },
			),
		);
		return asPage(`${deal.organization} availability`, "pipeline", (p) =>
			p
				.section("Availability poll", (s) =>
					s
						.caption(
							"Choose a confirmed delivery slot; this creates a workshop run and promotes the opportunity to Scheduled.",
						)
						.view(widget.schedule.availabilityPoll(poll, (b) => b.readOnly())),
				)
				.section("Confirm a delivery date", (s) => s.view(form))
				.section("Opportunity", (s) =>
					s.view(
						widget.ui.button(
							"Back to opportunity",
							act.navigate(`/pages/opportunity?deal=${deal.id}`),
						),
					),
				),
		);
	}

	function runsPage() {
		const runs = store.allRuns();
		const events = runs.map((run) => ({
			id: String(run.id),
			title: run.title,
			startISO: run.startISO,
			endISO: run.endISO,
			styleKey: "focus",
		}));
		return asPage("Workshop runs", "runs", (p) =>
			p
				.section("Scheduled workshops", (s) =>
					s
						.caption(
							"Confirmed delivery runs are persisted separately from pipeline opportunities.",
						)
						.view(
							runs.length
								? widget.time.week(events, (w) =>
										w
											.range(widget.time.range.week(events[0].startISO.slice(0, 10)))
											.hours(8, 18)
											.viewportHeight(420),
									)
								: widget.ui.emptyState(
										"No scheduled runs",
										"Confirm an availability slot from an opportunity.",
									),
						),
				)
				.section("Run records", (s) =>
					s.view(
						widget.ui.stack(
							{ gap: "sm" },
							...runs.map((run) =>
								widget.ui.button(
									`${run.title} · ${run.startISO.slice(0, 10)}`,
									act.navigate(`/pages/opportunity?deal=${run.dealId}`),
								),
							),
						),
					),
				),
		);
	}

	return { indexPage, pipelinePage, leadPage, opportunityPage, availabilityPage, runsPage };
}

module.exports = { createPages };
