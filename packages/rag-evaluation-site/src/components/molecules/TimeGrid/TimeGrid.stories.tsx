import type { Meta, StoryObj } from "@storybook/react-vite";
import { useState } from "react";
import type { ContextStyleSet } from "../../../context";
import { DateTile } from "../../atoms";
import { Stack } from "../../layout";
import { type TimeGridBlock, TimeGrid } from "./TimeGrid";

const eventStyleSet: ContextStyleSet = {
	id: "events",
	styles: {
		meeting: {
			fill: "color-mix(in srgb, var(--mac-accent) 22%, var(--mac-surface))",
			line: "var(--mac-accent)",
			labelColor: "var(--mac-text)",
		},
		focus: {
			fill: "color-mix(in srgb, var(--mac-green) 22%, var(--mac-surface))",
			line: "var(--mac-green)",
			labelColor: "var(--mac-text)",
		},
		personal: {
			fill: "color-mix(in srgb, var(--mac-amber) 26%, var(--mac-surface))",
			line: "var(--mac-amber)",
			labelColor: "var(--mac-text)",
		},
	},
	legend: [],
};

const WEEK = ["2026-07-06", "2026-07-07", "2026-07-08", "2026-07-09", "2026-07-10"];

const dayColumn = (dateISO: string) => ({
	dayISO: dateISO,
	header: (
		<Stack gap="xs" align="center">
			<DateTile dateISO={dateISO} size="sm" />
		</Stack>
	),
});

const BLOCKS: TimeGridBlock[] = [
	{
		id: "b1",
		dayISO: "2026-07-06",
		startISO: "2026-07-06T09:00",
		endISO: "2026-07-06T09:30",
		styleKey: "meeting",
		label: "Standup",
	},
	{
		id: "b2",
		dayISO: "2026-07-06",
		startISO: "2026-07-06T11:00",
		endISO: "2026-07-06T12:30",
		styleKey: "focus",
		label: "Focus block",
	},
	{
		id: "b3",
		dayISO: "2026-07-07",
		startISO: "2026-07-07T10:00",
		endISO: "2026-07-07T11:00",
		styleKey: "meeting",
		label: "1:1 with Chen",
	},
	// Overlapping trio to exercise lane packing.
	{
		id: "b4",
		dayISO: "2026-07-08",
		startISO: "2026-07-08T10:00",
		endISO: "2026-07-08T11:30",
		styleKey: "meeting",
		label: "Design review",
	},
	{
		id: "b5",
		dayISO: "2026-07-08",
		startISO: "2026-07-08T10:30",
		endISO: "2026-07-08T11:00",
		styleKey: "focus",
		label: "Sync",
	},
	{
		id: "b6",
		dayISO: "2026-07-08",
		startISO: "2026-07-08T10:45",
		endISO: "2026-07-08T12:00",
		styleKey: "personal",
		label: "Call",
	},
	{
		id: "b7",
		dayISO: "2026-07-09",
		startISO: "2026-07-09T14:00",
		endISO: "2026-07-09T15:00",
		styleKey: "meeting",
		label: "Team sync",
	},
	{
		id: "b8",
		dayISO: "2026-07-10",
		startISO: "2026-07-10T16:00",
		endISO: "2026-07-10T17:30",
		styleKey: "personal",
		label: "Gym",
	},
];

const meta = {
	title: "Component Library/Molecules/TimeGrid",
	component: TimeGrid,
	args: {
		days: WEEK.map(dayColumn),
		blocks: BLOCKS,
		styleSet: eventStyleSet,
		hourStart: 8,
		hourEnd: 18,
		style: { maxWidth: 640, maxHeight: 520 },
	},
} satisfies Meta<typeof TimeGrid>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Week: Story = {};

export const WithNowIndicator: Story = {
	args: { nowISO: "2026-07-09T13:20" },
};

export const SingleDay: Story = {
	args: {
		days: [dayColumn("2026-07-08")],
		style: { maxWidth: 260, maxHeight: 520 },
	},
};

export const Interactive: Story = {
	render: (args) => {
		const [selected, setSelected] = useState<string | undefined>("b7");
		const [note, setNote] = useState("Click a block or an empty slot.");
		return (
			<Stack gap="sm">
				<span style={{ font: "var(--rag-font-role-metadata)" }}>{note}</span>
				<TimeGrid
					{...args}
					selectedBlockId={selected}
					onBlockSelect={(id) => {
						setSelected(id);
						setNote(`Selected block ${id}`);
					}}
					onSlotCreate={({ dayISO, hour }) => setNote(`Create ${dayISO} at ${hour}:00`)}
				/>
			</Stack>
		);
	},
};

export const CustomBlockRenderer: Story = {
	args: {
		renderBlock: ({ block, onSelect }) => (
			<button
				type="button"
				onClick={onSelect}
				style={{
					width: "100%",
					height: "100%",
					border: "1px dashed var(--mac-accent)",
					background: "var(--mac-surface)",
					font: "var(--rag-font-role-metadata)",
					cursor: "pointer",
				}}
			>
				✦ {block.label}
			</button>
		),
	},
};

export const Empty: Story = {
	args: { blocks: [] },
};
