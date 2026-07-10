import type { Meta, StoryObj } from "@storybook/react-vite";
import { sampleWeekEvents } from "../scheduling";
import { defaultWidgetRegistry } from "./defaultRegistry";
import { component, type WidgetNode } from "./ir";
import { monthCalendar, weekCalendar } from "./presets/scheduling";
import { WidgetRenderer } from "./WidgetRenderer";

const meta = {
	title: "Widget IR/Renderer/Calendar",
	component: WidgetRenderer,
	args: {
		registry: defaultWidgetRegistry,
		onAction: (action, context) => {
			// eslint-disable-next-line no-console
			console.log("[widget action]", action, context);
		},
	},
} satisfies Meta<typeof WidgetRenderer>;

export default meta;
type Story = StoryObj<typeof meta>;

/** Month heatmap of event density, from the `monthCalendar` preset. */
export const MonthDensity: Story = {
	args: { node: monthCalendar(sampleWeekEvents, "2026-07") },
};

/** Week time-grid from the `weekCalendar` preset. */
export const Week: Story = {
	args: {
		node: component("Panel", { title: "Jul 6 – 10", density: "condensed" }, [
			weekCalendar(sampleWeekEvents, [
				"2026-07-06",
				"2026-07-07",
				"2026-07-08",
				"2026-07-09",
				"2026-07-10",
			]),
		]) as WidgetNode,
	},
};
