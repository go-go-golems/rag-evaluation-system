import type { ReactNode } from 'react';
import { ErrorCallout } from '../components/atoms';
import { Caption } from '../components/foundation';
import { AppShell, Panel } from '../components/layout';
import { AppNav } from '../components/molecules';
import { WidgetRenderer } from '../widgets/WidgetRenderer';
import type { ActionSpec, AppNavItemSpec, RenderableValue, WidgetNode } from '../widgets/ir';
import { dispatchWidgetAction, type WidgetActionContext } from '../widgets/actions';
import { useWidgetPage, type WidgetPageResponse } from '../hooks/useWidgetPage';
import './app.css';

export interface RagEvaluationSiteAppProps {
  apiBase?: string;
  defaultPageId?: string;
}

type PageShellMode = 'auto' | 'none' | 'app';
type PageMaxWidth = 'none' | 'content' | 'wide';

interface WidgetPageMeta {
  shell?: PageShellMode;
  activeNavItemId?: string;
  navItems?: AppNavItemSpec[];
  maxWidth?: PageMaxWidth;
}

const DEFAULT_NAV_ITEMS: AppNavItemSpec[] = [
  { id: 'index', label: 'Overview' },
  { id: 'demo', label: 'Demo' },
  { id: 'actions', label: 'Actions' },
];

export function RagEvaluationSiteApp({ apiBase = '/api/widget', defaultPageId = 'index' }: RagEvaluationSiteAppProps) {
  const pageId = readPageIdFromLocation(defaultPageId);
  const cleanApiBase = apiBase.replace(/\/$/, '');
  const { page, loading, error, refresh } = useWidgetPage(`${cleanApiBase}/pages/${encodeURIComponent(pageId)}`);

  async function handleAction(action: ActionSpec, context: WidgetActionContext): Promise<void> {
    if (action.kind !== 'server') {
      dispatchWidgetAction(action, context);
      return;
    }
    const response = await fetch(`${cleanApiBase}/actions/${encodeURIComponent(action.name)}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ payload: action.payload ?? {}, context }),
    });
    if (!response.ok) {
      throw new Error(`Widget action failed: ${response.status} ${response.statusText}`);
    }
    const result = await response.json() as { refresh?: boolean; toast?: string };
    if (result.toast && typeof window !== 'undefined') {
      window.dispatchEvent(new CustomEvent('widget:toast', { detail: result }));
    }
    if (result.refresh) refresh();
  }

  if (loading) {
    return <Panel className="rag-evaluation-site-state" title="RAG Evaluation Site" density="condensed"><Caption>Loading Widget IR…</Caption></Panel>;
  }

  if (error) {
    return <ErrorCallout className="rag-evaluation-site-state">Failed to load Widget IR page: {error.message}</ErrorCallout>;
  }

  if (!page) {
    return <Panel className="rag-evaluation-site-state" title="RAG Evaluation Site" density="condensed"><Caption>No Widget IR returned.</Caption></Panel>;
  }

  return renderPage(page, pageId, (action, context) => { void handleAction(action, context); });
}

function renderPage(page: WidgetPageResponse, pageId: string, onAction: (action: ActionSpec, context: WidgetActionContext) => void): ReactNode {
  const meta = normalizeMeta(page.meta);
  const renderedRoot = <WidgetRenderer node={page.root} onAction={onAction} />;
  const rootClassName = ['rag-evaluation-site-root', shouldUseDefaultShell(page, meta) ? 'rag-evaluation-site-root--shell' : 'rag-evaluation-site-root--raw']
    .filter(Boolean)
    .join(' ');

  if (!shouldUseDefaultShell(page, meta)) {
    return (
      <div className={rootClassName} data-rag-page="RagEvaluationSiteApp" data-page-id={page.id} data-rag-shell="none">
        {renderedRoot}
      </div>
    );
  }

  const navItems = meta.navItems?.length ? meta.navItems : DEFAULT_NAV_ITEMS;
  const activeNavItemId = meta.activeNavItemId ?? pageId;

  return (
    <div className={rootClassName} data-rag-page="RagEvaluationSiteApp" data-page-id={page.id} data-rag-shell="default">
      <AppShell
        className="rag-evaluation-site-shell"
        header={(
          <AppNav
            brand={<span className="rag-evaluation-site-brand">{page.title || 'RAG Evaluation Site'}</span>}
            items={navItems.map((item) => ({ id: item.id, label: renderNavLabel(item.label) }))}
            activeItemId={activeNavItemId}
            onItemSelect={(itemId) => {
              const item = navItems.find((candidate) => candidate.id === itemId);
              if (item?.action) {
                onAction(item.action, { value: itemId, componentType: 'AppNav' });
                return;
              }
              navigateToPage(itemId);
            }}
          />
        )}
      >
        <div className={contentClassName(meta.maxWidth)} data-rag-layout="PageContent">
          {renderedRoot}
        </div>
      </AppShell>
    </div>
  );
}

function shouldUseDefaultShell(page: WidgetPageResponse, meta: WidgetPageMeta): boolean {
  if (meta.shell === 'none') return false;
  if (isAppShellNode(page.root)) return false;
  return true;
}

function isAppShellNode(node: WidgetNode): boolean {
  return node.kind === 'component' && node.type === 'AppShell';
}

function normalizeMeta(meta: Record<string, unknown> | undefined): WidgetPageMeta {
  if (!meta) return {};
  return {
    shell: isPageShellMode(meta.shell) ? meta.shell : undefined,
    activeNavItemId: typeof meta.activeNavItemId === 'string' ? meta.activeNavItemId : undefined,
    navItems: Array.isArray(meta.navItems) ? meta.navItems.filter(isAppNavItemSpec) : undefined,
    maxWidth: isPageMaxWidth(meta.maxWidth) ? meta.maxWidth : undefined,
  };
}

function isPageShellMode(value: unknown): value is PageShellMode {
  return value === 'auto' || value === 'none' || value === 'app';
}

function isPageMaxWidth(value: unknown): value is PageMaxWidth {
  return value === 'none' || value === 'content' || value === 'wide';
}

function isAppNavItemSpec(value: unknown): value is AppNavItemSpec {
  if (!value || typeof value !== 'object') return false;
  const candidate = value as Partial<AppNavItemSpec>;
  return typeof candidate.id === 'string' && candidate.label !== undefined;
}

function renderNavLabel(label: RenderableValue): ReactNode {
  if (label && typeof label === 'object' && 'kind' in label) {
    return <WidgetRenderer node={label as WidgetNode} />;
  }
  return String(label ?? '');
}

function contentClassName(maxWidth: PageMaxWidth | undefined): string {
  const width = maxWidth ?? 'wide';
  return ['rag-evaluation-site-content', `rag-evaluation-site-content--${width}`].join(' ');
}

function navigateToPage(pageId: string): void {
  if (typeof window === 'undefined') return;
  const url = new URL(window.location.href);
  url.pathname = `/pages/${encodeURIComponent(pageId)}`;
  url.searchParams.delete('page');
  window.location.assign(url.toString());
}

function readPageIdFromLocation(defaultPageId: string): string {
  if (typeof window === 'undefined') return defaultPageId;
  const url = new URL(window.location.href);
  const queryPage = url.searchParams.get('page');
  if (queryPage) return queryPage;
  const parts = url.pathname.split('/').filter(Boolean);
  if (parts[0] === 'pages' && parts[1]) return parts[1];
  return defaultPageId;
}
