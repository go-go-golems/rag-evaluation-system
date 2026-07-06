import type { Meta, StoryObj } from "@storybook/react-vite";
import { Stack } from "../../layout";
import { MeterBar } from "./MeterBar";

const meta = {
	title: "Design System/Atoms/MeterBar",
	component: MeterBar,
	args: { value: 0.62, style: { width: 240 } },
} satisfies Meta<typeof MeterBar>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Progress: Story = {
	render: () => (
		<Stack gap="sm" style={{ width: 240 }}>
			<MeterBar value={0} />
			<MeterBar value={0.25} />
			<MeterBar value={0.5} />
			<MeterBar value={1} />
		</Stack>
	),
};

export const Tones: Story = {
	render: () => (
		<Stack gap="sm" style={{ width: 240 }}>
			<MeterBar value={0.62} tone="accent" />
			<MeterBar value={1} tone="success" />
			<MeterBar value={0.18} tone="danger" />
		</Stack>
	),
};

export const WithLabel: Story = {
	render: () => (
		<Stack gap="sm" style={{ width: 280 }}>
			<MeterBar value={0.62} label="62%" />
			<MeterBar value={0.62} tone="danger" label="413 error" />
		</Stack>
	),
};
