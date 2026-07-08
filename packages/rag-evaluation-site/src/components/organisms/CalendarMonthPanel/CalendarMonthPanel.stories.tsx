import type { Meta, StoryObj } from "@storybook/react-vite";
import { useState } from "react";
import { sampleWeekEvents } from "../../../scheduling";
import { CalendarMonthPanel } from "./CalendarMonthPanel";

const meta = {
	title: "Component Library/Organisms/CalendarMonthPanel",
	component: CalendarMonthPanel,
	args: {
		monthISO: "2026-07",
		events: sampleWeekEvents,
		todayISO: "2026-07-06",
	},
} satisfies Meta<typeof CalendarMonthPanel>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const Interactive: Story = {
	render: (args) => {
		const [monthISO, setMonthISO] = useState("2026-07");
		const [selected, setSelected] = useState<string | undefined>("2026-07-06");
		return (
			<CalendarMonthPanel
				{...args}
				monthISO={monthISO}
				selectedDateISO={selected}
				onDaySelect={setSelected}
				onMonthChange={setMonthISO}
			/>
		);
	},
};
