import type { Meta, StoryObj } from "@storybook/react-vite";
import { sampleTeamSyncPoll, sampleTeamSyncTallies } from "../../../scheduling";
import { PollResultsPanel } from "./PollResultsPanel";

const meta = {
	title: "Component Library/Organisms/PollResultsPanel",
	component: PollResultsPanel,
	args: {
		poll: sampleTeamSyncPoll,
		tallies: sampleTeamSyncTallies,
		invited: 6,
		pending: ["Dana", "Erin"],
		style: { maxWidth: 480 },
	},
} satisfies Meta<typeof PollResultsPanel>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const Finalized: Story = {
	args: { poll: { ...sampleTeamSyncPoll, status: "finalized" }, pending: [] },
};

export const AllResponded: Story = {
	args: { invited: 4, pending: [] },
};
