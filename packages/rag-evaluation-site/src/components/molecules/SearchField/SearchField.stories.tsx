import type { Meta, StoryObj } from "@storybook/react-vite";
import { useState } from "react";
import { Caption } from "../../foundation";
import { Stack } from "../../layout";
import { SearchField } from "./SearchField";

const meta = {
	title: "Component Library/Molecules/SearchField",
	component: SearchField,
	args: { value: "", onValueChange: () => {}, style: { width: 220 } },
} satisfies Meta<typeof SearchField>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Empty: Story = {};

export const WithValue: Story = {
	args: { value: "context window" },
};

export const Disabled: Story = {
	args: { value: "context window", disabled: true },
};

export const Interactive: Story = {
	render: () => {
		const [value, setValue] = useState("");
		const [submitted, setSubmitted] = useState("");
		return (
			<Stack gap="sm" style={{ width: 260 }}>
				<SearchField value={value} onValueChange={setValue} onSubmit={setSubmitted} />
				<Caption>submitted: {submitted || "(none)"}</Caption>
			</Stack>
		);
	},
};
