import type { Meta, StoryObj } from "@storybook/react-vite";
import {
	ACTIVITY_GLYPHS,
	activityStyleSet,
	sampleActivities,
	sampleContacts,
	tagStyleSet,
} from "../../../crm";
import { Button } from "../../atoms";
import { Caption, Text } from "../../foundation";
import { Panel, Stack } from "../../layout";
import { ActivityFeed, RecordFieldList, type RecordFieldListSection } from "../../molecules";
import { RecordShell } from "./RecordShell";

const contact = sampleContacts[0]!;

const sections: RecordFieldListSection[] = [
	{
		label: "Details",
		fields: [
			{ key: "email", type: "email", label: "Email" },
			{ key: "phone", type: "phone", label: "Phone" },
			{ key: "ownerId", type: "user", label: "Owner", relatedObject: "user" },
			{
				key: "segment",
				type: "select",
				label: "Segment",
				styleSet: tagStyleSet,
				options: [
					{ value: "enterprise", label: "Enterprise", colorKey: "enterprise" },
					{ value: "mid-market", label: "Mid-Market", colorKey: "mid-market" },
				],
			},
		],
	},
	{
		label: "Custom",
		fields: [
			{ key: "nps", type: "number", label: "NPS" },
			{ key: "renewalPct", type: "percent", label: "Renewal" },
		],
	},
];

const refs: Record<string, { label: string }> = { "u-you": { label: "You" } };

const activity = (
	<ActivityFeed
		activities={sampleActivities.map((a) => ({
			id: a.id,
			kind: a.kind,
			title: a.title,
			body: a.body,
			atISO: a.atISO,
			actor: a.actor,
		}))}
		glyphs={ACTIVITY_GLYPHS}
		styleSet={activityStyleSet}
	/>
);

const relatedDeals = (
	<Panel
		title="Deals (2)"
		density="condensed"
		actions={
			<Button size="compact" variant="default">
				+ Deal
			</Button>
		}
	>
		<Stack gap="xs">
			<div style={{ display: "flex", justifyContent: "space-between", gap: 8 }}>
				<Text size="compact">Acme renewal</Text>
				<Caption>$8,000 · Qualified</Caption>
			</div>
			<div style={{ display: "flex", justifyContent: "space-between", gap: 8 }}>
				<Text size="compact">Acme expansion</Text>
				<Caption>$25,000 · Proposal</Caption>
			</div>
		</Stack>
	</Panel>
);

const details = (
	<RecordFieldList
		values={contact.fields}
		sections={sections}
		mode="read"
		resolveRef={(id) => refs[id]}
	/>
);

const meta = {
	title: "Component Library/Organisms/RecordShell",
	component: RecordShell,
	args: {
		identity: {
			name: contact.name,
			subtitle: "VP Sales · Acme Corp · 🏷 enterprise",
			avatarText: "DW",
		},
		actions: (
			<>
				<Button size="compact">Edit</Button>
				<Button size="compact" variant="default">
					Log ▾
				</Button>
			</>
		),
		details,
		activity,
		activityActions: (
			<Button size="compact" variant="default">
				+ Note
			</Button>
		),
		related: relatedDeals,
	},
} satisfies Meta<typeof RecordShell>;

export default meta;
type Story = StoryObj<typeof meta>;

/** Full contact record: header, fields, timeline, related deals. */
export const ContactRecord: Story = {};

/** Deal record — same shell, different field list, no related panel. */
export const DealRecord: Story = {
	args: {
		identity: {
			name: "Acme expansion",
			subtitle: "$25,000 · Proposal · closes Aug 15",
			avatarText: "AE",
		},
		details: (
			<RecordFieldList
				values={{ amount: 25000, ownerId: "u-lee", closeDateISO: "2026-08-15", priority: "high" }}
				sections={[
					{
						label: "Details",
						fields: [
							{ key: "amount", type: "currency", label: "Amount", unit: "USD" },
							{ key: "ownerId", type: "user", label: "Owner" },
							{ key: "closeDateISO", type: "date", label: "Close date" },
						],
					},
				]}
				mode="read"
				resolveRef={(id) => ({ "u-lee": { label: "Lee Ortiz" } })[id]}
			/>
		),
		related: undefined,
	},
};

/** Edit mode — the fields flip to controls; timeline stays read-only. */
export const EditMode: Story = {
	args: {
		details: (
			<RecordFieldList
				values={contact.fields}
				sections={sections}
				mode="edit"
				resolveRef={(id) => refs[id]}
			/>
		),
	},
};

/** No activity yet — right column shows only related lists. */
export const NoActivity: Story = {
	args: { activity: <ActivityFeed activities={[]} />, related: relatedDeals },
};
