import type { HTMLAttributes } from "react";
import type { ContextStyleSet } from "../../../context";
import { availabilityStyleSet, type MeetingPoll, type SlotTally } from "../../../scheduling";
import { Button } from "../../atoms";
import { Caption, Divider } from "../../foundation";
import { Inline, Panel, Stack } from "../../layout";
import { SegmentedBar } from "../../molecules";

const MONTHS = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];

function formatSlot(startISO: string): string {
	const [, month, day] = startISO.slice(0, 10).split("-");
	return `${MONTHS[Number(month) - 1] ?? "?"} ${Number(day)} · ${startISO.slice(11, 16)}`;
}

export interface PollResultsPanelProps extends HTMLAttributes<HTMLDivElement> {
	poll: MeetingPoll;
	/** Ranked tallies (isBest flagged); rendered in the given order. */
	tallies: SlotTally[];
	invited?: number;
	pending?: string[];
	styleSet?: ContextStyleSet;
	onPick?: (optionId: string) => void;
	onFinalize?: () => void;
	onRemind?: () => void;
}

export function PollResultsPanel({
	poll,
	tallies,
	invited,
	pending,
	styleSet = availabilityStyleSet,
	onPick,
	onFinalize,
	onRemind,
	className,
	...rest
}: PollResultsPanelProps) {
	const responded = poll.responses.length;
	const invitedCount = invited ?? responded;
	const finalized = poll.status === "finalized";

	return (
		<div className={className} data-rag-organism="PollResultsPanel" {...rest}>
			<Panel
				title={`Results · ${poll.title}`}
				density="condensed"
				actions={
					!finalized && onFinalize ? (
						<Button variant="primary" onClick={() => onFinalize()}>
							Finalize time
						</Button>
					) : finalized ? (
						<Caption tone="success">Finalized</Caption>
					) : undefined
				}
			>
				<Stack gap="md">
					{tallies.map((tally) => {
						const option = poll.options.find((o) => o.id === tally.optionId);
						return (
							<Inline key={tally.optionId} gap="md" justify="between">
								<Stack gap="xs" style={{ flex: "1 1 auto", minWidth: 0 }}>
									<Caption tone={tally.isBest ? "accent" : "muted"}>
										{option ? formatSlot(option.slot.startISO) : tally.optionId}
										{tally.isBest ? " ★ best" : ""}
									</Caption>
									<SegmentedBar
										styleSet={styleSet}
										showCounts
										segments={[
											{ value: tally.yes, styleKey: "yes", label: "yes" },
											{ value: tally.ifneedbe, styleKey: "ifneedbe", label: "maybe" },
											{ value: tally.no, styleKey: "no", label: "no" },
										]}
									/>
								</Stack>
								{!finalized && onPick ? (
									<Button size="compact" onClick={() => onPick(tally.optionId)}>
										Pick
									</Button>
								) : null}
							</Inline>
						);
					})}

					<Divider />

					<Inline gap="sm" justify="between">
						<Caption tone="muted">
							{`${responded} of ${invitedCount} responded`}
							{pending && pending.length ? ` · pending: ${pending.join(", ")}` : ""}
						</Caption>
						{pending && pending.length && onRemind ? (
							<Button size="compact" onClick={() => onRemind()}>
								Remind
							</Button>
						) : null}
					</Inline>
				</Stack>
			</Panel>
		</div>
	);
}
