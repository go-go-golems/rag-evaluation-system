import type { Meta, StoryObj } from "@storybook/react-vite";
import { contextVisualStyleToCssVars } from "../../../context";
import { stageStyleSet } from "../../../crm";
import { Inline, Stack } from "../../layout";
import { DealCard } from "./DealCard";

const qualifiedAccent = contextVisualStyleToCssVars(stageStyleSet.styles.qualified!);
const proposalAccent = contextVisualStyleToCssVars(stageStyleSet.styles.proposal!);

const meta = {
	title: "Component Library/Molecules/DealCard",
	component: DealCard,
	args: {
		title: "Acme renewal",
		subtitle: "$8,000",
		meta: "🧑 Dana · closes Jul 31",
		status: "open",
		accentStyle: qualifiedAccent,
	},
	// A card is ~220px wide inside a board column; frame it so the story reads true.
	decorators: [
		(Story) => (
			<div style={{ width: 220 }}>
				<Story />
			</div>
		),
	],
} satisfies Meta<typeof DealCard>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const Selected: Story = { args: { selected: true } };

export const Dragging: Story = { args: { dragging: true } };

export const Won: Story = {
	args: { title: "Umbrella renewal", subtitle: "$30,000", status: "won", meta: "🧑 You" },
};

export const Lost: Story = {
	args: { title: "Initech pilot", subtitle: "$12,000", status: "lost", meta: "🧑 You" },
};

/** Long title + rich meta — proves truncation and the accent bar hold up. */
export const Overflow: Story = {
	args: {
		title: "Globex platform-wide multi-year enterprise expansion",
		subtitle: "$88,000",
		meta: "🧑 Priya · 🏷 enterprise · closes Aug 1",
		accentStyle: proposalAccent,
	},
};

/** Minimal card — title only, no subtitle/meta. */
export const TitleOnly: Story = {
	args: { subtitle: undefined, meta: undefined },
};

/** The set of states side by side, as they appear in a column. */
export const Gallery: Story = {
	render: () => (
		<Stack gap="sm">
			<Inline gap="sm">
				<div style={{ width: 200 }}>
					<DealCard
						title="Default"
						subtitle="$8,000"
						meta="🧑 Dana"
						accentStyle={qualifiedAccent}
					/>
				</div>
				<div style={{ width: 200 }}>
					<DealCard
						title="Selected"
						subtitle="$25,000"
						meta="🧑 Lee"
						selected
						accentStyle={proposalAccent}
					/>
				</div>
			</Inline>
			<Inline gap="sm">
				<div style={{ width: 200 }}>
					<DealCard title="Won" subtitle="$30,000" status="won" meta="🧑 You" />
				</div>
				<div style={{ width: 200 }}>
					<DealCard
						title="Dragging"
						subtitle="$40,000"
						dragging
						meta="🧑 Priya"
						accentStyle={proposalAccent}
					/>
				</div>
			</Inline>
		</Stack>
	),
};
