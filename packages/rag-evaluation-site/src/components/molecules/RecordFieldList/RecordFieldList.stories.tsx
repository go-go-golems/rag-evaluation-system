import type { Meta, StoryObj } from "@storybook/react-vite";
import { useState } from "react";
import { tagStyleSet } from "../../../crm";
import type { FieldValue } from "../../../crm/types";
import { RecordFieldList, type RecordFieldListSection } from "./RecordFieldList";

const sections: RecordFieldListSection[] = [
	{
		label: "Details",
		fields: [
			{ key: "email", type: "email", label: "Email" },
			{ key: "phone", type: "phone", label: "Phone" },
			{ key: "ownerId", type: "user", label: "Owner", relatedObject: "user" },
			{ key: "companyId", type: "relation", label: "Company", relatedObject: "company" },
		],
	},
	{
		label: "Custom",
		fields: [
			{
				key: "segment",
				type: "select",
				label: "Segment",
				styleSet: tagStyleSet,
				options: [
					{ value: "enterprise", label: "Enterprise", colorKey: "enterprise" },
					{ value: "mid-market", label: "Mid-Market", colorKey: "mid-market" },
					{ value: "smb", label: "SMB", colorKey: "default" },
				],
			},
			{ key: "nps", type: "number", label: "NPS" },
			{ key: "renewalPct", type: "percent", label: "Renewal likelihood" },
			{ key: "tags", type: "tags", label: "Tags", styleSet: tagStyleSet },
		],
	},
];

const values: Record<string, FieldValue> = {
	email: "dana@acme.com",
	phone: "+1 555 0142",
	ownerId: "u-you",
	companyId: "co-acme",
	segment: "mid-market",
	nps: 9,
	renewalPct: 82,
	tags: ["champion", "enterprise"],
};

const refs: Record<string, { label: string }> = {
	"u-you": { label: "You" },
	"co-acme": { label: "Acme Corp" },
};

const meta = {
	title: "Component Library/Molecules/RecordFieldList",
	component: RecordFieldList,
	args: {
		values,
		sections,
		mode: "read",
		resolveRef: (id: string) => refs[id],
	},
	decorators: [
		(Story) => (
			<div style={{ maxWidth: 420 }}>
				<Story />
			</div>
		),
	],
} satisfies Meta<typeof RecordFieldList>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Read: Story = {};

export const Edit: Story = { args: { mode: "edit" } };

/** Stacked layout — label above control, for narrow columns. */
export const Stacked: Story = { args: { mode: "read", rowLayout: "stacked" } };

/** Edit mode with a validation error on one field. */
export const WithInvalid: Story = { args: { mode: "edit", invalidKeys: ["email"] } };

/** A record with no values yet — every field reads as empty. */
export const Empty: Story = { args: { values: {} } };

/** Interactive edit — changes are applied to local state and reflected in read view. */
export const Interactive: Story = {
	render: () => {
		const [record, setRecord] = useState<Record<string, FieldValue>>(values);
		return (
			<div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 24, maxWidth: 720 }}>
				<RecordFieldList
					values={record}
					sections={sections}
					mode="edit"
					resolveRef={(id) => refs[id]}
					onFieldChange={(key, next) => setRecord((r) => ({ ...r, [key]: next }))}
				/>
				<RecordFieldList
					values={record}
					sections={sections}
					mode="read"
					resolveRef={(id) => refs[id]}
				/>
			</div>
		);
	},
};
