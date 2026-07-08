import type { Meta, StoryObj } from "@storybook/react-vite";
import { Stack } from "../../layout";
import { RatioBadge } from "./RatioBadge";

const meta = {
	title: "Design System/Atoms/RatioBadge",
	component: RatioBadge,
	args: { count: 5, total: 8 },
} satisfies Meta<typeof RatioBadge>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const DerivedTones: Story = {
	render: () => (
		<Stack gap="sm">
			<RatioBadge count={8} total={8} />
			<RatioBadge count={5} total={8} />
			<RatioBadge count={2} total={8} />
			<RatioBadge count={0} total={8} />
		</Stack>
	),
};

export const WithLabelAndTrack: Story = {
	render: () => (
		<Stack gap="sm">
			<RatioBadge count={6} total={8} label="yes" showTrack />
			<RatioBadge count={1} total={8} label="maybe" showTrack />
		</Stack>
	),
};

export const ForcedTone: Story = {
	args: { count: 2, total: 8, tone: "positive" },
};

export const ZeroTotal: Story = {
	args: { count: 0, total: 0 },
};
