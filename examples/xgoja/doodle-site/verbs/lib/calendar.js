function createCalendarHelpers({ widget, act, emptyState, nowISO }) {
	function monthIndex(monthName) {
		const months = {
			jan: 1,
			january: 1,
			feb: 2,
			february: 2,
			mar: 3,
			march: 3,
			apr: 4,
			april: 4,
			may: 5,
			jun: 6,
			june: 6,
			jul: 7,
			july: 7,
			aug: 8,
			august: 8,
			sep: 9,
			sept: 9,
			september: 9,
			oct: 10,
			october: 10,
			nov: 11,
			november: 11,
			dec: 12,
			december: 12,
		};
		return months[String(monthName || "").toLowerCase()] || 0;
	}

	function pad2(n) {
		return String(n).padStart(2, "0");
	}

	function parseSlotLabel(label, fallbackYear) {
		const text = String(label || "").trim();
		let match = text.match(/(\d{4})-(\d{1,2})-(\d{1,2})(?:[ T]+(\d{1,2}):(\d{2}))?/);
		if (match) {
			const hour = Number(match[4] || 9);
			const minute = Number(match[5] || 0);
			return {
				dayISO: `${match[1]}-${pad2(match[2])}-${pad2(match[3])}`,
				startISO: `${match[1]}-${pad2(match[2])}-${pad2(match[3])}T${pad2(hour)}:${pad2(minute)}:00`,
				endISO: `${match[1]}-${pad2(match[2])}-${pad2(match[3])}T${pad2(hour + 1)}:${pad2(minute)}:00`,
				hour,
			};
		}

		match = text.match(
			/(?:mon|tue|wed|thu|fri|sat|sun)?\w*\s*([A-Za-z]{3,9})\s+(\d{1,2}).*?(\d{1,2}):(\d{2})/i,
		);
		if (!match) return null;
		const month = monthIndex(match[1]);
		if (!month) return null;
		const year = Number(fallbackYear || new Date().getFullYear());
		const day = Number(match[2]);
		const hour = Number(match[3]);
		const minute = Number(match[4]);
		return {
			dayISO: `${year}-${pad2(month)}-${pad2(day)}`,
			startISO: `${year}-${pad2(month)}-${pad2(day)}T${pad2(hour)}:${pad2(minute)}:00`,
			endISO: `${year}-${pad2(month)}-${pad2(day)}T${pad2(hour + 1)}:${pad2(minute)}:00`,
			hour,
		};
	}

	function participantNamesFor(voteMap, participants, optionId, values) {
		return participants
			.filter((pt) => values.indexOf(voteMap[pt.id] && voteMap[pt.id][optionId]) >= 0)
			.map((pt) => pt.name);
	}

	function calendarEventsForPoll(poll, options, participants, voteMap, tally) {
		const fallbackYear = String(poll.created_at || nowISO()).slice(0, 4);
		return options
			.map((opt) => {
				const parsed = parseSlotLabel(opt.label, fallbackYear);
				if (!parsed) return null;
				const yesNames = participantNamesFor(voteMap, participants, opt.id, ["yes"]);
				const maybeNames = participantNamesFor(voteMap, participants, opt.id, ["maybe"]);
				const noNames = participantNamesFor(voteMap, participants, opt.id, ["no"]);
				const counts = tally.find((t) => t.option.id === opt.id) || { yes: 0, maybe: 0, no: 0 };
				const summary = `${counts.yes} yes · ${counts.maybe} maybe · ${counts.no} no`;
				return {
					id: `slot-${opt.id}`,
					title: opt.label,
					label: `${opt.label} — ${summary}`,
					dayISO: parsed.dayISO,
					startISO: parsed.startISO,
					endISO: parsed.endISO,
					styleKey: counts.yes > 0 ? "focus" : counts.maybe > 0 ? "personal" : "meeting",
					meta: {
						yes: yesNames,
						maybe: maybeNames,
						no: noNames,
						score: counts.yes * 2 + counts.maybe,
					},
					hour: parsed.hour,
				};
			})
			.filter(Boolean);
	}

	function calendarHourRange(_events) {
		return { start: 8, end: 22 };
	}

	function calendarMarkersForEvents(events) {
		const markers = {};
		events.forEach((event) => {
			const current = markers[event.dayISO] || { count: 0, styleKey: event.styleKey };
			current.count += 1;
			if (event.styleKey === "focus") current.styleKey = "focus";
			else if (event.styleKey === "personal" && current.styleKey !== "focus") {
				current.styleKey = "personal";
			}
			markers[event.dayISO] = current;
		});
		return markers;
	}

	function selectedCalendarState(query, calendarEvents) {
		const selectedSlotId = String((query && query.slot) || "");
		const selectedEvent = calendarEvents.find((event) => event.id === selectedSlotId);
		if (selectedEvent) return { selectedDay: selectedEvent.dayISO, selectedSlotId };

		const requestedDay = String((query && query.day) || "");
		if (requestedDay && calendarEvents.some((event) => event.dayISO === requestedDay)) {
			return { selectedDay: requestedDay, selectedSlotId: "" };
		}

		return {
			selectedDay: calendarEvents.length ? calendarEvents[0].dayISO : "",
			selectedSlotId: "",
		};
	}

	function dayDetailsView(selectedDay, calendarEvents) {
		const dayEvents = calendarEvents.filter((event) => event.dayISO === selectedDay);
		if (!selectedDay || dayEvents.length === 0) {
			return emptyState("No slots on this day", "Select a marked day to see the offered times.");
		}
		return widget.ui.stack(
			{ gap: "sm" },
			widget.ui.caption(`Selected day: ${selectedDay}`),
			...dayEvents.map((event) =>
				widget.ui.card(
					{ title: event.title, density: "condensed" },
					widget.ui.metadata({
						Time: widget.time.slotLabel(event.startISO, event.endISO),
						Yes: event.meta.yes.length ? event.meta.yes.join(", ") : "—",
						Maybe: event.meta.maybe.length ? event.meta.maybe.join(", ") : "—",
						No: event.meta.no.length ? event.meta.no.join(", ") : "—",
						Score: String(event.meta.score),
					}),
				),
			),
		);
	}

	function calendarView(poll, calendarEvents, query) {
		if (calendarEvents.length === 0) {
			return emptyState(
				"No calendar-ready slots",
				'Use labels like "Thu Jul 9 · 19:00" or "2026-07-09 19:00" to render this poll on the calendar.',
			);
		}

		const calendarHours = calendarHourRange(calendarEvents);
		const calendarState = selectedCalendarState(query, calendarEvents);
		const calendarMonth = widget.time.month(
			{
				monthISO: calendarState.selectedDay.slice(0, 7),
				markers: calendarMarkersForEvents(calendarEvents),
			},
			(m) =>
				m
					.selected(calendarState.selectedDay)
					.onSelect(act.navigate(`/pages/poll?poll=${poll.id}&day=$dateISO`)),
		);
		const dayDetails = dayDetailsView(calendarState.selectedDay, calendarEvents);
		const calendarWeek = widget.time.week(calendarEvents, (w) =>
			w
				.range(widget.time.range.week(calendarState.selectedDay))
				.hours(calendarHours.start, calendarHours.end)
				.hourHeight(48)
				.viewportHeight(360)
				.selected(calendarState.selectedSlotId)
				.onSelect(act.navigate(`/pages/poll?poll=${poll.id}&slot=$blockId`)),
		);

		return widget.ui.stack(
			{ gap: "md" },
			widget.ui.splitPane(calendarMonth, dayDetails, {
				ratio: "leftNarrow",
				gutter: "md",
				divider: true,
			}),
			calendarWeek,
		);
	}

	return { calendarEventsForPoll, calendarView };
}

module.exports = { createCalendarHelpers };
