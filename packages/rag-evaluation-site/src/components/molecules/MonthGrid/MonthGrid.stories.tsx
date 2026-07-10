import type { Meta, StoryObj } from "@storybook/react-vite";
import { useState } from "react";
import type { ContextStyleSet } from "../../../context";
import { MonthGrid } from "./MonthGrid";

const densityStyleSet: ContextStyleSet = {
	id: "density",
	styles: {
		low: { fill: "var(--mac-surface-2)", labelColor: "var(--mac-text)" },
		mid: { fill: "var(--mac-amber)", labelColor: "var(--mac-text)" },
		high: { fill: "var(--mac-green)", labelColor: "var(--mac-text-inv)" },
	},
	legend: [],
};

const meta = {
	title: "Component Library/Molecules/MonthGrid",
	component: MonthGrid,
	args: { monthISO: "2026-07", todayISO: "2026-07-06" },
} satisfies Meta<typeof MonthGrid>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const WithSelection: Story = {
	render: (args) => {
		const [monthISO, setMonthISO] = useState("2026-07");
		const [selected, setSelected] = useState<string | undefined>("2026-07-09");
		return (
			<MonthGrid
				{...args}
				monthISO={monthISO}
				selectedDateISO={selected}
				onDaySelect={setSelected}
				onMonthChange={setMonthISO}
			/>
		);
	},
};

export const AvailabilityMarkers: Story = {
	args: {
		styleSet: densityStyleSet,
		markers: {
			"2026-07-08": { styleKey: "low", label: "•" },
			"2026-07-09": { styleKey: "high", label: "•" },
			"2026-07-10": { styleKey: "mid", label: "•" },
			"2026-07-16": { styleKey: "mid", label: "•" },
		},
	},
};

export const CountMarkers: Story = {
	args: {
		markers: {
			"2026-07-09": { count: 3 },
			"2026-07-10": { count: 1 },
			"2026-07-22": { count: 5 },
		},
	},
};

export const BookingBounds: Story = {
	args: { minDateISO: "2026-07-06", maxDateISO: "2026-07-24" },
};

export const SundayStart: Story = {
	args: { weekStartsOn: 0 },
};

export const NoHeader: Story = {
	args: { showHeader: false },
};

export const FebruaryLeapYear: Story = {
	args: { monthISO: "2028-02", todayISO: undefined },
};
