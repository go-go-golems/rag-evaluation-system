import type { Meta, StoryObj } from "@storybook/react-vite";
import { cmsTagSuggestions } from "../../../cms";
import { Inline } from "../../layout";
import { Tag } from "./Tag";

const meta = {
	title: "Design System/Atoms/Tag",
	component: Tag,
	args: { label: "context-window" },
} satisfies Meta<typeof Tag>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const Selected: Story = {
	args: { selected: true },
};

export const Removable: Story = {
	render: () => (
		<Inline gap="xs">
			<Tag label="course" onRemove={() => {}} />
			<Tag label="diagram" selected onRemove={() => {}} />
		</Inline>
	),
};

export const Disabled: Story = {
	args: { disabled: true, onRemove: () => {} },
};

export const OverflowRow: Story = {
	render: () => (
		<div style={{ maxWidth: 320 }}>
			<Inline gap="xs">
				{[...cmsTagSuggestions, "a-really-long-tag-name-that-should-ellipsize", "extras"].map(
					(tag) => (
						<Tag key={tag} label={tag} onRemove={() => {}} />
					),
				)}
			</Inline>
		</div>
	),
};
