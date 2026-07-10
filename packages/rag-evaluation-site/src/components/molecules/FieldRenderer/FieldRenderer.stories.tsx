import type { Meta, StoryObj } from "@storybook/react-vite";
import { useState } from "react";
import { tagStyleSet } from "../../../crm";
import type { FieldOption, FieldType, FieldValue } from "../../../crm/types";
import { FieldRenderer } from "./FieldRenderer";

const meta = {
	title: "Component Library/Molecules/FieldRenderer",
	component: FieldRenderer,
	args: {
		fieldKey: "email",
		type: "email",
		value: "dana@acme.com",
		mode: "read",
		label: "Email",
	},
} satisfies Meta<typeof FieldRenderer>;

export default meta;
type Story = StoryObj<typeof meta>;

const segmentOptions: FieldOption[] = [
	{ value: "enterprise", label: "Enterprise", colorKey: "enterprise" },
	{ value: "mid-market", label: "Mid-Market", colorKey: "mid-market" },
	{ value: "smb", label: "SMB", colorKey: "default" },
];

const refs: Record<string, { label: string; avatarUrl?: string }> = {
	"u-you": { label: "You" },
	"u-dana": { label: "Dana Whitmore" },
};

/** One example per FieldType — the read/edit appearance table, made concrete. */
const SAMPLES: Array<{
	type: FieldType;
	value: FieldValue;
	extra?: Partial<Parameters<typeof FieldRenderer>[0]>;
}> = [
	{ type: "text", value: "Acme Corp" },
	{ type: "longtext", value: "Wants annual billing; decision expected by end of Q3." },
	{ type: "email", value: "dana@acme.com" },
	{ type: "phone", value: "+1 555 0142" },
	{ type: "url", value: "https://acme.com/pricing" },
	{ type: "number", value: 9 },
	{ type: "currency", value: 8000, extra: { unit: "USD" } },
	{ type: "percent", value: 62 },
	{ type: "date", value: "2026-07-31" },
	{ type: "datetime", value: "2026-07-31T14:00" },
	{ type: "boolean", value: true },
	{
		type: "select",
		value: "mid-market",
		extra: { options: segmentOptions, styleSet: tagStyleSet },
	},
	{
		type: "multiselect",
		value: ["enterprise", "mid-market"],
		extra: { options: segmentOptions, styleSet: tagStyleSet },
	},
	{ type: "tags", value: ["champion", "enterprise"], extra: { styleSet: tagStyleSet } },
	{ type: "user", value: "u-dana", extra: { resolveRef: (id: string) => refs[id] } },
	{
		type: "relation",
		value: "co-acme",
		extra: { resolveRef: () => ({ label: "Acme Corp", href: "#acme" }) },
	},
	{ type: "address", value: "123 Market St\nSan Francisco, CA" },
];

const gridStyle: React.CSSProperties = {
	display: "grid",
	gridTemplateColumns: "120px 1fr 1fr",
	gap: "8px 16px",
	alignItems: "center",
	maxWidth: 640,
};

/** The whole type table, read vs edit side by side. */
export const AllTypes: Story = {
	render: () => (
		<div style={gridStyle}>
			<strong style={{ font: "var(--rag-font-role-compact)" }}>Type</strong>
			<strong style={{ font: "var(--rag-font-role-compact)" }}>Read</strong>
			<strong style={{ font: "var(--rag-font-role-compact)" }}>Edit</strong>
			{SAMPLES.map((s) => (
				<div key={s.type} style={{ display: "contents" }}>
					<code style={{ font: "var(--rag-font-role-compact)", color: "var(--mac-text-dim)" }}>
						{s.type}
					</code>
					<FieldRenderer fieldKey={s.type} type={s.type} value={s.value} mode="read" {...s.extra} />
					<FieldRenderer fieldKey={s.type} type={s.type} value={s.value} mode="edit" {...s.extra} />
				</div>
			))}
		</div>
	),
};

export const Read: Story = { args: { mode: "read" } };

export const Edit: Story = { args: { mode: "edit" } };

/** Empty values render an em dash in read mode. */
export const Empty: Story = {
	render: () => (
		<div style={gridStyle}>
			{(["email", "currency", "select", "user", "tags"] as FieldType[]).map((type) => (
				<div key={type} style={{ display: "contents" }}>
					<code style={{ font: "var(--rag-font-role-compact)", color: "var(--mac-text-dim)" }}>
						{type}
					</code>
					<FieldRenderer fieldKey={type} type={type} value={null} mode="read" />
					<span />
				</div>
			))}
		</div>
	),
};

/** Invalid edit state — the control gets an error outline. */
export const Invalid: Story = {
	args: { type: "email", value: "not-an-email", mode: "edit", invalid: true },
};

/** Interactive: editing a currency field, value echoed below. */
export const Interactive: Story = {
	render: () => {
		const [value, setValue] = useState<FieldValue>(8000);
		return (
			<div style={{ display: "flex", flexDirection: "column", gap: 8, maxWidth: 320 }}>
				<FieldRenderer
					fieldKey="amount"
					type="currency"
					unit="USD"
					label="Amount"
					value={value}
					mode="edit"
					onChange={setValue}
				/>
				<div style={{ font: "var(--rag-font-role-compact)", color: "var(--mac-text-dim)" }}>
					value = {JSON.stringify(value)}
				</div>
				<FieldRenderer fieldKey="amount" type="currency" unit="USD" value={value} mode="read" />
			</div>
		);
	},
};
