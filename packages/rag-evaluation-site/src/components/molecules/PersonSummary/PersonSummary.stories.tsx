import type { Meta, StoryObj } from "@storybook/react-vite";
import { PersonSummary } from "./PersonSummary";
const meta = {
	title: "Component Library/Molecules/PersonSummary",
	component: PersonSummary,
	args: {
		name: "M. Calder",
		subtitle: "Applied ML, Developer Tools",
		bio: "Builds agent runtimes; has spent more nights staring at token budgets than is healthy.",
	},
} satisfies Meta<typeof PersonSummary>;
export default meta;
type Story = StoryObj<typeof meta>;
export const Instructor: Story = {};
