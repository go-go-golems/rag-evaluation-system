import type { Meta, StoryObj } from "@storybook/react-vite";
import { ErrorCallout } from "./ErrorCallout";

const meta = { title: "Design System/Atoms/ErrorCallout", component: ErrorCallout } satisfies Meta<
	typeof ErrorCallout
>;
export default meta;
type Story = StoryObj<typeof meta>;

export const ApiError: Story = {
	args: {
		children: "Failed to load result: provider returned 429",
	},
};

export const Preformatted: Story = {
	args: {
		children: '{\n  "error": "missing source_id"\n}',
	},
};
