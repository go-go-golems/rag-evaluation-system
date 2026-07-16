import { type ReactNode, useCallback, useEffect, useState } from "react";
import { ErrorCallout } from "../components/atoms";
import { Caption } from "../components/foundation";
import { AppShell, Panel, SidebarShell } from "../components/layout";
import { AppNav, KeyboardShortcutHelp, SidebarNav } from "../components/molecules";
import { CourseStudioShell } from "../components/organisms";
import { ariaKeyShortcuts, formatShortcutChord } from "../hooks/pageShortcuts.logic";
import { usePageShortcuts } from "../hooks/usePageShortcuts";
import {
	type PageNavigationItemSpec,
	type PageShellSpec,
	type PageShortcutSpec,
	useWidgetPage,
	type WidgetPageResponse,
} from "../hooks/useWidgetPage";
import {
	confirmWidgetAction,
	dispatchWidgetAction,
	resolveActionPayload,
	type WidgetActionContext,
} from "../widgets/actions";
import { defaultWidgetRegistry } from "../widgets/defaultRegistry";
import type {
	ActionSpec,
	AppNavItemSpec,
	ComponentNode,
	CourseStudioShellWidgetProps,
	RenderableValue,
	WidgetNode,
} from "../widgets/ir";
import { WidgetRenderer } from "../widgets/WidgetRenderer";
import "./app.css";

export interface RagEvaluationSiteAppProps {
	apiBase?: string;
	defaultPageId?: string;
}

type PageShellMode = "auto" | "none" | "app";
type PageMaxWidth = "none" | "content" | "wide";

interface WidgetPageMeta {
	shell?: PageShellMode;
	activeNavItemId?: string;
	navItems?: AppNavItemSpec[];
	maxWidth?: PageMaxWidth;
}

const SHORTCUT_PREFERENCE_KEY = "rag-evaluation-site:page-shortcuts-enabled";
const EMPTY_PAGE_SHORTCUTS: PageShortcutSpec[] = [];

function readShortcutPreference(): boolean {
	if (typeof window === "undefined") return true;
	try {
		return window.localStorage.getItem(SHORTCUT_PREFERENCE_KEY) !== "false";
	} catch {
		return true;
	}
}

const DEFAULT_NAV_ITEMS: AppNavItemSpec[] = [
	{ id: "index", label: "Overview" },
	{ id: "demo", label: "Demo" },
	{ id: "actions", label: "Actions" },
];

export function RagEvaluationSiteApp({
	apiBase = "/api/widget",
	defaultPageId = "index",
}: RagEvaluationSiteAppProps) {
	const [locationVersion, setLocationVersion] = useState(0);
	const [shortcutsEnabled, setShortcutsEnabled] = useState(readShortcutPreference);

	useEffect(() => {
		const handleLocationChange = () => setLocationVersion((version) => version + 1);
		window.addEventListener("popstate", handleLocationChange);
		return () => window.removeEventListener("popstate", handleLocationChange);
	}, []);

	const pageId = readPageIdFromLocation(defaultPageId);
	const pageSearch = readSearchFromLocation(locationVersion);
	const cleanApiBase = apiBase.replace(/\/$/, "");
	const { page, loading, error, refresh } = useWidgetPage(
		`${cleanApiBase}/pages/${encodeURIComponent(pageId)}${pageSearch}`,
	);

	const handleAction = useCallback(
		async (action: ActionSpec, context: WidgetActionContext): Promise<void> => {
			if (!confirmWidgetAction(action, context)) {
				return;
			}
			if (action.kind !== "server") {
				dispatchWidgetAction(action, context);
				return;
			}
			const response = await fetch(`${cleanApiBase}/actions/${encodeURIComponent(action.name)}`, {
				method: "POST",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify({ payload: resolveActionPayload(action.payload, context), context }),
			});
			const result = (await response.json().catch(() => ({
				ok: false,
				error: `Widget action failed: ${response.status} ${response.statusText}`,
			}))) as {
				ok?: boolean;
				refresh?: boolean;
				toast?: string;
				error?: string;
				fieldErrors?: Record<string, string>;
			};
			if (typeof window !== "undefined") {
				window.dispatchEvent(
					new CustomEvent("widget:action-result", {
						detail: { action, context, responseOk: response.ok, result },
					}),
				);
				if (result.toast || result.error)
					window.dispatchEvent(
						new CustomEvent("widget:toast", {
							detail: {
								message: result.toast ?? result.error,
								tone: response.ok && result.ok !== false ? "success" : "danger",
							},
						}),
					);
			}
			if (response.ok && result.refresh) refresh();
		},
		[cleanApiBase, refresh],
	);

	const shortcutBindings = page?.shortcuts?.bindings ?? EMPTY_PAGE_SHORTCUTS;
	usePageShortcuts({
		pageId: page?.id ?? pageId,
		bindings: shortcutBindings,
		enabled: shortcutsEnabled,
		onAction: handleAction,
	});

	const updateShortcutsEnabled = (enabled: boolean) => {
		setShortcutsEnabled(enabled);
		if (typeof window !== "undefined") {
			try {
				window.localStorage.setItem(SHORTCUT_PREFERENCE_KEY, String(enabled));
			} catch {
				// The in-memory preference still applies when storage is unavailable.
			}
		}
	};

	if (loading && !page) {
		return (
			<Panel className="rag-evaluation-site-state" title="RAG Evaluation Site" density="condensed">
				<Caption>Loading Widget IR…</Caption>
			</Panel>
		);
	}

	if (error && !page) {
		return (
			<ErrorCallout className="rag-evaluation-site-state">
				Failed to load Widget IR page: {error.message}
			</ErrorCallout>
		);
	}

	if (!page) {
		return (
			<Panel className="rag-evaluation-site-state" title="RAG Evaluation Site" density="condensed">
				<Caption>No Widget IR returned.</Caption>
			</Panel>
		);
	}

	return (
		<>
			{renderPage(page, pageId, shortcutsEnabled, (action, context) => {
				void handleAction(action, context);
			})}
			<KeyboardShortcutHelp
				items={shortcutBindings.map((binding) => ({
					id: binding.id,
					label: binding.label,
					chord: formatShortcutChord(binding),
				}))}
				enabled={shortcutsEnabled}
				onEnabledChange={updateShortcutsEnabled}
			/>
			{loading && <RoutePendingIndicator />}
			{error && (
				<ErrorCallout className="rag-evaluation-site-inline-error">
					Failed to refresh Widget IR page: {error.message}
				</ErrorCallout>
			)}
		</>
	);
}

function RoutePendingIndicator(): ReactNode {
	return (
		<div
			className="rag-evaluation-site-route-pending"
			role="status"
			aria-live="polite"
			data-rag-route-pending="true"
		>
			<span className="rag-evaluation-site-route-pending__bar" />
			<span className="rag-evaluation-site-route-pending__label">Loading next page…</span>
		</div>
	);
}

function renderPage(
	page: WidgetPageResponse,
	pageId: string,
	shortcutsEnabled: boolean,
	onAction: (action: ActionSpec, context: WidgetActionContext) => void,
): ReactNode {
	const meta = normalizeMeta(page.meta);
	const shell = resolvePageShell(page, pageId, meta);
	const shortcutKeys = shortcutsEnabled
		? ariaKeyShortcuts(page.shortcuts?.bindings ?? EMPTY_PAGE_SHORTCUTS)
		: undefined;
	const usesFramedAppShell = shell.kind === "app" && shell.navigation.placement === "top";
	const rootClassName = [
		"rag-evaluation-site-root",
		usesFramedAppShell ? "rag-evaluation-site-root--shell" : "rag-evaluation-site-root--raw",
	]
		.filter(Boolean)
		.join(" ");

	// Preserve historical course pages while typed root-owned pages use the
	// ordinary registry path shared by every other full-viewport organism.
	const legacyCourseShellNode = page.shell ? null : findRootCourseStudioShellNode(page.root);
	if (legacyCourseShellNode) {
		return (
			<div
				className={rootClassName}
				data-rag-page="RagEvaluationSiteApp"
				data-page-id={page.id}
				data-rag-shell="course-studio"
				aria-keyshortcuts={shortcutKeys}
			>
				{renderCourseStudioShellPage(legacyCourseShellNode, onAction)}
			</div>
		);
	}

	const renderedRoot = (
		<WidgetRenderer node={page.root} registry={defaultWidgetRegistry} onAction={onAction} />
	);

	if (shell.kind !== "app") {
		return (
			<div
				className={rootClassName}
				data-rag-page="RagEvaluationSiteApp"
				data-page-id={page.id}
				data-rag-shell={shell.kind}
				aria-keyshortcuts={shortcutKeys}
			>
				{renderedRoot}
			</div>
		);
	}

	const navigation = shell.navigation;
	const activeItemId = navigation.activeItemId ?? pageId;
	const dispatchNavigation = (itemId: string, componentType: "AppNav" | "SidebarNav") => {
		const item = navigation.sections
			.flatMap((section) => section.items)
			.find((candidate) => candidate.id === itemId);
		if (item?.action) {
			dispatchWidgetAction(item.action, { itemId, value: itemId, componentType }, onAction);
			return;
		}
		navigateToPage(itemId);
	};
	const content = (
		<div
			className={contentClassName(shell.content?.maxWidth)}
			data-rag-layout="PageContent"
			data-rag-scroll={shell.content?.scroll ?? "page"}
		>
			{renderedRoot}
		</div>
	);

	if (navigation.placement === "sidebar") {
		return (
			<div
				className={rootClassName}
				data-rag-page="RagEvaluationSiteApp"
				data-page-id={page.id}
				data-rag-shell="app-sidebar"
				aria-keyshortcuts={shortcutKeys}
			>
				<SidebarShell
					sidebarWidth={navigation.sidebarWidth ?? 188}
					sidebarAriaLabel={navigation.ariaLabel ?? "Primary navigation"}
					narrowMode={navigation.narrowMode ?? "stack"}
					contentPadding={shell.content?.padding ?? "md"}
					header={
						<span className="rag-evaluation-site-brand">
							{renderRenderableValue(navigation.brand ?? page.title, onAction)}
						</span>
					}
					sidebar={
						<SidebarNav
							ariaLabel={navigation.ariaLabel ?? "Primary navigation"}
							sections={navigation.sections.map((section) => ({
								...section,
								label: renderRenderableValue(section.label, onAction),
								items: section.items.map((item) => ({
									...item,
									label: renderRenderableValue(item.label, onAction),
									icon: renderRenderableValue(item.icon, onAction),
									badge: renderRenderableValue(item.badge, onAction),
								})),
							}))}
							activeItemId={activeItemId}
							onItemSelect={(itemId) => dispatchNavigation(itemId, "SidebarNav")}
						/>
					}
				>
					{content}
				</SidebarShell>
			</div>
		);
	}

	const items = navigation.sections.flatMap((section) => section.items);
	return (
		<div
			className={rootClassName}
			data-rag-page="RagEvaluationSiteApp"
			data-page-id={page.id}
			data-rag-shell="app-top"
			aria-keyshortcuts={shortcutKeys}
		>
			<AppShell
				className="rag-evaluation-site-shell"
				header={
					<AppNav
						brand={
							<span className="rag-evaluation-site-brand">
								{renderRenderableValue(navigation.brand ?? page.title, onAction)}
							</span>
						}
						items={items.map((item) => ({
							id: item.id,
							label: renderPageNavigationLabel(item, onAction),
							disabled: item.disabled,
						}))}
						activeItemId={activeItemId}
						onItemSelect={(itemId) => dispatchNavigation(itemId, "AppNav")}
					/>
				}
			>
				{content}
			</AppShell>
		</div>
	);
}

function resolvePageShell(
	page: WidgetPageResponse,
	pageId: string,
	meta: WidgetPageMeta,
): PageShellSpec {
	if (page.shell) return page.shell;
	if (meta.shell === "none" || isAppShellNode(page.root)) return { kind: "none" };
	if (findRootCourseStudioShellNode(page.root)) return { kind: "root-owned" };
	const items = meta.navItems?.length ? meta.navItems : DEFAULT_NAV_ITEMS;
	return {
		kind: "app",
		navigation: {
			placement: "top",
			brand: page.title || "RAG Evaluation Site",
			activeItemId: meta.activeNavItemId ?? pageId,
			ariaLabel: "Primary",
			sections: [{ id: "primary", label: "Primary", items }],
		},
		content: { maxWidth: meta.maxWidth ?? "wide", padding: "none", scroll: "page" },
	};
}

function renderCourseStudioShellPage(
	node: ComponentNode,
	onAction: (action: ActionSpec, context: WidgetActionContext) => void,
): ReactNode {
	const props = (node.props ?? {}) as CourseStudioShellWidgetProps;
	return (
		<CourseStudioShell
			className={props.className}
			sections={(props.sections ?? []).map((section) => ({
				...section,
				label: renderRenderableValue(section.label, onAction),
				items: section.items.map((item) => ({
					...item,
					label: renderRenderableValue(item.label, onAction),
					icon: renderRenderableValue(item.icon, onAction),
					badge: renderRenderableValue(item.badge, onAction),
				})),
			}))}
			activeItemId={props.activeItemId}
			onNavigate={(itemId) => {
				if (props.onNavigateAction)
					dispatchWidgetAction(
						props.onNavigateAction,
						{
							itemId,
							item: { id: itemId },
							value: itemId,
							componentType: "CourseStudioShell",
						},
						onAction,
					);
			}}
			title={renderRenderableValue(props.title, onAction)}
			subtitle={renderRenderableValue(props.subtitle, onAction)}
			sidebarFooter={
				props.sidebarFooter ? (
					<WidgetRenderer
						node={props.sidebarFooter}
						registry={defaultWidgetRegistry}
						onAction={onAction}
					/>
				) : undefined
			}
			contentPadding={props.contentPadding}
		>
			{(node.children ?? []).map((child, index) => (
				<WidgetRenderer
					key={widgetNodeKey(child, index)}
					node={child}
					registry={defaultWidgetRegistry}
					onAction={onAction}
				/>
			))}
		</CourseStudioShell>
	);
}

function isAppShellNode(node: WidgetNode): boolean {
	return node.kind === "component" && node.type === "AppShell";
}

function findRootCourseStudioShellNode(node: WidgetNode): ComponentNode | null {
	if (node.kind !== "component") return null;
	if (node.type === "CourseStudioShell") return node;
	const children = node.children ?? [];
	const onlyChild = children[0];
	if (children.length !== 1 || !onlyChild) return null;
	return isCourseStudioShellNode(onlyChild) ? onlyChild : null;
}

function isCourseStudioShellNode(node: WidgetNode): node is ComponentNode {
	return node.kind === "component" && node.type === "CourseStudioShell";
}

function normalizeMeta(meta: Record<string, unknown> | undefined): WidgetPageMeta {
	if (!meta) return {};
	return {
		shell: isPageShellMode(meta.shell) ? meta.shell : undefined,
		activeNavItemId: typeof meta.activeNavItemId === "string" ? meta.activeNavItemId : undefined,
		navItems: Array.isArray(meta.navItems) ? meta.navItems.filter(isAppNavItemSpec) : undefined,
		maxWidth: isPageMaxWidth(meta.maxWidth) ? meta.maxWidth : undefined,
	};
}

function isPageShellMode(value: unknown): value is PageShellMode {
	return value === "auto" || value === "none" || value === "app";
}

function isPageMaxWidth(value: unknown): value is PageMaxWidth {
	return value === "none" || value === "content" || value === "wide";
}

function isAppNavItemSpec(value: unknown): value is AppNavItemSpec {
	if (!value || typeof value !== "object") return false;
	const candidate = value as Partial<AppNavItemSpec>;
	return typeof candidate.id === "string" && candidate.label !== undefined;
}

function renderPageNavigationLabel(
	item: PageNavigationItemSpec,
	onAction: (action: ActionSpec, context: WidgetActionContext) => void,
): ReactNode {
	return (
		<>
			{item.icon ? (
				<span aria-hidden="true">{renderRenderableValue(item.icon, onAction)}</span>
			) : null}
			{renderRenderableValue(item.label, onAction)}
			{item.badge ? <span>{renderRenderableValue(item.badge, onAction)}</span> : null}
		</>
	);
}

function renderRenderableValue(
	value: RenderableValue | undefined,
	onAction?: (action: ActionSpec, context: WidgetActionContext) => void,
): ReactNode {
	if (value && typeof value === "object" && "kind" in value) {
		return (
			<WidgetRenderer
				node={value as WidgetNode}
				registry={defaultWidgetRegistry}
				onAction={onAction}
			/>
		);
	}
	return value == null ? null : String(value);
}

function widgetNodeKey(node: WidgetNode, fallback: number): string | number {
	if (node.kind === "component") {
		const props = (node.props ?? {}) as { id?: unknown; key?: unknown };
		if (typeof props.key === "string" || typeof props.key === "number") return props.key;
		if (typeof props.id === "string" || typeof props.id === "number") return props.id;
		return `${node.type}-${fallback}`;
	}
	if (node.kind === "element") {
		const attrs = (node.attrs ?? {}) as { id?: unknown; key?: unknown };
		if (typeof attrs.key === "string" || typeof attrs.key === "number") return attrs.key;
		if (typeof attrs.id === "string" || typeof attrs.id === "number") return attrs.id;
		return `${node.tag}-${fallback}`;
	}
	return fallback;
}

function contentClassName(maxWidth: PageMaxWidth | undefined): string {
	const width = maxWidth ?? "wide";
	return ["rag-evaluation-site-content", `rag-evaluation-site-content--${width}`].join(" ");
}

function navigateToPage(pageId: string): void {
	if (typeof window === "undefined") return;
	const url = new URL(window.location.href);
	url.pathname = `/pages/${encodeURIComponent(pageId)}`;
	url.searchParams.delete("page");
	window.history.pushState({}, "", url.toString());
	window.dispatchEvent(new PopStateEvent("popstate"));
}

function readSearchFromLocation(_locationVersion: number): string {
	if (typeof window === "undefined") return "";
	const url = new URL(window.location.href);
	const parts = url.pathname.split("/").filter(Boolean);
	if (parts[0] === "print" && parts[1] === "handouts" && parts[2] && !url.searchParams.has("doc")) {
		url.searchParams.set("doc", parts[2]);
		return url.search;
	}
	if (
		parts[0] === "present" &&
		parts[1] === "slides" &&
		parts[2] &&
		!url.searchParams.has("slide")
	) {
		url.searchParams.set("slide", parts[2]);
		return url.search;
	}
	return window.location.search || "";
}

function readPageIdFromLocation(defaultPageId: string): string {
	if (typeof window === "undefined") return defaultPageId;
	const url = new URL(window.location.href);
	const parts = url.pathname.split("/").filter(Boolean);
	if (parts[0] === "pages" && parts[1]) return parts[1];
	const queryPage = url.searchParams.get("page");
	if (queryPage) return queryPage;
	if (parts[0] === "print" && parts[1] === "handouts") return "print-handout";
	if (parts[0] === "present" && parts[1] === "slides") return "present-slide";
	return defaultPageId;
}
