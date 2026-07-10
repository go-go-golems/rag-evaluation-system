import type { Meta, StoryObj } from "@storybook/react-vite";
import { useState } from "react";
import { sampleBookableDays, sampleBookableSlots, sampleBookingType } from "../../../scheduling";
import { BookingPagePanel } from "./BookingPagePanel";

const meta = {
	title: "Component Library/Organisms/BookingPagePanel",
	component: BookingPagePanel,
	args: {
		bookingType: sampleBookingType,
		monthISO: "2026-07",
		days: sampleBookableDays,
		slots: sampleBookableSlots,
		tz: "Europe/Berlin (CET)",
		style: { maxWidth: 720 },
	},
} satisfies Meta<typeof BookingPagePanel>;

export default meta;
type Story = StoryObj<typeof meta>;

export const NoDaySelected: Story = {
	args: { selectedDateISO: undefined, slots: [] },
};

export const Interactive: Story = {
	render: (args) => {
		const [date, setDate] = useState<string | undefined>("2026-07-09");
		const [slot, setSlot] = useState<string | undefined>();
		const [confirmed, setConfirmed] = useState<string | null>(null);
		return (
			<div>
				<BookingPagePanel
					{...args}
					selectedDateISO={date}
					selectedSlotId={slot}
					slots={date ? args.slots : []}
					onDaySelect={(d) => {
						setDate(d);
						setSlot(undefined);
					}}
					onSlotSelect={setSlot}
					onConfirm={() => setConfirmed(`Booked ${slot} on ${date}`)}
				/>
				{confirmed ? (
					<p style={{ font: "var(--rag-font-role-metadata)", marginTop: 8 }}>{confirmed}</p>
				) : null}
			</div>
		);
	},
};
