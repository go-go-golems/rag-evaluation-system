import { useEffect } from "react";
import type { WidgetActionHandler } from "../widgets/actions";
import { matchPageShortcut, shouldIgnorePageShortcutEvent } from "./pageShortcuts.logic";
import type { PageShortcutSpec } from "./useWidgetPage";

export interface UsePageShortcutsOptions {
	pageId: string;
	bindings: PageShortcutSpec[];
	enabled?: boolean;
	onAction: WidgetActionHandler;
}

function isEditableOrModalTarget(target: EventTarget | null): boolean {
	if (!(target instanceof HTMLElement)) return false;
	return (
		target.isContentEditable ||
		Boolean(
			target.closest(
				"input, textarea, select, [contenteditable='true'], dialog[open], [role='dialog']",
			),
		)
	);
}

function isOwnedByNestedKeyboardScope(target: EventTarget | null): boolean {
	return target instanceof HTMLElement && Boolean(target.closest("[data-rag-keyboard-scope]"));
}

function hasActiveModal(): boolean {
	return (
		typeof document !== "undefined" &&
		Boolean(document.querySelector("dialog[open], [role='dialog'][aria-modal='true']"))
	);
}

export function usePageShortcuts({
	pageId,
	bindings,
	enabled = true,
	onAction,
}: UsePageShortcutsOptions): void {
	useEffect(() => {
		if (!enabled || bindings.length === 0 || typeof window === "undefined") return;

		const handleKeyDown = (event: KeyboardEvent) => {
			if (
				shouldIgnorePageShortcutEvent(event, {
					enabled,
					blockedTarget: isEditableOrModalTarget(event.target) || hasActiveModal(),
					nestedKeyboardScope: isOwnedByNestedKeyboardScope(event.target),
				})
			)
				return;
			const binding = matchPageShortcut(event, bindings);
			if (!binding) return;
			if (binding.preventDefault !== false) event.preventDefault();
			void onAction(binding.action, {
				componentType: "PageShortcut",
				pageId,
				shortcutId: binding.id,
				key: event.key,
			});
		};

		window.addEventListener("keydown", handleKeyDown);
		return () => window.removeEventListener("keydown", handleKeyDown);
	}, [bindings, enabled, onAction, pageId]);
}
