import { type ReactNode, useEffect, useState } from "react";
import { ErrorCallout } from "../components/atoms";
import { Caption } from "../components/foundation";
import { AppShell, Panel } from "../components/layout";
import { AppNav } from "../components/molecules";
import { CourseStudioShell } from "../components/organisms";
import { useWidgetPage, type WidgetPageResponse } from "../hooks/useWidgetPage";
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

	async function handleAction(action: ActionSpec, context: WidgetActionContext): Promise<void> {
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
		if (!response.ok) {
			throw new Error(`Widget action failed: ${response.status} ${response.statusText}`);
		}
		const result = (await response.json()) as { refresh?: boolean; toast?: string };
		if (result.toast && typeof window !== "undefined") {
			window.dispatchEvent(new CustomEvent("widget:toast", { detail: result }));
		}
		if (result.refresh) refresh();
	}

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
			{renderPage(page, pageId, (action, context) => {
				void handleAction(action, context);
			})}
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
	onAction: (action: ActionSpec, context: WidgetActionContext) => void,
): ReactNode {
	const meta = normalizeMeta(page.meta);
	const rootClassName = [
		"rag-evaluation-site-root",
		shouldUseDefaultShell(page, meta)
			? "rag-evaluation-site-root--shell"
			: "rag-evaluation-site-root--raw",
	]
		.filter(Boolean)
		.join(" ");

	const courseShellNode = findRootCourseStudioShellNode(page.root);
	if (courseShellNode) {
		return (
			<div
				className={rootClassName}
				data-rag-page="RagEvaluationSiteApp"
				data-page-id={page.id}
				data-rag-shell="course-studio"
			>
				{renderCourseStudioShellPage(courseShellNode, onAction)}
			</div>
		);
	}

	const renderedRoot = (
		<WidgetRenderer node={page.root} registry={defaultWidgetRegistry} onAction={onAction} />
	);

	if (!shouldUseDefaultShell(page, meta)) {
		return (
			<div
				className={rootClassName}
				data-rag-page="RagEvaluationSiteApp"
				data-page-id={page.id}
				data-rag-shell="none"
			>
				{renderedRoot}
			</div>
		);
	}

	const navItems = meta.navItems?.length ? meta.navItems : DEFAULT_NAV_ITEMS;
	const activeNavItemId = meta.activeNavItemId ?? pageId;

	return (
		<div
			className={rootClassName}
			data-rag-page="RagEvaluationSiteApp"
			data-page-id={page.id}
			data-rag-shell="default"
		>
			<AppShell
				className="rag-evaluation-site-shell"
				header={
					<AppNav
						brand={
							<span className="rag-evaluation-site-brand">
								{page.title || "RAG Evaluation Site"}
							</span>
						}
						items={navItems.map((item) => ({ id: item.id, label: renderNavLabel(item.label) }))}
						activeItemId={activeNavItemId}
						onItemSelect={(itemId) => {
							const item = navItems.find((candidate) => candidate.id === itemId);
							if (item?.action) {
								dispatchWidgetAction(
									item.action,
									{ value: itemId, componentType: "AppNav" },
									onAction,
								);
								return;
							}
							navigateToPage(itemId);
						}}
					/>
				}
			>
				<div className={contentClassName(meta.maxWidth)} data-rag-layout="PageContent">
					{renderedRoot}
				</div>
			</AppShell>
		</div>
	);
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

function shouldUseDefaultShell(page: WidgetPageResponse, meta: WidgetPageMeta): boolean {
	if (meta.shell === "none") return false;
	if (isAppShellNode(page.root)) return false;
	if (findRootCourseStudioShellNode(page.root)) return false;
	return true;
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

function renderNavLabel(label: RenderableValue): ReactNode {
	return renderRenderableValue(label);
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
	const queryPage = url.searchParams.get("page");
	if (queryPage) return queryPage;
	const parts = url.pathname.split("/").filter(Boolean);
	if (parts[0] === "pages" && parts[1]) return parts[1];
	if (parts[0] === "print" && parts[1] === "handouts") return "print-handout";
	if (parts[0] === "present" && parts[1] === "slides") return "present-slide";
	return defaultPageId;
}
