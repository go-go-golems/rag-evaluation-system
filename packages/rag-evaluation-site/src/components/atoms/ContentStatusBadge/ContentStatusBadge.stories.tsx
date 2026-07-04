import type { Meta, StoryObj } from "@storybook/react-vite";
import { Text } from "../../foundation";
import { Inline } from "../../layout";
import { ContentStatusBadge } from "./ContentStatusBadge";

const meta = {
	title: "Design System/Atoms/ContentStatusBadge",
	component: ContentStatusBadge,
	args: { status: "draft" },
} satisfies Meta<typeof ContentStatusBadge>;

export default meta;
type Story = StoryObj<typeof meta>;

export const AllStatuses: Story = {
	render: () => (
		<Inline gap="sm">
			<ContentStatusBadge status="draft" />
			<ContentStatusBadge status="published" />
			<ContentStatusBadge status="scheduled" />
			<ContentStatusBadge status="archived" />
		</Inline>
	),
};

export const NoIcon: Story = {
	render: () => (
		<Inline gap="sm">
			<ContentStatusBadge status="draft" icon={false} />
			<ContentStatusBadge status="published" icon={false} />
		</Inline>
	),
};

export const InTableRow: Story = {
	render: () => (
		<Inline gap="md">
			<Text size="compact" as="span">
				The Context Window — Field Guide
			</Text>
			<ContentStatusBadge status="published" />
			<Text size="metadata" tone="muted" as="span">
				2026-06-21
			</Text>
		</Inline>
	),
};
