import type { Meta, StoryObj } from "@storybook/react-vite";
import { useState } from "react";
import { cmsTagSuggestions } from "../../../cms";
import { TagListInput } from "./TagListInput";

const meta = {
	title: "Component Library/Molecules/TagListInput",
	component: TagListInput,
	args: {
		tags: ["course", "context-window"],
		onAdd: () => {},
		onRemove: () => {},
	},
} satisfies Meta<typeof TagListInput>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const Empty: Story = {
	args: { tags: [] },
};

export const ManyTags: Story = {
	args: {
		tags: [...cmsTagSuggestions, "edge-case", "a-really-long-tag-name-that-should-still-wrap"],
	},
	render: (args) => (
		<div style={{ maxWidth: 360 }}>
			<TagListInput {...args} />
		</div>
	),
};

export const WithSuggestions: Story = {
	args: { suggestions: cmsTagSuggestions },
};

export const Disabled: Story = {
	args: { disabled: true },
};

export const ReadOnly: Story = {
	args: { onAdd: undefined, onRemove: undefined },
};

export const Interactive: Story = {
	render: () => {
		const [tags, setTags] = useState(["course"]);
		return (
			<TagListInput
				tags={tags}
				suggestions={cmsTagSuggestions}
				onAdd={(tag) => setTags((current) => [...current, tag])}
				onRemove={(tag) => setTags((current) => current.filter((entry) => entry !== tag))}
			/>
		);
	},
};
