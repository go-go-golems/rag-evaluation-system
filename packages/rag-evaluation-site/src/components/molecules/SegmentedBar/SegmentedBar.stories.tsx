import type { Meta, StoryObj } from "@storybook/react-vite";
import type { ContextStyleSet } from "../../../context";
import { Stack } from "../../layout";
import { SegmentedBar } from "./SegmentedBar";

const availabilityStyleSet: ContextStyleSet = {
	id: "availability",
	styles: {
		yes: { fill: "var(--mac-green)", labelColor: "var(--mac-text-inv)" },
		ifneedbe: { fill: "var(--mac-amber)", labelColor: "var(--mac-text)" },
		no: { fill: "var(--mac-accent-2)", labelColor: "var(--mac-text-inv)" },
	},
	legend: [
		{ id: "yes", label: "Yes", styleKey: "yes" },
		{ id: "ifneedbe", label: "If need be", styleKey: "ifneedbe" },
		{ id: "no", label: "No", styleKey: "no" },
	],
};

const meta = {
	title: "Component Library/Molecules/SegmentedBar",
	component: SegmentedBar,
	args: {
		styleSet: availabilityStyleSet,
		segments: [
			{ value: 6, styleKey: "yes", label: "yes" },
			{ value: 1, styleKey: "ifneedbe", label: "maybe" },
			{ value: 1, styleKey: "no", label: "no" },
		],
		style: { width: 320 },
	},
} satisfies Meta<typeof SegmentedBar>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const WithCounts: Story = {
	args: { showCounts: true },
};

export const PollResults: Story = {
	render: (args) => (
		<Stack gap="sm" style={{ width: 340 }}>
			<SegmentedBar
				{...args}
				segments={[
					{ value: 6, styleKey: "yes", label: "yes" },
					{ value: 1, styleKey: "ifneedbe", label: "maybe" },
					{ value: 1, styleKey: "no", label: "no" },
				]}
				markers={[{ at: 6, styleKey: "yes", label: "★ best" }]}
			/>
			<SegmentedBar
				{...args}
				segments={[
					{ value: 4, styleKey: "yes" },
					{ value: 3, styleKey: "ifneedbe" },
					{ value: 1, styleKey: "no" },
				]}
			/>
			<SegmentedBar
				{...args}
				segments={[
					{ value: 2, styleKey: "yes" },
					{ value: 0, styleKey: "ifneedbe" },
					{ value: 6, styleKey: "no" },
				]}
			/>
		</Stack>
	),
};

export const Sizes: Story = {
	render: (args) => (
		<Stack gap="sm" style={{ width: 320 }}>
			<SegmentedBar {...args} size="sm" />
			<SegmentedBar {...args} size="md" />
			<SegmentedBar {...args} size="lg" />
		</Stack>
	),
};

export const FixedTotalWithHeadroom: Story = {
	args: {
		total: 12,
		segments: [
			{ value: 4, styleKey: "yes", label: "yes" },
			{ value: 2, styleKey: "ifneedbe", label: "maybe" },
		],
		showCounts: true,
	},
};

export const Empty: Story = {
	args: { segments: [] },
};
