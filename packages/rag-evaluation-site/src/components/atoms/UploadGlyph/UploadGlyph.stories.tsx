import type { Meta, StoryObj } from "@storybook/react-vite";
import { UploadGlyph } from "./UploadGlyph";

const meta = {
	title: "Design System/Atoms/UploadGlyph",
	component: UploadGlyph,
	args: { title: "Upload file" },
} satisfies Meta<typeof UploadGlyph>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const Muted: Story = {
	args: { style: { color: "var(--mac-text-dim)" } },
};
