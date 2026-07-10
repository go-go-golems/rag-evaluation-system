import type { Meta, StoryObj } from "@storybook/react-vite";
import { component, type ActionSpec, type WidgetNode } from "./ir";
import { defaultWidgetRegistry } from "./defaultRegistry";
import { WidgetRenderer, type WidgetRendererProps } from "./WidgetRenderer";

const meta = {
	title: "Widget IR/Renderer/ShareLink",
	component: WidgetRenderer,
	args: { registry: defaultWidgetRegistry },
} satisfies Meta<typeof WidgetRenderer>;

export default meta;
type Story = StoryObj<typeof meta>;

type WidgetActionContext = Parameters<NonNullable<WidgetRendererProps["onAction"]>>[1];

const shareLinkNode: WidgetNode = component("ShareLink", {
	label: "Share link",
	href: "/pages/poll?poll=1",
	description: "Copy this URL for poll participants.",
	copyAction: { kind: "copy", value: "/pages/poll?poll=1" } satisfies ActionSpec,
});

export const CopyableShareLink: Story = {
	args: {
		node: shareLinkNode,
		onAction: (action: ActionSpec, context: WidgetActionContext) => {
			console.log("share-link action", action, context);
		},
	},
};

export const LongShareLink: Story = {
	args: {
		node: component("ShareLink", {
			label: "Share link",
			href: "/pages/poll?poll=123456789&day=2026-07-11&slot=slot-987654321",
			description: "Long links truncate but remain copyable.",
			copyAction: {
				kind: "copy",
				value: "/pages/poll?poll=123456789&day=2026-07-11&slot=slot-987654321",
			} satisfies ActionSpec,
		}),
	},
};
