const widget = require("widget.dsl");

const fields = widget.crm.fields("Workshop opportunity", (f) =>
	f
		.text("organization", { label: "Organization", group: "Customer" })
		.email("buyerEmail", { label: "Buyer email", group: "Customer" })
		.currency("amount", { label: "Expected value", group: "Commercial", unit: "USD" })
		.select("format", {
			label: "Format",
			group: "Workshop",
			options: [{ value: "onsite-2d", label: "2-day onsite", colorKey: "onsite" }],
		}),
);

const pipeline = widget.crm.pipeline("AI engineering workshops", (p) =>
	p
		.stage("lead", "New lead", { colorKey: "lead", probability: 0.05 })
		.stage("proposal", "Proposal", { colorKey: "proposal", probability: 0.45 })
		.stage("won", "Won / scheduled", { colorKey: "won", probability: 1 }),
);

const deals = [
	{
		id: "deal-acme",
		title: "Acme Robotics workshop",
		amount: 18000,
		stageId: "proposal",
		ownerId: "manuel",
		status: "open",
		fields: { organization: "Acme Robotics", buyerEmail: "maya@acme.example", amount: 18000 },
	},
];

const activities = [
	{
		id: "activity-1",
		kind: "email",
		title: "Sent draft agenda",
		atISO: "2026-07-09T10:00:00Z",
		actor: { id: "manuel", name: "Manuel" },
	},
];

const page = widget.page("Workshop CRM pipeline", (p) =>
	p
		.section("Pipeline", (s) =>
			s.view(
				widget.crm.pipelineBoard(pipeline, deals, (b) =>
					b
						.summaries([
							{ stageId: "lead", amountTotal: 6000, count: 1 },
							{ stageId: "proposal", amountTotal: 18000, count: 1 },
							{ stageId: "won", amountTotal: 0, count: 0 },
						])
						.selected("deal-acme")
						.onMove(widget.crm.intent.moveDeal("${dealId}", "${toStage}"))
						.onOpen(widget.crm.intent.openDeal("${dealId}")),
				),
			),
		)
		.section("Opportunity", (s) =>
			s.view(
				widget.crm.recordFields(deals[0].fields, fields, (r) =>
					r.mode("edit").onChange(widget.crm.intent.updateField("deal-acme", "${key}", "${value}")),
				),
			),
		)
		.section("Activity", (s) => s.view(widget.crm.activityFeed(activities))),
);
