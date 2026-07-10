export interface MonthGridLogicCell {
	dateISO: string;
	dayOfMonth: number;
	inMonth: boolean;
	isToday: boolean;
	selected: boolean;
	disabled: boolean;
}

export function parseMonth(monthISO: string): { year: number; month: number } {
	const [y, m] = monthISO.slice(0, 7).split("-");
	return { year: Number(y), month: Number(m) - 1 };
}

export function pad2(n: number): string {
	return String(n).padStart(2, "0");
}

export function isoDate(year: number, month: number, day: number): string {
	// month is 0-based; Date normalizes over/underflow so adjacent-month days work.
	const d = new Date(Date.UTC(year, month, day));
	return `${d.getUTCFullYear()}-${pad2(d.getUTCMonth() + 1)}-${pad2(d.getUTCDate())}`;
}

export function shiftMonth(year: number, month: number, delta: number): string {
	const d = new Date(Date.UTC(year, month + delta, 1));
	return `${d.getUTCFullYear()}-${pad2(d.getUTCMonth() + 1)}`;
}

export function buildMonthGridCells({
	monthISO,
	weekStartsOn = 1,
	todayISO,
	selectedDateISO,
	minDateISO,
	maxDateISO,
}: {
	monthISO: string;
	weekStartsOn?: 0 | 1;
	todayISO?: string;
	selectedDateISO?: string;
	minDateISO?: string;
	maxDateISO?: string;
}): MonthGridLogicCell[] {
	const { year, month } = parseMonth(monthISO);
	const firstWeekday = new Date(Date.UTC(year, month, 1)).getUTCDay();
	const leading = (firstWeekday - weekStartsOn + 7) % 7;
	const daysInMonth = new Date(Date.UTC(year, month + 1, 0)).getUTCDate();
	const weekCount = Math.ceil((leading + daysInMonth) / 7);
	const cellCount = weekCount * 7;

	const cells: MonthGridLogicCell[] = [];
	for (let i = 0; i < cellCount; i++) {
		const dayNumber = i - leading + 1;
		const dateISO = isoDate(year, month, dayNumber);
		const inMonth = dayNumber >= 1 && dayNumber <= daysInMonth;
		const outOfRange =
			(minDateISO != null && dateISO < minDateISO) || (maxDateISO != null && dateISO > maxDateISO);
		const d = new Date(Date.UTC(year, month, dayNumber));
		cells.push({
			dateISO,
			dayOfMonth: d.getUTCDate(),
			inMonth,
			isToday: todayISO != null && dateISO === todayISO,
			selected: selectedDateISO != null && dateISO === selectedDateISO,
			disabled: !inMonth || outOfRange,
		});
	}
	return cells;
}
