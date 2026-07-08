import type { Meta, StoryObj } from "@storybook/react-vite";
import { useState } from "react";
import { sampleWeekEvents } from "../../../scheduling";
import { CalendarWeekPanel } from "./CalendarWeekPanel";

const WEEK = ["2026-07-06", "2026-07-07", "2026-07-08", "2026-07-09", "2026-07-10"];

const meta = {
	title: "Component Library/Organisms/CalendarWeekPanel",
	component: CalendarWeekPanel,
	args: {
		days: WEEK,
		events: sampleWeekEvents,
		title: "Jul 6 – 10",
		style: { maxWidth: 680 },
	},
} satisfies Meta<typeof CalendarWeekPanel>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const WithNow: Story = {
	args: { nowISO: "2026-07-09T13:20" },
};

export const Interactive: Story = {
	render: (args) => {
		const [selected, setSelected] = useState<string | undefined>("e3");
		return <CalendarWeekPanel {...args} selectedEventId={selected} onEventSelect={setSelected} />;
	},
};
