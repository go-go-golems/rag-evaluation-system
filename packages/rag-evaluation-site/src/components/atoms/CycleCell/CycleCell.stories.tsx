import type { Meta, StoryObj } from "@storybook/react-vite";
import { useState } from "react";
import type { ContextStyleSet } from "../../../context";
import { Inline } from "../../layout";
import { CycleCell } from "./CycleCell";

const AVAILABILITY_STATES = ["yes", "ifneedbe", "no", "unknown"];

const AVAILABILITY_GLYPHS: Record<string, string> = {
	yes: "✓",
	ifneedbe: "~",
	no: "✕",
	unknown: "·",
};

// Story-local palette: the guideline pattern is to map a simple palette arg to a
// ContextStyleSet in render rather than exposing nested style objects as controls.
const availabilityStyleSet: ContextStyleSet = {
	id: "availability",
	styles: {
		yes: { fill: "var(--mac-green)", line: "var(--mac-border)", labelColor: "var(--mac-text-inv)" },
		ifneedbe: {
			fill: "var(--mac-amber)",
			line: "var(--mac-border)",
			labelColor: "var(--mac-text)",
		},
		no: {
			fill: "var(--mac-accent-2)",
			line: "var(--mac-border)",
			labelColor: "var(--mac-text-inv)",
		},
		unknown: {
			fill: "var(--mac-surface)",
			line: "var(--mac-border)",
			labelColor: "var(--mac-text-dim)",
		},
	},
	legend: [
		{ id: "yes", label: "Yes", styleKey: "yes" },
		{ id: "ifneedbe", label: "If need be", styleKey: "ifneedbe" },
		{ id: "no", label: "No", styleKey: "no" },
	],
};

const meta = {
	title: "Design System/Atoms/CycleCell",
	component: CycleCell,
	args: {
		value: "yes",
		states: AVAILABILITY_STATES,
		glyphs: AVAILABILITY_GLYPHS,
		styleSet: availabilityStyleSet,
	},
} satisfies Meta<typeof CycleCell>;

export default meta;
type Story = StoryObj<typeof meta>;

export const AllStates: Story = {
	render: (args) => (
		<Inline gap="sm">
			{AVAILABILITY_STATES.map((state) => (
				<CycleCell key={state} {...args} value={state} />
			))}
		</Inline>
	),
};

export const Interactive: Story = {
	render: (args) => {
		const [value, setValue] = useState("unknown");
		return <CycleCell {...args} value={value} onCycle={setValue} />;
	},
};

export const Selected: Story = {
	args: { value: "yes", selected: true },
};

export const ReadOnly: Story = {
	args: { value: "no", readOnly: true },
};

export const Small: Story = {
	render: (args) => (
		<Inline gap="sm">
			{AVAILABILITY_STATES.map((state) => (
				<CycleCell key={state} {...args} size="sm" value={state} />
			))}
		</Inline>
	),
};

export const WithoutPalette: Story = {
	args: { styleSet: undefined, value: "yes" },
};
