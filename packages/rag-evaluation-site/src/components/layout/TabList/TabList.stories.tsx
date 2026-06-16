import { useState } from "react";
import type { Meta, StoryObj } from "@storybook/react-vite";
import { TabList } from "./TabList";

const meta = { title: "Design System/Layout/TabList", component: TabList } satisfies Meta<
	typeof TabList
>;
export default meta;
type Story = StoryObj;

export const Interactive: Story = {
	render: () => {
		const [tab, setTab] = useState("detail");
		return (
			<TabList
				items={[
					{ id: "detail", label: "detail" },
					{ id: "chunk", label: "chunk" },
					{ id: "document", label: "document" },
				]}
				activeId={tab}
				onChange={setTab}
			/>
		);
	},
};
