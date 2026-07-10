import type { Meta, StoryObj } from "@storybook/react-vite";
import { Inline } from "../../layout";
import { DateTile } from "./DateTile";

const meta = {
	title: "Design System/Atoms/DateTile",
	component: DateTile,
	args: { dateISO: "2026-07-09" },
} satisfies Meta<typeof DateTile>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const Emphasis: Story = {
	render: () => (
		<Inline gap="sm">
			<DateTile dateISO="2026-07-09" emphasis="default" />
			<DateTile dateISO="2026-07-09" emphasis="muted" />
			<DateTile dateISO="2026-07-09" emphasis="accent" />
		</Inline>
	),
};

export const Sizes: Story = {
	render: () => (
		<Inline gap="sm">
			<DateTile dateISO="2026-07-09" size="sm" />
			<DateTile dateISO="2026-07-09" size="md" />
			<DateTile dateISO="2026-07-09" size="lg" />
		</Inline>
	),
};

export const WithoutWeekday: Story = {
	args: { hideWeekday: true },
};

export const FullIsoInput: Story = {
	args: { dateISO: "2026-12-31T14:00:00Z" },
};

export const InvalidDate: Story = {
	args: { dateISO: "not-a-date" },
};
