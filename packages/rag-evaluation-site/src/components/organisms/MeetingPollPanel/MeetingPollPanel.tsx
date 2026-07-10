import type { HTMLAttributes } from "react";
import {
	AVAILABILITY_GLYPHS,
	AVAILABILITY_STATES,
	availabilityStyleSet,
	type AvailabilityState,
	type MeetingPoll,
	type ParticipantResponse,
	type PollOption,
	type SlotTally,
} from "../../../scheduling";
import type { ContextStyleSet } from "../../../context";
import { Button, CycleCell, DateTile, RatioBadge, TextInput } from "../../atoms";
import { Caption, Text } from "../../foundation";
import { Inline, Panel, Stack } from "../../layout";
import { type MatrixColumnSpec, KeyValueStrip, MatrixGrid } from "../../molecules";

const MONTHS = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];

function slotDate(startISO: string): { month: string; day: number } {
	const [, month, day] = startISO.slice(0, 10).split("-");
	return { month: MONTHS[Number(month) - 1] ?? "?", day: Number(day) };
}

export interface MeetingPollCellToggle {
	responseId: string;
	optionId: string;
	state: AvailabilityState;
}

export interface MeetingPollPanelProps extends Omit<HTMLAttributes<HTMLDivElement>, "onSubmit"> {
	poll: MeetingPoll;
	tallies?: SlotTally[];
	/** Response id whose row is editable (the "You" row). */
	currentResponseId?: string;
	styleSet?: ContextStyleSet;
	readOnly?: boolean;
	draftName?: string;
	draftComment?: string;
	onCellToggle?: (toggle: MeetingPollCellToggle) => void;
	onNameChange?: (name: string) => void;
	onCommentChange?: (comment: string) => void;
	onSubmit?: () => void;
}

function columnHeader(option: PollOption) {
	const { day } = slotDate(option.slot.startISO);
	const time = option.slot.startISO.slice(11, 16);
	return (
		<Stack gap="xs" align="center">
			<DateTile dateISO={option.slot.startISO.slice(0, 10)} size="sm" />
			<Caption>{`${time}`}</Caption>
			{option.note ? <Caption tone="muted">{option.note}</Caption> : null}
			<Text size="metric" tone="muted">
				{String(day)}
			</Text>
		</Stack>
	);
}

export function MeetingPollPanel({
	poll,
	tallies,
	currentResponseId,
	styleSet = availabilityStyleSet,
	readOnly = false,
	draftName,
	draftComment,
	onCellToggle,
	onNameChange,
	onCommentChange,
	onSubmit,
	className,
	...rest
}: MeetingPollPanelProps) {
	const tallyByOption = new Map((tallies ?? []).map((t) => [t.optionId, t]));
	const total = poll.responses.length;
	const canEdit = currentResponseId != null && !readOnly;
	const editing = canEdit;

	const columns: MatrixColumnSpec[] = poll.options.map((option) => ({
		id: option.id,
		header: columnHeader(option),
		meta: { option },
	}));

	const handleCell = canEdit
		? ({ rowKey, colId, value }: { rowKey: string; colId: string; value: unknown }) =>
				onCellToggle?.({
					responseId: rowKey,
					optionId: colId,
					state: value as AvailabilityState,
				})
		: undefined;

	const deadline = poll.settings.deadlineISO ? slotDate(poll.settings.deadlineISO) : null;

	return (
		<div className={className} data-rag-organism="MeetingPollPanel" {...rest}>
			<Panel title={poll.title} density="condensed">
				<Stack gap="sm">
					<KeyValueStrip
						items={[
							...(poll.location ? [{ key: "Location", value: poll.location }] : []),
							{ key: "Organizer", value: poll.organizer.name },
						]}
					/>
					<Caption tone="muted">
						{deadline ? `Closes ${deadline.month} ${deadline.day} · ` : ""}
						{`${total} responded`}
						{poll.status === "finalized" ? " · finalized" : ""}
					</Caption>

					<MatrixGrid<ParticipantResponse>
						ariaLabel={poll.title}
						rows={poll.responses}
						columns={columns}
						getRowKey={(row) => row.id}
						renderRowHeader={(row) => (
							<Inline gap="xs">
								<Text size="compact">{row.name}</Text>
								{row.id === currentResponseId ? <Caption tone="accent">✎</Caption> : null}
								{row.comment ? <Caption tone="muted">💬</Caption> : null}
							</Inline>
						)}
						valueAt={(row, col) => row.cells[col.id] ?? "unknown"}
						editableRowKey={canEdit ? currentResponseId : undefined}
						renderCell={(p) => (
							<CycleCell
								value={String(p.value)}
								states={AVAILABILITY_STATES}
								glyphs={AVAILABILITY_GLYPHS}
								styleSet={styleSet}
								readOnly={!p.editable}
								selected={p.selected}
								onCycle={(next) => p.onAction({ value: next })}
							/>
						)}
						onCell={handleCell}
						footer={{
							header: <Caption tone="muted">yes</Caption>,
							render: (col) => {
								const tally = tallyByOption.get(col.id);
								if (!tally) return null;
								return (
									<Inline gap="xs" justify="start">
										<RatioBadge count={tally.yes} total={total} />
										{tally.isBest ? <Caption tone="accent">★</Caption> : null}
									</Inline>
								);
							},
						}}
					/>

					{editing ? (
						<Inline gap="sm" wrap>
							<TextInput
								placeholder="Your name"
								value={draftName ?? ""}
								onChange={(e) => onNameChange?.(e.target.value)}
								aria-label="Your name"
							/>
							<TextInput
								placeholder="Comment (optional)"
								value={draftComment ?? ""}
								onChange={(e) => onCommentChange?.(e.target.value)}
								aria-label="Comment"
							/>
							<Button variant="primary" onClick={() => onSubmit?.()}>
								Submit
							</Button>
						</Inline>
					) : null}
				</Stack>
			</Panel>
		</div>
	);
}
