import type { PageShortcutSpec, ShortcutModifier } from "./useWidgetPage";

export interface ShortcutKeyboardEvent {
	key: string;
	altKey: boolean;
	ctrlKey: boolean;
	metaKey: boolean;
	shiftKey: boolean;
	repeat: boolean;
	isComposing: boolean;
	defaultPrevented: boolean;
}

const MODIFIER_ORDER: ShortcutModifier[] = ["Alt", "Control", "Meta", "Shift"];

function normalizeShortcutKey(key: string): string {
	return /^[A-Z]$/i.test(key) ? key.toLowerCase() : key;
}

export function canonicalShortcutChord(key: string, modifiers: ShortcutModifier[] = []): string {
	const present = new Set(modifiers);
	return [...MODIFIER_ORDER.filter((modifier) => present.has(modifier)), normalizeShortcutKey(key)]
		.filter(Boolean)
		.join("+");
}

export function eventShortcutChord(event: ShortcutKeyboardEvent): string {
	const modifiers: ShortcutModifier[] = [];
	if (event.altKey) modifiers.push("Alt");
	if (event.ctrlKey) modifiers.push("Control");
	if (event.metaKey) modifiers.push("Meta");
	if (event.shiftKey) modifiers.push("Shift");
	return canonicalShortcutChord(event.key, modifiers);
}

export interface PageShortcutEventGuards {
	enabled: boolean;
	blockedTarget: boolean;
	nestedKeyboardScope: boolean;
}

export function shouldIgnorePageShortcutEvent(
	event: ShortcutKeyboardEvent,
	guards: PageShortcutEventGuards,
): boolean {
	return (
		!guards.enabled ||
		event.defaultPrevented ||
		event.isComposing ||
		guards.blockedTarget ||
		guards.nestedKeyboardScope
	);
}

export function matchPageShortcut(
	event: ShortcutKeyboardEvent,
	bindings: PageShortcutSpec[],
): PageShortcutSpec | undefined {
	if (event.defaultPrevented || event.isComposing) return undefined;
	const chord = eventShortcutChord(event);
	return bindings.find(
		(binding) =>
			canonicalShortcutChord(binding.key, binding.modifiers) === chord &&
			(!event.repeat || binding.allowRepeat === true),
	);
}

export function formatShortcutChord(binding: Pick<PageShortcutSpec, "key" | "modifiers">): string {
	return canonicalShortcutChord(binding.key, binding.modifiers);
}

export function ariaKeyShortcuts(bindings: PageShortcutSpec[]): string | undefined {
	if (bindings.length === 0) return undefined;
	return bindings.map(formatShortcutChord).join(" ");
}
