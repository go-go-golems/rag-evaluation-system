import type { Meta, StoryObj } from "@storybook/react-vite";
import { useState } from "react";
import { contextCourseFixture } from "../../../context";
import { StepList } from "./StepList";

const items = contextCourseFixture.agenda.map((item) => ({
	id: item.id,
	index: item.number,
	title: item.title,
	description: item.description,
	meta: item.duration,
}));
const meta = {
	title: "Component Library/Molecules/StepList",
	component: StepList,
	args: { items },
} satisfies Meta<typeof StepList>;
export default meta;
type Story = StoryObj<typeof meta>;
export const MinimalSteps: Story = {
	render: () => {
		const [active, setActive] = useState(items[0]?.id);
		return <StepList items={items} activeItemId={active} onItemSelect={setActive} />;
	},
};

export const ReadOnly: Story = {
	render: () => {
		return <StepList items={items} selectable={false} />;
	},
};

export const ReadOnlyWithActive: Story = {
	render: () => {
		return <StepList items={items} selectable={false} activeItemId={items[2]?.id} />;
	},
};
