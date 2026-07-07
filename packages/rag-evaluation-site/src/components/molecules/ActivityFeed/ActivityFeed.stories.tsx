import type { Meta, StoryObj } from "@storybook/react-vite";
import { ACTIVITY_GLYPHS, activityStyleSet, sampleActivities } from "../../../crm";
import type { Activity } from "../../../crm/types";
import { ActivityFeed, type ActivityFeedItem } from "./ActivityFeed";

function toItems(activities: Activity[]): ActivityFeedItem[] {
	return activities.map((a) => ({
		id: a.id,
		kind: a.kind,
		title: a.title,
		body: a.body,
		atISO: a.atISO,
		actor: { name: a.actor.name, avatarUrl: a.actor.avatarUrl },
	}));
}

const meta = {
	title: "Component Library/Molecules/ActivityFeed",
	component: ActivityFeed,
	args: {
		activities: toItems(sampleActivities),
		glyphs: ACTIVITY_GLYPHS,
		styleSet: activityStyleSet,
	},
	decorators: [
		(Story) => (
			<div style={{ maxWidth: 460 }}>
				<Story />
			</div>
		),
	],
} satisfies Meta<typeof ActivityFeed>;

export default meta;
type Story = StoryObj<typeof meta>;

/** The full timeline, grouped by day with the connective spine. */
export const Default: Story = {};

/** With "load earlier" and clickable rows. */
export const WithActions: Story = {
	args: {
		onOpen: (id) => console.log("open", id),
		onLoadMore: () => console.log("load more"),
	},
};

/** Flat (no day grouping). */
export const Flat: Story = { args: { groupByDay: false } };

/** Empty state. */
export const Empty: Story = { args: { activities: [] } };

/** Dense — many activities across several days exercises the spine + grouping. */
export const Dense: Story = {
	args: {
		activities: toItems([
			...sampleActivities,
			{
				id: "b1",
				kind: "meeting",
				actor: { id: "u-you", name: "You" },
				atISO: "2026-06-27T11:00:00",
				subjectId: "c-dana",
				title: "Discovery call",
				body: "45 min",
			},
			{
				id: "b2",
				kind: "task",
				actor: { id: "u-lee", name: "Lee Ortiz" },
				atISO: "2026-06-27T09:30:00",
				subjectId: "c-dana",
				title: "Task · Send deck",
			},
			{
				id: "b3",
				kind: "field_change",
				actor: { id: "u-you", name: "You" },
				atISO: "2026-06-26T16:20:00",
				subjectId: "c-dana",
				title: "Amount → $8,000",
			},
		]),
	},
};

/** Single kind — one long note, wraps in the body. */
export const SingleNote: Story = {
	args: {
		activities: toItems([
			{
				id: "n1",
				kind: "note",
				actor: { id: "u-you", name: "You" },
				atISO: "2026-07-07T09:00:00",
				subjectId: "c-dana",
				title: "Note",
				body: "Customer wants annual billing and a security review before signing; decision expected end of Q3.",
			},
		]),
	},
};
