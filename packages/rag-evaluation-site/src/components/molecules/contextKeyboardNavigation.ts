import type { KeyboardEvent } from "react";

export type ContextArrowOrientation = "horizontal" | "vertical" | "both";

export function handleContextPartKeyDown(
	event: KeyboardEvent<HTMLElement>,
	focusedPartId: string,
	orderedPartIds: string[],
	onPartSelect: ((partId: string) => void) | undefined,
	orientation: ContextArrowOrientation,
	navigationOriginPartId?: string,
) {
	if (!onPartSelect) return;
	if (event.key === "Enter" || event.key === " ") {
		event.preventDefault();
		onPartSelect(focusedPartId);
		return;
	}

	const previousKeys =
		orientation === "horizontal"
			? ["ArrowLeft"]
			: orientation === "vertical"
				? ["ArrowUp"]
				: ["ArrowLeft", "ArrowUp"];
	const nextKeys =
		orientation === "horizontal"
			? ["ArrowRight"]
			: orientation === "vertical"
				? ["ArrowDown"]
				: ["ArrowRight", "ArrowDown"];
	const originPartId = navigationOriginPartId ?? focusedPartId;
	const currentIndex = orderedPartIds.indexOf(originPartId);
	if (currentIndex < 0) return;

	if (previousKeys.includes(event.key)) {
		event.preventDefault();
		onPartSelect(orderedPartIds[Math.max(0, currentIndex - 1)] ?? originPartId);
		return;
	}
	if (nextKeys.includes(event.key)) {
		event.preventDefault();
		onPartSelect(
			orderedPartIds[Math.min(orderedPartIds.length - 1, currentIndex + 1)] ?? originPartId,
		);
		return;
	}
	if (event.key === "Home") {
		event.preventDefault();
		onPartSelect(orderedPartIds[0] ?? originPartId);
		return;
	}
	if (event.key === "End") {
		event.preventDefault();
		onPartSelect(orderedPartIds[orderedPartIds.length - 1] ?? originPartId);
	}
}
