import type { Meta, StoryObj } from "@storybook/react-vite";
import { TextInput } from "../../atoms";
import { FormRow } from "../FormRow";
import { FieldGrid } from "./FieldGrid";

const meta = {
	title: "Design System/Layout/FieldGrid",
	component: FieldGrid,
	args: {
		columns: 2,
		gap: "md",
	},
} satisfies Meta<typeof FieldGrid>;
export default meta;
type Story = StoryObj<typeof meta>;

export const TwoColumns: Story = {
	render: (args) => (
		<FieldGrid {...args}>
			<FormRow label="When" control={<TextInput defaultValue="Meetup Club Med · atelier LLM" />} />
			<FormRow label="Where" control={<TextInput defaultValue="Sur place + démo remote" />} />
			<FormRow label="Format" control={<TextInput defaultValue="3h30 · démos + exercices" />} />
			<FormRow label="Kicker" control={<TextInput placeholder="Workshop · GenAI" />} />
		</FieldGrid>
	),
};

export const ThreeColumns: Story = {
	args: { columns: 3, gap: "sm" },
	render: (args) => (
		<FieldGrid {...args}>
			<FormRow label="Time" control={<TextInput defaultValue="14h30" />} />
			<FormRow label="Duration" control={<TextInput defaultValue="15 min" />} />
			<FormRow label="ID" control={<TextInput defaultValue="agenda-foundations" />} />
		</FieldGrid>
	),
};
