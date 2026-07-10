import type { Activity, Company, Contact, CrmUser, Deal, FieldDef, Pipeline, Task } from "./types";

/**
 * Sample CRM data so every widget has something realistic to render in
 * Storybook and in presets. The server would supply these in production; they
 * are hard-coded here exactly like `src/scheduling/fixtures.ts`.
 */

export const sampleUsers: CrmUser[] = [
	{ id: "u-you", name: "You", email: "you@company.com" },
	{ id: "u-dana", name: "Dana Whitmore", email: "dana@acme.com" },
	{ id: "u-lee", name: "Lee Ortiz", email: "lee@company.com" },
	{ id: "u-priya", name: "Priya Raman", email: "priya@company.com" },
];

export const sampleSalesPipeline: Pipeline = {
	id: "pipe-sales",
	name: "Sales",
	stages: [
		{ id: "st-lead", name: "Lead", order: 0, colorKey: "lead", probability: 0.1 },
		{ id: "st-qualified", name: "Qualified", order: 1, colorKey: "qualified", probability: 0.3 },
		{ id: "st-proposal", name: "Proposal", order: 2, colorKey: "proposal", probability: 0.6 },
		{
			id: "st-negotiation",
			name: "Negotiation",
			order: 3,
			colorKey: "negotiation",
			probability: 0.8,
		},
		{ id: "st-won", name: "Won", order: 4, colorKey: "won", probability: 1 },
	],
};

/** Field definitions for a Contact record (drives RecordFieldList + tables). */
export const contactFieldDefs: FieldDef[] = [
	{ key: "email", label: "Email", type: "email", group: "Details" },
	{ key: "phone", label: "Phone", type: "phone", group: "Details" },
	{ key: "ownerId", label: "Owner", type: "user", relatedObject: "user", group: "Details" },
	{
		key: "companyId",
		label: "Company",
		type: "relation",
		relatedObject: "company",
		group: "Details",
	},
	{
		key: "segment",
		label: "Segment",
		type: "select",
		group: "Custom",
		options: [
			{ value: "enterprise", label: "Enterprise", colorKey: "enterprise" },
			{ value: "mid-market", label: "Mid-Market", colorKey: "mid-market" },
			{ value: "smb", label: "SMB", colorKey: "default" },
		],
	},
	{ key: "nps", label: "NPS", type: "number", group: "Custom" },
	{ key: "renewalPct", label: "Renewal likelihood", type: "percent", group: "Custom" },
	{ key: "tags", label: "Tags", type: "tags", group: "Custom" },
];

/** Field definitions for a Deal record. */
export const dealFieldDefs: FieldDef[] = [
	{ key: "amount", label: "Amount", type: "currency", unit: "USD", group: "Details" },
	{ key: "ownerId", label: "Owner", type: "user", relatedObject: "user", group: "Details" },
	{ key: "closeDateISO", label: "Close date", type: "date", group: "Details" },
	{
		key: "companyId",
		label: "Company",
		type: "relation",
		relatedObject: "company",
		group: "Details",
	},
	{
		key: "priority",
		label: "Priority",
		type: "select",
		group: "Custom",
		options: [
			{ value: "low", label: "Low", colorKey: "default" },
			{ value: "med", label: "Medium", colorKey: "mid-market" },
			{ value: "high", label: "High", colorKey: "churn_risk" },
		],
	},
];

export const sampleContacts: Contact[] = [
	{
		id: "c-dana",
		name: "Dana Whitmore",
		title: "VP Sales",
		companyId: "co-acme",
		ownerId: "u-you",
		tags: ["enterprise"],
		updatedAtISO: "2026-07-06T15:03:00",
		fields: {
			email: "dana@acme.com",
			phone: "+1 555 0142",
			ownerId: "u-you",
			companyId: "co-acme",
			segment: "mid-market",
			nps: 9,
			renewalPct: 82,
			tags: ["enterprise", "champion"],
		},
	},
	{
		id: "c-marcus",
		name: "Marcus Vale",
		title: "CTO",
		companyId: "co-globex",
		ownerId: "u-lee",
		tags: ["mid-market"],
		updatedAtISO: "2026-07-02T09:20:00",
		fields: {
			email: "marcus@globex.com",
			phone: "+1 555 0199",
			ownerId: "u-lee",
			companyId: "co-globex",
			segment: "enterprise",
			nps: 7,
			renewalPct: 61,
			tags: ["mid-market"],
		},
	},
];

export const sampleCompanies: Company[] = [
	{
		id: "co-acme",
		name: "Acme Corp",
		domain: "acme.com",
		ownerId: "u-you",
		tags: ["enterprise"],
		fields: {},
	},
	{
		id: "co-globex",
		name: "Globex",
		domain: "globex.com",
		ownerId: "u-lee",
		tags: ["mid-market"],
		fields: {},
	},
];

export const sampleDeals: Deal[] = [
	{
		id: "d-acme-renew",
		title: "Acme renewal",
		amount: 8000,
		currency: "USD",
		stageId: "st-qualified",
		pipelineId: "pipe-sales",
		companyId: "co-acme",
		contactIds: ["c-dana"],
		ownerId: "u-dana",
		closeDateISO: "2026-07-31",
		status: "open",
		fields: { priority: "high" },
	},
	{
		id: "d-acme-expand",
		title: "Acme expansion",
		amount: 25000,
		currency: "USD",
		stageId: "st-proposal",
		pipelineId: "pipe-sales",
		companyId: "co-acme",
		contactIds: ["c-dana"],
		ownerId: "u-lee",
		closeDateISO: "2026-08-15",
		status: "open",
		fields: { priority: "med" },
	},
	{
		id: "d-globex",
		title: "Globex platform",
		amount: 40000,
		currency: "USD",
		stageId: "st-proposal",
		pipelineId: "pipe-sales",
		companyId: "co-globex",
		ownerId: "u-priya",
		closeDateISO: "2026-08-01",
		status: "open",
		fields: { priority: "high" },
	},
	{
		id: "d-initech",
		title: "Initech pilot",
		amount: 12000,
		currency: "USD",
		stageId: "st-lead",
		pipelineId: "pipe-sales",
		ownerId: "u-you",
		status: "open",
		fields: {},
	},
	{
		id: "d-stark",
		title: "Stark Industries",
		amount: 88000,
		currency: "USD",
		stageId: "st-negotiation",
		pipelineId: "pipe-sales",
		ownerId: "u-lee",
		closeDateISO: "2026-07-20",
		status: "open",
		fields: { priority: "high" },
	},
	{
		id: "d-umbrella",
		title: "Umbrella renewal",
		amount: 30000,
		currency: "USD",
		stageId: "st-won",
		pipelineId: "pipe-sales",
		ownerId: "u-you",
		closeDateISO: "2026-07-01",
		status: "won",
		fields: {},
	},
];

export const sampleActivities: Activity[] = [
	{
		id: "a1",
		kind: "email",
		actor: { id: "u-you", name: "You" },
		atISO: "2026-07-07T10:24:00",
		subjectId: "c-dana",
		title: 'Email sent · "Q3 proposal"',
		meta: { attachments: 2 },
	},
	{
		id: "a2",
		kind: "stage_change",
		actor: { id: "u-you", name: "You" },
		atISO: "2026-07-06T09:12:00",
		subjectId: "c-dana",
		title: "Stage changed",
		meta: { from: "Lead", to: "Qualified" },
	},
	{
		id: "a3",
		kind: "call",
		actor: { id: "u-lee", name: "Lee Ortiz" },
		atISO: "2026-07-06T08:40:00",
		subjectId: "c-dana",
		title: "Call · 12 min",
		body: "left voicemail",
		meta: { durationMin: 12 },
	},
	{
		id: "a4",
		kind: "note",
		actor: { id: "u-you", name: "You" },
		atISO: "2026-06-28T15:03:00",
		subjectId: "c-dana",
		title: "Note",
		body: "Wants annual billing, decision by Q3",
	},
];

export const sampleTasks: Task[] = [
	{
		id: "t1",
		title: "Send Acme contract",
		dueISO: "2026-07-07",
		status: "open",
		assigneeId: "u-you",
		relatedId: "d-acme-renew",
		priority: "high",
	},
	{
		id: "t2",
		title: "Follow up with Globex",
		dueISO: "2026-07-08",
		status: "open",
		assigneeId: "u-lee",
		relatedId: "d-globex",
		priority: "med",
	},
	{
		id: "t3",
		title: "Prep Stark demo",
		dueISO: "2026-07-09",
		status: "open",
		assigneeId: "u-you",
		relatedId: "d-stark",
		priority: "high",
	},
	{
		id: "t4",
		title: "Log Umbrella win notes",
		status: "done",
		assigneeId: "u-you",
		relatedId: "d-umbrella",
		priority: "low",
	},
];

/** Server-computed pipeline metrics; hard-coded here for stories. */
export interface StageSummary {
	stageId: string;
	count: number;
	amountTotal: number;
}

export const sampleStageSummaries: StageSummary[] = [
	{ stageId: "st-lead", count: 1, amountTotal: 12000 },
	{ stageId: "st-qualified", count: 1, amountTotal: 8000 },
	{ stageId: "st-proposal", count: 2, amountTotal: 65000 },
	{ stageId: "st-negotiation", count: 1, amountTotal: 88000 },
	{ stageId: "st-won", count: 1, amountTotal: 30000 },
];
