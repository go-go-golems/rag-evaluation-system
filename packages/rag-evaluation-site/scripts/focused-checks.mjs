import assert from "node:assert/strict";

const { packTimeGridColumn } = await import(
	"../src/components/molecules/TimeGrid/TimeGrid.logic.ts"
);
const { buildMonthGridCells, shiftMonth } = await import(
	"../src/components/molecules/MonthGrid/MonthGrid.logic.ts"
);
const { resolveStyleByStyle } = await import("../src/widgets/styleBy.logic.ts");
const {
	ariaKeyShortcuts,
	canonicalShortcutChord,
	matchPageShortcut,
	shouldIgnorePageShortcutEvent,
} = await import("../src/hooks/pageShortcuts.logic.ts");

function ids(packed) {
	return packed.map((p) => ({ id: p.block.id, lane: p.lane, lanes: p.lanes }));
}

{
	const packed = packTimeGridColumn(
		[
			{ id: "a", startISO: "2026-07-09T09:00:00Z", endISO: "2026-07-09T10:00:00Z" },
			{ id: "b", startISO: "2026-07-09T09:30:00Z", endISO: "2026-07-09T10:30:00Z" },
			{ id: "c", startISO: "2026-07-09T10:30:00Z", endISO: "2026-07-09T11:00:00Z" },
		],
		8 * 60,
		4 * 60,
	);
	assert.deepEqual(ids(packed), [
		{ id: "a", lane: 0, lanes: 2 },
		{ id: "b", lane: 1, lanes: 2 },
		{ id: "c", lane: 0, lanes: 1 },
	]);
}

{
	const cells = buildMonthGridCells({
		monthISO: "2026-02",
		weekStartsOn: 1,
		minDateISO: "2026-02-03",
		maxDateISO: "2026-02-27",
		todayISO: "2026-02-14",
		selectedDateISO: "2026-02-20",
	});
	assert.equal(cells.length, 35);
	assert.equal(cells[0].dateISO, "2026-01-26");
	assert.equal(cells.at(-1).dateISO, "2026-03-01");
	assert.equal(cells.find((c) => c.dateISO === "2026-02-02").disabled, true);
	assert.equal(cells.find((c) => c.dateISO === "2026-02-03").disabled, false);
	assert.equal(cells.find((c) => c.dateISO === "2026-02-14").isToday, true);
	assert.equal(cells.find((c) => c.dateISO === "2026-02-20").selected, true);
	assert.equal(shiftMonth(2026, 0, -1), "2025-12");
}

{
	const fallback = { pattern: "fallback", fill: "gray" };
	const styleSet = {
		id: "test",
		name: "Test",
		styles: {
			yes: { pattern: "yes", fill: "green" },
			maybe: { pattern: "maybe", fill: "yellow" },
		},
		fallbackStyle: fallback,
	};
	assert.equal(resolveStyleByStyle({ styleSet }, "yes")?.fill, "green");
	assert.equal(resolveStyleByStyle({ styleSet, map: { y: "yes" } }, "y")?.fill, "green");
	assert.equal(
		resolveStyleByStyle({ styleSet, fallbackStyleKey: "maybe" }, "missing")?.fill,
		"yellow",
	);
	assert.equal(resolveStyleByStyle({ styleSet }, "missing"), fallback);
}

{
	const bindings = [
		{
			id: "accept",
			key: "y",
			label: "Yes",
			action: { kind: "server", name: "triage.accept" },
			preventDefault: true,
		},
		{
			id: "save",
			key: "s",
			modifiers: ["Control"],
			label: "Save",
			action: { kind: "server", name: "triage.save" },
		},
	];
	const keyEvent = (overrides = {}) => ({
		key: "y",
		altKey: false,
		ctrlKey: false,
		metaKey: false,
		shiftKey: false,
		repeat: false,
		isComposing: false,
		defaultPrevented: false,
		...overrides,
	});

	assert.equal(canonicalShortcutChord("Y"), "y");
	assert.equal(canonicalShortcutChord("s", ["Control"]), "Control+s");
	assert.equal(matchPageShortcut(keyEvent(), bindings)?.id, "accept");
	assert.equal(matchPageShortcut(keyEvent({ key: "Y" }), bindings)?.id, "accept");
	assert.equal(matchPageShortcut(keyEvent({ key: "s", ctrlKey: true }), bindings)?.id, "save");
	assert.equal(matchPageShortcut(keyEvent({ repeat: true }), bindings), undefined);
	assert.equal(matchPageShortcut(keyEvent({ isComposing: true }), bindings), undefined);
	assert.equal(matchPageShortcut(keyEvent({ defaultPrevented: true }), bindings), undefined);
	assert.equal(ariaKeyShortcuts(bindings), "y Control+s");
	assert.equal(
		shouldIgnorePageShortcutEvent(keyEvent(), {
			enabled: true,
			blockedTarget: true,
			nestedKeyboardScope: false,
		}),
		true,
	);
	assert.equal(
		shouldIgnorePageShortcutEvent(keyEvent(), {
			enabled: true,
			blockedTarget: false,
			nestedKeyboardScope: true,
		}),
		true,
	);
	assert.equal(
		shouldIgnorePageShortcutEvent(keyEvent(), {
			enabled: false,
			blockedTarget: false,
			nestedKeyboardScope: false,
		}),
		true,
	);
}

console.log("focused checks passed");
