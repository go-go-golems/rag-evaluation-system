import { useCallback, useEffect, useState } from "react";
import type { ActionSpec, RenderableValue, WidgetNode } from "../widgets/ir";

export interface PageNavigationItemSpec {
	id: string;
	label: RenderableValue;
	action?: ActionSpec;
	icon?: RenderableValue;
	badge?: RenderableValue;
	disabled?: boolean;
}

export interface PageNavigationSectionSpec {
	id: string;
	label: RenderableValue;
	items: PageNavigationItemSpec[];
}

export interface PageNavigationSpec {
	placement: "top" | "sidebar";
	brand?: RenderableValue;
	ariaLabel?: string;
	activeItemId?: string;
	sidebarWidth?: number;
	narrowMode?: "stack";
	sections: PageNavigationSectionSpec[];
}

export interface PageContentViewportSpec {
	maxWidth?: "none" | "content" | "wide";
	padding?: "none" | "md" | "lg";
	scroll?: "page" | "main";
}

export type PageShellSpec =
	| { kind: "none" | "root-owned" }
	| { kind: "app"; navigation: PageNavigationSpec; content?: PageContentViewportSpec };

export interface WidgetPageResponse {
	id: string;
	title: string;
	shell?: PageShellSpec;
	root: WidgetNode;
	meta?: Record<string, unknown>;
}

export interface UseWidgetPageOptions {
	enabled?: boolean;
	fetcher?: typeof fetch;
}

export interface UseWidgetPageResult {
	page: WidgetPageResponse | null;
	loading: boolean;
	error: Error | null;
	refresh: () => void;
}

export function useWidgetPage(
	url: string,
	options: UseWidgetPageOptions = {},
): UseWidgetPageResult {
	const { enabled = true, fetcher = fetch } = options;
	const [page, setPage] = useState<WidgetPageResponse | null>(null);
	const [loading, setLoading] = useState(Boolean(enabled && url));
	const [error, setError] = useState<Error | null>(null);
	const [version, setVersion] = useState(0);

	const refresh = useCallback(() => setVersion((current) => current + 1), []);

	useEffect(() => {
		if (!enabled || !url) {
			setLoading(false);
			return;
		}

		const controller = new AbortController();
		setLoading(true);
		setError(null);

		fetcher(url, { signal: controller.signal })
			.then(async (response) => {
				if (!response.ok) {
					throw new Error(`Widget page request failed: ${response.status} ${response.statusText}`);
				}
				return response.json() as Promise<WidgetPageResponse>;
			})
			.then((nextPage) => {
				if (!controller.signal.aborted) {
					setPage(nextPage);
				}
			})
			.catch((err: unknown) => {
				if (controller.signal.aborted) return;
				setError(err instanceof Error ? err : new Error(String(err)));
			})
			.finally(() => {
				if (!controller.signal.aborted) {
					setLoading(false);
				}
			});

		return () => controller.abort();
	}, [enabled, fetcher, url, version]);

	return { page, loading, error, refresh };
}
