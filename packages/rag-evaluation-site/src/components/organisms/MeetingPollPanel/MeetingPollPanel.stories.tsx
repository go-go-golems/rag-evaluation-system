import type { Meta, StoryObj } from "@storybook/react-vite";
import { useState } from "react";
import {
	type AvailabilityState,
	type MeetingPoll,
	sampleTeamSyncPoll,
	sampleTeamSyncTallies,
} from "../../../scheduling";
import { MeetingPollPanel } from "./MeetingPollPanel";

const meta = {
	title: "Component Library/Organisms/MeetingPollPanel",
	component: MeetingPollPanel,
	args: {
		poll: sampleTeamSyncPoll,
		tallies: sampleTeamSyncTallies,
		style: { maxWidth: 560 },
	},
} satisfies Meta<typeof MeetingPollPanel>;

export default meta;
type Story = StoryObj<typeof meta>;

/** Read-only view — no editable row, no submit. */
export const ReadOnly: Story = {
	args: { readOnly: true },
};

/** Interactive: the "You" row cycles and a draft name/comment can be entered. */
export const Respond: Story = {
	render: (args) => {
		const [poll, setPoll] = useState<MeetingPoll>(sampleTeamSyncPoll);
		const [name, setName] = useState("You");
		const [comment, setComment] = useState("");
		const [submitted, setSubmitted] = useState<string | null>(null);
		return (
			<div>
				<MeetingPollPanel
					{...args}
					poll={poll}
					currentResponseId="you"
					draftName={name}
					draftComment={comment}
					onCellToggle={({ responseId, optionId, state }) =>
						setPoll((prev) => ({
							...prev,
							responses: prev.responses.map((r) =>
								r.id === responseId
									? { ...r, cells: { ...r.cells, [optionId]: state as AvailabilityState } }
									: r,
							),
						}))
					}
					onNameChange={setName}
					onCommentChange={setComment}
					onSubmit={() => setSubmitted(`${name} submitted (${comment || "no comment"})`)}
				/>
				{submitted ? (
					<p style={{ font: "var(--rag-font-role-metadata)", marginTop: 8 }}>{submitted}</p>
				) : null}
			</div>
		);
	},
};

/** Finalized poll: best slot starred, no editing. */
export const Finalized: Story = {
	args: {
		poll: { ...sampleTeamSyncPoll, status: "finalized" },
		readOnly: true,
	},
};

/** Empty poll (no responses yet). */
export const NoResponses: Story = {
	args: {
		poll: { ...sampleTeamSyncPoll, responses: [] },
		tallies: [],
	},
};
